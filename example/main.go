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
}
