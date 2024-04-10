package myGraph

import (
	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"os"
	"strconv"
)

func (g *MyGraph) GenerateDOT() {
	if g.size > 5 {
		print("MyGraph is too big to visualize")
		return
	}

	gr := graph.New(graph.IntHash, graph.Directed(), graph.Weighted())
	for i := 0; i < g.size; i++ {
		_ = gr.AddVertex(i + 1)
	}

	edges, _ := g.graph.Edges()
	for _, edge := range edges {
		_ = gr.AddEdge(edge.Source+1, edge.Target+1, graph.EdgeAttribute("label", strconv.Itoa(edge.Properties.Weight)))
	}

	file, _ := os.Create("./my-graph.gv")
	_ = draw.DOT(gr, file)
}
