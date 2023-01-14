package graph

// Greedy implementation of the maximum independent set algorithm.
func (g *Graph[T, N]) MaximumIndependentSet() []NodeID {
	if len(g.nodes) == 0 {
		return []NodeID{}
	}

	return maximumIndependentSet(g)
}

// ms(G) = max(1 + ms(G - N(B)), ms(G - B)).
func maximumIndependentSet[T any, N Number](g *Graph[T, N]) []NodeID {
	if len(g.nodes) == 0 {
		return []NodeID{}
	}

	if len(g.nodes) == 1 {
		for _, node := range g.nodes {
			return []NodeID{node.ID}
		}
	}

	mis := []NodeID{}
	for _, node := range g.nodes {
		g1 := newGraphWithoutNodeAndNeighbours(g, node.ID)
		g2 := newGraphWithoutNode(g, node.ID)
		m := max(
			append(maximumIndependentSet(g1), node.ID),
			maximumIndependentSet(g2))
		mis = max(m, mis)
	}
	return mis
}

func newGraphWithoutNode[T any, N Number](g *Graph[T, N], nodeID NodeID) *Graph[T, N] {
	newGraph := New[T, N](g.directed)
	for _, node := range g.nodes {
		if node.ID != nodeID {
			newGraph.AddNode(node.ID, node.Value)
		}
	}
	for _, node := range g.nodes {
		if node.ID != nodeID {
			for _, edge := range node.edgesOut {
				newGraph.AddEdge(node.ID, edge.end.ID, edge.weight)
			}
		}
	}
	return newGraph
}

func newGraphWithoutNodeAndNeighbours[T any, N Number](g *Graph[T, N], nodeID NodeID) *Graph[T, N] {
	newGraph := New[T, N](g.directed)
	for _, node := range g.nodes {
		if !g.isNodeOrNeighbour(nodeID, node.ID) {
			newGraph.AddNode(node.ID, node.Value)
		}
	}
	for _, node := range g.nodes {
		if !g.isNodeOrNeighbour(nodeID, node.ID) {
			for _, edge := range node.edgesOut {
				newGraph.AddEdge(node.ID, edge.end.ID, edge.weight)
			}
		}
	}
	return newGraph
}

func (g *Graph[T, N]) isNodeOrNeighbour(nodeID, neighbourID NodeID) bool {
	return nodeID == neighbourID || g.hasNeighbour(nodeID, neighbourID)
}

func max[T any](x, y []T) []T {
	if len(x) > len(y) {
		return x
	}
	return y
}
