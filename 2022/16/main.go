// https://adventofcode.com/2022/day/16

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

type Valve struct {
	id   string
	mask int

	flow          int
	shortestPaths map[*Valve]int
}

type Path struct {
	pressure int
	mask     int
}

// getPossiblePaths returns the list of all feasible paths in the given time
func getPossiblePaths(time int, start *Valve) []Path {
	return getPathsRec(0, time, 0, 0, start)
}

// getPathsRec returns the list of all feasible paths in the given time recursively
// it uses a cache of opened valves and current path valve ids
func getPathsRec(pressure int, remaining int, opened int, path int, node *Valve) []Path {
	paths := []Path{{pressure: pressure, mask: path}}
	for next, cost := range node.shortestPaths {
		if opened&next.mask != 0 || next.flow == 0 { // opened, skip
			continue
		}

		remaining := remaining - cost - 1 // go to valve and open it

		if remaining <= 0 { // no time, skip
			continue
		}

		paths = append(paths, getPathsRec(pressure+(remaining*next.flow), remaining, opened|next.mask, path|next.mask, next)...)
	}
	return paths
}

// parseValves parses input and return the map of valves
func parseValves(inputs []string) (valves map[string]*Valve) {
	valves = make(map[string]*Valve, len(inputs))

	neighborsID := make(map[*Valve][]string)
	for i, input := range inputs {
		valve := &Valve{mask: 1 << i}

		parts := strings.SplitAfter(input, "; ")
		fmt.Sscanf(parts[0], "Valve %s has flow rate=%d", &valve.id, &valve.flow)

		neighborsID[valve] = strings.Split(strings.Trim(parts[1], "tunelsadov "), ", ")
		valves[valve.id] = valve
	}

	neighbors := utils.MapMap(neighborsID, func(v *Valve, neighborsID []string) (*Valve, []*Valve) {
		return v, utils.ArrayMap(neighborsID, func(id string) *Valve { return valves[id] })
	})

	// for each valve, find the shortest path to other valves
	for id, start := range valves {
		costs := map[*Valve]int{start: 0}
		openList := utils.NewOrderedArray(start)
		for len(openList) > 0 {
			current := openList.Remove(0)
			costSoFar := costs[current]
			for _, neighbor := range neighbors[current] {
				if _, visited := costs[neighbor]; visited {
					continue
				}
				costs[neighbor] = costSoFar + 1
				openList.Insert(neighbor)
			}
		}
		valves[id].shortestPaths = costs
	}

	return
}

func main() {
	fmt.Println("--- 2022 Day 16: Proboscidea Volcanium ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	valves := parseValves(inputs)

	////////////////////////////////////////

	pressureMax, paths := 0, getPossiblePaths(30, valves["AA"])
	for _, path := range paths {
		pressureMax = utils.Max(pressureMax, path.pressure)
	}

	// 1651
	fmt.Println("Part 1:", pressureMax)

	////////////////////////////////////////

	pressureMax, paths = 0, getPossiblePaths(26, valves["AA"])
	for i, me := range paths {
		for _, elephant := range paths[i:] {
			if me.mask&elephant.mask == 0 { // ensure no common part
				pressureMax = utils.Max(pressureMax, me.pressure+elephant.pressure)
			}
		}
	}

	// 1707
	fmt.Println("Part 2:", pressureMax)
}
