// https://adventofcode.com/

package main

import (
	"fmt"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

func main() {
	fmt.Println("--- ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	for _, input := range inputs {
		fmt.Println(input)
	}

	////////////////////////////////////////

	// 0
	fmt.Println("Part 1:", 0)

	////////////////////////////////////////

	// 0
	fmt.Println("Part 2:", 0)
}
