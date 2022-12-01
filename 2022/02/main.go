// https://adventofcode.com/2022/day/2

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

const (
	rock = iota
	paper
	scissors
)

const (
	draw = iota
	loss
	win
)

// convert converts a choice into the corresponding figure
func convert(choice string) int {
	switch choice {
	case "A", "X":
		return rock
	case "B", "Y":
		return paper
	case "C", "Z":
		return scissors
	}
	return -1
}

// points return the number of points scored for this play
func points(opponent, me int) int {
	score := 0

	switch me {
	case rock:
		score += 1
	case paper:
		score += 2
	case scissors:
		score += 3
	}

	switch play(opponent-me, 3) {
	case draw:
		score += 3
	case loss:
		score += 0
	case win:
		score += 6
	}

	return score
}

// play executes a round, and return whether it results in a win, draw or loss
func play(i, n int) int {
	return ((i % n) + n) % n
}

func main() {
	fmt.Println("--- 2022 Day 2: Rock Paper Scissors ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	////////////////////////////////////////

	score := 0
	for _, input := range inputs {
		plays := strings.Fields(input)
		score += points(convert(plays[0]), convert(plays[1]))
	}

	// 15
	fmt.Println("Part 1:", score)

	////////////////////////////////////////

	score = 0
	for _, input := range inputs {
		plays := strings.Fields(input)
		switch plays[1] {
		case "X": // need to lose
			score += points(convert(plays[0]), play(convert(plays[0])-1, 3))
		case "Y": // need to draw
			score += points(convert(plays[0]), convert(plays[0]))
		case "Z": // need to win
			score += points(convert(plays[0]), play(convert(plays[0])+1, 3))
		}
	}

	// 12
	fmt.Println("Part 2:", score)
}
