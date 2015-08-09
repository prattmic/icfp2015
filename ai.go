package main

import (
	"errors"
	"log"
	"math/rand"
)

var (
	errTreeExhausted = errors.New("Tree not deep enough!")
	errNoMoves       = errors.New("Out of moves!")
)

// MovePermuter provides each of the possible moves in a random order.
// TODO)prattmic): Better name...
type MovePermuter struct {
	possible []Direction
}

// NewMovePermuter builds a MovePermuter with all of the legal game moves.
func NewMovePermuter() MovePermuter {
	return MovePermuter{
		// All the possible moves.
		possible: []Direction{E, W, SE, SW, CCW, CW},
	}
}

// Next returns the next move to attempt, or an error if none remain.
func (m *MovePermuter) Next() (Command, error) {
	if len(m.possible) == 0 {
		return 0, errors.New("No more possible moves")
	}

	i := rand.Intn(len(m.possible))
	move := m.possible[i]
	m.possible = append(m.possible[:i], m.possible[i+1:]...)

	// Pick a random command for this direction
	c := directionToCommands[move][rand.Intn(len(directionToCommands[move]))]

	return c, nil
}

type TreeDescender struct {
	root *Node
}

func (t *TreeDescender) Next() (Command, error) {
	if t.root.IsDead() {
		return 0, errNoMoves
	}

	if t.root.IsLeaf() {
		return 0, errTreeExhausted
	}

	next := t.root.BestMove()
	t.root = next

	// TODO(myenik) First command is *best* command!
	return directionToCommands[next.d][0], nil
}

func NewTreeDescender(g *Game) *TreeDescender {
	// TODO(myenik) paramterize depth
	depth := 4

	// This bullshit is to avoid an error on start.
	var startNode *Node
	bestStart := 0.0
	for _, d := range dirs {
		thisgame := g.Fork()
		node := BuildScoreTree(d, thisgame, depth, 0)
		if node.score > bestStart {
			bestStart = node.score
			startNode = node
		}
	}

	return &TreeDescender{root: startNode}
}

// AI wins the game!
type AI struct {
	Game *Game
}

// NewAI builds a new AI.
func NewAI(g *Game) AI {
	return AI{Game: g}
}

type aiResult struct {
	command Command
	game *Game
	score float64
	done bool
	err error
}

func (a *AI) Next() (bool, error) {
	var coms Commands
	for _, d := range []Direction{E, W, SE, SW, CCW, CW} {
		coms = append(coms, directionToCommands[d][0])
	}

	// Try each command, use the highest score.
	ch := make(chan aiResult, len(coms))

	runner := func(c Command) {
		g := a.Game.Fork()
		done, err := g.Update(c)

		ch <- aiResult{
			command: c,
			game: g,
			score: g.Score(),
			done: done,
			err: err,
		}
	}

	// Compute moves
	for _, c := range coms {
		go runner(c)
	}

	// Collect results
	var ret []aiResult
	for range coms {
		ret = append(ret, <-ch)
	}

	bestScore := -1.0
	var best aiResult
	for _, r := range ret {
		if r.err != nil {
			log.Printf("skipping: %+v", r)
			continue
		}

		if r.score > bestScore {
			log.Printf("best by score: %+v", r)
			bestScore = r.score
			best = r
		}

		if best.done && !r.done {
			log.Printf("best by not done: %+v", r)
			bestScore = r.score
			best = r
		}
	}

	log.Printf("results: %+v", ret)

	if best.game == nil {
		// No best move! Just use first.
		best = ret[0]
		log.Printf("best by best is nil: %+v", best)
	}

	log.Printf("best: %+v", best)

	a.Game = best.game
	return best.done, best.err
}

type TreeAI struct {
	Game *Game
}

func NewTreeAI(g *Game) *TreeAI {
	return &TreeAI{
		Game: g,
	}
}

// Next steps the game, returning true when the game is done.
func (a *TreeAI) Next() (bool, error) {
	t := NewTreeDescender(a.Game)
	c, err := t.Next()
	if err == errNoMoves {
		// No possible moves, we are stuck!
		return false, err
	}

	if err == errTreeExhausted {
		// Need to build new tree.
		return false, err
	}

	done, err := a.Game.Update(c)
	log.Printf("Update(%s) -> %v, %v", c, done, err)
	return done, err
}
