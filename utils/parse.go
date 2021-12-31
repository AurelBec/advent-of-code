package utils

import (
	"regexp"
	"strconv"
	"strings"
)

var extractNumbers = regexp.MustCompile(`-?\d+`)
var extractWords = regexp.MustCompile(`[a-zA-Z]+`)

// Numbers return the list of numbers contained in the input string
func Numbers(s string) []string {
	return extractNumbers.FindAllString(s, -1)
}

// FastNumbers return the list of numbers contained in the input string
// input string need to contains only digits separated by blanks
func FastNumbers(s string) []int {
	return ArrayMap(strings.Fields(s), MustInt)
}

// Words return the list of words contained in the input string
func Words(s string) []string {
	return extractWords.FindAllString(s, -1)
}

// MustInt interprets a string s as an int. No error check is done
func MustInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
