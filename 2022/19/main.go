// https://adventofcode.com/2022/day/19

package main

import (
	"fmt"
	"math"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

const (
	ore = iota
	clay
	obsidian
	geode
	N
)

type State struct {
	time      int
	robots    [N]int
	resources [N]int
}

func (s State) hash() string {
	return fmt.Sprint(s.time, s.resources[ore], s.resources[clay], s.resources[obsidian], s.resources[geode], s.robots[ore], s.robots[clay], s.robots[obsidian], s.robots[geode])
}

func (state *State) forward(time int) {
	state.time += time
	for resource, number := range state.robots {
		state.resources[resource] += number * time
	}
}

func (state State) getMaxUntilEnd(resource int, end int) int {
	dt := end - state.time
	return state.resources[resource] + state.robots[resource]*dt + (dt*(dt-1))/2
}

type Blueprint struct {
	ID         int
	robotCosts [N][N]int
	maxRobots  [N]int
}

func (blueprint Blueprint) getMax(resource int, timeLimit int) int {
	max := 0
	cache := make(map[string]struct{})

	var explore func(State)
	explore = func(state State) {
		if state.time == timeLimit {
			max = utils.Max(max, state.resources[resource])
			return
		}

		for _, next := range blueprint.nextStates(state, timeLimit) {
			hash := next.hash()
			if _, visited := cache[hash]; visited {
				continue
			}
			cache[hash] = struct{}{}

			if max > 0 && next.getMaxUntilEnd(resource, timeLimit) <= max {
				continue
			}

			explore(next)
		}
	}

	start := State{}
	start.robots[ore] = 1
	explore(start)
	return max
}

func (blueprint Blueprint) getQualityLevel(resource int, timeLimit int) int {
	return blueprint.ID * blueprint.getMax(resource, timeLimit)
}

func (blueprint Blueprint) nextStates(current State, end int) (nextStates []State) {
	if current.time >= end {
		return
	}

	// iterate for each kind of resources if we can buy a robot
robots:
	for _, robot := range []int{geode, obsidian, clay, ore} {
		// no need to create a new robot for this resource
		if robot != geode && current.robots[robot] >= blueprint.maxRobots[robot] {
			continue
		}

		nextState := current

		// try to find time until robot is created
		timeUntil := 0
		for resource, cost := range blueprint.robotCosts[robot] {
			// no robot producing resource, can not buy
			if current.robots[resource] == 0 && cost > 0 {
				continue robots
			}

			// if there is no time remaining, ignore
			timeNeeded := int(math.Ceil(float64(cost-current.resources[resource]) / float64(current.robots[resource])))
			if current.time+timeNeeded >= end-1 {
				continue robots
			}

			// take time and cost into account
			timeUntil = utils.Max(timeUntil, timeNeeded)
			nextState.resources[resource] -= cost
		}

		// forward time to when robot is ready
		nextState.forward(timeUntil + 1)
		nextState.robots[robot]++
		nextStates = append(nextStates, nextState)
	}

	// forward to the end if no robots built
	if len(nextStates) == 0 {
		current.forward(end - current.time)
		nextStates = append(nextStates, current)
	}

	return
}

type Blueprints []Blueprint

func (blueprints Blueprints) getQualityLevels(timeLimit int) int {
	sum := 0
	for _, blueprint := range blueprints {
		sum += blueprint.getQualityLevel(geode, timeLimit)
	}
	return sum
}

func (blueprints Blueprints) getMaxMultiplied(n int, timeLimit int) int {
	mul := 1
	for i := 0; i < utils.Min(n, len(blueprints)); i++ {
		mul *= blueprints[i].getMax(geode, timeLimit)
	}
	return mul
}

func parseBlueprints(inputs []string) Blueprints {
	blueprints := make([]Blueprint, len(inputs))

	for i, input := range inputs {
		fmt.Sscanf(
			input,
			"Blueprint %v: Each ore robot costs %v ore. Each clay robot costs %v ore. Each obsidian robot costs %v ore and %v clay. Each geode robot costs %v ore and %v obsidian.",
			&blueprints[i].ID,
			&blueprints[i].robotCosts[ore][ore],
			&blueprints[i].robotCosts[clay][ore],
			&blueprints[i].robotCosts[obsidian][ore],
			&blueprints[i].robotCosts[obsidian][clay],
			&blueprints[i].robotCosts[geode][ore],
			&blueprints[i].robotCosts[geode][obsidian],
		)

		for cost := 0; cost < N; cost++ {
			for resource := 0; resource < N; resource++ {
				blueprints[i].maxRobots[cost] = utils.Max(blueprints[i].maxRobots[cost], blueprints[i].robotCosts[resource][cost])
			}
		}
	}

	return blueprints
}

func main() {
	fmt.Println("--- 2022 Day 19: Not Enough Minerals ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	blueprints := parseBlueprints(inputs)

	////////////////////////////////////////

	// 33
	fmt.Println("Part 1:", blueprints.getQualityLevels(24))

	////////////////////////////////////////

	// 3472
	fmt.Println("Part 2:", blueprints.getMaxMultiplied(3, 32))
}
