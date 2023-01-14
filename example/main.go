package main

import (
	"fmt"
	"os"

	"github.com/trelore/go-graph"
)

func main() {
	g := graph.ImportG6([]byte(os.Args[1]))
	mis := g.MaximumIndependentSet()
	fmt.Println(mis)

	dot := g.ExportDot()
	fmt.Println(string(dot))

	root, err := g.Node(0)
	if err != nil {
		panic(err)
	}
	rootBag := g.ComputeTreeDecomposition(root)
	fmt.Println(rootBag)
	for _, n := range rootBag.Neighbors {
		fmt.Println(n)
	}
}
