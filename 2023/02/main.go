// https://adventofcode.com/2023/day/2

package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

type game struct {
	id               int
	red, green, blue int
}

func (g game) isValid(red, green, blue int) bool {
	return g.red <= red && g.green <= green && g.blue <= blue
}

func (g game) power() int {
	return g.red * g.green * g.blue
}

func getGameIDs(games []game, red, green, blue int) int {
	gameIDs := 0
	for _, game := range games {
		if game.isValid(red, green, blue) {
			gameIDs += game.id
		}
	}
	return gameIDs
}

func getGamePowers(games []game) int {
	gamePowers := 0
	for _, game := range games {
		gamePowers += game.power()
	}
	return gamePowers
}

func parseGames(inputs []string) []game {
	games := make([]game, len(inputs))
	for i, input := range inputs {
		information := strings.Split(input, ": ")
		fmt.Sscanf(information[0], "Game %v", &games[i].id)
		for _, subset := range strings.Split(information[1], "; ") {
			for _, cubes := range strings.Split(subset, ", ") {
				colors := strings.Split(cubes, " ")
				n, _ := strconv.Atoi(colors[0])
				switch colors[1] {
				case "red":
					games[i].red = max(n, games[i].red)
				case "green":
					games[i].green = max(n, games[i].green)
				case "blue":
					games[i].blue = max(n, games[i].blue)
				}
			}
		}
	}
	return games
}

func main() {
	fmt.Println("--- 2023 Day 2: Cube Conundrum ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	games := parseGames(inputs)

	////////////////////////////////////////

	// 8
	fmt.Println("Part 1:", getGameIDs(games, 12, 13, 14))

	////////////////////////////////////////

	// 2286
	fmt.Println("Part 2:", getGamePowers(games))
}
