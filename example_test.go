package graph_test

import (
	"fmt"

	"github.com/trelore/go-graph"
)

func ExampleGraph() {
	g := graph.New[int, int](false)
	id := g.AddNode(graph.NodeID(1), 1)
	id2 := g.AddNode(graph.NodeID(2), 2)
	g.AddEdge(id, id2, 1)
	g.AddEdge(id2, id, 1)
	fmt.Println(g.Node(id))
	fmt.Println(g.Node(id2))
	// Output:
	// 1 <nil>
	// 2 <nil>
}

func ExampleImportG6() {
	g := graph.ImportG6([]byte("DQc"))
	dot := g.ExportDot()
	fmt.Println(string(dot))
	// Output:
	// digraph {
	// 	3 -> 4 [label=1];
	// 	0 -> 2 [label=1];
	// 	0 -> 4 [label=1];
	// 	1 -> 3 [label=1];
	// }
}
