// https://adventofcode.com/2022/day/4

package main

import (
	"fmt"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

var inputs = [...]string{
	"2-4,6-8",
	"2-3,4-5",
	"5-7,7-9",
	"2-8,3-7",
	"6-6,4-6",
	"2-6,4-8",
}

// parseAssignments returns the section IDs affected to a pair
func parseAssignments(inputs []string) (assignments [][2]utils.Interval[int]) {
	assignments = make([][2]utils.Interval[int], len(inputs))
	for i, pair := range inputs {
		fmt.Sscanf(pair, "%d-%d,%d-%d", &assignments[i][0].Min, &assignments[i][0].Max, &assignments[i][1].Min, &assignments[i][1].Max)
	}
	return
}

func main() {
	defer func(start time.Time) { fmt.Println("took:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init

	assignments := parseAssignments(inputs[:])

	////////////////////////////////////////

	fullyContained := 0
	for _, assignment := range assignments {
		if assignment[0].Contains(assignment[1]) || assignment[1].Contains(assignment[0]) {
			fullyContained++
		}
	}

	// 2
	fmt.Println("part1:", fullyContained)

	////////////////////////////////////////

	overlapsAtAll := 0
	for _, assignment := range assignments {
		if assignment[0].Overlaps(assignment[1]) {
			overlapsAtAll++
		}
	}

	// 4
	fmt.Println("part2:", overlapsAtAll)
}
