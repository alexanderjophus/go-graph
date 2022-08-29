package graph

import "fmt"

// Greedy implementation of the maximum independent set algorithm.
func (g *Graph[T, N]) MaximumIndependentSet() ([]NodeID, error) {
	if len(g.nodes) == 0 {
		return []NodeID{}, nil
	}

	return nil, fmt.Errorf("not implemented")
}
