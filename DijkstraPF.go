// Package DijkstraPF provides tools to find the shortest path between nodes in a graph.

// The Dijkstra algorithm is a search algorithm used to find the shortest path between two nodes
// in a graph. This package provides an easy way to create and update a graph, find the shortest
// path between a start and a target node and to print said graph. Grid and graph are used somewhat
// intechangably throuough this package.

// It adds the all vertices to a list of unvisited nodes, sets the distance from the starting node
// to all other nodes to infinity (in this case a number big enough to assure no path can be of that
// length), removing the current node from the queue, then proceeds to visits all neighbors of the
// least distant node, updating their distances, adding each previous node to a list and repeating.
// Thus, this list of previous nodes can be back traced to form a path of the shortest distance.
// The actual path finding algorithm runs in O(V + E), where V denotes the number of vertices and E
// the number of edges. The graph is stored in an adjecency list.
// It's a list of data structures, in this case a list of lists, that store
// the neighbors of each node. This way, the memory consumption is optimized, as no
// "non neighbors" are stored, as opposed to an adjecency matrix.

package DijkstraPF

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type node struct {
	category  string
	set_index int
}

type Graph struct {
	set           []node
	adjecencyList [][]node
	start_node    int
	goal_node     int
	dist          []int
	prev          []int
	gridMatrix    [][]string
	grid_width    int
	grid_height   int
}

// Takes a width and a height of the grid of boxes that represents the graph. Returns empty graph.
// x is greater towards the right adn y is greater toward the bottom of the graph.
func NewGraph(width, height int) (graph Graph) {
	dist, prev := make([]int, width*height), make([]int, width*height)
	graphSet, adjecencyList := make([]node, width*height), make([][]node, width*height)
	graph = Graph{set: graphSet, adjecencyList: adjecencyList, start_node: -1, goal_node: -1, dist: dist, prev: prev, grid_width: width, grid_height: height}
	graph.NewGrid()
	return graph
}

// Method of the Graph class. Prints the graph represented by a grid to the command line.
func (graph *Graph) PrintGrid() {
	for i := 0; i < len(graph.gridMatrix); i++ {
		fmt.Println(graph.gridMatrix[i])
	}
}

// Makes a wall block in the grid at coordinate (x, y).
func (graph *Graph) MakeWallBlock(x, y int) {
	graph.gridMatrix[y][x] = "#"
	graph.fillAdjecencyList()
}

// Makes a wall in the grid between (x1, y1) and (x2, y2).
func (graph *Graph) MakeWall(x1, y1, x2, y2 int) {
	if x1 != x2 && y1 != y2 {
		fmt.Println("Coordinate choice does not make a line. Try again")
	} else if x1 == x2 {
		for i := y1; i <= y2; i++ {
			graph.gridMatrix[i][x1] = "#"
		}
	} else if y1 == y2 {
		for i := x1; i <= x2; i++ {
			graph.gridMatrix[y1][i] = "#"
		}
	}
	graph.fillAdjecencyList()
}

// Places start at (x, y).
func (graph *Graph) PlaceStart(x, y int) {
	graph.gridMatrix[y][x] = "s"
	graph.fillAdjecencyList()
}

// Places goal at (x, y).
func (graph *Graph) PlaceGoal(x, y int) {
	graph.gridMatrix[y][x] = "g"
	graph.fillAdjecencyList()
}

// Graph method. When called the user i prompted by questions to edit the graph in the
// command line. A visual representation of the changes made are printed as they are made.
// This is also called when initiating a new graph with the NewGraph function.
func (graph *Graph) EditGraph() {
	scanner := bufio.NewScanner(os.Stdin)
	startPos := make([]int, 2) // (x, y)
	goalPos := make([]int, 2)  // (x, y)
	start_chosen := false
	goal_chosen := false
questionLoop:
	for {
		fmt.Println("Current grid:")
		graph.PrintGrid()
		fmt.Print("Wall block (b), wall (w), start (s), goal (g) or clear (c)? Type 'exit' when done \n")
		scanner.Scan()
		input := scanner.Text()

		// Block input
		switch input {
		case "c":
			graph.NewGrid()
		case "b":
			x, y := coordinateInput()
			graph.MakeWallBlock(x, y)

			// Line input
		case "w":

			// (x1, y1)
			fmt.Println("First point:")
			x1, y1 := coordinateInput()

			// (x2, y2)
			fmt.Println("Second point:")
			x2, y2 := coordinateInput()

			graph.MakeWall(x1, y1, x2, y2)

		case "s":
			x, y := coordinateInput()
			x_prev, y_prev := startPos[0], startPos[1]
			if graph.gridMatrix[y_prev][x_prev] != "g" { // don't clear g from origin if s is not already placed
				graph.gridMatrix[y_prev][x_prev] = " "
			}
			startPos[0], startPos[1] = x, y
			graph.gridMatrix[y][x] = "s"
			start_chosen = true

		case "g":
			x, y := coordinateInput()
			x_prev, y_prev := goalPos[0], goalPos[1]
			if graph.gridMatrix[y_prev][x_prev] != "s" { // don't clear s from origin if g is not already placed
				graph.gridMatrix[y_prev][x_prev] = " "
			}
			goalPos[0], goalPos[1] = x, y
			graph.gridMatrix[y][x] = "g"
			goal_chosen = true

		case "exit":
			if start_chosen && goal_chosen {
				break questionLoop
			} else {
				fmt.Println("You must choose both start and goal.")
			}
		default:
			fmt.Println("Invalid choice. Try again")
		}
	}
	graph.fillAdjecencyList()
}

func coordinateInput() (x, y int) {
	scanner := bufio.NewScanner(os.Stdin)
	// (x, y)
	fmt.Print("x coordinate:\n")
	scanner.Scan()
	x_input := scanner.Text()
	fmt.Print("y coordinate:\n")
	scanner.Scan()
	y_input := scanner.Text()
	// input check
	x, errx := strconv.Atoi(x_input)
	if errx != nil {
		fmt.Println("Error, input an integer")
	}
	y, erry := strconv.Atoi(y_input)
	if erry != nil {
		fmt.Println("Error, input an integer")
	}

	return x, y
}

// Clears the grid
func (graph *Graph) NewGrid() {
	w := graph.grid_width
	h := graph.grid_height
	var gridMatrix = make([][]string, h) // gridMatrix[y][x]
	var gridRow = make([]string, w)

	for i := 0; i < w; i++ {
		gridRow[i] = " "
	}
	for i := 0; i < h; i++ {
		var tmp = make([]string, w)
		copy(tmp, gridRow)
		gridMatrix[i] = tmp
	}
	graph.gridMatrix = gridMatrix
	graph.fillAdjecencyList()
}

func (graph *Graph) fillAdjecencyList() {

	var list_idx int
	w := graph.grid_width
	h := graph.grid_height

	for y := 0; y < h; y++ {

		for x := 0; x < w; x++ {

			neighbors := make([]node, 5)
			category := graph.gridMatrix[y][x] // current node
			set_index := y*h + x
			graph.set[set_index] = node{category, set_index} // add to set
			number_of_neighboring_walls := 0
			neighbor_index := 0

			switch graph.gridMatrix[y][x] {
			case "s":
				graph.start_node = list_idx
			case "g":
				graph.goal_node = list_idx
			}

			// no neighbors if current node i a wall
			if graph.gridMatrix[y][x] == "#" {
				neighbors = neighbors[:0]

				// left grid edge
			} else if x == 0 {

				category := graph.gridMatrix[y][x+1] // right of node
				if category != "#" {                 // only add neighbor if not a wall
					neighbors[neighbor_index] = node{category, y*h + x + 1}
					neighbor_index++
				} else {
					number_of_neighboring_walls++
				}

				if y == 0 { // top left corner
					category = graph.gridMatrix[y+1][x] // below node
					if category != "#" {                // only add neighbor if not a wall
						neighbors[neighbor_index] = node{category, (y+1)*h + x}
						neighbor_index++
					} else {
						number_of_neighboring_walls++
					}

					neighbors = neighbors[:2-number_of_neighboring_walls]
					number_of_neighboring_walls = 0

				} else if y == 4 { // bottom left corner
					category = graph.gridMatrix[y-1][x] // above node
					if category != "#" {                // only add neighbor if not a wall
						neighbors[neighbor_index] = node{category, (y-1)*h + x}
						neighbor_index++
					} else {
						number_of_neighboring_walls++
					}

					neighbors = neighbors[:2-number_of_neighboring_walls]
					number_of_neighboring_walls = 0

				} else { // left grid edge
					category := graph.gridMatrix[y+1][x] // below node
					if category != "#" {                 // only add neighbor if not a wall
						neighbors[neighbor_index] = node{category, (y+1)*h + x}
						neighbor_index++
					} else {
						number_of_neighboring_walls++
					}

					category = graph.gridMatrix[y-1][x] // above node
					if category != "#" {                // only add neighbor if not a wall
						neighbors[neighbor_index] = node{category, (y-1)*h + x}
						neighbor_index++
					} else {
						number_of_neighboring_walls++
					}

					neighbors = neighbors[:3-number_of_neighboring_walls]
					number_of_neighboring_walls = 0
				}

				// right grid edge. x != 0
			} else if x == 4 {

				category := graph.gridMatrix[y][x-1] // left of node
				if category != "#" {                 // only add neighbor if not a wall
					neighbors[neighbor_index] = node{category, y*h + x - 1}
					neighbor_index++
				} else {
					number_of_neighboring_walls++
				}

				if y == 0 { // top right corner
					category = graph.gridMatrix[y+1][x] // below node
					if category != "#" {                // only add neighbor if not a wall
						neighbors[neighbor_index] = node{category, (y+1)*h + x}
						neighbor_index++
					} else {
						number_of_neighboring_walls++
					}

					neighbors = neighbors[:2-number_of_neighboring_walls]
					number_of_neighboring_walls = 0

				} else if y == 4 { // bottom right corner
					category = graph.gridMatrix[y-1][x] // above node
					if category != "#" {                // only add neighbor if not a wall
						neighbors[neighbor_index] = node{category, (y-1)*h + x}
						neighbor_index++
					} else {
						number_of_neighboring_walls++
					}

					neighbors = neighbors[:2-number_of_neighboring_walls]
					number_of_neighboring_walls = 0

				} else { // right grid edge
					category := graph.gridMatrix[y+1][x] // below node
					if category != "#" {                 // only add neighbor if not a wall
						neighbors[neighbor_index] = node{category, (y+1)*h + x}
						neighbor_index++
					} else {
						number_of_neighboring_walls++
					}

					category = graph.gridMatrix[y-1][x] // above node
					if category != "#" {                // only add neighbor if not a wall
						neighbors[neighbor_index] = node{category, (y-1)*h + x}
						neighbor_index++
					} else {
						number_of_neighboring_walls++
					}

					neighbors = neighbors[:3-number_of_neighboring_walls]
					number_of_neighboring_walls = 0
				}

				// top grid edge. x != 0, x != 4
			} else if y == 0 {
				category := graph.gridMatrix[y+1][x] // below node
				if category != "#" {                 // only add neighbor if not a wall
					neighbors[neighbor_index] = node{category, (y+1)*h + x}
					neighbor_index++
				} else {
					number_of_neighboring_walls++
				}

				category = graph.gridMatrix[y][x-1] // left of node
				if category != "#" {                // only add neighbor if not a wall
					neighbors[neighbor_index] = node{category, y*h + x - 1}
					neighbor_index++
				} else {
					number_of_neighboring_walls++
				}

				category = graph.gridMatrix[y][x+1] // right of node
				if category != "#" {                // only add neighbor if not a wall
					neighbors[neighbor_index] = node{category, y*h + x + 1}
					neighbor_index++
				} else {
					number_of_neighboring_walls++
				}

				neighbors = neighbors[:3-number_of_neighboring_walls]
				number_of_neighboring_walls = 0

				// bottom edge. x != 0, x != 4
			} else if y == 4 {
				category := graph.gridMatrix[y-1][x] // above node
				if category != "#" {                 // only add neighbor if not a wall
					neighbors[neighbor_index] = node{category, (y-1)*h + x}
					neighbor_index++
				} else {
					number_of_neighboring_walls++
				}

				category = graph.gridMatrix[y][x-1] // left of node
				if category != "#" {                // only add neighbor if not a wall
					neighbors[neighbor_index] = node{category, y*h + x + 1}
					neighbor_index++
				} else {
					number_of_neighboring_walls++
				}

				category = graph.gridMatrix[y][x+1] // right of node
				if category != "#" {                // only add neighbor if not a wall
					neighbors[neighbor_index] = node{category, y*h + x + 1}
					neighbor_index++
				} else {
					number_of_neighboring_walls++
				}

				neighbors = neighbors[:3-number_of_neighboring_walls]
				number_of_neighboring_walls = 0

				// non edge node
			} else {
				category := graph.gridMatrix[y-1][x] // above node
				if category != "#" {                 // only add neighbor if not a wall
					neighbors[neighbor_index] = node{category, (y-1)*h + x}
					neighbor_index++
				} else {
					number_of_neighboring_walls++
				}

				category = graph.gridMatrix[y][x-1] // left of node
				if category != "#" {                // only add neighbor if not a wall
					neighbors[neighbor_index] = node{category, y*h + x - 1}
					neighbor_index++
				} else {
					number_of_neighboring_walls++
				}

				category = graph.gridMatrix[y][x+1] // right of node
				if category != "#" {                // only add neighbor if not a wall
					neighbors[neighbor_index] = node{category, y*h + x + 1}
					neighbor_index++
				} else {
					number_of_neighboring_walls++
				}

				category = graph.gridMatrix[y+1][x] // below node
				if category != "#" {                // only add neighbor if not a wall
					neighbors[neighbor_index] = node{category, (y+1)*h + x}
					neighbor_index++
				} else {
					number_of_neighboring_walls++
				}

				neighbors = neighbors[:4-number_of_neighboring_walls]
				number_of_neighboring_walls = 0
			}
			graph.adjecencyList[list_idx] = neighbors
			list_idx++
		}
	}
}

// Method of the Graph class. Finds the shortest path between two nodes specified in the EditGraph()
// function. Returns shortest distance if found and -1 otherwise. Use PrintGraph()
// to see visual representation of the path.
func (graph *Graph) PathFinder() (shortest int) {

	var unvisitedNodes = make([]node, len(graph.adjecencyList))
	inf := len(graph.set)

	for index, node := range graph.set {
		unvisitedNodes[index] = node
		graph.dist[index] = inf
	}

	// set distance from start to itself to 0
	graph.dist[graph.start_node] = 0

	for len(unvisitedNodes) > 0 {

		// find min
		min := inf
		min_index := 0
		for index, current_node := range unvisitedNodes {

			if graph.dist[current_node.set_index] < min {
				min_index = index
				min = graph.dist[current_node.set_index]
			}
		}

		curr_node := unvisitedNodes[min_index]

		// remove current node from queue
		unvisitedNodes = append(unvisitedNodes[:min_index], unvisitedNodes[min_index+1:]...)

		// break if unreachable
		if graph.dist[curr_node.set_index] == inf {
			break
		}

		// for each neighbor of current node
		for _, neighbor := range graph.adjecencyList[curr_node.set_index] {

			alt := graph.dist[curr_node.set_index] + 1 // distance from start to current node + distance from current node to neighbor

			if alt < graph.dist[neighbor.set_index] { // distance from start to current neighbor
				graph.dist[neighbor.set_index] = alt                 // update distance
				graph.prev[neighbor.set_index] = curr_node.set_index // update path
			}
		}
	}

	// update grid to show path
	i := graph.prev[graph.goal_node]
	for i > 0 {
		x := i % 5
		y := (i - x) / 5
		graph.gridMatrix[y][x] = "."
		i = graph.prev[i]
	}
	if graph.dist[graph.goal_node] != inf {
		return graph.dist[graph.goal_node]
	} else {
		return -1
	}

}
