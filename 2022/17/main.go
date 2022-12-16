// https://adventofcode.com/2022/day/17

package main

import (
	"fmt"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

var rocks = [...][][]byte{
	{
		{'#', '#', '#', '#'},
	},
	{
		{' ', '#', ' '},
		{'#', '#', '#'},
		{' ', '#', ' '},
	},
	{
		{' ', ' ', '#'},
		{' ', ' ', '#'},
		{'#', '#', '#'},
	},
	{
		{'#'},
		{'#'},
		{'#'},
		{'#'},
	},
	{
		{'#', '#'},
		{'#', '#'},
	},
}

type Rock struct {
	visible bool
	layout  [][]byte
	pos     utils.Location2D[int]
}

type Cave struct {
	width  int      // cave width
	layout [][]byte // current cave layout
	offset int      // cave height offset

	heights []int       // cave height after each rock fall
	rows    []int       // cave top row state after each rock fall
	states  map[int]int // cache of all top rows state + rock + jet after each rock fall and index it was encountered

	currentRock Rock // current rock falling

	rocks utils.CyclicArray[Rock] // list of rocks falling
	jets  utils.CyclicArray[byte] // list of steam jets
}

// height returns the current cave height
func (cave Cave) height() int {
	return len(cave.layout) + cave.offset
}

// maxUnreachableHeight returns the index of the first top row that can no longer be reached
func (cave Cave) maxUnreachableHeight() int {
	max := 0
	for x := 0; x < cave.width; x++ {
		for y, row := range cave.layout {
			if y == len(cave.layout)-1 {
				return len(cave.layout)
			}
			if row[x] == '#' {
				max = utils.Max(max, y)
				break
			}
		}
	}

	// keep extra 8 as safety
	return utils.Min(max+8, len(cave.layout)-1)
}

// String prints the first 10 rows of the cave
func (cave Cave) String() string {
	layout := make([][]byte, len(cave.layout))
	for i, row := range cave.layout {
		layout[i] = make([]byte, len(row))
		copy(layout[i], row)
	}

	if cave.currentRock.visible {
		for i, row := range cave.currentRock.layout {
			for j, block := range row {
				if block != ' ' {
					layout[i+cave.currentRock.pos.Y][j+cave.currentRock.pos.X] = '@'
				}
			}
		}
	}

	str := ""
	for i, row := range layout {
		if i > 10 {
			return str
		}
		str += fmt.Sprintf("|%s|\n", row)
	}

	return str
}

// currentRockCanGoesLeft returns whether the rock can go left or not
func (cave Cave) currentRockCanGoesLeft() bool {
	return cave.currentRock.pos.X > 0 && cave.isCurrentRockValid(-1, 0)
}

// currentRockCanGoesRight returns whether the rock can go right or not
func (cave Cave) currentRockCanGoesRight() bool {
	return cave.currentRock.pos.X+len(cave.currentRock.layout[0]) < cave.width && cave.isCurrentRockValid(+1, 0)
}

// currentRockCanGoesDown returns whether the rock can go down or not
func (cave Cave) currentRockCanGoesDown() bool {
	return len(cave.currentRock.layout)+cave.currentRock.pos.Y < len(cave.layout) && cave.isCurrentRockValid(0, +1)
}

// isCurrentRockValid returns whether the rock will be valid given the next state
func (cave Cave) isCurrentRockValid(dx, dy int) bool {
	px := cave.currentRock.pos.X + dx
	py := cave.currentRock.pos.Y + dy
	for i, row := range cave.currentRock.layout {
		for j, block := range row {
			if block == '#' && cave.layout[i+py][j+px] == '#' {
				return false
			}
		}
	}
	return true
}

// stopCurrentRock registers the rock at its current location
func (cave *Cave) stopCurrentRock() {
	cave.currentRock.visible = false
	for i, row := range cave.currentRock.layout {
		for j, block := range row {
			if block != ' ' {
				cave.layout[i+cave.currentRock.pos.Y][j+cave.currentRock.pos.X] = '#'
			}
		}
	}
}

// topLineMask returns a binary representation of the top line
func (cave Cave) topLineMask() (mask int) {
	if len(cave.layout) == 0 {
		return (1 << cave.width) - 1
	}

	for i, r := range cave.layout[0] {
		if r == '#' {
			mask |= 1 << i
		}
	}
	return
}

// topLineState returns a unique value describing the top row and the next piece and jet
func (cave Cave) topLineState() int {
	r := len(cave.rocks.Values())
	j := len(cave.jets.Values())
	return cave.topLineMask()*r*j + cave.rocks.Index()*j + cave.jets.Index()
}

// simulateFall simulates the fall of n rocks, and return the height of the cave after that
func (cave *Cave) simulateFall(n int) int {
	for i := len(cave.heights); i < n+len(cave.heights); i++ {
		// get a new rock at pos (0, 2)
		cave.currentRock = cave.rocks.Next()
		cave.currentRock.visible = true
		cave.currentRock.pos.X = 2
		cave.currentRock.pos.Y = 0

		// free enough space to allow rock to appear
		for i := 0; i < 3+len(cave.currentRock.layout); i++ {
			cave.layout = append([][]byte{make([]byte, cave.width)}, cave.layout...)
		}

		// move rock by jet
		for {
			switch cave.jets.Next() {
			case '>':
				// move rock right
				if cave.currentRockCanGoesRight() {
					cave.currentRock.pos.X++
				}
			case '<':
				// move rock left
				if cave.currentRockCanGoesLeft() {
					cave.currentRock.pos.X--
				}
			}

			// move rock down
			if cave.currentRockCanGoesDown() {
				cave.currentRock.pos.Y++
			} else {
				break
			}
		}

		// stop the rock
		cave.stopCurrentRock()

		// remove free rows
		for cave.topLineMask() == 0 {
			cave.layout = cave.layout[1:]
		}

		// remove rows no longer useful from bottom
		max := cave.maxUnreachableHeight()
		cave.offset += len(cave.layout) - max
		cave.layout = cave.layout[:max]

		// save heights
		cave.heights = append(cave.heights, cave.height())

		// save row mask
		cave.rows = append(cave.rows, cave.topLineMask())

		// save state for cycle detection
		state := cave.topLineState()
		prev, cycleDetected := cave.states[state]
		cave.states[state] = i

		// in case of detected cycle
		if cycleDetected {
			// ensure the cycle is valid comparing all previous state in period
			cycleValid := true
			for reverse := 1; reverse <= prev && i-reverse > prev; reverse++ {
				if cave.rows[prev-reverse] != cave.rows[i-reverse] {
					cycleValid = false
					break
				}
			}

			if cycleValid {
				// get current height
				currentHeight := cave.heights[i]

				// add complete remaining periods
				period := i - prev
				remainingRocks := (n - i - 1)
				completePeriodsRemaining := remainingRocks / period
				heightPerPeriod := cave.heights[i] - cave.heights[prev]
				currentHeight += completePeriodsRemaining * heightPerPeriod

				// add remaining rocks
				remainingRocks %= period
				currentHeight += cave.heights[prev+remainingRocks] - cave.heights[prev]

				return currentHeight
			}
		}
	}

	return cave.height()
}

// getNewCave initialises a new cave from lists of rocks and jets
func getNewCave(width int, rocks [][][]byte, input string) *Cave {
	return &Cave{
		width:  width,
		rocks:  utils.NewCyclicArray(utils.ArrayMap(rocks, func(rock [][]byte) Rock { return Rock{layout: rock} })...),
		jets:   utils.NewCyclicArray([]byte(input)...),
		states: make(map[int]int),
	}
}

func main() {
	fmt.Println("--- 2022 Day 17: Pyroclastic Flow ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	input := inputs[0]

	////////////////////////////////////////

	// 3068
	fmt.Println("Part 1:", getNewCave(7, rocks[:], input).simulateFall(2022))

	////////////////////////////////////////

	// 1514285714288
	fmt.Println("Part 2:", getNewCave(7, rocks[:], input).simulateFall(1_000_000_000_000))
}
