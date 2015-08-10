package main

import ()

type Node struct {
	id       int
	score    float64
	dead     bool // Game ending leaf.
	d        Direction
	children []*Node
	game     *Game
}

var (
	dirs = []Direction{SE, SW, E, W, CW, CCW}
	nary = len(dirs)

	depthWeight = 100.0

	uniqueId int
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

func (n *Node) IsLeaf() bool {
	return n.children == nil
}

func (n *Node) IsDead() bool {
	return n.dead
}

func BuildScoreTree(d Direction, g *Game, depth int, height int) *Node {
	n := &Node{d: d, id: uniqueId}
	uniqueId++

	c := directionToCommands[d][0]

	n.game = g.Fork()
	locked, done, err := n.game.Update(c)
	if err != nil {
		// NO POINTS FOR U
		n.score = -1000000000
		n.dead = true
		return n
	}

	if done {
		// NO POINTS FOR U
		n.score = -1000000000
		n.dead = true
		return n
	}

	if depth == 0 {
		midY := 0.0
		for _, c := range n.game.currUnit.Members {
			midY += float64(c.Y)
		}
		midY /= float64(len(n.game.currUnit.Members))

		n.score = n.game.Score() + depthWeight*(midY+float64(height))
		return n
	}

	n.children = make([]*Node, nary)
	for i := range n.children {
		n.children[i] = BuildScoreTree(dirs[i], n.game, depth-1, height+1)
	}

	if locked {
		n.score = -1000
	}

	n.score += n.BestMove().score

	return n
}
