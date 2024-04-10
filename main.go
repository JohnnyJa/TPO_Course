package main

import "TPO/myGraph"

type test struct {
}

func (*test) Case1() [][]int {
	testCase1 := [][]int{
		{0, 4, 7, 2, 8, 1, 6, 3},
		{0, 0, 9, 0, 0, 6, 8, 1},
		{0, 0, 0, 0, 10, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 8, 6, 0, 5, 0, 0},
		{0, 0, 0, 0, 6, 0, 0, 0},
		{0, 4, 0, 0, 0, 7, 0, 0},
		{0, 0, 3, 0, 0, 0, 0, 0}}
	return testCase1
}

func (*test) Case2() [][]int {
	testCase3 := [][]int{
		{0, 8, 0, 0, 0, 3},
		{4, 0, 7, 0, 6, 0},
		{0, 5, 0, 1, 0, 0},
		{0, 0, 2, 0, 5, 0},
		{0, 9, 0, 7, 0, 2},
		{9, 0, 0, 0, 9, 0},
	}

	return testCase3
}

func (*test) Case3() [][]int {
	testCase3 := [][]int{
		{0, 0, 0, 0, 0, 343, 0, 1435, 464, 0},
		{0, 0, 0, 0, 0, 879, 954, 811, 0, 524},
		{0, 0, 0, 0, 1364, 1054, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 433, 0, 0, 1053},
		{0, 0, 1364, 0, 0, 1106, 0, 0, 0, 766},
		{343, 879, 1054, 0, 1106, 0, 0, 0, 0, 0},
		{0, 954, 0, 433, 0, 0, 0, 837, 0, 0},
		{1435, 811, 0, 0, 0, 0, 837, 0, 0, 0},
		{464, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 524, 0, 1053, 766, 0, 0, 0, 0, 0},
	}

	return testCase3
}

func main() {
	test := test{}
	initMatrix := test.Case3()
	g := myGraph.New()
	g.Set(initMatrix)
	g.Print()
	println()
	println()

	res, _ := g.ParallelWFI(4)
	a, _ := res.AdjacencyMap()
	for i := 0; i < len(initMatrix); i++ {
		for j := 0; j < len(initMatrix); j++ {
			if _, ok := a[i][j]; !ok {
				print("0 ")
				continue
			}
			print(a[i][j].Properties.Weight, " ")
		}
		println()
	}

	g.GenerateDOT()
}
