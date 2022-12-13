// https://adventofcode.com/2022/day/14

package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

const (
	void   = '§'
	air    = ' '
	trace  = '.'
	rock   = '#'
	sand   = 'o'
	ground = '='
)

const (
	dimY = 0
	dimX = 1
)

type Cave struct {
	dim          [2]utils.Interval[int] // 0: Y range, 1: X range
	layout       [][]byte
	floorEnabled bool
}

// String returns a graphical representation of the cave's layout
func (cave Cave) String() (str string) {
	lines := make([]string, len(cave.layout))
	for l, line := range cave.layout {
		lines[l] = string(line)
	}
	return strings.Join(lines, "\n")
}

// clear resets all blocks that are not obstacles to air
func (cave *Cave) clear() {
	for y := range cave.layout {
		for x, block := range cave.layout[y] {
			if block != ground && block != rock {
				cave.layout[y][x] = air
			}
		}
	}
}

// enableFloor adds an infinite bottom layer of solid ground
func (cave *Cave) enableFloor() {
	cave.floorEnabled = true
	cave.layout = append(cave.layout, bytes.Repeat([]byte{air}, cave.dim[dimX].Len()))
	cave.layout = append(cave.layout, bytes.Repeat([]byte{ground}, cave.dim[dimX].Len()))
	cave.dim[dimY].Max = len(cave.layout)
}

// set sets the block at location x,y to the given type
func (cave *Cave) set(location utils.Location2D[int], value byte) {
	cave.layout[location.Y-cave.dim[dimY].Min][location.X-cave.dim[dimX].Min] = value
}

// get returns the block type at location x,y
// if the infinite floor is enable, the cave will be extended, and a type ground will be returned
// else it's a void
func (cave *Cave) get(location utils.Location2D[int]) byte {
	// out or range and infinite flor is disabled, return void
	if !cave.floorEnabled && (!cave.dim[dimX].Contains(location.X) || !cave.dim[dimY].Contains(location.Y)) {
		return void
	}

	if cave.floorEnabled {
		if location.X < cave.dim[dimX].Min { // extend the cave to the left
			for y, line := range cave.layout {
				cave.layout[y] = append([]byte{air}, line...)
			}
			cave.dim[dimX].Min--
			cave.set(utils.NewLocation2D(cave.dim[dimX].Min, cave.dim[dimY].Max-1), ground)
		} else if location.X >= cave.dim[dimX].Max { // extend the cave to the right
			for y, line := range cave.layout {
				cave.layout[y] = append(line, air)
			}
			cave.dim[dimX].Max++
			cave.set(utils.NewLocation2D(cave.dim[dimX].Max, cave.dim[dimY].Max-1), ground)
		}
	}

	return cave.layout[location.Y-cave.dim[dimY].Min][location.X-cave.dim[dimX].Min]
}

// putSand puts a sand unit at the given x,y location
// the sand unit will fall to a rest location
// it returns true if a rest location is found, false if the unit could not be put or fall into the void
func (cave *Cave) putSand(start utils.Location2D[int]) bool {
	// if the start location is not free, return failure
	if block := cave.get(start); block != air && block != trace {
		return false
	}

	// add a sand trace
	cave.set(start, trace)

	for _, next := range [...]utils.Location2D[int]{
		{X: start.X, Y: start.Y + 1},     // always falls down one step if possible
		{X: start.X - 1, Y: start.Y + 1}, // else, move one step down and to the left
		{X: start.X + 1, Y: start.Y + 1}, // else, move one step down and to the right
	} {
		switch cave.get(next) {
		// if next location is free, continue
		case air, trace:
			return cave.putSand(next)
		// if next location is void, return failure
		case void:
			return false
		}
	}

	// if next location is free, put sand unit here
	cave.set(start, sand)
	return true
}

// caveFromLines creates a cave from the a layout described by lines
func caveFromLines(lines [][]utils.Location2D[int]) (cave Cave) {
	// get cave dimensions first
	cave.dim = [2]utils.Interval[int]{{Min: 0, Max: -1}, {Min: -1, Max: -1}}
	for _, line := range lines {
		for _, c := range line {
			if c.X < cave.dim[dimX].Min || cave.dim[dimX].Min < 0 {
				cave.dim[dimX].Min = c.X
			}
			if c.X > cave.dim[dimX].Max || cave.dim[dimX].Max < 0 {
				cave.dim[dimX].Max = c.X
			}
			if c.Y < cave.dim[dimY].Min || cave.dim[dimY].Min < 0 {
				cave.dim[dimY].Min = c.Y
			}
			if c.Y > cave.dim[dimY].Max || cave.dim[dimY].Max < 0 {
				cave.dim[dimY].Max = c.Y
			}
		}
	}

	// set blank layout
	cave.layout = make([][]byte, cave.dim[dimY].Len())
	for y := range cave.layout {
		cave.layout[y] = bytes.Repeat([]byte{air}, cave.dim[dimX].Len())
	}

	// add rocks
	for _, line := range lines {
		from := line[0]
		for i := 1; i < len(line); i++ {
			to := line[i]
			dy, dx := 1, 1
			if from.Y > to.Y {
				dy = -1
			}
			if from.X > to.X {
				dx = -1
			}
			rockLocation := from
			for ; (rockLocation.Y >= to.Y && dy < 0) || (rockLocation.Y <= to.Y && dy > 0); rockLocation.Y += dy {
				for rockLocation.X = from.X; (rockLocation.X >= to.X && dx < 0) || (rockLocation.X <= to.X && dx > 0); rockLocation.X += dx {
					cave.set(rockLocation, rock)
				}
			}
			from = to
		}
	}

	return
}

// parseLines parses input and returns lines layout
func parseLines(inputs []string) (lines [][]utils.Location2D[int]) {
	lines = make([][]utils.Location2D[int], len(inputs))
	for i, input := range inputs {
		coordinates := strings.Split(input, " -> ")
		lines[i] = make([]utils.Location2D[int], len(coordinates))
		for c, coordinate := range coordinates {
			fmt.Sscanf(coordinate, "%v,%v", &lines[i][c].X, &lines[i][c].Y)
		}
	}

	return
}

func main() {
	fmt.Println("--- 2022 Day 14: Regolith Reservoir ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	cave := caveFromLines(parseLines(inputs))
	sandHole := utils.NewLocation2D(500, 0)

	////////////////////////////////////////

	units := 0
	for ; cave.putSand(sandHole); units++ {
		// fmt.Print("\033[H\033[2J")
		// fmt.Println(cave)
		// time.Sleep(50 * time.Millisecond)
	}

	// 24
	fmt.Println("Part 1:", units)

	////////////////////////////////////////

	cave.enableFloor()
	cave.clear()

	units = 0
	for ; cave.putSand(sandHole); units++ {
		// fmt.Print("\033[H\033[2J")
		// fmt.Println(cave)
		// time.Sleep(10 * time.Millisecond)
	}

	// 93
	fmt.Println("Part 2:", units)
}
