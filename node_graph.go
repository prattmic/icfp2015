// +build graph

package main

import (
	"fmt"
	"strconv"

	"github.com/awalterschulze/gographviz"
)

func (n *Node) addToGraph(g *gographviz.Graph, parent string) {
	id := strconv.Itoa(n.id)

	label := "root"
	if n.game != nil {
		label = fmt.Sprintf(`"Score: %d
Weights: %v
Direction: %s
Unit: %+v"`, int(n.score), n.weights, n.d, n.game.currUnit)
	}

	attrs := map[string]string{
		"label": label,
	}

	g.AddNode(parent, id, attrs)
	g.AddEdge(parent, id, true, nil)

	for _, c := range n.children {
		c.addToGraph(g, id)
	}
}

func (n *Node) Graph() string {
	g := gographviz.NewGraph()
	g.SetName("Tree")
	g.SetDir(true)

	g.AddNode("Tree", "root", nil)

	n.addToGraph(g, "root")
	return g.String()
}
