package graph

import "sync"

// algo to convert a graph to a tree decomposition
// https://en.wikipedia.org/wiki/Tree_decomposition

type TreeDecomposition[T any] struct {
	sync.RWMutex

	nodes map[NodeID]*tdNode[[]T]
	width int
}

type tdNode[T any] struct {
	Value    T
	id       NodeID
	edgesOut map[NodeID]*edge[T, int] // hardcoded weights to 1 as they're meaningless in TD
	edgesIn  map[NodeID]*edge[T, int] // hardcoded weights to 1 as they're meaningless in TD
}

func (g *Graph[T, N]) ComputeTreeDecomposition() *TreeDecomposition[T] {
	// todo
	return nil
}
