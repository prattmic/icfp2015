package main

import (
	_ "fmt"
)

type Node struct {
	score    float64
	dead     bool // Game ending leaf.
	d        Direction
	children []*Node
}

var (
	dirs = []Direction{SE, SW, E, W, CW, CCW}
	nary = len(dirs)

	depthWeight = 100.0
)

// TODO(mgyenik) make this and d Diection in Node a Command
// and look for phrases of power.
// This returns the best move to make from this node, it should
// not be called on a leaf node.
func (n *Node) BestMove() *Node {
	// For debugging
	if n.IsLeaf() {
		panic("Can't find best move on leaf node")
	}

	best := n.children[0]
	hiscore := 0.0
	for _, c := range n.children {
		if c.score > hiscore {
			hiscore = c.score
			best = c
		}
	}

	return best
}

func (n *Node) IsLeaf() bool {
	return n.children == nil
}

func (n *Node) IsDead() bool {
	return n.dead
}

func BuildScoreTree(d Direction, g *Game, depth int, height int) *Node {
	n := &Node{d: d}
	thisgame := g.Fork()
	c := directionToCommands[d][0]
	done, err := thisgame.Update(c)
	if err != nil {
		// NO POINTS FOR U
		n.score = 0
		n.dead = true
		return n
	}

	if done {
		// NO POINTS FOR U
		n.score = 0
		n.dead = true
		return n
	}

	if depth == 0 {
		midY := 0.0
		for _, c := range g.currUnit.Members {
			midY += float64(c.Y)
		}
		midY /= float64(len(g.currUnit.Members))

		n.score = g.Score() + depthWeight*(midY+100.0*float64(height))
		n.dead = true
		return n
	}

	n.children = make([]*Node, nary)
	for i := range n.children {
		n.children[i] = BuildScoreTree(dirs[i], thisgame, depth-1, height+1)
	}

	n.score = n.BestMove().score

	return n
}
