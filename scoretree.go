package main

type Node struct {
	score    int
	d        Direction
	children []*Node
}

var (
	dirs = []Direction{E, SE, SW, W, CW, CCW}
	nary = len(dirs)
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

	var best *Node
	hiscore := 0
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

// TODO(myenik) We just pick the best of the first moves and discard the rest...
func BuildScoreTree(d Direction, g *Game, depth int) *Node {
	n := Node{d: d}
	thisgame := g.Fork()
	done, err := thisgame.Update(d)
	if err != nil {
		// NO POINTS FOR U
		n.score = 0
		return n
	}

	if done || depth == 0 {
		n.score = g.Score()
		return n
	}

	n.children = make([]*Node, nary)
	for i := range n.children {
		n.children[i] = BuildScoreTree(dirs[i], thisgame, depth-1)
	}

	n.score = n.BestMove().score
	return n
}
