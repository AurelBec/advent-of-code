// https://adventofcode.com/2022/day/5

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

// getStacks parse the first part of the input, and return the stacks
func getStacks(inputs []string) []string {
	// creates the right number of stacks using number of number on last line
	stacks := make([]string, len(strings.Fields(inputs[len(inputs)-1])))

	// parse input to fill stacks
	for _, input := range inputs[:len(inputs)-1] {
		for i := range stacks {
			if len(input) <= i*4+1 {
				break
			}
			if item := input[i*4+1]; item != ' ' {
				stacks[i] += string(item)
			}
		}
	}

	return stacks
}

// getMoves parse the second part of the input, and return the moves to execute
func getMoves(inputs []string) [][3]int {
	moves := make([][3]int, len(inputs))
	for i, input := range inputs {
		fmt.Sscanf(input, "move %v from %v to %v", &moves[i][0], &moves[i][1], &moves[i][2])
	}
	return moves
}

// reverse reverse a string
func reverse(in string) (out string) {
	for _, r := range in {
		out = string(r) + out
	}
	return
}

// firstChar returns the first char of input string
func firstChar(s string, _ ...int) string {
	return string(s[0])
}

func main() {
	fmt.Println("--- 2022 Day 5: Supply Stacks ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	sep := 0
	for ; inputs[sep] != ""; sep++ {
	}

	moves := getMoves(inputs[sep+1:])

	////////////////////////////////////////

	// separate stacks description and rearrangement procedure

	stacks := getStacks(inputs[:sep])
	for _, move := range moves {
		n, from, to := move[0], move[1]-1, move[2]-1
		stacks[to] = reverse(stacks[from][:n]) + stacks[to]
		stacks[from] = stacks[from][n:]
	}

	// CMZ
	fmt.Println("Part 1:", utils.SumFunc(stacks, firstChar))

	////////////////////////////////////////

	stacks = getStacks(inputs[:sep])
	for _, move := range moves {
		n, from, to := move[0], move[1]-1, move[2]-1
		stacks[to] = stacks[from][:n] + stacks[to]
		stacks[from] = stacks[from][n:]
	}

	// MCD
	fmt.Println("Part 2:", utils.SumFunc(stacks, firstChar))
}
