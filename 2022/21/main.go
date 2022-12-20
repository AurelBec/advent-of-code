// https://adventofcode.com/2022/day/21

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

const (
	add = '+'
	sub = '-'
	mul = '*'
	div = '/'
)

var inverse = map[byte]byte{
	add: sub,
	sub: add,
	mul: div,
	div: mul,
}

func calculate(lhs int, operator byte, rhs int) int {
	switch operator {
	case add:
		return lhs + rhs
	case sub:
		return lhs - rhs
	case mul:
		return lhs * rhs
	case div:
		return lhs / rhs
	default:
		panic(fmt.Sprintf("unknown operator %v", operator))
	}
}

type Monkey struct {
	name        string
	parent      *Monkey
	left, right *Monkey
	operation   byte
	value       int
}

func (m *Monkey) yell() int {
	if m.value == 0 {
		m.value = calculate(m.left.yell(), m.operation, m.right.yell())
	}
	return m.value
}

func (m *Monkey) waitOther() int {
	if m.parent.left == m {
		return m.parent.right.yell()
	} else if m.parent.right == m {
		return m.parent.left.yell()
	}
	panic("unknown dependency")
}

func (m *Monkey) yellFor(target *Monkey) int {
	// when target is found, evaluate other branch to get equality
	if m.parent == target {
		return m.waitOther()
	}

	// else, get the parent evaluation, and the other branch value
	// inverse operand for sub and div operations on right side
	if m.parent.right == m && (m.parent.operation == sub || m.parent.operation == div) {
		return calculate(m.waitOther(), m.parent.operation, m.parent.yellFor(target))
	} else {
		return calculate(m.parent.yellFor(target), inverse[m.parent.operation], m.waitOther())
	}
}

func parseMonkeys(inputs []string) map[string]*Monkey {
	args := utils.ArrayMap(inputs, strings.Fields)

	monkeys := make(map[string]*Monkey, len(inputs))
	for _, monkey := range args {
		name := strings.TrimRight(monkey[0], ":")
		monkeys[name] = &Monkey{name: name}
	}

	for _, monkey := range args {
		name := strings.TrimRight(monkey[0], ":")
		switch args := monkey[1:]; {
		case len(args) == 1:
			monkeys[name].value = utils.MustInt(args[0])
		default:
			monkeys[name].left = monkeys[args[0]]
			monkeys[name].operation = args[1][0]
			monkeys[name].right = monkeys[args[2]]

			monkeys[name].left.parent = monkeys[name]
			monkeys[name].right.parent = monkeys[name]
		}
	}
	return monkeys
}

func main() {
	fmt.Println("--- 2022 Day 21: Monkey Math ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	monkeys := parseMonkeys(inputs)

	root := monkeys["root"]
	humn := monkeys["humn"]

	////////////////////////////////////////

	// 152
	fmt.Println("Part 1:", root.yell())

	////////////////////////////////////////

	// 301
	fmt.Println("Part 2:", humn.yellFor(root))
}
