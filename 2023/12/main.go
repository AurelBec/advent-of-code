// https://adventofcode.com/2023/day/12

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

// see https://github.com/jonathanpaulson/AdventOfCode/blob/master/2023/12.py

var cache map[tuple]int

type tuple struct{ r, b, current int }

type Record struct {
	record string
	blocks []int
}

func (r Record) getArrangementsCount(...int) int {
	cache = make(map[tuple]int, len(r.record)*len(r.record)*len(r.blocks))
	return countArrangements(r.record, r.blocks, 0, 0, 0)
}

func (r Record) getUnfoldedArrangementsCount(...int) int {
	cache = make(map[tuple]int, len(r.record)*len(r.record)*len(r.blocks))

	var record string
	var blocks []int
	for i := 0; i < 5; i++ {
		record += "?" + r.record
		blocks = append(blocks, r.blocks...)
	}
	record = record[1:]

	return countArrangements(record, blocks, 0, 0, 0)
}

func countArrangements(record string, blocks []int, recordPos, blockPos, current int) int {
	key := tuple{recordPos, blockPos, current}
	if cache, found := cache[key]; found {
		return cache
	}

	if recordPos == len(record) {
		if blockPos == len(blocks) && current == 0 {
			return 1
		} else if blockPos == len(blocks)-1 && blocks[blockPos] == current {
			return 1
		} else {
			return 0
		}
	}

	count := 0
	for _, c := range []byte{'.', '#'} {
		if record[recordPos] != c && record[recordPos] != '?' {
			continue
		}

		if c == '.' && current == 0 {
			count += countArrangements(record, blocks, recordPos+1, blockPos, 0)
		} else if c == '.' && current > 0 && blockPos < len(blocks) && blocks[blockPos] == current {
			count += countArrangements(record, blocks, recordPos+1, blockPos+1, 0)
		} else if c == '#' {
			count += countArrangements(record, blocks, recordPos+1, blockPos, current+1)
		}
	}
	cache[key] = count
	return count
}

func parseRecords(inputs []string) []Record {
	records := make([]Record, len(inputs))
	for i, input := range inputs {
		record, list, _ := strings.Cut(input, " ")
		records[i].record = record
		records[i].blocks = utils.FastNumbers(strings.ReplaceAll(list, ",", " "))
	}
	return records
}

func main() {
	fmt.Println("--- 2023 Day 12: Hot Springs ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	records := parseRecords(inputs)

	////////////////////////////////////////

	// 21
	fmt.Println("Part 1:", utils.SumFunc(records, Record.getArrangementsCount))

	////////////////////////////////////////

	// 525152
	fmt.Println("Part 2:", utils.SumFunc(records, Record.getUnfoldedArrangementsCount))
}
