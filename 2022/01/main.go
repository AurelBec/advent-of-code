// https://adventofcode.com/2022/day/1

package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

func main() {
	fmt.Println("--- 2022 Day 1: Calorie Counting ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	sums, current := make([]int, 0), 0
	for _, food := range inputs {
		if food == "" {
			sums, current = append(sums, current), 0
		} else {
			current += utils.MustInt(food)
		}
	}
	sums = append(sums, current)

	sort.Sort(sort.Reverse(sort.IntSlice(sums)))

	////////////////////////////////////////

	// 24000
	fmt.Println("Part 1:", sums[0])

	////////////////////////////////////////

	// 45000
	fmt.Println("Part 2:", sums[0]+sums[1]+sums[2])
}
