// https://adventofcode.com/2023/day/7

package main

import (
	"fmt"
	"slices"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

const (
	HighCard = iota
	OnePair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

var cardValues = map[rune][2]int{
	'A': {13, 13},
	'K': {12, 12},
	'Q': {11, 11},
	'J': {10, 00},
	'T': {9, 9},
	'9': {8, 8},
	'8': {7, 7},
	'7': {6, 6},
	'6': {5, 5},
	'5': {4, 4},
	'4': {3, 3},
	'3': {2, 2},
	'2': {1, 1},
}

type Hand struct {
	cards  string
	values [2]int
	kind   [2]int
	bid    int
}

// compare should return:
// a negative number when a < b
// a positive number when a > b
// and zero when a == b
func (a Hand) compareJack(b Hand) int  { return a.compare(b, 0) }
func (a Hand) compareJoker(b Hand) int { return a.compare(b, 1) }
func (a Hand) compare(b Hand, i int) int {
	switch {
	case a.kind[i] > b.kind[i]:
		return 1
	case a.kind[i] < b.kind[i]:
		return -1
	case a.values[i] > b.values[i]:
		return 1
	case a.values[i] < b.values[i]:
		return -1
	default:
		return 0
	}
}

func (hand Hand) totalWinning(rank ...int) int {
	return (rank[0] + 1) * hand.bid
}

func parseHands(inputs []string) []Hand {
	hands := make([]Hand, len(inputs))
	for i, input := range inputs {
		fmt.Sscanf(input, "%s %v", &hands[i].cards, &hands[i].bid)

		cardsMap := make(map[rune]int)
		maxRepeatingCards := 0
		for _, card := range hands[i].cards {
			cardsMap[card]++
			maxRepeatingCards = max(maxRepeatingCards, cardsMap[card])
			hands[i].values[0] = hands[i].values[0]<<4 + cardValues[card][0]
			hands[i].values[1] = hands[i].values[1]<<4 + cardValues[card][1]
		}

		n, j := len(cardsMap), cardsMap['J']
		switch {
		case maxRepeatingCards == 5 && n == 1:
			hands[i].kind = [2]int{FiveOfAKind, FiveOfAKind}
		case maxRepeatingCards == 4 && n == 2:
			hands[i].kind = [2]int{FourOfAKind, FiveOfAKind}
		case maxRepeatingCards == 3 && n == 2:
			hands[i].kind = [2]int{FullHouse, FiveOfAKind}
		case maxRepeatingCards == 3 && n == 3:
			hands[i].kind = [2]int{ThreeOfAKind, FourOfAKind}
		case maxRepeatingCards == 2 && n == 3 && j < 2:
			hands[i].kind = [2]int{TwoPair, FullHouse}
		case maxRepeatingCards == 2 && n == 3 && j == 2:
			hands[i].kind = [2]int{TwoPair, FourOfAKind}
		case maxRepeatingCards == 2 && n == 4:
			hands[i].kind = [2]int{OnePair, ThreeOfAKind}
		case maxRepeatingCards == 1 && n == 5:
			hands[i].kind = [2]int{HighCard, OnePair}
		}
		if j == 0 {
			hands[i].kind[1] = hands[i].kind[0]
		}
	}
	return hands
}

func main() {
	fmt.Println("--- 2023 Day 7: Camel Cards ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	hands := parseHands(inputs)

	////////////////////////////////////////

	slices.SortFunc(hands, Hand.compareJack)

	// 6440
	fmt.Println("Part 1:", utils.SumFunc(hands, Hand.totalWinning))

	////////////////////////////////////////

	slices.SortFunc(hands, Hand.compareJoker)

	// 5905
	fmt.Println("Part 2:", utils.SumFunc(hands, Hand.totalWinning))
}
