// https://adventofcode.com/2022/day/1

package main

import (
	"fmt"
	"strconv"
	"time"
)

var inputs = [...]string{
	"1000",
	"2000",
	"3000",
	"",
	"4000",
	"",
	"5000",
	"6000",
	"",
	"7000",
	"8000",
	"9000",
	"",
	"10000",
}

func main() {
	defer func(start time.Time) { fmt.Println("took:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	////////////////////////////////////////

	maximum := 0
	current := 0

	updateMaximum := func(sum int) {
		if sum > maximum {
			maximum = sum
		}
	}

	for _, food := range inputs {
		if food == "" {
			updateMaximum(current)
			current = 0
			continue
		}

		cal, _ := strconv.Atoi(food)
		current += cal
	}
	updateMaximum(current)

	// 24000
	fmt.Println("part1:", maximum)

	////////////////////////////////////////

	size := 3
	current = 0
	maximums := make([]int, size)

	updateMaximums := func(sum int) {
		for i := 0; i < size; i++ {
			if sum > maximums[i] {
				for j := size - 1; j > i; j-- {
					maximums[j] = maximums[j-1]
				}
				maximums[i] = sum
				return
			}
		}
	}

	for _, food := range inputs {
		if food == "" {
			updateMaximums(current)
			current = 0
			continue
		}

		cal, _ := strconv.Atoi(food)
		current += cal
	}
	updateMaximums(current)

	// 45000
	fmt.Println("part2:", maximums, "->", func() (sum int) {
		for _, v := range maximums {
			sum += v
		}
		return
	}())
}
