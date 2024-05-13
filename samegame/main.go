package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	SIZE = 15
	BEAM = 2
)

type coord struct {
	x int
	y int
}

func findCoordNeighbors(c coord) []coord {
	var result []coord

	if c.x > 0 {
		result = append(result, coord{c.x - 1, c.y})
	}

	if c.y > 0 {
		result = append(result, coord{c.x, c.y - 1})
	}

	if c.y < SIZE-1 {
		result = append(result, coord{c.x, c.y + 1})
	}

	if c.x < SIZE-1 {
		result = append(result, coord{c.x + 1, c.y})
	}

	return result
}

type board struct {
	data  [][]int
	moves []coord
	score int
}

func createBoard() board {
	result := board{
		data: make([][]int, SIZE),
	}
	for i := 0; i < SIZE; i++ {
		result.data[i] = make([]int, SIZE)
	}
	return result
}

func (b board) solved() bool {
	for r := 0; r < SIZE; r++ {
		for c := 0; c < SIZE; c++ {
			if b.data[r][c] != -1 {
				return false
			}
		}
	}
	return true
}

func (b board) clone() board {
	result := board{
		data:  make([][]int, SIZE),
		moves: make([]coord, len(b.moves)),
		score: b.score,
	}
	for i := range b.data {
		result.data[i] = make([]int, SIZE)
		copy(result.data[i], b.data[i])
	}
	copy(result.moves, b.moves)
	return result
}

func (b board) findValidMoves() []coord {
	visited := make([][]bool, SIZE)
	for r := 0; r < SIZE; r++ {
		visited[r] = make([]bool, SIZE)
	}

	var q []coord = make([]coord, 0, SIZE*SIZE)
	for r := 0; r < SIZE; r++ {
		for c := 0; c < SIZE; c++ {
			if b.data[r][c] != -1 {
				q = append(q, coord{
					x: c,
					y: r,
				})
			}
		}
	}

	floodfill := func(c coord) int {
		count := 0
		q := []coord{c}
		for len(q) > 0 {
			e := q[0]
			q = q[1:]

			if visited[e.y][e.x] {
				continue
			}

			visited[e.y][e.x] = true
			count++

			neighbors := findCoordNeighbors(e)
			for _, n := range neighbors {
				if !visited[n.y][n.x] && b.data[n.y][n.x] == b.data[c.y][c.x] {
					q = append(q, n)
				}
			}
		}
		return count
	}

	var possible []coord
	for len(q) > 0 {
		e := q[0]
		q = q[1:]

		if visited[e.y][e.x] {
			continue
		}

		count := floodfill(e)
		if count == 1 {
			continue
		}

		possible = append(possible, e)
	}

	return possible
}

func (b *board) move(c coord) {
	b.moves = append(b.moves, c)

	// remove the region
	visited := make([][]bool, SIZE)
	for r := 0; r < SIZE; r++ {
		visited[r] = make([]bool, SIZE)
	}
	count := 0
	var xmin, xmax, ymin, ymax int = SIZE, 0, SIZE, 0
	q := []coord{c}
	for len(q) > 0 {
		e := q[0]
		q = q[1:]

		if visited[e.y][e.x] {
			continue
		} else {
			visited[e.y][e.x] = true
			count++
			if e.x < xmin {
				xmin = e.x
			}
			if e.x > xmax {
				xmax = e.x
			}
			if e.y < ymin {
				ymin = e.y
			}
			if e.y > ymax {
				ymax = e.y
			}
		}

		neighbors := findCoordNeighbors(e)
		for _, n := range neighbors {
			if !visited[n.y][n.x] && b.data[n.y][n.x] == b.data[e.y][e.x] {
				q = append(q, n)
			}
		}

		b.data[e.y][e.x] = -1
	}

	if ymin > 0 {
		ymin--
	}

	// from top to bottom
	// collapse cells downward
	for r := ymin; r < ymax; r++ {
		for c := xmin; c <= xmax; c++ {
			if b.data[r][c] != -1 && b.data[r+1][c] == -1 {
				// move downward

				// determine the gap size
				gsize := 0
				for cr := r + 1; cr < SIZE && b.data[cr][c] == -1; cr++ {
					gsize++
				}

				// move everything down gsize
				for cr := r; cr >= 0 && b.data[cr][c] != -1; cr-- {
					b.data[cr+gsize][c] = b.data[cr][c]
					b.data[cr][c] = -1
				}
			}
		}
	}

	if xmin > 0 {
		xmin--
	}

	// if ymax is the last row
	// then move empty rows to the left
	if ymax == SIZE-1 {
		last := SIZE - 1
		for c := xmax; c > xmin; c-- {
			if b.data[last][c] != -1 && b.data[last][c-1] == -1 {
				// move leftward

				// determine the gap size
				gsize := 0
				for cc := c - 1; cc >= 0 && b.data[last][cc] == -1; cc-- {
					gsize++
				}

				// move everything left gsize
				for cc := c; cc < SIZE && b.data[last][cc] != -1; cc++ {
					b.data[last][cc-gsize] = b.data[last][cc]
					b.data[last][cc] = -1
				}
			}
		}
	}

	b.score += (count - 2) * (count - 2)
}

func (b board) solve() ([]coord, int) {
	moves := b.findValidMoves()
	if len(moves) == 0 {
		if b.solved() {
			b.score += 1000
		}
		return b.moves, b.score
	}

	children := make([]board, len(moves))
	for mi, m := range moves {
		children[mi] = b.clone()
		children[mi].move(m)
	}
	sort.Slice(children, func(i, j int) bool {
		// descending order
		return children[i].score > children[j].score
	})

	var bestMove []coord
	bestScore := -1
	for c := 0; c < len(children) && c < BEAM; c++ {
		cmoves, cscore := children[c].solve()
		if cscore > bestScore {
			bestMove = cmoves
			bestScore = cscore
		}
	}
	return bestMove, bestScore
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)
	var inputs []string

	var actions []coord
	msg := ""

	for {
		base := createBoard()

		for i := 0; i < 15; i++ {
			scanner.Scan()
			inputs = strings.Split(scanner.Text(), " ")
			for j := 0; j < 15; j++ {
				fmt.Sscan(inputs[j], &base.data[i][j])
			}
			fmt.Fprintln(os.Stderr, strings.Join(inputs, " "))
		}

		if len(actions) == 0 {
			moves, _ := base.solve()
			actions = append(actions, moves...)
		}

		if len(actions) == 0 {
			panic("ran out of legal actions!")
		}

		act := actions[0]
		actions = actions[1:]

		fmt.Printf("%d %d %s\n", act.x, (SIZE-1)-act.y, msg)
	}
}
