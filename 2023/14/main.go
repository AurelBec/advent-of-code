// https://adventofcode.com/2023/day/14

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

const (
	roundedRock = 'O'
	cubeRock    = '#'
)

type Platform struct {
	N     int
	rocks [][]byte
}

func (p *Platform) rollCycles(n int) {
	cache := make(map[string]int)
	for ; n > 0; n-- {
		for i := 0; i < 4; i++ {
			p.rollNorth()
			p.rotate()
		}

		state := p.String()
		if last, found := cache[state]; !found {
			cache[state] = n
		} else {
			n %= (last - n)
		}
	}
}

func (p *Platform) rollNorth() {
	for x := 0; x < p.N; x++ {
		freeY := p.N
		for y := 0; y < p.N; y++ {
			switch p.rocks[y][x] {
			case roundedRock:
				if freeY != y && freeY < p.N {
					p.rocks[y][x], p.rocks[freeY][x] = p.rocks[freeY][x], p.rocks[y][x]
					freeY++
				}
			case cubeRock:
				freeY = p.N
			default:
				freeY = min(freeY, y)
			}
		}
	}
}

func (p *Platform) rotate() {
	for x := 0; x < p.N/2; x++ {
		for y := x; y < p.N-x-1; y++ {
			temp := p.rocks[x][y]
			p.rocks[x][y] = p.rocks[p.N-1-y][x]
			p.rocks[p.N-1-y][x] = p.rocks[p.N-1-x][p.N-1-y]
			p.rocks[p.N-1-x][p.N-1-y] = p.rocks[y][p.N-1-x]
			p.rocks[y][p.N-1-x] = temp
		}
	}
}

func (p Platform) getNorthLoad() int {
	totalLoad := 0
	for y, row := range p.rocks {
		for _, rock := range row {
			if rock == roundedRock {
				totalLoad += p.N - y
			}
		}
	}
	return totalLoad
}

func (p Platform) String() string {
	return strings.Join(utils.ArrayMap(p.rocks, func(b []byte) string { return string(b) }), "\n")
}

func parsePlatform(inputs []string) Platform {
	platform := Platform{
		N:     len(inputs),
		rocks: make([][]byte, len(inputs)),
	}

	for y, input := range inputs {
		platform.rocks[y] = []byte(input)
	}

	return platform
}

func main() {
	fmt.Println("--- 2023 Day 14: Parabolic Reflector Dish ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	platform := parsePlatform(inputs)

	////////////////////////////////////////

	platform.rollNorth()

	// 136
	fmt.Println("Part 1:", platform.getNorthLoad())

	////////////////////////////////////////

	platform.rollCycles(1_000_000_000)

	// 64
	fmt.Println("Part 2:", platform.getNorthLoad())
}
