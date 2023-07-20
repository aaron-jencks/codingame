package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTest1(t *testing.T) {
	sdata, err := os.ReadFile("./test1.txt")
	if err != nil {
		assert.True(t, false, "failed to open test1.txt")
	}

	odata, err := os.ReadFile("./test1.out.txt")
	if err != nil {
		assert.True(t, false, "failed to open test1.out.txt")
	}

	lines := strings.Split(string(sdata), "\n")

	start := State{
		Solution: make([][]rune, SIZE),
	}

	for i, line := range lines {
		start.Solution[i] = []rune(line)
	}

	fmt.Fprint(os.Stderr, start.String())

	solution := FindSolution(start)

	// fmt.Fprintln(os.Stderr, "Debug messages...")
	fmt.Print(solution.String()) // Write answer to stdout

	assert.Equal(t, string(odata), solution.String(), "expected solutions to match")
}

func TestTest3(t *testing.T) {
	sdata, err := os.ReadFile("./test3.txt")
	if err != nil {
		assert.True(t, false, "failed to open test1.txt")
	}

	// odata, err := os.ReadFile("./test1.out.txt")
	// if err != nil {
	// 	assert.True(t, false, "failed to open test1.out.txt")
	// }

	lines := strings.Split(string(sdata), "\n")

	start := State{
		Solution: make([][]rune, SIZE),
	}

	for i, line := range lines {
		start.Solution[i] = []rune(line)
	}

	fmt.Fprint(os.Stderr, start.String())

	solution := FindSolution(start)

	// fmt.Fprintln(os.Stderr, "Debug messages...")
	fmt.Print(solution.String()) // Write answer to stdout

	// assert.Equal(t, string(odata), solution.String(), "expected solutions to match")
}
