// https://adventofcode.com/2022/day/3

package main

import (
	"fmt"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

// getItems returns a map of all items contained in a compartment
func getItems(compartment string) map[rune]rune {
	items := make(map[rune]rune, len(compartment))
	for _, item := range compartment {
		items[item] = item
	}
	return items
}

// getCommon return the first common element between compartments
func getCommon(compartments ...string) rune {
	if len(compartments) < 2 {
		return -1
	}

	itemsPerCompartments := make([]map[rune]rune, 0, len(compartments))
	for _, compartment := range compartments {
		itemsPerCompartments = append(itemsPerCompartments, getItems(compartment))
	}

	// iterate on items in first compartment
	// and ensure it is present in others
	for item := range itemsPerCompartments[0] {
		isCommon := true
		for _, itemsPerCompartment := range itemsPerCompartments[1:] {
			if _, found := itemsPerCompartment[item]; !found {
				isCommon = false
				break
			}
		}
		if isCommon {
			return item
		}
	}

	return -1
}

// getPriority returns the item priority
func getPriority(item rune) int {
	switch {
	case 'a' <= item && item <= 'z':
		return int(item-'a') + 1
	case 'A' <= item && item <= 'Z':
		return int(item-'A') + 27
	default:
		return -1
	}
}

func main() {
	fmt.Println("--- 2022 Day 3: Rucksack Reorganization ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	////////////////////////////////////////

	priorities := 0
	for _, rucksack := range inputs {
		sep := len(rucksack) / 2
		common := getCommon(rucksack[:sep], rucksack[sep:])
		priorities += getPriority(common)
	}

	// 157
	fmt.Println("Part 1:", priorities)

	////////////////////////////////////////

	priorities = 0
	for i := 0; i < len(inputs); i += 3 {
		common := getCommon(inputs[i], inputs[i+1], inputs[i+2])
		priorities += getPriority(common)
	}

	// 70
	fmt.Println("Part 2:", priorities)
}
