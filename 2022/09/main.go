// https://adventofcode.com/2022/day/9

package main

import (
	"fmt"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

// displacement needed to move closer to target the center
var displacement = map[int]map[int]struct{ x, y int }{
	+2: {-2: {+1, -1}, -1: {+1, -1}, 00: {00, -1}, +1: {-1, -1}, +2: {-1, -1}},
	+1: {-2: {+1, -1}, -1: {00, 00}, 00: {00, 00}, +1: {00, 00}, +2: {-1, -1}},
	00: {-2: {+1, 00}, -1: {00, 00}, 00: {00, 00}, +1: {00, 00}, +2: {-1, 00}},
	-1: {-2: {+1, +1}, -1: {00, 00}, 00: {00, 00}, +1: {00, 00}, +2: {-1, +1}},
	-2: {-2: {+1, +1}, -1: {+1, +1}, 00: {00, +1}, +1: {-1, +1}, +2: {-1, +1}},
}

type knot struct {
	x, y    int
	visited map[int]map[int]struct{}
}

// visits returns the number of different position visited
func (k knot) visits() (visits int) {
	for _, x := range k.visited {
		visits += len(x)
	}
	return
}

// visit saves the current position into a map
func (k *knot) visit() {
	switch {
	case k.visited == nil:
		k.visited = map[int]map[int]struct{}{k.x: {k.y: {}}}
	case k.visited[k.x] == nil:
		k.visited[k.x] = map[int]struct{}{k.y: {}}
	default:
		k.visited[k.x][k.y] = struct{}{}
	}
}

// move moves a knot into the given direction
func (k *knot) move(direction string) {
	k.visit()
	defer k.visit()

	switch direction {
	case "U":
		k.y++
	case "D":
		k.y--
	case "L":
		k.x--
	case "R":
		k.x++
	}
}

// moveCloserTo moves a knot closer to the target, following displacement rules
func (k *knot) moveCloserTo(target knot) {
	k.visit()
	defer k.visit()

	displacement := displacement[k.y-target.y][k.x-target.x]
	k.x += displacement.x
	k.y += displacement.y
}

type rope []knot

// move moves the head of the rope into the given direction, and move also the following knots
func (r rope) move(direction string) {
	for i := range r {
		if i == 0 {
			r[i].move(direction)
		} else {
			r[i].moveCloserTo(r[i-1])
		}
	}
}

// tail returns the tail knot of the rope
func (r rope) tail() knot {
	return r[len(r)-1]
}

// parseMotionMove parses a motion move and return the direction with number of steps
func parseMotionMove(motionMove string) (direction string, step int) {
	fmt.Sscanf(motionMove, "%s %v", &direction, &step)
	return
}

func main() {
	fmt.Println("--- 2022 Day 9: Rope Bridge ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs1 := utils.MustReadInput("example1.txt")
	inputs2 := utils.MustReadInput("example2.txt")

	////////////////////////////////////////

	head := knot{x: 0, y: 0}
	tail := knot{x: 0, y: 0}
	for _, input := range inputs1 {
		direction, step := parseMotionMove(input)
		for i := 0; i < step; i++ {
			head.move(direction)
			tail.moveCloserTo(head)
		}
	}

	// 13
	fmt.Println("Part 1:", tail.visits())

	////////////////////////////////////////

	rope := make(rope, 10)
	for _, input := range inputs2 {
		direction, step := parseMotionMove(input)
		for i := 0; i < step; i++ {
			rope.move(direction)
		}
	}

	// 36
	fmt.Println("Part 2:", rope.tail().visits())
}
