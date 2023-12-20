// https://adventofcode.com/2023/day/21

package main

import (
	"fmt"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

const (
	gardenPlot = '.'
	rock       = '#'
	start      = 'S'
)

var dirs = [4][2]int{
	{-1, 0},
	{1, 0},
	{0, -1},
	{0, 1},
}

type Grid struct {
	X, Y  int
	tiles [][]rune
}

func (g Grid) isTileValid(tile [2]int, infinite bool) bool {
	if tile[0] < 0 || tile[1] < 0 || tile[0] >= g.X || tile[1] >= g.Y {
		if !infinite {
			return false
		}
		tile[0] = utils.Mod(tile[0], g.X)
		tile[1] = utils.Mod(tile[1], g.Y)
	}
	return g.tiles[tile[0]][tile[1]] != rock
}

func (g Grid) getReachableTiles(start [2]int, steps int, infinite bool) int {
	// express steps as steps=n*P+R
	p := g.X
	r := steps % p

	points := []utils.Location2D[int]{}
	visited := map[int]map[[2]int]struct{}{0: {start: {}}}

	for s := 0; s <= steps; s++ {
		visited[s+1] = map[[2]int]struct{}{}
		for tile := range visited[s] {
			for _, dir := range dirs {
				next := [2]int{tile[0] + dir[0], tile[1] + dir[1]}
				if !g.isTileValid(next, infinite) {
					continue
				}
				visited[s+1][next] = struct{}{}
			}
		}

		// find the first 3 interpolation points if grid is infinite
		if s%p == r && infinite {
			points = append(points, utils.NewLocation2D(s, len(visited[s])))
			if l := len(points); l > 3 {
				x0, x1, x2, y0, y1, y2 := points[l-3].X, points[l-2].X, points[l-1].X, points[l-3].Y, points[l-2].Y, points[l-1].Y
				// ensure interpolation is valid by looking to previous value
				prev := utils.Interpolation(steps, points[l-4].X, x0, x1, points[l-4].Y, y0, y1)
				current := utils.Interpolation(steps, x0, x1, x2, y0, y1, y2)
				if prev == current {
					return current
				}
			}
		}

		// stop after 2 periods iterations if grid is not infinite
		if s > 2*(p-1) && !infinite {
			return len(visited[2*(p-1)-steps%2])
		}
	}

	return len(visited[steps])
}

func parseGrid(inputs []string) (grid Grid, startPos [2]int) {
	grid.tiles = make([][]rune, len(inputs))
	for x, row := range inputs {
		grid.tiles[x] = make([]rune, len(row))
		for y, tile := range row {
			if tile == start {
				startPos = [2]int{x, y}
				tile = gardenPlot
			}
			grid.tiles[x][y] = tile
		}
	}
	grid.X = len(grid.tiles)
	grid.Y = len(grid.tiles[0])
	return grid, startPos
}

func main() {
	fmt.Println("--- 2023 Day 21: Step Counter ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	grid, start := parseGrid(inputs)

	////////////////////////////////////////

	// 6: 16
	// 64: 42
	fmt.Println("Part 1:", grid.getReachableTiles(start, 64, false))

	////////////////////////////////////////

	// 6: 16
	// 10: 50
	// 50: 1594
	// 100: 6536
	// 500: 167004
	// 1000: 668697
	// 5000: 16733044
	// 26501365: 470149643712804
	fmt.Println("Part 2:", grid.getReachableTiles(start, 26501365, true))
}
