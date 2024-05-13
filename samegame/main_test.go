package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readBoardDataFromFile(t *testing.T, f string) board {
	data, err := os.ReadFile(f)
	assert.NoError(t, err, "unexpected error while reading file")

	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	scanner.Buffer(make([]byte, 1000000), 1000000)
	var inputs []string

	base := createBoard()

	for i := 0; i < 15; i++ {
		scanner.Scan()
		inputs = strings.Split(scanner.Text(), " ")
		for j := 0; j < 15; j++ {
			fmt.Sscan(inputs[j], &base.data[i][j])
		}
	}

	return base
}

func TestHorizontal(t *testing.T) {
	b := readBoardDataFromFile(t, "./horizontal.txt")
	moves, score := b.solve()
	fmt.Printf("solved with %d moves and a score of %d\n", len(moves), score)
}

func TestVertical(t *testing.T) {
	b := readBoardDataFromFile(t, "./vertical.txt")
	moves, score := b.solve()
	fmt.Printf("solved with %d moves and a score of %d\n", len(moves), score)
}
