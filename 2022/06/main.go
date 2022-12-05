// https://adventofcode.com/2022/day/6

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

// hasDuplicates return whether a character is present more than once in a string
func hasDuplicates(s string) bool {
	for _, r := range s {
		if strings.Count(s, string(r)) > 1 {
			return true
		}
	}
	return false
}

func getBuffer(input string, window int) (buffer int) {
	for buffer = window; buffer < len(input)-1 && hasDuplicates(input[buffer-window:buffer]); buffer++ {
	}
	return
}

func main() {
	fmt.Println("--- 2022 Day 6: Tuning Trouble ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	////////////////////////////////////////

	// 7
	fmt.Println("Part 1:", getBuffer(inputs[0], 4))

	////////////////////////////////////////

	// 19
	fmt.Println("Part 2:", getBuffer(inputs[0], 14))
}
