// https://adventofcode.com/2023/day/22

package main

import (
	"fmt"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
	"github.com/aurelbec/advent-of-code/utils/collections"
)

type Brick struct {
	id            int
	finalLocation bool
	start, end    utils.Location3D[int]

	isUnder   []*Brick
	isOver    []*Brick
	inCascade int
}

func (b *Brick) isSafeToDisintegrated() bool {
	return b.getDependentBricks() == 0
}

func (b *Brick) getDependentBricks() int {
	if b.inCascade >= 0 {
		return b.inCascade
	}

	fallen := make(map[*Brick]bool)
	queue := collections.NewQueue(b)
	for !queue.IsEmpty() {
		current, _ := queue.Dequeue()
		fallen[current] = true

	over:
		for _, over := range current.isUnder {
			for _, underOver := range over.isOver {
				if !fallen[underOver] {
					continue over
				}
			}
			queue.Enqueue(over)
		}
	}

	b.inCascade = len(fallen) - 1
	return b.inCascade
}

func parseBricks(inputs []string) []*Brick {
	X, Y, Z := 0, 0, 0
	bricks := make([]*Brick, len(inputs))
	for i, input := range inputs {
		brick := Brick{id: i + 1, inCascade: -1}
		fmt.Sscanf(input, "%v,%v,%v~%v,%v,%v", &brick.start.X, &brick.start.Y, &brick.start.Z, &brick.end.X, &brick.end.Y, &brick.end.Z)
		bricks[i] = &brick
		X = max(X, brick.start.X, brick.end.X)
		Y = max(Y, brick.start.Y, brick.end.Y)
		Z = max(Z, brick.start.Z, brick.end.Z)
	}

	// init space
	X, Y, Z = X+1, Y+1, Z+1
	cubes := make([][][]*Brick, Z)
	for z := range cubes {
		cubes[z] = make([][]*Brick, Y)
		for y := range cubes[z] {
			cubes[z][y] = make([]*Brick, X)
		}
	}

	for _, brick := range bricks {
		for x := brick.start.X; x <= brick.end.X; x++ {
			for y := brick.start.Y; y <= brick.end.Y; y++ {
				for z := brick.start.Z; z <= brick.end.Z; z++ {
					cubes[z][y][x] = brick
				}
			}
		}
	}

	// initiate fall
	for z := 1; z < Z; z++ {
		for y := 0; y < Y; y++ {
			for x := 0; x < X; x++ {
				brick := cubes[z][y][x]
				if brick == nil {
					continue
				} else if brick.finalLocation {
					continue
				} else if brick.start.Z < z {
					continue
				}

				blockers := make(map[*Brick]bool)
				for dZ := 1; dZ <= z && !brick.finalLocation; dZ++ {
					for x := brick.start.X; x <= brick.end.X; x++ {
						for y := brick.start.Y; y <= brick.end.Y; y++ {
							blocker := cubes[z-dZ][y][x]
							if dZ != z {
								if blocker == nil {
									continue
								} else if blockers[blocker] {
									continue
								}
								blockers[blocker] = true
							}

							if brick.finalLocation || dZ == 1 {
								brick.finalLocation = true
								continue
							}

							for x := brick.start.X; x <= brick.end.X; x++ {
								for y := brick.start.Y; y <= brick.end.Y; y++ {
									for z := brick.start.Z; z <= brick.end.Z; z++ {
										cubes[z-dZ+1][y][x] = brick
										cubes[z][y][x] = nil
									}
								}
							}
							brick.start.Z = brick.start.Z - dZ + 1
							brick.end.Z = brick.end.Z - dZ + 1
							brick.finalLocation = true
						}
					}
				}
				for blocker := range blockers {
					brick.isOver = append(brick.isOver, blocker)
					blocker.isUnder = append(blocker.isUnder, brick)
				}
				brick.finalLocation = true
			}
		}
	}

	return bricks
}

func main() {
	fmt.Println("--- 2023 Day 22: Sand Slabs ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	bricks := parseBricks(inputs)

	////////////////////////////////////////

	// 5
	fmt.Println("Part 1:", utils.SumFunc(bricks, func(b *Brick, _ ...int) int {
		if b.isSafeToDisintegrated() {
			return 1
		}
		return 0
	}))

	////////////////////////////////////////

	// 7
	fmt.Println("Part 2:", utils.SumFunc(bricks, func(b *Brick, _ ...int) int {
		return b.getDependentBricks()
	}))
}
