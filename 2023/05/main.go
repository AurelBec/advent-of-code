// https://adventofcode.com/2023/day/5

package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

type Almanac struct {
	seeds    []int
	mappings []Mapping
}

type Mapping struct {
	name   string
	ranges []Range
}

type Range struct {
	source utils.Interval[int]
	offset int
}

func (almanac Almanac) getLowestLocation(sources utils.Intervals[int]) int {
	for _, mapping := range almanac.mappings {
		n := len(mapping.ranges)
		if n == 0 {
			continue
		}

		destinations := utils.Intervals[int]{}
		for _, source := range sources.Values() {
			// project values before first range
			destination := utils.NewInterval(source.Min, min(mapping.ranges[0].source.Min-1, source.Max))
			if destination.Len() > 0 {
				destinations.Insert(destination)
			}

			// project values through ranges
			for _, r := range mapping.ranges {
				destination = source.Intersection(r.source)
				if destination.Len() > 0 {
					destination.Shift(r.offset)
					destinations.Insert(destination)
				}
			}

			// project values after last range
			destination = utils.NewInterval(max(mapping.ranges[n-1].source.Max, source.Min), source.Max)
			if destination.Len() > 0 {
				destinations.Insert(destination)
			}
		}
		sources = destinations
	}

	res := math.MaxInt
	for _, sub := range sources.Values() {
		res = min(res, sub.Min)
	}
	return res
}

func parseAlmanac(inputs []string) Almanac {
	destination, source, length := 0, 0, 0
	almanac := Almanac{}

	for _, line := range inputs {
		switch {
		// get seeds list
		case len(almanac.seeds) == 0 && strings.HasPrefix(line, "seeds:"):
			almanac.seeds = utils.ArrayMap(utils.Numbers(line), utils.MustInt)

		// register new mapping
		case strings.HasSuffix(line, "map:"):
			almanac.mappings = append(almanac.mappings, Mapping{})
			fmt.Sscanf(line, "%s map:", &almanac.mappings[len(almanac.mappings)-1].name)

		// end map, do nothing
		case line == "":

		// add range to the last mapping
		default:
			fmt.Sscanf(line, "%d %d %d", &destination, &source, &length)
			almanac.mappings[len(almanac.mappings)-1].ranges = append(
				almanac.mappings[len(almanac.mappings)-1].ranges,
				Range{source: utils.NewInterval(source, source+length-1), offset: destination - source},
			)
		}
	}

	// order mapping ranges by their source
	for m := range almanac.mappings {
		sort.Slice(almanac.mappings[m].ranges, func(i, j int) bool {
			return almanac.mappings[m].ranges[i].source.Min < almanac.mappings[m].ranges[j].source.Min
		})
	}

	return almanac
}

func main() {
	fmt.Println("--- 2023 Day 5: If You Give A Seed A Fertilizer ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	almanac := parseAlmanac(inputs)

	////////////////////////////////////////

	intervals := utils.Intervals[int]{}
	for i := 0; i < len(almanac.seeds); i += 1 {
		intervals.Insert(utils.NewInterval(almanac.seeds[i], almanac.seeds[i]))
	}

	// 35
	fmt.Println("Part 1:", almanac.getLowestLocation(intervals))

	////////////////////////////////////////

	intervals = utils.Intervals[int]{}
	for i := 0; i < len(almanac.seeds); i += 2 {
		intervals.Insert(utils.NewInterval(almanac.seeds[i], almanac.seeds[i]+almanac.seeds[i+1]-1))
	}

	// 46
	fmt.Println("Part 2:", almanac.getLowestLocation(intervals))
}
