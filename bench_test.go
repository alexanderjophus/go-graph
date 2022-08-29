package graph_test

import (
	"testing"

	"github.com/trelore/go-graph"
)

func BenchNaiveMaximalIndependentSet(b *testing.B) {
	g := graph.New[int, int](false)
	for i := 0; i < b.N; i++ {
		g.MaximumIndependentSet()
	}
}
