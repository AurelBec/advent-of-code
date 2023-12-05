// https://adventofcode.com/2023/day/6

package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

type Race struct {
	time     int
	distance int
}

func (race Race) numberOfWaysToWin(...int) int {
	T := float64(race.time)
	D := float64(race.distance)
	// distance equation: D(held) = held * (T-held)
	// -> solve -h²+hT-D > 0
	// -> h = (T +/- sqrt(T²-4D)) / 2
	d := math.Sqrt(T*T - 4*D)
	return int(math.Ceil((T+d)/2)-math.Floor((T-d)/2)) - 1
}

func concatenateRaces(races []Race) Race {
	time, distance := "", ""
	for _, race := range races {
		time += fmt.Sprint(race.time)
		distance += fmt.Sprint(race.distance)
	}
	return Race{time: utils.MustInt(time), distance: utils.MustInt(distance)}
}

func parseRaces(inputs []string) []Race {
	times := utils.ArrayMap(utils.Numbers(strings.TrimPrefix(inputs[0], "Time:")), utils.MustInt)
	distances := utils.ArrayMap(utils.Numbers(strings.TrimPrefix(inputs[1], "Distance:")), utils.MustInt)
	n := min(len(times), len(distances))
	races := make([]Race, n)
	for i := 0; i < n; i++ {
		races[i].time = times[i]
		races[i].distance = distances[i]
	}
	return races
}

func main() {
	fmt.Println("--- 2023 Day 6: Wait For It ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	races := parseRaces(inputs)

	////////////////////////////////////////

	// 288
	fmt.Println("Part 1:", utils.MultiplyFunc(races, Race.numberOfWaysToWin))

	////////////////////////////////////////

	// 71503
	fmt.Println("Part 2:", concatenateRaces(races).numberOfWaysToWin(0))
}
