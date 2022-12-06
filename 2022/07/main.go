// https://adventofcode.com/2022/day/7

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/aurelbec/advent-of-code/utils"
)

type directory struct {
	name        string
	parent      *directory
	directories map[string]*directory
	files       map[string]int
}

func (d *directory) cd(dir string) *directory {
	if d == nil {
		return d
	}

	switch dir {
	case "/":
		if d.parent == nil {
			return d
		} else {
			return d.parent.cd(dir)
		}

	case "..":
		if d.parent == nil {
			return d
		} else {
			return d.parent
		}

	default:
		if directory, exists := d.directories[dir]; exists {
			return directory
		}
		return d.mkdir(dir)
	}
}

func (d *directory) copy(items []string) {
	for _, item := range items {
		if strings.HasPrefix(item, "dir ") { // register dir
			d.cd(item[4:])
		} else { //register file
			var size int
			var name string
			fmt.Sscanf(item, "%v %s", &size, &name)
			d.files[name] = size
		}
	}
}

func (d *directory) mkdir(dir string) *directory {
	directory := newDirectory()
	directory.parent = d
	directory.name = dir
	d.directories[dir] = directory
	return directory
}

func (d *directory) pwd() string {
	if d == nil || d.parent == nil {
		return ""
	}
	return d.parent.pwd() + "/" + d.name
}

func (d directory) size() int {
	total := 0
	for _, size := range d.files {
		total += size
	}
	for _, directory := range d.directories {
		total += directory.size()
	}
	return total
}

func (d directory) print(prefix ...string) {
	if len(prefix) == 0 {
		prefix = []string{""}
		fmt.Printf("/ (dir, size=%v)\n", d.size())
	}

	prefix[0] += "  "
	for name, size := range d.files {
		fmt.Printf("%s%s (file, size=%v)\n", prefix[0], name, size)
	}
	for name, directory := range d.directories {
		fmt.Printf("%s%s (dir, size=%v)\n", prefix[0], name, directory.size())
		directory.print(prefix[0])
	}
}

func (d *directory) walk(callback func(*directory)) {
	callback(d)
	for _, directory := range d.directories {
		directory.walk(callback)
	}
}

func newDirectory() *directory {
	return &directory{
		parent:      nil,
		name:        "",
		directories: make(map[string]*directory),
		files:       make(map[string]int),
	}
}

func main() {
	fmt.Println("--- 2022 Day 7: No Space Left On Device ---")
	defer func(start time.Time) { fmt.Println("Total time:", time.Since(start).Round(time.Microsecond)) }(time.Now())

	// init
	inputs := utils.MustReadInput("example.txt")

	pwd := newDirectory()
	for i := 0; i < len(inputs); i++ {
		switch inputs[i][2:4] {
		case "cd":
			pwd = pwd.cd(inputs[i][5:])
		case "ls":
			items := []string{}
			for i = i + 1; i < len(inputs) && inputs[i][0] != '$'; i++ {
				items = append(items, inputs[i])
			}
			pwd.copy(items)
			i--
		}
	}
	pwd = pwd.cd("/")

	////////////////////////////////////////

	size := 0
	pwd.walk(func(d *directory) {
		if s := d.size(); s < 100_000 {
			size += s
		}
	})

	// 95437
	fmt.Println("Part 1:", size)

	////////////////////////////////////////

	max := 70_000_000
	current := pwd.size()
	remain := max - current
	need := 30_000_000 - remain

	size = max
	pwd.walk(func(d *directory) {
		if s := d.size(); s > need && s < size {
			size = s
		}
	})

	// 24933642
	fmt.Println("Part 2:", size)
}
