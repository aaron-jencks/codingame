package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	WATER = '~'
	LAND  = '#'
)

type coord struct {
	x int
	y int
}

func (c coord) equals(o coord) bool {
	return c.x == o.x && c.y == o.y
}

func (c coord) waterNeighborCount(visited coordMap, limit int, grid []string) int {
	result := 0
	if c.x > 0 && grid[c.y][c.x-1] == WATER && !visited.contains(coord{c.x - 1, c.y}) {
		visited.insert(coord{c.x - 1, c.y})
		result++
	}
	if c.x < limit-1 && grid[c.y][c.x+1] == WATER && !visited.contains(coord{c.x + 1, c.y}) {
		visited.insert(coord{c.x + 1, c.y})
		result++
	}
	if c.y > 0 && grid[c.y-1][c.x] == WATER && !visited.contains(coord{c.x, c.y - 1}) {
		visited.insert(coord{c.x, c.y - 1})
		result++
	}
	if c.y < limit-1 && grid[c.y+1][c.x] == WATER && !visited.contains(coord{c.x, c.y + 1}) {
		visited.insert(coord{c.x, c.y + 1})
		result++
	}
	return result
}

func (c coord) findNeighborLand(limit int, grid []string) []coord {
	var result []coord

	if c.x > 0 && grid[c.y][c.x-1] == LAND {
		result = append(result, coord{c.x - 1, c.y})
	}

	if c.y > 0 && grid[c.y-1][c.x] == LAND {
		result = append(result, coord{c.x, c.y - 1})
	}
	if c.y < limit-1 && grid[c.y+1][c.x] == LAND {
		result = append(result, coord{c.x, c.y + 1})
	}

	if c.x < limit-1 && grid[c.y][c.x+1] == LAND {
		result = append(result, coord{c.x + 1, c.y})
	}

	return result
}

type coordMap map[coord]bool

func (vm coordMap) contains(c coord) bool {
	_, ok := vm[c]
	return ok
}

func (vm coordMap) insert(c coord) {
	vm[c] = true
}

type island struct {
	index     int
	origin    coord
	coastline int
}

func stackContains(stack []coord, e coord) bool {
	for _, se := range stack {
		if se.equals(e) {
			return true
		}
	}
	return false
}

func mapIsland(index int, origin coord, visited coordMap, limit int, grid []string) island {
	result := island{
		index:  index,
		origin: origin,
	}

	coastVisited := coordMap{}
	stack := []coord{origin}
	for len(stack) > 0 {
		e := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited.contains(e) {
			continue
		}

		visited.insert(e)

		// count water edges
		result.coastline += e.waterNeighborCount(coastVisited, limit, grid)

		// append children
		for _, child := range e.findNeighborLand(limit, grid) {
			if !(visited.contains(child) || stackContains(stack, child)) {
				stack = append(stack, child)
			}
		}
	}

	return result
}

func mainIo(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	var n int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &n)
	fmt.Fprintln(os.Stderr, n)

	grid := make([]string, n)
	for i := 0; i < n; i++ {
		scanner.Scan()
		grid[i] = scanner.Text()
		fmt.Fprintln(os.Stderr, grid[i])
	}

	var islands []island
	visited := coordMap{}
	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {
			if visited.contains(coord{x, y}) {
				continue
			}

			if grid[y][x] == LAND {
				islands = append(islands, mapIsland(len(islands)+1, coord{x, y}, visited, n, grid))
			}

			visited.insert(coord{x, y})
		}
	}

	max := island{
		index: 1,
	}
	for _, isle := range islands {
		if isle.coastline > max.coastline {
			max = isle
		}
	}

	fmt.Fprintf(out, "%d %d\n", max.index, max.coastline)
}

func main() {
	mainIo(os.Stdin, os.Stdout)
}
