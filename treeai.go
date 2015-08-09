package main

import (
	"errors"
	"log"
)

var (
	errTreeExhausted = errors.New("Tree not deep enough!")
	errNoMoves       = errors.New("Out of moves!")
)

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
	height := 0

	// Fake root, there is no direction here.
	root := &Node{}

	root.children = make([]*Node, nary)
	for i := range root.children {
		root.children[i] = BuildScoreTree(dirs[i], g, depth-1, height+1)
	}

	root.score = root.BestMove().score

	return &TreeDescender{root: root}
}

// TreeAI implements AI.
type TreeAI struct {
	game *Game
}

func NewTreeAI(g *Game) AI {
	return &TreeAI{
		game: g,
	}
}

func (a *TreeAI) Game() *Game {
	return a.game
}

// Next steps the game, returning true when the game is done.
func (a *TreeAI) Next() (bool, error) {
	t := NewTreeDescender(a.game)
	c, err := t.Next()
	if err == errNoMoves {
		// No possible moves, we are stuck!
		return false, err
	}

	if err == errTreeExhausted {
		// Need to build new tree.
		return false, err
	}

	_, done, err := a.game.Update(c)
	log.Printf("Update(%s) -> %v, %v", c, done, err)
	return done, err
}