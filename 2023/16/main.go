// https://adventofcode.com/2023/day/16

package main

import (
	"fmt"
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

var next = map[rune]map[int][]int{
	'.':  {U: []int{U}, R: []int{R}, D: []int{D}, L: []int{L}},
	'|':  {U: []int{U}, R: []int{U, D}, D: []int{D}, L: []int{U, D}},
	'-':  {U: []int{L, R}, R: []int{R}, D: []int{L, R}, L: []int{L}},
	'\\': {U: []int{L}, R: []int{D}, D: []int{R}, L: []int{U}},
	'/':  {U: []int{R}, R: []int{U}, D: []int{L}, L: []int{D}},
}

type node struct {
	x, y int
	dir  int
}

type Layout struct {
	N     int
	cells [][]rune
}

func (l Layout) energize(x, y int, dir int) int {
	explored := make(map[node]bool, l.N*l.N*4)

	stack := collections.NewStack(node{x, y, dir})
	for !stack.IsEmpty() {
		step, _ := stack.Pop()

		if 0 > step.x || step.x >= l.N || 0 > step.y || step.y >= l.N {
			continue
		}

		if explored[step] {
			continue
		}
		explored[step] = true

		for _, next := range next[l.cells[step.x][step.y]][step.dir] {
			switch next {
			case U:
				stack.Push(node{step.x, step.y - 1, next})
			case R:
				stack.Push(node{step.x + 1, step.y, next})
			case D:
				stack.Push(node{step.x, step.y + 1, next})
			case L:
				stack.Push(node{step.x - 1, step.y, next})
			}
		}
	}

	energizedCells := make(map[int]bool, l.N*l.N)
	for c := range explored {
		energizedCells[c.x+c.y*l.N] = true
	}
	return len(energizedCells)
}

func parseLayout(inputs []string) Layout {
	layout := Layout{
		N:     len(inputs),
		cells: make([][]rune, len(inputs)),
	}
	for x := 0; x < layout.N; x++ {
		layout.cells[x] = make([]rune, layout.N)
	}
	for y, row := range inputs {
		for x, cell := range row {
			layout.cells[x][y] = cell
		}
	}
	return layout
}

func main() {
	fmt.Println("--- 2023 Day 16: The Floor Will Be Lava ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	layout := parseLayout(inputs)

	////////////////////////////////////////

	// 46
	fmt.Println("Part 1:", layout.energize(0, 0, R))

	////////////////////////////////////////

	// GHDM
	bruteForce := 0
	for x := 0; x < layout.N; x++ {
		for y := 0; y < layout.N; y++ {
			for d := 0; d < 4; d++ {
				bruteForce = max(bruteForce, layout.energize(x, y, d))
			}
		}
	}

	// 51
	fmt.Println("Part 2:", bruteForce)
}
