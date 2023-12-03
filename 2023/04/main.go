// https://adventofcode.com/2023/day/4

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

type pile struct {
	cards []card
}

type card struct {
	id             int
	numbers        []string
	winningNumbers map[string]bool
}

func (c card) getMatchingNumbers() int {
	n := 0
	for _, number := range c.numbers {
		if c.winningNumbers[number] {
			n++
		}
	}
	return n
}

func (c card) getPoints() int {
	if n := c.getMatchingNumbers(); n == 0 {
		return 0
	} else {
		return 1 << (n - 1)
	}
}

func (p pile) getPoints() int {
	points := 0
	for _, card := range p.cards {
		points += card.getPoints()
	}
	return points
}

func (p pile) getTotalCards() int {
	total := 0
	copies := make(map[int]int, len(p.cards))
	for _, card := range p.cards {
		n := copies[card.id] + 1
		for id := card.id + 1; id <= card.id+card.getMatchingNumbers(); id++ {
			copies[id] += n
		}
		total += n
	}
	return total
}

func parsePile(inputs []string) pile {
	pile := pile{
		cards: make([]card, 0, len(inputs)),
	}

	for id, input := range inputs {
		information := strings.Split(strings.Split(input, ":")[1], "|")
		winningNumbers := utils.Numbers(information[0])
		numbers := utils.Numbers(information[1])

		card := card{
			id:             id + 1,
			winningNumbers: make(map[string]bool, len(winningNumbers)),
			numbers:        numbers,
		}

		for _, winningNumber := range winningNumbers {
			card.winningNumbers[winningNumber] = true
		}

		pile.cards = append(pile.cards, card)
	}

	return pile
}

func main() {
	fmt.Println("--- 2023 Day 4: Scratchcards ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	pile := parsePile(inputs)

	////////////////////////////////////////

	// 13
	fmt.Println("Part 1:", pile.getPoints())
	//

	////////////////////////////////////////

	// 30
	fmt.Println("Part 2:", pile.getTotalCards())
}
