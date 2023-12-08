// https://adventofcode.com/2023/day/9

package main

import (
	"fmt"
	"slices"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

func getNextValue(history []int, _ ...int) (extrapolation int) {
	history = slices.Clone(history)
	for {
		allZeros := true
		extrapolation += history[len(history)-1]
		for i := 1; i < len(history); i++ {
			history[i-1] = history[i] - history[i-1]
			allZeros = allZeros && history[i-1] == 0
		}
		if allZeros {
			return
		}
		history = history[:len(history)-1]
	}
}

func getPastValue(history []int, _ ...int) (extrapolation int) {
	history = slices.Clone(history)
	slices.Reverse(history)
	return getNextValue(history)
}

func parseHistories(inputs []string) [][]int {
	return utils.ArrayMap(inputs, utils.FastNumbers)
}

func main() {
	fmt.Println("--- 2023 Day 9: Mirage Maintenance ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	histories := parseHistories(inputs)

	////////////////////////////////////////

	// 114
	fmt.Println("Part 1:", utils.SumFunc(histories, getNextValue))

	////////////////////////////////////////

	// 2
	fmt.Println("Part 2:", utils.SumFunc(histories, getPastValue))
}
