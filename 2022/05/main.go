// https://adventofcode.com/2022/day/5

package main

import (
	"fmt"
	"strings"
	"time"
)

var inputs = [...]string{
	"    [D]    ",
	"[N] [C]    ",
	"[Z] [M] [P]",
	" 1   2   3 ",
	"",
	"move 1 from 2 to 1",
	"move 3 from 1 to 3",
	"move 2 from 2 to 1",
	"move 1 from 1 to 2",
}

// getStacks parse the first part of the input, and return the stacks
func getStacks(inputs []string) []string {
	// creates the right number of stacks using number of number on last line
	stacks := make([]string, len(strings.Fields(inputs[len(inputs)-1])))

	// parse input to fill stacks
	for _, input := range inputs[:len(inputs)-1] {
		for i := range stacks {
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

func main() {
	defer func(start time.Time) { fmt.Println("took:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	////////////////////////////////////////

	result := ""

	// separate stacks description and rearrangement procedure
	sep := 0
	for inputs[sep] != "" {
		sep++
	}

	stacks := getStacks(inputs[:sep])
	moves := getMoves(inputs[sep+1:])
	for _, move := range moves {
		n, from, to := move[0], move[1]-1, move[2]-1
		stacks[to] = reverse(stacks[from][:n]) + stacks[to]
		stacks[from] = stacks[from][n:]
	}

	for _, stack := range stacks {
		result += string(stack[0])
	}

	// CMZ
	fmt.Println("part1:", result)

	////////////////////////////////////////

	result = ""

	stacks = getStacks(inputs[:sep])
	for _, move := range moves {
		n, from, to := move[0], move[1]-1, move[2]-1
		stacks[to] = stacks[from][:n] + stacks[to]
		stacks[from] = stacks[from][n:]
	}

	for _, stack := range stacks {
		result += string(stack[0])
	}

	// MCD
	fmt.Println("part2:", result)
}
