// https://adventofcode.com/2023/day/25

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
	"github.com/aurelbec/advent-of-code/utils/collections"
)

type Node struct {
	neighbors map[*Node]*Edge
	group     int
}

type Edge struct {
	*Edge // avoid empty struct, mess-up with language
}

type Graph struct {
	nodes map[*Node]struct{}
	edges map[*Edge]struct{}
}

func (g *Graph) bfs(source, dest *Node) bool {
	queue := collections.NewQueue(source)
	visited := make(map[*Node]*Node, len(g.nodes))
	visited[source] = nil
	for !queue.IsEmpty() {
		current, _ := queue.Dequeue()

		// tag group if BFS is executed blindly
		if dest == nil && visited[current] != nil {
			current.group = visited[current].group
		}

		// reconstruct path if dest is found, forbidden edges for next step
		if current == dest {
			for prev := visited[current]; prev != nil; current, prev = prev, visited[prev] {
				g.edges[current.neighbors[prev]] = struct{}{}
			}
			return true
		}

		// explore next nodes if possible
		for next, edge := range current.neighbors {
			if _, visited := visited[next]; visited {
				continue
			} else if _, visited := g.edges[edge]; visited {
				continue
			}

			visited[next] = current
			queue.Enqueue(next)
		}
	}

	return false
}

func (g *Graph) cutPaths(source, dest *Node, n int) bool {
	g.edges = make(map[*Edge]struct{})

	for i := 0; i <= n; i++ {
		if !g.bfs(source, dest) {
			return false
		}
	}

	return true
}

func (g *Graph) split(cuts int) (int, int) {
	g.edges = make(map[*Edge]struct{})

	var source *Node
	for node := range g.nodes {
		if len(node.neighbors) > cuts {
			source = node
			source.group = 1
			break
		}
	}

	for dest := range g.nodes {
		if dest.group > 0 {
			continue
		}

		if !g.cutPaths(source, dest, cuts) {
			// use disconnected graph to categorize as many nodes as possible
			g.bfs(source, nil)
			dest.group = 2
			g.bfs(dest, nil)
		}
	}

	lhs, rhs := 0, 0
	for node := range g.nodes {
		if node.group == 1 {
			lhs++
		} else {
			rhs++
		}
	}
	return lhs, rhs
}

func parseGraph(inputs []string) *Graph {
	nodes := make(map[string]*Node, len(inputs))
	for _, input := range inputs {
		names := strings.Fields(strings.ReplaceAll(input, ":", ""))
		from := names[0]
		if _, found := nodes[from]; !found {
			nodes[from] = &Node{neighbors: make(map[*Node]*Edge)}
		}

		for _, to := range names[1:] {
			if _, found := nodes[to]; !found {
				nodes[to] = &Node{neighbors: make(map[*Node]*Edge)}
			}

			edge := new(Edge)
			nodes[from].neighbors[nodes[to]] = edge
			nodes[to].neighbors[nodes[from]] = edge
		}
	}

	graph := Graph{nodes: make(map[*Node]struct{}, len(nodes))}
	for _, node := range nodes {
		graph.nodes[node] = struct{}{}
	}
	return &graph
}

func main() {
	fmt.Println("--- 2023 Day 25: Snowverload ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	graph := parseGraph(inputs)

	////////////////////////////////////////

	lhs, rhs := graph.split(3)

	// 54
	fmt.Println("Part 1:", lhs*rhs)

	////////////////////////////////////////

	// 0
	fmt.Println("Part 2:", 0)
}
