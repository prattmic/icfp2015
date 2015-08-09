package main

import ()

type Node struct {
	id       int
	score    float64
	dead     bool // Game ending leaf.
	d        Direction
	children []*Node
	game     *Game
	weights  map[string]float64
}

var (
	dirs = []Direction{SE, SW, E, W, CW, CCW}
	nary = len(dirs)

	depthWeight = 10.0

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
	n := &Node{
		d:       d,
		id:      uniqueId,
		weights: make(map[string]float64),
	}
	uniqueId++

	c := directionToCommands[d][0]

	n.game = g.Fork()

	unit := n.game.currUnit.DeepCopy()
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

	midY := 0.0
	for _, c := range n.game.currUnit.Members {
		midY += float64(c.Y)
	}
	midY /= float64(len(n.game.currUnit.Members))

	n.weights["gameScore"] = n.game.Score()
	n.weights["depth"] = depthWeight * (midY + float64(height))

	n.score = n.weights["gameScore"] + n.weights["depth"]

	if depth == 0 {
		return n
	}

	n.children = make([]*Node, nary)
	for i := range n.children {
		n.children[i] = BuildScoreTree(dirs[i], n.game, depth-1, height+1)
	}

	if locked {
		if n.game.B.GapBelowAny(unit) {
			n.weights["locked"] = -10000
		} else {
			n.weights["locked"] = 10000
		}
		n.score += n.weights["locked"]
	}

	n.weights["bestMove"] = n.BestMove().score
	n.score += n.weights["bestMove"]

	return n
}
