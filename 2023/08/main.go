// https://adventofcode.com/2023/day/8

package main

import (
	"fmt"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

const (
	Start = 'A'
	Goal  = 'Z'
	Left  = 'L'
	Right = 'R'
)

type node struct {
	last byte
	next [2]*node
}

func getDirection(b byte) int {
	if b == Left {
		return 0
	} else {
		return 1
	}
}

func getCosts(nodes map[string]*node, sequence string) map[string]int {
	steps := make(map[string]int)
	directions := utils.NewCyclicArray(utils.ArrayMap([]byte(sequence), getDirection)...)
	for start, next := range nodes {
		// ignore nodes that are not valid starts
		if next.last != Start {
			continue
		}
		// reset directions
		directions.Get(0)
		// explore until we reach a valid destination
		for next.last != Goal {
			steps[start]++
			next = next.next[directions.Next()]
		}
	}
	return steps
}

func parseNodes(inputs []string) map[string]*node {
	nodes := make(map[string]*node, len(inputs))

	// create nodes from inputs
	links := make([][3]string, 0, len(inputs))
	for _, input := range inputs {
		if input := utils.Words(input); len(input) == 3 {
			links = append(links, [3]string(input))
			nodes[input[0]] = &node{last: input[0][2]}
			nodes[input[1]] = &node{last: input[1][2]}
			nodes[input[2]] = &node{last: input[2][2]}
		}
	}

	// create connections from inputs
	for _, link := range links {
		nodes[link[0]].next[getDirection(Left)] = nodes[link[1]]
		nodes[link[0]].next[getDirection(Right)] = nodes[link[2]]
	}

	return nodes
}

func main() {
	fmt.Println("--- 2023 Day 8: Haunted Wasteland ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	costs := getCosts(parseNodes(inputs), inputs[0])

	////////////////////////////////////////

	// 2
	fmt.Println("Part 1:", utils.LCM(costs["AAA"]))

	////////////////////////////////////////

	// 6
	fmt.Println("Part 2:", utils.LCM(utils.MapValues(costs)...))
}
