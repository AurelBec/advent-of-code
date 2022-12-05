// https://adventofcode.com/2022/day/6

package main

import (
	"fmt"
	"strings"
	"time"
)

var input = "mjqjpqmgbljsphdztnvjfqwrcgsmlb"

// hasDuplicates return whether a character is present more than once in a string
func hasDuplicates(s string) bool {
	for _, r := range s {
		if strings.Count(s, string(r)) > 1 {
			return true
		}
	}
	return false
}

func main() {
	defer func(start time.Time) { fmt.Println("took:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	////////////////////////////////////////

	buffer := 0

	window := 4
	for buffer = window; buffer < len(input)-1 && hasDuplicates(input[buffer-window:buffer]); buffer++ {
	}

	// 7
	fmt.Println("part1:", buffer)

	////////////////////////////////////////

	buffer = 0

	window = 14
	for buffer = window; buffer < len(input)-1 && hasDuplicates(input[buffer-window:buffer]); buffer++ {
	}

	// 19
	fmt.Println("part2:", buffer)
}
