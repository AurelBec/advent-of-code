// https://adventofcode.com/2023/day/20

package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
	"github.com/aurelbec/advent-of-code/utils/collections"
)

type Pulse int

const (
	low Pulse = iota
	high
)

type Propagation struct {
	from, to *Module
	pulse    Pulse
}

type Module struct {
	name         string
	kind         byte
	state        Pulse
	first        [2]int
	inputs       map[*Module]Pulse
	destinations []*Module
}

func (m *Module) String() string {
	return string(m.kind) + m.name
}

func (m *Module) generate(pulse Pulse) []Propagation {
	return utils.ArrayMap(m.destinations, func(next *Module) Propagation { return Propagation{from: m, to: next, pulse: pulse} })
}

func (m *Module) propagatePulse(n int, from *Module, pulse Pulse) []Propagation {
	// fmt.Printf("%s -%v-> %s\n", from, pulse, m)

	if m.first[pulse] == 0 {
		m.first[pulse] = n
	}

	switch m.kind {
	// flip-flop
	case '%':
		// do nothing on high
		if pulse == high {
			return nil
		}

		// else, flip and propagate
		m.state = 1 - m.state
		return m.generate(m.state)

	// conjunction
	case '&':
		// update memory
		m.inputs[from] = pulse

		// check if at least one low is present
		for _, pulse := range m.inputs {
			if pulse == low {
				return m.generate(high)
			}
		}

		// propagate low if all high
		return m.generate(low)

	// broadcaster
	default:
		return m.generate(pulse)
	}
}

func (m *Module) getConjunctionValue(p string) (int, bool) {
	if m == nil {
		return 0, false
	}

	if m.kind == '%' {
		return m.first[low], m.first[low] > 0
	}

	allFlipFlop := true
	allConjunctions := true
	values := make([]int, 0, len(m.inputs))
	for input := range m.inputs {
		value, valid := input.getConjunctionValue(p + "  ")
		if !valid {
			return 0, false
		}
		allFlipFlop = allFlipFlop && input.kind == '%'
		allConjunctions = allConjunctions && input.kind == '&'
		values = append(values, value)
	}

	if allFlipFlop {
		return utils.Sum(values), true
	} else if allConjunctions {
		return utils.LCM(values...), true
	} else {
		return 0, false
	}
}

type Button struct {
	pushes      int
	broadcaster *Module
}

func (b *Button) getSignalsCountAfter(pushes int) int {
	queue := collections.NewQueue[Propagation]()
	counts := [2]int{0, 0}
	for i := 0; i < pushes; i++ {
		b.pushes++
		queue.Enqueue(Propagation{to: b.broadcaster, pulse: low})
		for !queue.IsEmpty() {
			propagation, _ := queue.Dequeue()
			counts[propagation.pulse]++
			queue.Enqueue(propagation.to.propagatePulse(b.pushes, propagation.from, propagation.pulse)...)
		}
	}
	return counts[low] * counts[high]
}

func (b *Button) getCountUntilConjunctionOn(target *Module) int {
	for pow := int(math.Log2(float64(b.pushes))); pow < 14; pow++ {
		b.getSignalsCountAfter((2 << pow) - b.pushes)
		if count, valid := target.getConjunctionValue(""); valid {
			return count
		}
	}
	return 0
}

func parseModules(inputs []string) map[string]*Module {
	modules := make(map[string]*Module, len(inputs))
	outputs := make(map[string][]string, len(inputs))

	for _, input := range inputs {
		module := &Module{
			state:  low,
			inputs: make(map[*Module]Pulse),
		}

		left, right, _ := strings.Cut(input, " -> ")
		switch kind := left[0]; kind {
		case '%', '&':
			module.name = left[1:]
			module.kind = kind
		default:
			module.name = left
		}

		modules[module.name] = module
		outputs[module.name] = strings.Split(right, ", ")
	}

	for name, module := range modules {
		for _, destination := range outputs[name] {
			dest, found := modules[destination]
			if !found {
				dest = &Module{name: destination, inputs: make(map[*Module]Pulse)}
				modules[destination] = dest
			}

			module.destinations = append(module.destinations, dest)
			dest.inputs[module] = low
		}
	}

	return modules
}

func main() {
	fmt.Println("--- 2023 Day 20: Pulse Propagation ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	modules := parseModules(inputs)

	rx := modules["rx"]
	broadcaster := modules["broadcaster"]
	button := Button{broadcaster: broadcaster}

	////////////////////////////////////////

	// 32000000
	fmt.Println("Part 1:", button.getSignalsCountAfter(1000))

	////////////////////////////////////////

	// undefined
	fmt.Println("Part 2:", button.getCountUntilConjunctionOn(rx))
}
