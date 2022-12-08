// https://adventofcode.com/2022/day/9

package main

import (
	"fmt"
	"time"
)

var inputs1 = [...]string{
	"R 4",
	"U 4",
	"L 3",
	"D 1",
	"R 4",
	"D 1",
	"L 5",
	"R 2",
}

var inputs2 = [...]string{
	"R 5",
	"U 8",
	"L 8",
	"D 3",
	"R 17",
	"D 10",
	"L 25",
	"U 20",
}

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
	defer func(start time.Time) { fmt.Println("took:", time.Since(start).Round(time.Microsecond)) }(time.Now())

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
	fmt.Println("part1:", tail.visits())

	////////////////////////////////////////

	rope := make(rope, 10)
	for _, input := range inputs2 {
		direction, step := parseMotionMove(input)
		for i := 0; i < step; i++ {
			rope.move(direction)
		}
	}

	// 36
	fmt.Println("part2:", rope.tail().visits())
}
