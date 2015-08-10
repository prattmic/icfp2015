package main

import (
	"log"
)

var (
	cary = len(defaultPhrases)
)

type Chant struct {
	id       int
	score    float64
	dead     bool // game ending leaf
	d        Commands
	children []*Chant
	game     *Game
}

func (n *Chant) BestMove() *Chant {
	if n.IsLeaf() {
		panic("Can't find best move on leaf chant")
	}
	best := n.children[0]
	hiscore := n.children[0].score
	for _, c := range n.children {
		if !c.IsDead() && c.score > hiscore {
			hiscore = c.score
			best = c
		}
	}
	if best == nil {
		panic("nil best")
	}
	return best
}

func (n *Chant) IsLeaf() bool {
	return n.children == nil
}

func (n *Chant) IsDead() bool {
	return n.dead
}

type ChanterDescender struct {
	root *Chant
}

func (t *ChanterDescender) Next() (Commands, error) {
	if t.root.IsDead() {
		return nil, errNoMoves
	}
	if t.root.IsLeaf() {
		return nil, errTreeExhausted
	}
	next := t.root.BestMove()
	t.root = next
	return next.d, nil
}

func NewChanterDescender(g *Game) *ChanterDescender {
	depth := 4
	height := 0
	root := &Chant{} // fake root
	root.children = make([]*Chant, cary)
	for i := range root.children {
		root.children[i] = BuildScoreChanter(normalizedCommands[i], g, depth-1, height+1)
	}
	root.score = root.BestMove().score
	return &ChanterDescender{root: root}
}

func BuildScoreChanter(d Commands, g *Game, depth int, height int) *Chant {
	n := &Chant{
		d:  d,
		id: uniqueId,
	}
	uniqueId++
	n.game = g.Fork()
	// XXX spinner update
	for _, c := range d {
		_, done, err := n.game.Update(c)
		if err != nil || done {
			n.score = -1000000000
			n.dead = true
			return n
		}
	}
	n.score = n.game.Score()
	if depth == 0 {
		return n
	}
	n.children = make([]*Chant, cary)
	for i := range n.children {
		n.children[i] = BuildScoreChanter(normalizedCommands[i], n.game, depth-1, height+1)
	}
	n.score += n.BestMove().score
	return n
}

type ChanterAI struct {
	index   int
	game    *Game
	current Commands
}

func NewChanterAI(g *Game) AI {
	return &ChanterAI{index: 0, game: g}
}

// Game returns the Game used by the AI.
// It may change after calls to Next().
func (ai *ChanterAI) Game() *Game {
	return ai.game
}

// Next steps the AI one step, returning true if the game is
// complete, or an error if the game cannot continue.
func (ai *ChanterAI) Next() (bool, error) {
	if ai.current == nil {
		t := NewChanterDescender(ai.game)
		current, err := t.Next()
		if err == errNoMoves {
			return false, err // no possible moves, we are stuck!
		}
		if err == errTreeExhausted {
			return false, err // need to build a new tree
		}
		ai.current = current
		log.Printf("Chanting word %v", current.String())
	}

	c := ai.current[0]
	if len(ai.current) > 1 {
		ai.current = ai.current[1:len(ai.current)]
	} else {
		ai.current = nil
	}

	locked, done, err := ai.game.Update(c)
	log.Printf("Update(%s) -> locked %v done %v, %v", c, locked, done, err)
	return done, err
}
