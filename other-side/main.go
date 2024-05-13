package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	SPACE = '+'
)

type coord struct {
	x int
	y int
}

type set[T comparable] map[T]bool

var visited = set[coord]{}

func dfs(start, h, w int, grid [][]rune) bool {
	if _, ok := visited[coord{0, start}]; ok {
		return visited[coord{0, start}]
	}

	var last int
	var e coord
	found := false
	tvisited := set[coord]{}
	stack := []coord{{0, start}}
	for len(stack) > 0 {
		last = len(stack) - 1
		e = stack[last]
		stack = stack[:last]

		if _, ok := visited[e]; ok {
			found = visited[e]
			break
		}

		if tvisited[e] {
			continue
		}
		tvisited[e] = true

		if e.x == w-1 {
			found = true
			break
		}

		if e.x > 0 && grid[e.y][e.x-1] == SPACE && !tvisited[coord{e.x - 1, e.y}] {
			stack = append(stack, coord{e.x - 1, e.y})
		}
		if e.x < w-1 && grid[e.y][e.x+1] == SPACE && !tvisited[coord{e.x + 1, e.y}] {
			stack = append(stack, coord{e.x + 1, e.y})
		}
		if e.y > 0 && grid[e.y-1][e.x] == SPACE && !tvisited[coord{e.x, e.y - 1}] {
			stack = append(stack, coord{e.x, e.y - 1})
		}
		if e.y < h-1 && grid[e.y+1][e.x] == SPACE && !tvisited[coord{e.x, e.y + 1}] {
			stack = append(stack, coord{e.x, e.y + 1})
		}
	}

	for k := range tvisited {
		visited[k] = found
	}

	return found
}

func main_io(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	var h int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &h)

	var w int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &w)

	grid := make([][]rune, h)

	for i := 0; i < h; i++ {
		scanner.Scan()
		row := strings.ReplaceAll(scanner.Text(), " ", "")
		grid[i] = []rune(row)
	}

	count := 0
	for i := 0; i < h; i++ {
		if grid[i][0] != SPACE {
			continue
		}
		if dfs(i, h, w, grid) {
			count++
		}
	}

	// fmt.Fprintln(os.Stderr, "Debug messages...")
	fmt.Fprintln(out, count) // Write answer to stdout
}

func main() {
	main_io(os.Stdin, os.Stdout)
}
