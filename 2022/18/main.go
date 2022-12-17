// https://adventofcode.com/2022/day/18

package main

import (
	"fmt"
	"math"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

const (
	X = iota
	Y
	Z
)

type Boulder struct {
	ranges [3]utils.Interval[int]  // coordinates ranges [0; max+2]
	cubes  []utils.Location3D[int] // solid cube list
	shape  [][][]bool              // boulder shape
}

// isSolid returns whether the given cube is solid or not
func (boulder Boulder) isSolid(cube utils.Location3D[int]) bool {
	return boulder.shape[cube.Z][cube.Y][cube.X]
}

// outOfRange returns whether the given point is in coordinates ranges or not
func (boulder Boulder) outOfRange(cube utils.Location3D[int]) bool {
	return !boulder.ranges[X].Contains(cube.X) || !boulder.ranges[Y].Contains(cube.Y) || !boulder.ranges[Z].Contains(cube.Z)
}

// cubeNeighbors returns the list of adjacent cube neighbors
func (boulder Boulder) cubeNeighbors(cube utils.Location3D[int]) [6]utils.Location3D[int] {
	return [6]utils.Location3D[int]{
		{X: cube.X - 1, Y: cube.Y, Z: cube.Z},
		{X: cube.X + 1, Y: cube.Y, Z: cube.Z},
		{X: cube.X, Y: cube.Y - 1, Z: cube.Z},
		{X: cube.X, Y: cube.Y + 1, Z: cube.Z},
		{X: cube.X, Y: cube.Y, Z: cube.Z - 1},
		{X: cube.X, Y: cube.Y, Z: cube.Z + 1},
	}
}

// getSurfaceArea returns the area that is touching air
func (boulder Boulder) getSurfaceArea() (surfaceArea int) {
	// for each solid cubes
	for _, cube := range boulder.cubes {
		// get neighbors
		for _, neighbor := range boulder.cubeNeighbors(cube) {
			// if the neighbor is out of range neither solid, then this current face is touching air
			if boulder.outOfRange(neighbor) || !boulder.isSolid(neighbor) {
				surfaceArea++
			}
		}
	}

	return
}

// getVisibleSurfaceArea returns the area touching air only on the outside
func (boulder Boulder) getVisibleSurfaceArea() (surfaceArea int) {
	visited := make(map[utils.Location3D[int]]bool)

	// start a flood-fill algorithm at location (0,0,0), clear by design
	flood := utils.NewUnorderedArray(utils.Location3D[int]{X: 0, Y: 0, Z: 0})
	for len(flood) > 0 {
		water := flood.Remove(0)

		// do not visit twice
		if visited[water] {
			continue
		}
		visited[water] = true

		// for each neighbor
		for _, neighbor := range boulder.cubeNeighbors(water) {
			// ensure in range
			if boulder.outOfRange(neighbor) {
				continue
			}

			// if it's a solid block, then this current face is on the external surface
			if boulder.isSolid(neighbor) {
				surfaceArea++
			} else {
				// else continue the flood by exploring the neighbor
				flood.Insert(neighbor)
			}
		}
	}

	return
}

// getBoulder parses input and create the corresponding boulder
func getBoulder(input []string) *Boulder {
	boulder := Boulder{}

	// initialize ranges
	for i := range boulder.ranges {
		boulder.ranges[i] = utils.NewInterval(0, -math.MaxInt)
	}

	// parse input
	boulder.cubes = make([]utils.Location3D[int], len(input))
	for i, input := range input {
		fmt.Sscanf(input, "%v,%v,%v", &boulder.cubes[i].X, &boulder.cubes[i].Y, &boulder.cubes[i].Z)
		boulder.ranges[X].Merge(boulder.cubes[i].X)
		boulder.ranges[Y].Merge(boulder.cubes[i].Y)
		boulder.ranges[Z].Merge(boulder.cubes[i].Z)
	}

	// add one extra layer on each face, so increase ranges by 2
	for i := range boulder.ranges {
		boulder.ranges[i].Max += 2
	}

	// initialize empty shape
	boulder.shape = make([][][]bool, boulder.ranges[Z].Len())
	for z := range boulder.shape {
		boulder.shape[z] = make([][]bool, boulder.ranges[Y].Len())
		for y := range boulder.shape[z] {
			boulder.shape[z][y] = make([]bool, boulder.ranges[X].Len())
		}
	}

	// add solid cubes
	for i, cube := range boulder.cubes {
		// move cube by one to add free extra layer behind
		cube.X += 1
		cube.Y += 1
		cube.Z += 1
		boulder.shape[cube.Z][cube.Y][cube.X] = true
		boulder.cubes[i] = cube
	}

	return &boulder
}

func main() {
	fmt.Println("--- 2022 Day 18: Boiling Boulders ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	boulder := getBoulder(inputs)

	////////////////////////////////////////

	// 64
	fmt.Println("Part 1:", boulder.getSurfaceArea())

	////////////////////////////////////////

	// 58
	fmt.Println("Part 2:", boulder.getVisibleSurfaceArea())
}
