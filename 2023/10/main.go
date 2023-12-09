// https://adventofcode.com/2023/day/10

package main

import (
	"fmt"
	"slices"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
	"github.com/aurelbec/advent-of-code/utils/collections"
)

// directions
const (
	N = 1 << iota
	E
	S
	W
)

// pipes
const (
	NE = N | E
	SE = E | S
	SW = S | W
	NW = W | N
	NS = N | S
	EW = E | W
)

var conversion = map[byte]int{'|': NS, '-': EW, '7': SW, 'J': NW, 'F': SE, 'L': NE, 'S': 15}

type Tile struct {
	kind int
	x, y int
	next []*Tile
}

type Network struct {
	dims utils.Location2D[int]

	tiles [][]*Tile
	start *Tile
}

func (n Network) getLongestLoop(start *Tile) []*Tile {
	type node struct {
		tile          *Tile
		parent        *node
		distFromStart int
	}

	var part1, part2 *node

	visitedTiles := make(map[*Tile]*node, n.dims.X*n.dims.Y)
	queue := collections.NewQueue(&node{start, nil, 0})
	for {
		current, found := queue.Dequeue()
		// if queue is empty, backtrace loop
		if !found || current.tile == nil {
			if part1 == nil || part2 == nil {
				return nil
			}

			loop := make([]*Tile, 0, part1.distFromStart+part2.distFromStart)
			// start backtracking one extremity
			for node := part1; node != nil; node = node.parent {
				loop = append(loop, node.tile)
			}
			slices.Reverse(loop)
			// append the second
			for node := part2.parent; node != nil; node = node.parent {
				loop = append(loop, node.tile)
			}
			return loop
		}

		// if node has not been visited, continue exploring
		previous, visited := visitedTiles[current.tile]
		if !visited {
			visitedTiles[current.tile] = current
			queue.Enqueue(utils.ArrayMap(current.tile.next, func(tile *Tile) *node { return &node{tile, current, current.distFromStart + 1} })...)
			continue
		}

		// else, if the visit is longer, keep track of the two paths
		if current.distFromStart <= previous.distFromStart {
			part1 = previous
			part2 = current
		}
	}
}

func (n Network) getTilesInLoop(loop []*Tile) []*Tile {
	isLoop := make(map[*Tile]bool, len(loop))
	for _, tile := range loop {
		isLoop[tile] = true
	}

	insideLoop := make([]*Tile, 0)
	for x := 0; x < n.dims.X; x++ {
		isInsideLoop := false
		lastBend := 0
		for y := 0; y < n.dims.Y; y++ {
			tile := n.tiles[x][y]

			// switch loop if pipe direction changed to up/down
			isInsideLoop = isInsideLoop != (isLoop[tile] &&
				((tile.kind == EW) || (lastBend == SW && tile.kind == NE) || (lastBend == SE && tile.kind == NW)))

			// keep trace of last corner in row
			if isLoop[tile] && (tile.kind != NS && tile.kind != EW) {
				lastBend = tile.kind
			}

			// tile inside loop found
			if isInsideLoop && !isLoop[tile] {
				insideLoop = append(insideLoop, tile)
			}
		}
	}

	return insideLoop
}

func (n Network) print() {
	beautify := map[int]string{NE: "└", SE: "┌", SW: "┐", NW: "┘", NS: "│", EW: "─"}

	longestLoop := n.getLongestLoop(n.start)
	isOnLoop := make(map[*Tile]bool, len(longestLoop))
	for _, tile := range longestLoop {
		isOnLoop[tile] = true
	}

	tilesInLoop := n.getTilesInLoop(longestLoop)
	isInLoop := make(map[*Tile]bool, len(tilesInLoop))
	for _, tile := range tilesInLoop {
		isInLoop[tile] = true
	}

	for y := 0; y < n.dims.Y; y++ {
		for x := 0; x < n.dims.X; x++ {
			tile := n.tiles[x][y]
			if tile == n.start {
				fmt.Print("S")
			} else if isOnLoop[tile] {
				fmt.Print(beautify[tile.kind])
			} else if isInLoop[tile] {
				fmt.Print("░")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func parseNetwork(inputs []string) Network {
	network := Network{dims: utils.NewLocation2D(len(inputs[0]), len(inputs))}

	network.tiles = make([][]*Tile, network.dims.X)
	for x := 0; x < network.dims.X; x++ {
		network.tiles[x] = make([]*Tile, network.dims.Y)
		for y := 0; y < network.dims.Y; y++ {
			tile := inputs[y][x]
			network.tiles[x][y] = &Tile{conversion[tile], x, y, nil}
			if tile == 'S' {
				network.start = network.tiles[x][y]
			}
		}
	}

	for x, line := range network.tiles {
		for y, tile := range line {
			if x > 0 {
				if left := network.tiles[x-1][y]; (tile.kind&W) > 0 && (left.kind&E) > 0 {
					tile.next = append(tile.next, left)
					left.next = append(left.next, tile)
				}
			}
			if y > 0 {
				if up := network.tiles[x][y-1]; (tile.kind&N) > 0 && (up.kind&S) > 0 {
					tile.next = append(tile.next, up)
					up.next = append(up.next, tile)
				}
			}
		}
	}

	// infer start kind
	network.start.kind = 0
	for _, next := range network.start.next {
		switch {
		case next.kind&S > 0 && network.start.y > next.y && network.start.x == next.x:
			network.start.kind |= N
		case next.kind&N > 0 && network.start.y < next.y && network.start.x == next.x:
			network.start.kind |= S
		case next.kind&E > 0 && network.start.x > next.x && network.start.y == next.y:
			network.start.kind |= W
		case next.kind&W > 0 && network.start.x < next.x && network.start.y == next.y:
			network.start.kind |= E
		}
	}

	return network
}

func main() {
	fmt.Println("--- 2023 Day 10: Pipe Maze ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	network := parseNetwork(inputs)

	////////////////////////////////////////

	longestLoop := network.getLongestLoop(network.start)

	// 80
	fmt.Println("Part 1:", (len(longestLoop)-1)/2)

	////////////////////////////////////////

	tilesInLoop := network.getTilesInLoop(longestLoop)

	// 10
	fmt.Println("Part 2:", len(tilesInLoop))

	network.print()
}
