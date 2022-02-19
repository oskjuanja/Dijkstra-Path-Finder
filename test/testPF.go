package main

import (
	"DijkstraPF"
	"fmt"
)

func main() {
	// Unit tests
	h := 5
	w := 5

	expGrid := make([][]string, h)
	expGridRow := make([]string, w)

	for i := 0; i < w; i++ {
		expGridRow[i] = " "
	}
	for i := 0; i < h; i++ {
		var tmp = make([]string, w)
		copy(tmp, expGridRow)
		expGrid[i] = tmp
	}

	Graph := DijkstraPF.NewGraph(w, h)

	// Test of PrintGrid()
	fmt.Println("Expected: ")
	for i := 0; i < len(expGrid); i++ {
		fmt.Println(expGrid[i])
	}

	fmt.Println("Got: ")
	Graph.PrintGrid()

	// Test of MakeWallBlock() och MakeWall()
	expGrid[0] = []string{"#", "#", "#", "#", "#"}
	expGrid[2][2] = "#"

	Graph.MakeWall(0, 0, 4, 0)
	Graph.MakeWallBlock(2, 2)

	fmt.Println("Expected: ")
	for i := 0; i < len(expGrid); i++ {
		fmt.Println(expGrid[i])
	}

	fmt.Println("Got: ")
	Graph.PrintGrid()

	// Test of PathFinder()
	expGrid[0] = []string{" ", " ", " ", " ", " "}
	expGrid[0][0] = "s"
	expGrid[0][4] = "g"
	expGrid[0][3] = "."
	for i := 1; i <= 4; i++ {
		expGrid[i][0] = "."
		expGrid[i-1][1] = "#"
		expGrid[i-1][2] = "."
		expGrid[4][i-1] = "."
	}
	expGrid[4][3] = " "

	Graph.NewGrid()
	Graph.PlaceStart(0, 0)
	Graph.PlaceGoal(4, 0)
	Graph.MakeWall(1, 0, 1, 3)
	shortest := Graph.PathFinder()

	fmt.Println("Expected: ")
	for i := 0; i < len(expGrid); i++ {
		fmt.Println(expGrid[i])
	}
	fmt.Println("expected shortest: 12")

	fmt.Println("Got: ")
	Graph.PrintGrid()
	fmt.Printf("got shortest: %v", shortest)
}
