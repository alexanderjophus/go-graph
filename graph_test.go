package graph_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trelore/go-graph"
)

func ExampleGraph() {
	g := graph.New[int, int](false)
	id := g.AddNode(1)
	id2 := g.AddNode(2)
	g.AddEdge(id, id2, 1)
	g.AddEdge(id2, id, 1)
	fmt.Println(g.Node(id))
	fmt.Println(g.Node(id2))
	// Output:
	// 1 <nil>
	// 2 <nil>
}

func TestImportG6(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		want *graph.Graph[int, int]
	}{
		{
			name: "empty",
			b:    []byte(``),
			want: graph.New[int, int](false),
		},
		{
			name: "simple",
			b:    []byte(`DQc`),
			want: func() *graph.Graph[int, int] {
				g := graph.New[int, int](false)
				g.AddNode(0)
				g.AddNode(1)
				g.AddNode(2)
				g.AddNode(3)
				g.AddNode(4)
				g.AddEdge(0, 2, 1)
				g.AddEdge(0, 4, 1)
				g.AddEdge(1, 3, 1)
				g.AddEdge(3, 4, 1)
				return g
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := graph.ImportG6(tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}
