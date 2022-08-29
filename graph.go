package graph

import (
	"fmt"
	"sync"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

type Graph[T any, N Number] struct {
	nodes    map[NodeID]*node[T, N]
	directed bool
	nextID   NodeID
	lock     sync.RWMutex

	misNodes []NodeID
}

type NodeID int

type node[T any, N Number] struct {
	Value    T
	id       NodeID
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
		nextID:   0,
		lock:     sync.RWMutex{},
		misNodes: make([]NodeID, 0),
	}
}

func ImportG6(b []byte) *Graph[int, int] {
	g := New[int, int](false)

	if len(b) == 0 {
		return g
	}

	// decode number of nodes in graph
	n := int(b[0] - 63)
	for i := 0; i < n; i++ {
		g.AddNode(i)
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

func (g *Graph[T, N]) Node(id NodeID) (T, error) {
	g.lock.RLock()
	defer g.lock.RUnlock()
	var value T
	if node, exists := g.nodes[id]; !exists {
		return value, ErrNodeNotFound(id)
	} else {
		return node.Value, nil
	}
}

func (g *Graph[T, N]) NodeIDs() []NodeID {
	g.lock.RLock()
	defer g.lock.RUnlock()
	ids := make([]NodeID, 0)
	for _, n := range g.nodes {
		ids = append(ids, n.id)
	}
	return ids
}

func (g *Graph[T, N]) AddNode(nodeData T) NodeID {
	g.lock.Lock()
	id := NodeID(g.nextID)
	g.nextID++
	g.nodes[id] = &node[T, N]{
		Value:    nodeData,
		id:       id,
		edgesOut: make(map[NodeID]*edge[T, N]),
		edgesIn:  make(map[NodeID]*edge[T, N]),
	}
	g.lock.Unlock()
	return id
}

func (g *Graph[T, N]) RemoveNode(id NodeID) error {
	g.lock.Lock()
	defer g.lock.Unlock()
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
	g.lock.Lock()
	defer g.lock.Unlock()
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
	g.lock.Lock()
	defer g.lock.Unlock()
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

func (g *Graph[T, N]) Neighbors(id NodeID) []NodeID {
	g.lock.RLock()
	defer g.lock.RUnlock()

	neighbors := make([]NodeID, 0)
	node := g.nodes[id]
	for _, edge := range node.edgesOut {
		neighbors = append(neighbors, edge.end.id)
	}
	return neighbors
}

// errors
var (
	ErrNodeNotFound = func(id NodeID) error {
		return fmt.Errorf("node not found: %d", id)
	}
)
