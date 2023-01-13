package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkMaximalIndependentSet(b *testing.B) {
	g := ImportG6([]byte("DQc"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.MaximumIndependentSet()
	}
}

func TestMaximalIndependentSet(t *testing.T) {
	tests := []struct {
		name    string
		in      []byte
		want    []NodeID
		wantErr bool
	}{
		{
			name: "empty",
			in:   []byte(``),
			want: []NodeID{},
		},
		{
			name: "simple",
			in:   []byte(`DQc`),
			want: []NodeID{1, 2, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ImportG6(tt.in).MaximumIndependentSet()
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func Test_newGraphWithoutNode(t *testing.T) {
	tests := []struct {
		name string
		g    *Graph[int, int]
		n    NodeID
		want *Graph[int, int]
	}{
		{
			name: "simple clique",
			g: func() *Graph[int, int] {
				g := New[int, int](false)
				g.AddNode(NodeID(0), 0)
				g.AddNode(NodeID(1), 1)
				g.AddNode(NodeID(2), 2)
				g.AddEdge(0, 1, 1)
				g.AddEdge(1, 2, 1)
				g.AddEdge(2, 0, 1)
				return g
			}(),
			n: 0,
			want: func() *Graph[int, int] {
				g := New[int, int](false)
				g.AddNode(NodeID(1), 1)
				g.AddNode(NodeID(2), 2)
				g.AddEdge(1, 2, 1)
				return g
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newGraphWithoutNode(tt.g, tt.n)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_newGraphWithoutNodeAndNeighbours(t *testing.T) {
	tests := []struct {
		name   string
		g      *Graph[int, int]
		nodeID NodeID
		want   *Graph[int, int]
	}{
		{
			name: "simple graph",
			g: func() *Graph[int, int] {
				g := New[int, int](false)
				g.AddNode(NodeID(0), 0)
				g.AddNode(NodeID(1), 1)
				g.AddNode(NodeID(2), 2)
				g.AddEdge(0, 1, 1)
				g.AddEdge(1, 2, 1)
				return g
			}(),
			nodeID: 0,
			want: func() *Graph[int, int] {
				g := New[int, int](false)
				g.AddNode(NodeID(2), 2)
				return g
			}(),
		},
		{
			name: "simple graph - remove last node",
			g: func() *Graph[int, int] {
				g := New[int, int](false)
				g.AddNode(NodeID(0), 0)
				g.AddNode(NodeID(1), 1)
				g.AddNode(NodeID(2), 2)
				g.AddEdge(0, 1, 1)
				g.AddEdge(1, 2, 1)
				return g
			}(),
			nodeID: 2,
			want: func() *Graph[int, int] {
				g := New[int, int](false)
				g.AddNode(NodeID(0), 0)
				return g
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newGraphWithoutNodeAndNeighbours(tt.g, tt.nodeID)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_maximumIndependentSet(t *testing.T) {
	tests := []struct {
		name    string
		g       *Graph[int, int]
		wantLen int
	}{
		{
			name: "single node",
			g: func() *Graph[int, int] {
				g := New[int, int](false)
				g.AddNode(NodeID(0), 0)
				return g
			}(),
			wantLen: 1,
		},
		{
			name: "2 node clique",
			g: func() *Graph[int, int] {
				g := New[int, int](false)
				g.AddNode(NodeID(0), 0)
				g.AddNode(NodeID(1), 1)
				g.AddEdge(0, 1, 1)
				return g
			}(),
			wantLen: 1,
		},
		{
			name: "3 node clique",
			g: func() *Graph[int, int] {
				g := New[int, int](false)
				g.AddNode(NodeID(0), 0)
				g.AddNode(NodeID(1), 1)
				g.AddNode(NodeID(2), 2)
				g.AddEdge(0, 1, 1)
				g.AddEdge(1, 2, 1)
				g.AddEdge(2, 0, 1)
				return g
			}(),
			wantLen: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maximumIndependentSet(tt.g)
			assert.Equal(t, tt.wantLen, len(got))
		})
	}
}
