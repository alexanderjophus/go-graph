package graph_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trelore/go-graph"
)

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
				g.AddNode(graph.NodeID(0), 0)
				g.AddNode(graph.NodeID(1), 1)
				g.AddNode(graph.NodeID(2), 2)
				g.AddNode(graph.NodeID(3), 3)
				g.AddNode(graph.NodeID(4), 4)
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
