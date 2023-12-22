// https://adventofcode.com/2023/day/23

package main

import (
	"fmt"
	"maps"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
	"github.com/aurelbec/advent-of-code/utils/collections"
)

var dirs = map[byte][2]int{
	'^': {0, -1},
	'>': {+1, 0},
	'v': {0, +1},
	'<': {-1, 0},
}

type node struct {
	origin string
	loc    [2]int
	dir    byte
	steps  int
}

func (n node) next(dir byte) node {
	n.dir = dir
	n.steps += 1
	n.loc[0] += dirs[dir][0]
	n.loc[1] += dirs[dir][1]
	return n
}

func isSlope(r byte) bool {
	return r == '^' || r == '>' || r == 'v' || r == '<'
}

func isCrossing(neighbors ...byte) bool {
	count := 0
	for _, neighbor := range neighbors {
		if isSlope(neighbor) {
			count++
			if count >= 2 {
				return true
			}
		}
	}
	return false
}

func getLongestDist(edges map[string]map[string]int, start string, end string) int {
	if dist, found := getLongestDistRec(edges, start, end, map[string]struct{}{start: {}}); found {
		return dist
	} else {
		return -1
	}
}

func getLongestDistRec(edges map[string]map[string]int, start string, end string, visited map[string]struct{}) (int, bool) {
	maxDist, found := 0, false
	for next, cost := range edges[start] {
		if next == end {
			return cost, true
		}

		if _, visited := visited[next]; visited {
			continue
		}
		visited := maps.Clone(visited)
		visited[next] = struct{}{}

		dist, ok := getLongestDistRec(edges, next, end, visited)
		found = found || ok
		maxDist = max(maxDist, dist+cost)
	}

	if !found {
		return -1, false
	}
	return maxDist, true
}

func parsePaths(inputs []string) map[string]map[string]int {
	nodes := make(map[[2]int]string)
	edges := make(map[string]map[string]int)

	queue := collections.NewStack[node]()
	for x, y := 1, 0; x < len(inputs[y])-1; x++ {
		if inputs[y][x] == '.' {
			nodes[[2]int{x, y}] = "S"
			queue.Push(node{origin: "S", steps: 0, loc: [2]int{x, y}})
			break
		}
	}
	for x, y := 1, len(inputs)-1; x < len(inputs[y])-1; x++ {
		if inputs[y][x] == '.' {
			nodes[[2]int{x, y}] = "E"
			break
		}
	}

	visited := make(map[[2]int]bool)
	for !queue.IsEmpty() {
		n, _ := queue.Pop()
		x, y := n.loc[0], n.loc[1]
		// check tile validity
		if !(x >= 0 && y >= 0 && y < len(inputs) && x < len(inputs[y])) {
			continue
		}
		if inputs[y][x] == '#' {
			continue
		}
		if inputs[y][x] != '.' && inputs[y][x] != n.dir {
			continue
		}

		// check if on crossing
		if x > 0 && y > 0 && y < len(inputs)-1 && x < len(inputs[y])-1 {
			if nodes[n.loc] == "" && isCrossing(inputs[y][x-1], inputs[y+1][x], inputs[y][x+1], inputs[y-1][x]) {
				nodes[[2]int{x, y}] = fmt.Sprint(len(nodes) - 1)
			}
		}
		if node := nodes[n.loc]; node != "" && n.steps > 2 {
			if edges[n.origin] == nil {
				edges[n.origin] = map[string]int{node: n.steps}
			} else {
				edges[n.origin][node] = n.steps
			}

			n.steps = 0
			n.origin = node
		}

		// do not visit twice
		if visited[n.loc] {
			continue
		}
		visited[n.loc] = true

		for slope := range dirs {
			queue.Push(n.next(slope))
		}
	}

	return edges
}

func main() {
	fmt.Println("--- 2023 Day 23: A Long Walk ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	edges := parsePaths(inputs)

	////////////////////////////////////////

	// 94
	fmt.Println("Part 1:", getLongestDist(edges, "S", "E"))

	////////////////////////////////////////

	// add edges in both direction
	for from, next := range edges {
		for to, dist := range next {
			if edges[to] == nil {
				edges[to] = map[string]int{from: dist}
			} else {
				edges[to][from] = dist
			}
		}
	}

	// 154
	fmt.Println("Part 2:", getLongestDist(edges, "S", "E"))
}
