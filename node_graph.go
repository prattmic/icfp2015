// +build graph

package main

import (
	"fmt"
	"strconv"

	"github.com/awalterschulze/gographviz"
)

func (n *Node) addToGraph(g *gographviz.Graph, parent string) {
	id := strconv.Itoa(n.id)

	label := fmt.Sprintf(`"Score %d\nDirection: %s"`, int(n.score), n.d)

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
