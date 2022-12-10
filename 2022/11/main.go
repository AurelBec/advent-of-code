// https://adventofcode.com/2022/day/11

package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

type monkeyID = int

type worryLevel = int

type Monkey struct {
	ID          monkeyID     // monkey ID
	items       []worryLevel // worry level list for each item currently holding in the order they will be inspected
	inspections int          // number of items inspected during all rounds

	testWorryLevel func(worryLevel) worryLevel // update worry level during item inspection

	test            func(worryLevel) bool // use worry level to decide where to throw an item next
	trueTestMonkey  monkeyID              // ID of the monkey which will receive item if the test is true
	falseTestMonkey monkeyID              // ID of the monkey which will receive item if the test is false

	reduceFactor worryLevel // reduction factor (product of all test divisors) to reduce worry level and avoid overflows
}

// inspect updates worry level about item during inspection, and returns to which monkey to item should be thrown
func (monkey *Monkey) inspect(item worryLevel, worried bool) (worryLevel, monkeyID) {
	// increase number of inspected items
	monkey.inspections++

	// update worry level after item inspection
	item = monkey.testWorryLevel(item)

	// if not worry by item inspection, reduce worry level
	if !worried {
		item /= 3
	}

	// trick to reduce worry level value, to avoid overflows
	item %= monkey.reduceFactor

	// determine to which monkey throw the item next
	if monkey.test(item) {
		return item, monkey.trueTestMonkey
	} else {
		return item, monkey.falseTestMonkey
	}
}

// addItem add a new item ath the end of the current list
func (monkey *Monkey) addItem(item worryLevel) {
	monkey.items = append(monkey.items, item)
}

type Monkeys []*Monkey

// playRound executes a complete run of monkeys inspecting and throwing items
func (mks Monkeys) playRound(worried bool) {
	for _, monkey := range mks {
		for _, item := range monkey.items {
			// let the monkey inspect item, and decide which to throw next to
			item, next := monkey.inspect(item, worried)
			// give the item to the selected monkey
			mks[next].addItem(item)
		}
		// consider that all items has been thrown, thus empty the list
		monkey.items = []worryLevel{}
	}
}

// getBusiness returns the business level of the group
func (mks Monkeys) getBusiness() int {
	sort.Slice(mks, func(i, j int) bool { return mks[i].inspections > mks[j].inspections })
	return mks[0].inspections * mks[1].inspections
}

// getMonkeys returns a list of monkeys created using description as inputs
func getMonkeys(inputs []string) (monkeys Monkeys) {
	reduceFactor := worryLevel(1)

	for _, input := range inputs {
		// add a new empty monkey to the list if empty or new line detected
		if input == "" || len(monkeys) == 0 {
			monkeys = append(monkeys, &Monkey{})
			continue
		}

		monkey := monkeys[len(monkeys)-1]

		input = strings.TrimSpace(input)
		switch {
		case strings.HasPrefix(input, "Monkey"):
			// set the monkey ID
			fmt.Sscanf(input, "Monkey %d:", &monkey.ID)

		case strings.HasPrefix(input, "Starting items:"):
			// set the list of monkey starting items
			for _, item := range strings.Split(strings.TrimPrefix(input, "Starting items: "), ", ") {
				monkey.items = append(monkey.items, utils.MustInt(item))
			}

		case strings.HasPrefix(input, "If true: throw to monkey"):
			// set the monkey target for a true test
			fmt.Sscanf(input, "If true: throw to monkey %d", &monkey.trueTestMonkey)

		case strings.HasPrefix(input, "If false: throw to monkey"):
			// set the monkey target for a false test
			fmt.Sscanf(input, "If false: throw to monkey %d", &monkey.falseTestMonkey)

		case strings.HasPrefix(input, "Test: divisible by"):
			var divisor worryLevel
			fmt.Sscanf(input, "Test: divisible by %d", &divisor)

			// set the test func to check if a worry level is divisible by the input
			monkey.test = func(item worryLevel) bool {
				return item%divisor == 0
			}

			// multiply every divisors by each other to get a common factor used to reduce worry level later
			reduceFactor *= divisor

		case strings.HasPrefix(input, "Operation:"):
			// set the operation func updating worry level during inspection
			args := strings.Fields(strings.TrimPrefix(input, "Operation: new = "))

			// get operation operands
			lhs, operator, rhs := utils.MustInt(args[0]), args[1], utils.MustInt(args[2])

			monkey.testWorryLevel = func(old worryLevel) worryLevel {
				// if left-hand operand is "old", use current value
				if args[0] == "old" {
					lhs = old
				}
				// if right-hand operand is "old", use current value
				if args[2] == "old" {
					rhs = old
				}

				// return operation result
				switch operator {
				case "+":
					return lhs + rhs
				case "-":
					return lhs - rhs
				case "*":
					return lhs * rhs
				case "/":
					return lhs / rhs
				default:
					return -1
				}
			}
		}
	}

	// set the reduce factor for all monkeys
	for _, monkey := range monkeys {
		monkey.reduceFactor = reduceFactor
	}

	return
}

func main() {
	fmt.Println("--- 2022 Day 11: Monkey in the Middle ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	////////////////////////////////////////

	monkeys := getMonkeys(inputs)

	// play 20 rounds, without being worried about items inspection
	for i := 0; i < 20; i++ {
		monkeys.playRound(false)
	}

	// 10605
	fmt.Println("Part 1:", monkeys.getBusiness())

	////////////////////////////////////////

	monkeys = getMonkeys(inputs)

	// play 10000 rounds, being worried about items inspection
	for i := 0; i < 10_000; i++ {
		monkeys.playRound(true)
	}

	// 2713310158
	fmt.Println("Part 2:", monkeys.getBusiness())
}
