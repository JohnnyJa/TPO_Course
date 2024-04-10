package algo_test

import (
	"TPO/myGraph"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	size = 4000
)

func TestParallel(t *testing.T) {
	t.Log("Test Parallel")

	for i := 0; i < 1000; i++ {

		g := myGraph.New()
		g.Generate(100, 10)

		mapStandard, _ := g.StandardWFI().AdjacencyMap()
		resParallel, err := g.ParallelWFI(4)
		assert.Nil(t, err)
		mapParallel, _ := resParallel.AdjacencyMap()

		assert.Equal(t, mapStandard, mapParallel)

		resParallel, err = g.ParallelWFI(16)
		assert.Nil(t, err)
		mapParallel, _ = resParallel.AdjacencyMap()

		assert.Equal(t, mapStandard, mapParallel)

		resParallel, err = g.ParallelWFI(25)
		assert.Nil(t, err)
		mapParallel, _ = resParallel.AdjacencyMap()

		assert.Equal(t, mapStandard, mapParallel)

		resParallel, err = g.ParallelWFI(100)
		assert.Nil(t, err)
		mapParallel, _ = resParallel.AdjacencyMap()

		assert.Equal(t, mapStandard, mapParallel)

	}
}

func BenchmarkStandard(b *testing.B) {
	g := myGraph.New()
	g.Generate(100, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.StandardWFI()
	}
}

func BenchmarkParallel(b *testing.B) {
	g := myGraph.New()
	g.Generate(4000, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.ParallelWFI(16)
	}
}
