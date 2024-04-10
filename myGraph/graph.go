package myGraph

import (
	"errors"
	"github.com/dominikbraun/graph"
	"math"
	"math/rand"
	"strconv"
	"sync"
)

type MyGraph struct {
	graph     graph.Graph[int, int]
	matrix    [][]int
	size      int
	maxWeight int
}

type Signal int

const (
	Wait Signal = iota
	Start
)

type Pair struct {
	X int
	Y int
}

func New() MyGraph {
	return MyGraph{
		graph: graph.New(graph.IntHash, graph.Directed(), graph.Weighted()),
	}
}

func (g *MyGraph) Generate(size int, maxWeight int) {
	matrix := GenerateWeightedDirectedGraph(size, maxWeight)

	for i := 0; i < g.size; i++ {
		for j := 0; j < g.size; j++ {
			_ = g.graph.AddEdge(i, j, graph.EdgeAttribute("label", strconv.Itoa(matrix[i][j])))
		}
	}

	g.size = size
	g.maxWeight = maxWeight
}

func (g *MyGraph) Set(matrix [][]int) {
	g.size = len(matrix)

	g.matrix = matrix

	g.maxWeight = matrix[0][0]

	for i := 0; i < g.size; i++ {
		_ = g.graph.AddVertex(i)
	}

	for i := 0; i < g.size; i++ {
		for j := 0; j < g.size; j++ {
			if matrix[i][j] != 0 {
				_ = g.graph.AddEdge(i, j, graph.EdgeWeight(matrix[i][j]))
			}
			if matrix[i][j] > g.maxWeight {
				g.maxWeight = matrix[i][j]
			}
		}
	}

}

func (g *MyGraph) Print() {

	edges, _ := g.graph.Edges()
	for _, edge := range edges {
		println(edge.Source, " ", edge.Target, " ", edge.Properties.Weight)
	}
}

func (g *MyGraph) StandardWFI() graph.Graph[int, int] {
	dist, _ := g.graph.Clone()

	adjacencyMap, _ := dist.AdjacencyMap()

	//Перехід від матриці суміжності до списку суміжності нічого не змінив - ми все перевіяємо всі можливі ребра
	for k := 0; k < g.size; k++ {
		for i := 0; i < g.size; i++ {
			for j := 0; j < g.size; j++ {
				if i == j {
					continue
				}

				edgeIK, ok := adjacencyMap[i][k]
				if !ok {
					continue
				}

				edgeKJ, ok := adjacencyMap[k][j]
				if !ok {
					continue
				}

				edgeIJ, ok := adjacencyMap[i][j]
				var err error
				if !ok {
					err = dist.AddEdge(i, j, graph.EdgeWeight(edgeIK.Properties.Weight+edgeKJ.Properties.Weight))
					adjacencyMap[i][j] = graph.Edge[int]{
						Source: i,
						Target: j,
						Properties: graph.EdgeProperties{
							Weight: edgeIK.Properties.Weight + edgeKJ.Properties.Weight,
						},
					}
				} else if edgeIK.Properties.Weight+edgeKJ.Properties.Weight < edgeIJ.Properties.Weight {
					err = dist.UpdateEdge(i, j, graph.EdgeWeight(edgeIK.Properties.Weight+edgeKJ.Properties.Weight))
					adjacencyMap[i][j] = graph.Edge[int]{
						Source: i,
						Target: j,
						Properties: graph.EdgeProperties{
							Weight: edgeIK.Properties.Weight + edgeKJ.Properties.Weight,
						},
					}
				}

				if err != nil {
					panic(err)
				}

			}
		}

	}
	return dist

}

func (g *MyGraph) ParallelWFI(numOfParallels int) (graph.Graph[int, int], error) {
	//Don't ask me what's going on here, it's a magic
	numOfParallels = (int)(math.Pow((float64)((int)(math.Sqrt((float64)(numOfParallels)))), 2))

	if g.size%(int)(math.Sqrt((float64)(numOfParallels))) != 0 {
		return nil, errors.New("sqrt of number of parallels should be divisible by graph size")
	}

	edgesInBlock := g.size / (int)(math.Sqrt((float64)(numOfParallels)))
	dist, _ := g.graph.Clone()
	adjacencyMap, _ := dist.AdjacencyMap()

	blocksInRaw := (int)(math.Sqrt((float64)(numOfParallels)))

	channels := make([][]chan Signal, blocksInRaw)
	for i := 0; i < blocksInRaw; i++ {
		channels[i] = make([]chan Signal, blocksInRaw)
	}
	for k := 0; k < g.size; k++ {
		var wg sync.WaitGroup
		mainBlockI := k / edgesInBlock
		for i := 0; i < blocksInRaw; i++ {
			for j := 0; j < blocksInRaw; j++ {
				if i == mainBlockI || j == mainBlockI {
					channels[i][j] = make(chan Signal)
				} else {
					channels[i][j] = make(chan Signal, 1)
				}
			}
		}

		//fmt.Printf("k = %d\n", k)

		for x := 0; x < blocksInRaw; x++ {
			for y := 0; y < blocksInRaw; y++ {
				wg.Add(1)
				go func(x int, y int) {
					//fmt.Printf("Block %dx%d inited\n", x, y)
					defer wg.Done()
					signalCounter := 0
					//fmt.Printf("Block %dx%d is waiting for signal\n", x, y)
					for range channels[x][y] { //wait for signal
						signalCounter++
						//fmt.Printf("Block %dx%d get %d signal\n", x, y, signalCounter)
						if (x == mainBlockI || y == mainBlockI) && signalCounter == 1 {
							close(channels[x][y])
							//fmt.Printf("Chan %dx%d closed \n", x, y)
						} else if signalCounter == 2 {
							close(channels[x][y])
							//fmt.Printf("Chan %dx%d closed \n", x, y)
						}

					}

					//fmt.Printf("Block %dx%d started\n", x, y)

					for i := x * edgesInBlock; i < (x+1)*edgesInBlock; i++ {
						for j := y * edgesInBlock; j < (y+1)*edgesInBlock; j++ {
							if i == j {
								continue
							}

							edgeIK, ok := adjacencyMap[i][k]
							if !ok {
								continue
							}

							edgeKJ, ok := adjacencyMap[k][j]
							if !ok {
								continue
							}

							edgeIJ, ok := adjacencyMap[i][j]
							var err error
							if !ok {
								err = dist.AddEdge(i, j, graph.EdgeWeight(edgeIK.Properties.Weight+edgeKJ.Properties.Weight))
								adjacencyMap[i][j] = graph.Edge[int]{
									Source: i,
									Target: j,
									Properties: graph.EdgeProperties{
										Weight: edgeIK.Properties.Weight + edgeKJ.Properties.Weight,
									},
								}
							} else if edgeIK.Properties.Weight+edgeKJ.Properties.Weight < edgeIJ.Properties.Weight {
								err = dist.UpdateEdge(i, j, graph.EdgeWeight(edgeIK.Properties.Weight+edgeKJ.Properties.Weight))
								adjacencyMap[i][j] = graph.Edge[int]{
									Source: i,
									Target: j,
									Properties: graph.EdgeProperties{
										Weight: edgeIK.Properties.Weight + edgeKJ.Properties.Weight,
									},
								}
							}

							if err != nil {
								panic(err)
							}
						}
					}

					for i := 0; i < blocksInRaw; i++ {
						if i == mainBlockI { //skip main block
							continue
						}

						if x == mainBlockI && y == mainBlockI { //send signal to all blocks if it's main block
							//fmt.Printf("Block %dx%d send signal to block %dx%d\n", x, y, x, i)
							channels[x][i] <- Start

							//fmt.Printf("Block %dx%d send signal to block %dx%d\n", x, y, i, y)
							channels[i][y] <- Start

						} else if y == mainBlockI && x != mainBlockI && i != mainBlockI && i != y { //send signal to all columns if it's main block raw
							//send signal to all columns if it's main block raw
							//fmt.Printf("Block %dx%d send signal to block %dx%d\n", x, y, x, i)
							channels[x][i] <- Start
						} else if x == mainBlockI && y != mainBlockI && i != mainBlockI && i != x { //send signal to all raws if it's main block column
							//fmt.Printf("Block %dx%d send signal to block %dx%d\n", x, y, i, y)
							channels[i][y] <- Start
						}
					}
					//fmt.Printf("Block %dx%d finished\n", x, y)

				}(x, y)

			}

		}
		channels[mainBlockI][mainBlockI] <- Start //send signal to block with [k][k] element
		//fmt.Printf("Send signal to block %dx%d\n", mainBlockI, mainBlockI)

		wg.Wait()
		//if k%100 == 0 {
		//	print("k = ", k, "\n")
		//
		//}
		//fmt.Printf("End %d iteration\n", k)
		//fmt.Printf("-----------------------\n")
	}

	return dist, nil

}

func GenerateWeightedDirectedGraph(numVertices int, maxWeight int) [][]int {

	// Ініціалізуємо матрицю суміжності з нульовими вагами
	graphMatrix := make([][]int, numVertices)
	for i := range graphMatrix {
		graphMatrix[i] = make([]int, numVertices)
	}

	// Генеруємо мінімальне остовне дерево за допомогою алгоритму Прима
	visited := make([]bool, numVertices)
	visited[0] = true
	j := 0
	edgesNum := 0
	for count := 0; count < numVertices-1; {
		i := j

		for {
			j = rand.Intn(numVertices)

			if i != j && graphMatrix[i][j] == 0 {
				graphMatrix[i][j] = rand.Intn(maxWeight) + 1
				edgesNum++
				break
			}
		}

		if !visited[j] {
			visited[j] = true
			count++
		}
	}

	randomEdgesNum := rand.Intn(numVertices*(numVertices-1)/2 - edgesNum + 1)
	for i := 0; i < randomEdgesNum; i++ {
		for {
			x := rand.Intn(numVertices)
			y := rand.Intn(numVertices)
			if x != y && graphMatrix[x][y] == 0 {
				graphMatrix[x][y] = rand.Intn(maxWeight) + 1
				break
			}
		}
	}
	return graphMatrix
}
