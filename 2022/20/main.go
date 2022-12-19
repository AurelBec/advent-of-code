// https://adventofcode.com/2022/day/20

package main

import (
	"fmt"
	"slices"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

type Mix struct {
	utils.CyclicArray[int]
	zeroIndex int
}

func (mix Mix) GetCoordinate(index int, _ ...int) int {
	return mix.Get(mix.zeroIndex + index)
}

func mixSequence(sequence []int, decryptionKey, steps int) Mix {
	// initialize arrays to keep trace of indexes
	atNowIs := utils.ArrayIota(0, len(sequence))
	isNowAt := slices.Clone(atNowIs)

	// mix the sequence
	for step := 0; step < steps; step++ {
		for originalIndex, offset := range sequence {
			from := isNowAt[originalIndex]
			to := utils.Mod(from+offset*decryptionKey, len(sequence)-1)
			delta := utils.Sign(to - from)

			// rotate indexes
			for i := from; i != to; i += delta {
				atNowIs[i] = atNowIs[i+delta]
				isNowAt[atNowIs[i+delta]] = i
			}
			atNowIs[to] = originalIndex
			isNowAt[originalIndex] = to
		}
	}

	mix := make([]int, len(sequence))
	for originalIndex, value := range sequence {
		mix[isNowAt[originalIndex]] = value * decryptionKey
	}
	return Mix{utils.NewCyclicArray(mix...), slices.Index(mix, 0)}
}

func parseSequence(inputs []string) []int {
	return utils.ArrayMap(inputs, utils.MustInt)
}

func main() {
	fmt.Println("--- 2022 Day 20: Grove Positioning System ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")
	groveCoordinates := []int{1000, 2000, 3000}

	sequence := parseSequence(inputs)

	////////////////////////////////////////

	mix := mixSequence(sequence, 1, 1)

	// 3
	fmt.Println("Part 1:", utils.SumFunc(groveCoordinates, mix.GetCoordinate))

	////////////////////////////////////////

	mix = mixSequence(sequence, 811589153, 10)

	// 1623178306
	fmt.Println("Part 2:", utils.SumFunc(groveCoordinates, mix.GetCoordinate))
}
