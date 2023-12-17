// https://adventofcode.com/2023/day/18

package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

type Instructions []Instruction

type Instruction struct {
	dir string
	n   int
}

func (i Instructions) getLagoonArea() int {
	// https://en.wikipedia.org/wiki/Shoelace_formula
	// https://en.wikipedia.org/wiki/Pick%27s_theorem
	area := 0
	perimeter := 0
	current := utils.NewLocation2D(0, 0)
	for _, instruction := range i {
		prev := current
		dist := instruction.n
		switch instruction.dir {
		case "U", "3":
			current = current.MovedBy(0, -dist)
		case "R", "0":
			current = current.MovedBy(+dist, 0)
		case "D", "1":
			current = current.MovedBy(0, +dist)
		case "L", "2":
			current = current.MovedBy(-dist, 0)
		}
		area += (prev.Y + current.Y) * (prev.X - current.X)
		perimeter += dist
	}
	return area/2 + perimeter/2 + 1
}

func parseInstructions(inputs []string) (Instructions, Instructions) {
	instructions := make(Instructions, len(inputs))
	correctedInstructions := make(Instructions, len(inputs))
	for i, input := range inputs {
		fmt.Sscanf(input, "%v %v (#%v)", &instructions[i].dir, &instructions[i].n, &input)

		n, _ := strconv.ParseInt(input[:5], 16, 0)
		correctedInstructions[i].n = int(n)
		correctedInstructions[i].dir = input[5:6]
	}
	return instructions, correctedInstructions
}

func main() {
	fmt.Println("--- 2023 Day 18: Lavaduct Lagoon ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	instructions, correctedInstructions := parseInstructions(inputs)

	////////////////////////////////////////

	// 62
	fmt.Println("Part 1:", instructions.getLagoonArea())

	////////////////////////////////////////

	// 952408144115
	fmt.Println("Part 2:", correctedInstructions.getLagoonArea())
}
