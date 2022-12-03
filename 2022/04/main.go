// https://adventofcode.com/2022/day/4

package main

import (
	"fmt"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

// parseAssignments returns the section IDs affected to a pair
func parseAssignments(inputs []string) (assignments [][2]utils.Interval[int]) {
	assignments = make([][2]utils.Interval[int], len(inputs))
	for i, pair := range inputs {
		fmt.Sscanf(pair, "%d-%d,%d-%d", &assignments[i][0].Min, &assignments[i][0].Max, &assignments[i][1].Min, &assignments[i][1].Max)
	}
	return
}

func main() {
	fmt.Println("--- 2022 Day 4: Camp Cleanup ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	assignments := parseAssignments(inputs)

	////////////////////////////////////////

	fullyContained := 0
	for _, assignment := range assignments {
		if assignment[0].Contains(assignment[1]) || assignment[1].Contains(assignment[0]) {
			fullyContained++
		}
	}

	// 2
	fmt.Println("Part 1:", fullyContained)

	////////////////////////////////////////

	overlapsAtAll := 0
	for _, assignment := range assignments {
		if assignment[0].Overlaps(assignment[1]) {
			overlapsAtAll++
		}
	}

	// 4
	fmt.Println("Part 2:", overlapsAtAll)
}
