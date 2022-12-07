// https://adventofcode.com/2022/day/8

package main

import (
	"fmt"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

const (
	U = iota // Up
	D        // Down
	L        // Left
	R        // Right
)

type tree struct {
	height            int    // tree height
	visibleFromBorder bool   // tell whether the tree is visible from forest boder, or not
	viewingDistances  [4]int // store the viewing distances by direction
}

// viewingScore compute the viewing score of a tree
func (t tree) viewingScore() int {
	return t.viewingDistances[L] * t.viewingDistances[R] * t.viewingDistances[U] * t.viewingDistances[D]
}

// setVisibilityOnLine iterate over a forest line to set the visibility of a tree from the border
// note: expect either (x1=x2 and dx=0) or (y1=y2 and dy=0)
func setVisibilityOnLine(x1, x2, dx, y1, y2, dy int, forest [][]tree) {
	forest[x1][y1].visibleFromBorder = true
	hMax := forest[x1][y1].height
	for x := x1 + dx; (x >= x2 && dx < 0) || (x < x2 && dx > 0) || (dx == 0); x += dx {
		for y := y1 + dy; (y >= y2 && dy < 0) || (y < y2 && dy > 0) || (dy == 0); y += dy {
			if forest[x][y].height > hMax {
				forest[x][y].visibleFromBorder = true
				hMax = forest[x][y].height
			}
			if dy == 0 {
				break
			}
		}
		if dx == 0 {
			break
		}
	}
}

// getViewingDistanceOnLine iterate over a forest line to get the viewing distances from a tree into all directions
// note: expect either (x1=x2 and dx=0) or (y1=y2 and dy=0)
func getViewingDistanceOnLine(x1, x2, dx, y1, y2, dy int, forest [][]tree) int {
	viewingDistance := 0
	hMax := forest[x1][y1].height
	for x := x1 + dx; (x >= x2 && dx < 0) || (x < x2 && dx > 0) || (dx == 0); x += dx {
		for y := y1 + dy; (y >= y2 && dy < 0) || (y < y2 && dy > 0) || (dy == 0); y += dy {
			viewingDistance++
			if hMax <= forest[x][y].height || dy == 0 {
				break
			}
		}
		if hMax <= forest[x][y1].height || dx == 0 {
			break
		}
	}
	return viewingDistance
}

func main() {
	fmt.Println("--- 2022 Day 8: Treetop Tree House ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	X, Y := len(inputs), 0
	forest := make([][]tree, X)
	for x := 0; x < X; x++ {
		Y = len(inputs[x])
		forest[x] = make([]tree, Y)
		for y := 0; y < Y; y++ {
			forest[x][y].height = int(inputs[x][y] - '0')
			forest[x][y].visibleFromBorder = false
			forest[x][y].viewingDistances = [4]int{0, 0, 0, 0}
		}
	}

	////////////////////////////////////////

	visible := 0
	for x := 0; x < X; x++ {
		setVisibilityOnLine(x, x, 0, 0, Y, 1, forest)
		setVisibilityOnLine(x, x, 0, Y-1, 0, -1, forest)
	}
	for y := 0; y < Y; y++ {
		setVisibilityOnLine(0, X, 1, y, y, 0, forest)
		setVisibilityOnLine(X-1, 0, -1, y, y, 0, forest)
	}
	for x := 0; x < X; x++ {
		for y := 0; y < Y; y++ {
			if forest[x][y].visibleFromBorder {
				visible++
			}
		}
	}

	// 21
	fmt.Println("Part 1:", visible)

	////////////////////////////////////////

	viewingScore := 0
	for x := 0; x < X; x++ {
		for y := 0; y < Y; y++ {
			tree := forest[x][y]
			tree.viewingDistances[L] = getViewingDistanceOnLine(x, x, 0, y, 0, -1, forest)
			tree.viewingDistances[R] = getViewingDistanceOnLine(x, x, 0, y, Y, 1, forest)
			tree.viewingDistances[U] = getViewingDistanceOnLine(x, 0, -1, y, y, 0, forest)
			tree.viewingDistances[D] = getViewingDistanceOnLine(x, X, 1, y, y, 0, forest)
			viewingScore = max(viewingScore, tree.viewingScore())
		}
	}

	// 8
	fmt.Println("Part 2:", viewingScore)
}
