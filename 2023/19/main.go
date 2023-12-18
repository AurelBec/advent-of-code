// https://adventofcode.com/2023/day/19

package main

import (
	"fmt"
	"maps"
	"strings"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

const (
	X = 'x'
	M = 'm'
	A = 'a'
	S = 's'
)

var getter = map[byte]func(Item) int{
	X: func(i Item) int { return i.x },
	M: func(i Item) int { return i.m },
	A: func(i Item) int { return i.a },
	S: func(i Item) int { return i.s },
}

var comparator = map[byte]func(int, int) bool{
	'<': func(a, b int) bool { return a < b },
	'>': func(a, b int) bool { return a > b },
}

type Item struct {
	x, m, a, s int
}

func (i Item) rating() int {
	return i.x + i.m + i.a + i.s
}

type ItemRanges map[byte]utils.Interval[int]

func NewItemRanges(min, max int) ItemRanges {
	return ItemRanges{
		X: utils.NewInterval(min, max),
		M: utils.NewInterval(min, max),
		A: utils.NewInterval(min, max),
		S: utils.NewInterval(min, max),
	}
}

func (ir ItemRanges) count() int {
	return ir[X].Len() * ir[M].Len() * ir[A].Len() * ir[S].Len()
}

type Workflows map[string]Workflow

type Workflow struct {
	rules []Rule
}

type Rule struct {
	field    byte
	operator byte
	value    int
	next     string
}

func (r Rule) eval(item Item) bool {
	return r.operator == 0 || comparator[r.operator](getter[r.field](item), r.value)
}

func (ws Workflows) getItemsRating(start string, items []Item) []int {
	return utils.ArrayMap(items, func(item Item) int { return ws.getItemRating(start, item) })
}

func (ws Workflows) getItemRating(start string, item Item) int {
	for _, rule := range ws[start].rules {
		if !rule.eval(item) {
			continue
		}
		switch rule.next {
		case "A":
			return item.rating()
		case "R":
			return 0
		default:
			return ws.getItemRating(rule.next, item)
		}
	}
	return 0
}

func (ws Workflows) getAcceptedItemsCount(start string, current ItemRanges) int {
	if start == "A" {
		return current.count()
	} else if start == "R" {
		return 0
	}

	count := 0
	for _, rule := range ws[start].rules {
		next := maps.Clone(current)
		if rule.operator == '>' {
			current[rule.field], next[rule.field] = current[rule.field].Split(rule.value, true)
		} else if rule.operator == '<' {
			next[rule.field], current[rule.field] = current[rule.field].Split(rule.value, false)
		}
		count += ws.getAcceptedItemsCount(rule.next, next)
	}
	return count
}

func parseSystem(inputs []string) (Workflows, []Item) {
	workflows := make(Workflows, len(inputs))
	items := make([]Item, 0, len(inputs))

	// parse workflows
	i := 0
	for ; i < len(inputs); i++ {
		if inputs[i] == "" {
			break
		}

		name, input, _ := strings.Cut(inputs[i], "{")
		rules := strings.Split(input, ",")
		rules[len(rules)-1] = strings.TrimSuffix(rules[len(rules)-1], "}")

		var workflow Workflow
		for _, rule := range rules {
			operation, result, found := strings.Cut(rule, ":")
			if !found {
				workflow.rules = append(workflow.rules, Rule{
					next: operation,
				})
			} else {
				workflow.rules = append(workflow.rules, Rule{
					field:    operation[0],
					operator: operation[1],
					value:    utils.MustInt(operation[2:]),
					next:     result,
				})
			}
		}

		workflows[name] = workflow
	}

	// parse items
	i = i + 1
	for ; i < len(inputs); i++ {
		var item Item
		fmt.Sscanf(inputs[i], "{x=%v,m=%v,a=%v,s=%v}", &item.x, &item.m, &item.a, &item.s)
		items = append(items, item)
	}

	return workflows, items
}

func main() {
	fmt.Println("--- 2023 Day 19: Aplenty ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	workflows, items := parseSystem(inputs)

	////////////////////////////////////////

	// 19114
	fmt.Println("Part 1:", utils.Sum(workflows.getItemsRating("in", items)))

	////////////////////////////////////////

	// 167409079868000
	fmt.Println("Part 2:", workflows.getAcceptedItemsCount("in", NewItemRanges(1, 4000)))
}
