// https://adventofcode.com/2022/day/10

package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

type CPU struct {
	x      int
	memory []int
}

// execute parses and executes an instruction
// it also keeps trace in memory of all x register values among cycles
func (cpu *CPU) execute(instruction string) {
	switch args := strings.Fields(instruction); args[0] {
	// noop takes one cycle to complete
	// it has no other effect
	case "noop":
		cpu.memory = append(cpu.memory, cpu.x)

	// addx V takes two cycles to complete
	// after two cycles, the x register is increased by the value V (V can be negative)
	case "addx":
		v, _ := strconv.Atoi(args[1])
		cpu.memory = append(cpu.memory, cpu.x, cpu.x)
		cpu.x += v
	}
}

// value returns the x register value at the n-th cycle
// note: cycle is expected to be >= 1
// note: x register start value is expected to be 1, so a unit is added to all values returned
func (cpu CPU) value(cycle int) int {
	if cycle <= 1 {
		return 1
	} else if cycle > len(cpu.memory) {
		return cpu.x + 1
	} else {
		return cpu.memory[cycle-1] + 1
	}
}

// strength returns the signal strength at the n-th cycle
// note: cycle is expected to be >= 1
func (cpu CPU) strength(cycle int, _ ...int) int {
	return cycle * cpu.value(cycle)
}

type CRT struct {
	w, h   int
	cpu    CPU
	buffer strings.Builder
}

// execute executes a single instruction
func (crt *CRT) execute(instruction string) {
	// execute CPU instruction
	crt.cpu.execute(instruction)

	// perform drawing operation
	for cycle := crt.buffer.Len() + 1; cycle <= len(crt.cpu.memory); cycle++ {
		// pixel currently drawn
		pixel := (cycle - 1) % crt.w
		// get sprite location
		sprite := crt.cpu.value(cycle)

		// if the pixel is in the 3 pixel wide window of the sprite, draw it
		if math.Abs(float64(pixel-sprite)) <= 1 {
			crt.buffer.WriteRune('#')
		} else {
			crt.buffer.WriteRune(' ')
		}
	}
}

func (crt CRT) String() string {
	buffer := crt.buffer.String()
	screen := strings.Builder{}
	screen.Grow((crt.w + 1) * crt.h)
	for y := 0; y < crt.h; y++ {
		if len(buffer) >= crt.w {
			fmt.Fprintf(&screen, "\n%s", buffer[:crt.w])
			buffer = buffer[crt.w:]
		} else {
			fmt.Fprintf(&screen, "\n%-[2]*[1]s", buffer, crt.w)
		}
	}
	return screen.String()
}

func main() {
	fmt.Println("--- 2022 Day 10: Cathode-Ray Tube ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	////////////////////////////////////////

	cpu := CPU{}
	for _, instruction := range inputs {
		cpu.execute(instruction)
	}

	cycles := []int{20, 60, 100, 140, 180, 220}

	// 13140
	fmt.Println("Part 1:", utils.SumFunc(cycles, cpu.strength))

	////////////////////////////////////////

	crt := CRT{w: 40, h: 6}
	for _, instruction := range inputs {
		crt.execute(instruction)
	}

	// ##..##..##..##..##..##..##..##..##..##..
	// ###...###...###...###...###...###...###.
	// ####....####....####....####....####....
	// #####.....#####.....#####.....#####.....
	// ######......######......######......####
	// #######.......#######.......#######.....
	fmt.Println("Part 2:", crt)
}
