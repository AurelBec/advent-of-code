// https://adventofcode.com/2022/day/12

package main

import (
	"fmt"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

type Node struct {
	elevation        rune
	neighbors        []*Node
	neighborsReverse []*Node
}

// connect connects a node to an other if the reachability condition is met
// the reverse connection is also created
// if the reverse param is true, the connection is tested to between the destination node and this one
func (from *Node) connect(to *Node, reverse bool) {
	if from == nil || to == nil {
		return
	}
	if from.elevation-to.elevation >= -1 {
		from.neighbors = append(from.neighbors, to)
		to.neighborsReverse = append(to.neighborsReverse, from)
	}
	if reverse {
		to.connect(from, false)
	}
}

// getShortestPathLength returns the length of the shorter path from start and the first goal reached
func getShortestPathLength(start *Node, neighbors func(*Node) []*Node, goalReached func(*Node) bool) int {
	openList := []*Node{start}
	costs := map[*Node]int{start: 0}

	for len(openList) != 0 {
		start, openList = openList[0], openList[1:] // pop the first node, it will have the lower cost
		currentCost := costs[start]

		if goalReached(start) {
			return currentCost
		}

		for _, neighbor := range neighbors(start) {
			if _, visited := costs[neighbor]; !visited {
				costs[neighbor] = currentCost + 1
				openList = append(openList, neighbor) // node are visited with BFS strategy, append them to the list
			}
		}
	}

	return -1
}

// parseWorld parses inputs to create a graph of Node, and returns the start and end
func parseWorld(inputs []string) (start *Node, end *Node) {
	var previousRow []*Node = make([]*Node, len(inputs[0])) // keep trace of the previous upper nodes
	for _, row := range inputs {
		var previousNode *Node = nil // keep trace of the previous left node
		for y, elevation := range row {
			node := &Node{}

			if elevation == 'S' { // found start
				node.elevation = 'a'
				start = node
			} else if elevation == 'E' { // found end
				node.elevation = 'z'
				end = node
			} else {
				node.elevation = elevation
			}

			node.connect(previousRow[y], true) // test the connection with the upper node, and the reverse
			node.connect(previousNode, true)   // test the connection with the left node, and the reverse

			previousNode = node
			previousRow[y] = previousNode
		}
	}

	return
}

func main() {
	fmt.Println("--- 2022 Day 12: Hill Climbing Algorithm ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	start, end := parseWorld(inputs)

	////////////////////////////////////////

	// 31
	fmt.Println("Part 1:", getShortestPathLength(
		// begin from start node
		start,
		// navigate through normal neighbors
		func(node *Node) []*Node { return node.neighbors },
		// stop at end
		func(node *Node) bool { return node == end },
	))

	////////////////////////////////////////

	// 29
	fmt.Println("Part 2:", getShortestPathLength(
		// begin from end node
		end,
		// navigate through reversed neighbors, as we start high to the target the lowest level
		func(node *Node) []*Node { return node.neighborsReverse },
		// stop at elevation level 'a'
		func(node *Node) bool { return node.elevation == 'a' },
	))
}
