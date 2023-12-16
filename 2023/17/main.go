// https://adventofcode.com/2023/day/17

package main

import (
	"fmt"
	"math"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
	"github.com/aurelbec/advent-of-code/utils/collections"
)

const (
	U = iota
	R
	D
	L
)

var offsets = map[int][2]int{
	U: {0, -1},
	R: {+1, 0},
	D: {0, +1},
	L: {-1, 0},
}

type node struct {
	loc utils.Location2D[int]
	dir int
}

type City struct {
	X, Y int
	loss [][]int
}

func (c City) getMinimalHeatLoss(start, end utils.Location2D[int], minStreak, maxStreak int) int {
	minHeatLoss := math.MaxInt

	visited := make(map[node]int, c.X*c.Y*4)
	for dir := range offsets {
		visited[node{loc: start, dir: dir}] = 0
	}

	queue := collections.NewQueue(utils.MapKeys(visited)...)
	for !queue.IsEmpty() {
		current, _ := queue.Dequeue()
		currentHeatLoss := visited[current]
		if current.loc == end {
			minHeatLoss = min(minHeatLoss, currentHeatLoss)
		}

		for dir, offset := range offsets {
			// ignore reverse
			if (current.dir+2)%4 == dir {
				continue
			}
			// ignore forward
			if current.dir == dir {
				continue
			}

			// explore all forward steps directly in direction
			nextHeatLoss := currentHeatLoss
			for i := 1; i <= maxStreak; i++ {
				next := node{loc: current.loc.MovedBy(i*offset[0], i*offset[1]), dir: dir}
				if 0 > next.loc.X || next.loc.X >= c.X || 0 > next.loc.Y || next.loc.Y >= c.Y {
					continue
				}

				nextHeatLoss += c.loss[next.loc.X][next.loc.Y]
				if prev, found := visited[next]; found && prev <= nextHeatLoss {
					continue
				}

				if i < minStreak {
					continue
				}

				visited[next] = nextHeatLoss
				queue.Enqueue(next)
			}
		}
	}

	return minHeatLoss
}

func parseCity(inputs []string) City {
	city := City{
		X:    len(inputs[0]),
		Y:    len(inputs),
		loss: make([][]int, len(inputs[0])),
	}
	for x := 0; x < city.X; x++ {
		city.loss[x] = make([]int, city.Y)
	}
	for y, row := range inputs {
		for x, loss := range row {
			city.loss[x][y] = int(loss - '0')
		}
	}
	return city
}

func main() {
	fmt.Println("--- 2023 Day 17: Clumsy Crucible ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	city := parseCity(inputs)

	////////////////////////////////////////

	// 102
	fmt.Println("Part 1:", city.getMinimalHeatLoss(
		utils.NewLocation2D(0, 0),
		utils.NewLocation2D(city.X-1, city.Y-1),
		1, 3,
	))

	////////////////////////////////////////

	// 94
	fmt.Println("Part 2:", city.getMinimalHeatLoss(
		utils.NewLocation2D(0, 0),
		utils.NewLocation2D(city.X-1, city.Y-1),
		4, 10,
	))
}
