package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func MustReadInput(name string) []string {
	input, err := ReadInput(name)
	if err != nil {
		panic(err)
	}
	return input
}

func ReadInput(name string) ([]string, error) {
	// retrieve caller file
	_, current, _, _ := runtime.Caller(0)
	file := ""
	for i, ok := 1, true; ok; i++ {
		_, file, _, ok = runtime.Caller(i)
		if !ok {
			return nil, fmt.Errorf("fail to get caller to read input")
		}
		if file != current {
			break
		}
	}

	// switch to input path
	file = filepath.Join(filepath.Dir(file), name)

	// open input
	input, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer input.Close()

	//scan the input's contents line by line
	lines := make([]string, 0, 1000)
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// return error if scanning is not done properly
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading from file: %v", err)
	}

	return lines, nil
}
