package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const EMPTY = '.'
const VALID = "ABCDEFGHIJKLMNOP"
const SIZE = 16
const BLOCK_SIZE = SIZE >> 2

type State struct {
	Solution [][]rune
}

func (s State) String() string {
	result := ""
	for _, row := range s.Solution {
		result += string(row) + "\n"
	}
	return result
}

func (s State) Done() bool {
	for _, row := range s.Solution {
		for _, col := range row {
			if col == EMPTY {
				return false
			}
		}
	}
	return true
}

func (s State) Clone() State {
	result := State{
		Solution: make([][]rune, len(s.Solution)),
	}

	for ri, row := range s.Solution {
		result.Solution[ri] = make([]rune, SIZE)
		copy(result.Solution[ri], row)
	}

	return result
}

func (s State) ValidPlacement(r, c int, ch rune) bool {
	for i := 0; i < SIZE; i++ {
		if s.Solution[r][i] == ch || s.Solution[i][c] == ch {
			return false
		}
	}

	rblock := r >> 2 << 2
	cblock := c >> 2 << 2

	for ri := rblock; ri < rblock+BLOCK_SIZE; ri++ {
		for ci := cblock; ci < cblock+BLOCK_SIZE; ci++ {
			if s.Solution[ri][ci] == ch {
				return false
			}
		}
	}

	return true
}

type Neighbors [][]State

func (n Neighbors) Len() int {
	return len(n)
}

func (n Neighbors) Less(i, j int) bool {
	return len(n[i]) < len(n[j])
}

func (n Neighbors) Swap(i, j int) {
	temp := n[i]
	n[i] = n[j]
	n[j] = temp
}

type FilteredNeighbors []State

func (n FilteredNeighbors) Len() int {
	return len(n)
}

func (n FilteredNeighbors) Less(i, j int) bool {
	return len(FindNeighbors(n[i])) < len(FindNeighbors(n[j]))
}

func (n FilteredNeighbors) Swap(i, j int) {
	temp := n[i]
	n[i] = n[j]
	n[j] = temp
}

func FindNeighbors(s State) FilteredNeighbors {
	var ns Neighbors

	for r := 0; r < 16; r++ {
		for c := 0; c < 16; c++ {
			if s.Solution[r][c] != EMPTY {
				continue
			}

			var cell []State

			for _, ch := range VALID {
				if s.ValidPlacement(r, c, ch) {
					ns := s.Clone()
					ns.Solution[r][c] = ch
					cell = append(cell, ns)
				}
			}

			ns = append(ns, cell)
		}
	}

	sort.Sort(ns)

	return ns[0]
}

func FindSolution(initial State) State {
	stack := []State{initial}
	for len(stack) > 0 {
		last := len(stack) - 1
		element := stack[last]
		stack = stack[:last]

		if element.Done() {
			return element
		}

		neighbors := FindNeighbors(element)
		// sort.Sort(neighbors)
		for _, n := range neighbors {
			if n.Done() {
				return n
			}
			stack = append(stack, n)
		}
	}

	panic("no solution found")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	start := State{
		Solution: make([][]rune, SIZE),
	}

	for i := 0; i < SIZE; i++ {
		scanner.Scan()
		row := scanner.Text()
		start.Solution[i] = []rune(row)
	}

	fmt.Fprint(os.Stderr, start.String())

	solution := FindSolution(start)

	// fmt.Fprintln(os.Stderr, "Debug messages...")
	fmt.Print(solution.String()) // Write answer to stdout
}
