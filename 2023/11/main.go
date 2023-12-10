// https://adventofcode.com/2023/day/11

package main

import (
	"fmt"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

type Universe struct {
	galaxies []utils.Location2D[int]

	usedRows map[int]bool
	usedCols map[int]bool
}

func (u Universe) getGalaxiesDistances(expansionFactor int) []int {
	distances := make([]int, 0, len(u.galaxies)*len(u.galaxies))
	for i := 0; i < len(u.galaxies); i++ {
		start := u.galaxies[i]
		for j := i + 1; j < len(u.galaxies); j++ {
			end := u.galaxies[j]
			dist := 0
			for x := start.X; x != end.X; x += utils.Sign(end.X - start.X) {
				if u.usedCols[x] {
					dist += 1
				} else {
					dist += expansionFactor
				}
			}
			for y := start.Y; y != end.Y; y += utils.Sign(end.Y - start.Y) {
				if u.usedRows[y] {
					dist += 1
				} else {
					dist += expansionFactor
				}
			}
			distances = append(distances, dist)
		}
	}
	return distances
}

func parseUniverse(inputs []string) Universe {
	universe := Universe{
		galaxies: make([]utils.Location2D[int], 0),
		usedRows: make(map[int]bool),
		usedCols: make(map[int]bool),
	}

	for y, input := range inputs {
		for x, data := range input {
			if data == '#' {
				universe.galaxies = append(universe.galaxies, utils.NewLocation2D(x, y))
				universe.usedCols[x] = true
				universe.usedRows[y] = true
			}
		}
	}

	return universe
}

func main() {
	fmt.Println("--- 2023 Day 11: Cosmic Expansion ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	universe := parseUniverse(inputs)

	////////////////////////////////////////

	// 374
	fmt.Println("Part 1:", utils.Sum(universe.getGalaxiesDistances(2)))

	////////////////////////////////////////

	// 10: 1030
	// 100: 8410
	// 1000000: 82000210
	fmt.Println("Part 2:", utils.Sum(universe.getGalaxiesDistances(1_000_000)))
}
