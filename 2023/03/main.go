// https://adventofcode.com/2023/day/3

package main

import (
	"fmt"
	"time"
	"unicode"

	"github.com/aurelbec/advent-of-code/utils"
)

type engineSchematic struct {
	symbols map[int]map[int]bool
	digits  map[int]map[int]*partNumber

	partNumber []*partNumber
	gears      [][2]int
}

type partNumber struct {
	value int
	valid bool
}

func (es engineSchematic) getPartNumbersSum() int {
	sum := 0
	for _, partNumber := range es.partNumber {
		if partNumber.valid {
			sum += partNumber.value
		}
	}
	return sum
}

func (es engineSchematic) getGearRatiosSum() int {
	sum := 0
	for _, gear := range es.gears {
		partNumbers := make(map[*partNumber]struct{}, 3)
		for i := gear[0] - 1; i <= gear[0]+1; i++ {
			for j := gear[1] - 1; j <= gear[1]+1; j++ {
				partNumbers[es.digits[i][j]] = struct{}{}
			}
		}

		delete(partNumbers, nil)
		if len(partNumbers) == 2 {
			ratio := 1
			for partNumber := range partNumbers {
				ratio *= partNumber.value
			}
			sum += ratio
		}
	}
	return sum
}

func (es engineSchematic) isAdjacentToSymbol(i, j int) bool {
	return es.symbols[i-1][j-1] || es.symbols[i-1][j] || es.symbols[i-1][j+1] ||
		es.symbols[i][j-1] || es.symbols[i][j+1] ||
		es.symbols[i+1][j-1] || es.symbols[i+1][j] || es.symbols[i+1][j+1]
}

func parseEngineSchematic(inputs []string) engineSchematic {
	engineSchematic := engineSchematic{
		symbols: make(map[int]map[int]bool, len(inputs)),
		digits:  make(map[int]map[int]*partNumber, len(inputs)),
	}

	for i, input := range inputs {
		engineSchematic.symbols[i] = make(map[int]bool, len(input))
		engineSchematic.digits[i] = make(map[int]*partNumber, len(input))
		for j, r := range input {
			if !unicode.IsDigit(r) && r != '.' {
				engineSchematic.symbols[i][j] = true
				if r == '*' {
					engineSchematic.gears = append(engineSchematic.gears, [2]int{i, j})
				}
			}
		}
	}

	pn := &partNumber{}
	for i, input := range inputs {
		for j, r := range input {
			if unicode.IsDigit(r) {
				pn.value = pn.value*10 + int(r-'0')
				pn.valid = pn.valid || engineSchematic.isAdjacentToSymbol(i, j)
				engineSchematic.digits[i][j] = pn
			} else if pn.value > 0 {
				engineSchematic.partNumber = append(engineSchematic.partNumber, pn)
				pn = &partNumber{}
			}
		}

		if pn.value > 0 {
			engineSchematic.partNumber = append(engineSchematic.partNumber, pn)
			pn = &partNumber{}
		}
	}

	return engineSchematic
}

func main() {
	fmt.Println("--- 2023 Day 3: Gear Ratios ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	engineSchematic := parseEngineSchematic(inputs)

	////////////////////////////////////////

	// 4361
	fmt.Println("Part 1:", engineSchematic.getPartNumbersSum())

	////////////////////////////////////////

	// 467835
	fmt.Println("Part 2:", engineSchematic.getGearRatiosSum())
}

//
