// https://adventofcode.com/2023/day/15

package main

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

type Lens struct {
	focal      int
	insertedAt int
}

func getHash(input string, _ ...int) int {
	hash := 0
	for _, r := range input {
		hash = (hash + int(r)) * 17 % 256
	}
	return hash
}

func getFocusingPower(boxes map[int]map[string]Lens) int {
	focusingPower := 0
	for box := 0; box < 256; box++ {
		lenses := utils.MapValues(boxes[box])
		slices.SortFunc(lenses, func(lhs, rhs Lens) int { return lhs.insertedAt - rhs.insertedAt })
		for slot, lens := range lenses {
			focusingPower += (box + 1) * (slot + 1) * lens.focal
		}
	}
	return focusingPower
}

func fillBoxes(steps []string) map[int]map[string]Lens {
	boxes := make(map[int]map[string]Lens)
	for i, step := range steps {
		values := strings.FieldsFunc(step, func(r rune) bool { return r == '=' || r == '-' })
		label := values[0]
		box := getHash(label)
		switch len(values) {
		case 1:
			delete(boxes[box], label)
		case 2:
			if len(boxes[box]) == 0 {
				boxes[box] = map[string]Lens{label: {focal: utils.MustInt(values[1]), insertedAt: i}}
			} else if lens, found := boxes[box][label]; !found {
				boxes[box][label] = Lens{focal: utils.MustInt(values[1]), insertedAt: i}
			} else {
				lens.focal = utils.MustInt(values[1])
				boxes[box][label] = lens
			}
		}
	}
	return boxes
}

func parseSteps(inputs []string) []string {
	steps := make([]string, 0)
	for _, input := range inputs {
		steps = append(steps, strings.Split(input, ",")...)
	}
	return steps
}

func main() {
	fmt.Println("--- 2023 Day 15: Lens Library ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	steps := parseSteps(inputs)

	////////////////////////////////////////

	// 1320
	fmt.Println("Part 1:", utils.SumFunc(steps, getHash))

	////////////////////////////////////////

	// 145
	fmt.Println("Part 2:", getFocusingPower(fillBoxes(steps)))
}
