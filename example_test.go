package graph_test

import (
	"fmt"
	"sort"

	"github.com/trelore/go-graph"
)

func ExampleGraph() {
	g := graph.New[int, int](false)
	id := g.AddNode(graph.NodeID(1), 1)
	id2 := g.AddNode(graph.NodeID(2), 2)
	g.AddEdge(id, id2, 1)
	g.AddEdge(id2, id, 1)
	n1, _ := g.Node(id)
	n2, _ := g.Node(id2)
	fmt.Println(n1.ID)
	fmt.Println(n2.ID)
	// Output:
	// 1
	// 2
}

type byID []graph.NodeID

func (b byID) Len() int           { return len(b) }
func (b byID) Less(i, j int) bool { return b[i] < b[j] }
func (b byID) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func ExampleImportG6() {
	g := graph.ImportG6([]byte("DQc"))
	nodeIDs := g.NodeIDs()
	sort.Sort(byID(nodeIDs))
	fmt.Println(nodeIDs)
	n, _ := g.Node(0)
	ne := n.Neighbours()
	for _, n := range ne {
		fmt.Println(n.ID)
	}
	// Output:
	// [0 1 2 3 4]
	// 2
	// 4
}
