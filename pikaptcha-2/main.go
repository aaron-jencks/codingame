package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	UP = iota
	RIGHT
	DOWN
	LEFT
)

type sides struct {
	t, r, b, l bool
}

type coord struct {
	x          int
	y          int
	dir        int
	prev_sides sides
}

func simulate(width, height int, grid [][]rune, origin coord, isLeft bool) [][]rune {
	grid[origin.y][origin.x] = '0'

	if origin.y > 0 {
		origin.prev_sides.t = grid[origin.y-1][origin.x] != '#'
	}
	if origin.y < height-1 {
		origin.prev_sides.b = grid[origin.y+1][origin.x] != '#'
	}
	if origin.x > 0 {
		origin.prev_sides.l = grid[origin.y][origin.x-1] != '#'
	}
	if origin.x < width-1 {
		origin.prev_sides.r = grid[origin.y][origin.x+1] != '#'
	}

	// simulate maze
	current := origin
	first := true
	for {
		if !first && current.x == origin.x && current.y == origin.y {
			break
		} else if first {
			first = false
		}

		var t, r, b, l bool
		if current.y > 0 {
			t = grid[current.y-1][current.x] != '#'
		}
		if current.y < height-1 {
			b = grid[current.y+1][current.x] != '#'
		}
		if current.x > 0 {
			l = grid[current.y][current.x-1] != '#'
		}
		if current.x < width-1 {
			r = grid[current.y][current.x+1] != '#'
		}

		if !(t || r || b || l) {
			// trapped
			break
		}

		grid[current.y][current.x]++

		switch current.dir {
		case UP:
			if isLeft && l && !current.prev_sides.l {
				// inner corner
				current.x--
				current.dir = LEFT
			} else if !isLeft && r && !current.prev_sides.r {
				current.x++
				current.dir = RIGHT
			} else if t {
				// passthrough
				current.y--
			} else if r {
				// corner
				current.x++
				current.dir = RIGHT
			} else if l {
				current.x--
				current.dir = LEFT
			} else {
				// dead end
				current.y++
				current.dir = DOWN
			}
		case RIGHT:
			if isLeft && t && !current.prev_sides.t {
				// inner corner
				current.y--
				current.dir = UP
			} else if !isLeft && b && !current.prev_sides.b {
				current.y++
				current.dir = DOWN
			} else if r {
				// passthrough
				current.x++
			} else if t {
				// corner
				current.y--
				current.dir = UP
			} else if b {
				current.y++
				current.dir = DOWN
			} else {
				// dead end
				current.x--
				current.dir = LEFT
			}
		case LEFT:
			if isLeft && b && !current.prev_sides.b {
				// inner corner
				current.y++
				current.dir = DOWN
			} else if !isLeft && t && !current.prev_sides.t {
				current.y--
				current.dir = UP
			} else if l {
				// passthrough
				current.x--
			} else if t {
				// corner
				current.y--
				current.dir = UP
			} else if b {
				current.y++
				current.dir = DOWN
			} else {
				// dead end
				current.x++
				current.dir = RIGHT
			}
		case DOWN:
			if isLeft && r && !current.prev_sides.r {
				// inner corner
				current.x++
				current.dir = RIGHT
			} else if !isLeft && l && !current.prev_sides.l {
				current.x--
				current.dir = LEFT
			} else if b {
				// passthrough
				current.y++
			} else if r {
				// corner
				current.x++
				current.dir = RIGHT
			} else if l {
				current.x--
				current.dir = LEFT
			} else {
				// dead end
				current.y--
				current.dir = UP
			}
		}
		current.prev_sides = sides{
			t, r, b, l,
		}
	}

	return grid
}

func main_io(in io.Reader, out io.Writer) {
	var width, height int
	fmt.Fscan(in, &width, &height)
	var origin coord
	grid := make([][]rune, height)
	for i := 0; i < height; i++ {
		var line string
		fmt.Fscan(in, &line)
		grid[i] = []rune(line)
		for j, r := range line {
			if strings.ContainsRune("<^>v", r) {
				origin = coord{
					x: j,
					y: i,
				}
				switch r {
				case '<':
					origin.dir = LEFT
				case '>':
					origin.dir = RIGHT
				case '^':
					origin.dir = UP
				case 'v':
					origin.dir = DOWN
				}
			}
		}
	}
	var side string
	fmt.Fscan(in, &side)
	isLeft := side[0] == 'L'

	grid = simulate(width, height, grid, origin, isLeft)

	for i := 0; i < height; i++ {
		fmt.Fprintln(out, string(grid[i]))
	}
}

func main() {
	main_io(os.Stdin, os.Stdout)
}
