package graph

import (
	"bytes"
	"fmt"
	"sync"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

type Graph[T any, N Number] struct {
	sync.RWMutex

	nodes    map[NodeID]*node[T, N]
	directed bool
}

type NodeID int

type node[T any, N Number] struct {
	Value    T
	ID       NodeID
	edgesOut map[NodeID]*edge[T, N]
	edgesIn  map[NodeID]*edge[T, N]
}

type edge[T any, N Number] struct {
	end    *node[T, N]
	weight N
}

func New[T any, N Number](directed bool) *Graph[T, N] {
	return &Graph[T, N]{
		nodes:    make(map[NodeID]*node[T, N]),
		directed: directed,
	}
}

// for interesting graphs visit https://houseofgraphs.org/
func ImportG6(b []byte) *Graph[int, int] {
	g := New[int, int](false)

	if len(b) == 0 {
		return g
	}

	// decode number of nodes in graph
	n := int(b[0] - 63)
	for i := 0; i < n; i++ {
		g.AddNode(NodeID(i), i)
	}

	n1, n2 := 0, 1
	rx := b[1:]
	for i := 0; i < len(rx); i++ {
		v := rx[i] - 63
		for j := 5; j >= 0; j-- {
			mask := byte(1 << uint(j))
			if (v & mask) != 0 {
				g.AddEdge(NodeID(n1), NodeID(n2), 1)
			}
			n1++
			if n1 == n2 {
				n1, n2 = 0, n2+1
			}
		}
	}

	return g
}

func (g *Graph[T, N]) Node(id NodeID) (*node[T, N], error) {
	g.RLock()
	defer g.RUnlock()
	if node, exists := g.nodes[id]; !exists {
		return nil, ErrNodeNotFound(id)
	} else {
		return node, nil
	}
}

func (g *Graph[T, N]) NodeIDs() []NodeID {
	g.RLock()
	defer g.RUnlock()
	ids := make([]NodeID, 0)
	for _, n := range g.nodes {
		ids = append(ids, n.ID)
	}
	return ids
}

func (g *Graph[T, N]) AddNode(id NodeID, nodeData T) NodeID {
	g.Lock()
	defer g.Unlock()
	g.nodes[id] = &node[T, N]{
		Value:    nodeData,
		ID:       id,
		edgesOut: make(map[NodeID]*edge[T, N]),
		edgesIn:  make(map[NodeID]*edge[T, N]),
	}
	return id
}

func (g *Graph[T, N]) RemoveNode(id NodeID) error {
	g.Lock()
	defer g.Unlock()
	_, exists := g.nodes[id]
	if !exists {
		return ErrNodeNotFound(id)
	}
	for _, endNode := range g.nodes {
		delete(endNode.edgesOut, id)
		delete(endNode.edgesIn, id)
	}
	delete(g.nodes, id)
	return nil
}

func (g *Graph[T, N]) AddEdge(startNodeID, endNodeID NodeID, weight N) error {
	g.Lock()
	defer g.Unlock()
	startNode, exists := g.nodes[startNodeID]
	if !exists {
		return ErrNodeNotFound(startNodeID)
	}
	endNode, exists := g.nodes[endNodeID]
	if !exists {
		return ErrNodeNotFound(endNodeID)
	}
	startNode.edgesOut[endNodeID] = &edge[T, N]{
		end:    endNode,
		weight: weight,
	}
	endNode.edgesIn[startNodeID] = &edge[T, N]{
		end:    startNode,
		weight: weight,
	}
	return nil
}

func (g *Graph[T, N]) RemoveEdge(startNodeID, endNodeID NodeID) error {
	g.Lock()
	defer g.Unlock()
	startNode, exists := g.nodes[startNodeID]
	if !exists {
		return ErrNodeNotFound(startNodeID)
	}
	endNode, exists := g.nodes[endNodeID]
	if !exists {
		return ErrNodeNotFound(endNodeID)
	}
	delete(startNode.edgesOut, endNodeID)
	delete(endNode.edgesOut, startNodeID)
	return nil
}

func (g *Graph[T, N]) Neighbours(id NodeID) []NodeID {
	g.RLock()
	defer g.RUnlock()

	neighbours := make([]NodeID, 0)
	node := g.nodes[id]
	for _, edge := range node.edgesOut {
		neighbours = append(neighbours, edge.end.ID)
	}
	return neighbours
}

func (n *node[T, N]) Neighbours() []*node[T, N] {
	neighbours := make([]*node[T, N], 0)
	for _, edge := range n.edgesOut {
		neighbours = append(neighbours, edge.end)
	}
	return neighbours
}

func (g *Graph[T, N]) ExportDot() []byte {
	g.RLock()
	defer g.RUnlock()

	var buf bytes.Buffer
	buf.WriteString("digraph {\n")
	for _, node := range g.nodes {
		for _, edge := range node.edgesOut {
			buf.WriteString(fmt.Sprintf("\t%d -> %d [label=%v];\n", node.ID, edge.end.ID, edge.weight))
		}
	}
	buf.WriteString("}\n")
	return buf.Bytes()
}

func (g *Graph[T, N]) hasNeighbour(id, neighbour NodeID) bool {
	g.RLock()
	defer g.RUnlock()

	node := g.nodes[id]
	for _, edge := range node.edgesOut {
		if edge.end.ID == neighbour {
			return true
		}
	}
	for _, edge := range node.edgesIn {
		if edge.end.ID == neighbour {
			return true
		}
	}
	return false
}

// errors
var (
	ErrNodeNotFound = func(id NodeID) error {
		return fmt.Errorf("node not found: %d", id)
	}
)
