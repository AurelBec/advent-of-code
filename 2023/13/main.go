// https://adventofcode.com/2023/day/13

package main

import (
	"fmt"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

type Pattern [][]byte

func (p Pattern) getSummary(tolerance int) int {
	if line := findReflection(p, tolerance); line > 0 {
		return 100 * line
	}
	if line := findReflection(utils.ArrayTranspose(p), tolerance); line > 0 {
		return line
	}
	return 0
}

func findReflection(pattern [][]byte, tolerance int) int {
next:
	for i := 1; i < len(pattern); i++ {
		tolerance := tolerance

		// find first adjacent match
		if !equals(pattern[i-1], pattern[i], &tolerance) {
			continue next
		}
		line := i

		// check next lines with the reflection
		for j := 1; 0 <= line-1-j && line+j < len(pattern); j++ {
			if !equals(pattern[line-1-j], pattern[line+j], &tolerance) {
				continue next
			}
		}

		if tolerance == 0 {
			return line
		}
	}
	return 0
}

func equals(a, b []byte, tolerance *int) bool {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			if *tolerance == 0 {
				return false
			}
			*tolerance--
		}
	}
	return true
}

func parsePatterns(inputs []string) []Pattern {
	patterns := make([]Pattern, 1)

	for _, input := range inputs {
		if input == "" {
			patterns = append(patterns, Pattern{})
		} else {
			patterns[len(patterns)-1] = append(patterns[len(patterns)-1], []byte(input))
		}
	}

	return patterns
}

func main() {
	fmt.Println("--- 2023 Day 13: Point of Incidence ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	patterns := parsePatterns(inputs)

	////////////////////////////////////////

	// 405
	fmt.Println("Part 1:", utils.SumFunc(patterns, func(p Pattern, _ ...int) int { return p.getSummary(0) }))

	////////////////////////////////////////

	// 400
	fmt.Println("Part 2:", utils.SumFunc(patterns, func(p Pattern, _ ...int) int { return p.getSummary(1) }))
}
