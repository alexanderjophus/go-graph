package graph

import (
	"bytes"
	"fmt"
	"sync"
)

// algo to convert a graph to a tree decomposition
// https://en.wikipedia.org/wiki/Tree_decomposition
//
// A tree decomoposition of a graph G is a pair (T, v) where;
// T is a tree
// v is a map of 'bags' to nodes in G
// 3 conditions must be met:
// V(G) = union of all bags in v
// for all edges (u, v) in G, there exists a node t in T such that u and v are in the same bag
// the intersection of Vt1 and Vt3 is  contained in Vt2 whenever t2 is on the path in T connecting t1 and t3
type TreeDecomposition[T any, N Number] struct {
	sync.RWMutex

	nodes map[NodeID]*bag[T, N]
	width int
}

type bag[T any, N Number] struct {
	Vertices  []*node[T, N]
	Neighbors []*bag[T, N]
}

func (b *bag[T, N]) String() string {
	var buf bytes.Buffer
	buf.WriteString("{")
	for i, v := range b.Vertices {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("%v", v.Value))
	}
	buf.WriteString("}")
	return buf.String()
}

// first attempt written by chatgpt 1/10
// Need to fix, currently broken - always saying things have a tree width of one
func (g *Graph[T, N]) ComputeTreeDecomposition(root *node[T, N]) *bag[T, N] {
	visited := make(map[*node[T, N]]bool)
	var dfs func(n *node[T, N], b *bag[T, N])
	dfs = func(n *node[T, N], b *bag[T, N]) {
		if visited[n] {
			return
		}
		visited[n] = true
		b.Vertices = append(b.Vertices, n)
		for _, neighbor := range n.Neighbours() {
			if !visited[neighbor] {
				newBag := &bag[T, N]{}
				b.Neighbors = append(b.Neighbors, newBag)
				dfs(neighbor, newBag)
			}
		}
	}

	rootBag := &bag[T, N]{}
	dfs(root, rootBag)
	return rootBag
}
