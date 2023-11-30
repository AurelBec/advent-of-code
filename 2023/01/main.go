// https://adventofcode.com/2023/day/1

package main

import (
	"fmt"
	"slices"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

func parseDigit(word string) (int, bool) {
	switch word {
	case "0", "zero", "orez":
		return 0, true
	case "1", "one", "eno":
		return 1, true
	case "2", "two", "owt":
		return 2, true
	case "3", "three", "eerht":
		return 3, true
	case "4", "four", "ruof":
		return 4, true
	case "5", "five", "evif":
		return 5, true
	case "6", "six", "xis":
		return 6, true
	case "7", "seven", "neves":
		return 7, true
	case "8", "eight", "thgie":
		return 8, true
	case "9", "nine", "enin":
		return 9, true
	default:
		return -1, false
	}
}

func findDigits(input []byte) [2]int {
	digits := [2]int{-1, -1}
	for i := 0; i < len(input); i++ {
		for j := i + 1; j <= min(i+5, len(input)); j++ {
			if digit, found := parseDigit(string(input[i:j])); found {
				if digits[1] < 0 {
					digits[1] = digit
				}
				if j == i+1 {
					digits[0] = digit
					return digits
				}
			}
		}
	}
	return digits
}

func parseCalibrationValues(inputs []string) [][2][2]int {
	calibrationValues := make([][2][2]int, len(inputs))
	for i, input := range inputs {
		input := []byte(input)
		calibrationValues[i][0] = findDigits(input)
		slices.Reverse(input)
		calibrationValues[i][1] = findDigits(input)
	}
	return calibrationValues
}

func getCalibrationValuesSum(calibrationValues [][2][2]int, digit int) int {
	sum := 0
	for _, calibrationValue := range calibrationValues {
		sum += calibrationValue[0][digit]*10 + calibrationValue[1][digit]
	}
	return sum
}

func main() {
	fmt.Println("--- 2023 Day 1: Trebuchet?! ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs1 := utils.MustReadInput("example1.txt")
	inputs2 := utils.MustReadInput("example2.txt")

	////////////////////////////////////////

	// 142
	fmt.Println("Part 1:", getCalibrationValuesSum(parseCalibrationValues(inputs1), 0))

	////////////////////////////////////////

	// 281
	fmt.Println("Part 2:", getCalibrationValuesSum(parseCalibrationValues(inputs2), 1))
}
