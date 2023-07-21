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

type Coord struct {
	Row int
	Col int
}

type Openings struct {
	Rows   []int
	Cols   []int
	Quad   [][]int
	Coords []Coord
}

func (o Openings) GetCount(r, c int) int {
	return o.Rows[r] + o.Cols[c] + o.Quad[r>>2][c>>2]
}

func (o Openings) Len() int {
	return len(o.Coords)
}

func (o Openings) Less(i, j int) bool {
	ic := o.Coords[i]
	jc := o.Coords[j]
	return o.GetCount(ic.Row, ic.Col) < o.GetCount(jc.Row, jc.Col)
}

func (n Openings) Swap(i, j int) {
	temp := n.Coords[i]
	n.Coords[i] = n.Coords[j]
	n.Coords[j] = temp
}

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
		if ri == r {
			continue
		}

		for ci := cblock; ci < cblock+BLOCK_SIZE; ci++ {
			if ci == c {
				continue
			}

			if s.Solution[ri][ci] == ch {
				return false
			}
		}
	}

	return true
}

func (s State) FindValidValues(r, c int) []rune {
	isValid := map[rune]bool{}
	for _, r := range VALID {
		isValid[r] = true
	}

	if s.Solution[r][c] != EMPTY && isValid[s.Solution[r][c]] {
		isValid[s.Solution[r][c]] = false
	}

	for i := 0; i < SIZE; i++ {
		if i != c && s.Solution[r][i] != EMPTY && isValid[s.Solution[r][i]] {
			isValid[s.Solution[r][i]] = false
		}
		if i != r && s.Solution[i][c] != EMPTY && isValid[s.Solution[i][c]] {
			isValid[s.Solution[i][c]] = false
		}
	}

	rblock := r >> 2 << 2
	cblock := c >> 2 << 2

	for ri := rblock; ri < rblock+BLOCK_SIZE; ri++ {
		if ri == r {
			continue
		}

		for ci := cblock; ci < cblock+BLOCK_SIZE; ci++ {
			if ci == c {
				continue
			}

			if s.Solution[ri][ci] != EMPTY && isValid[s.Solution[ri][ci]] {
				isValid[s.Solution[ri][ci]] = false
			}
		}
	}

	result := make([]rune, 0, len(isValid))
	for k, v := range isValid {
		if v {
			result = append(result, k)
		}
	}
	return result
}

func (s State) FindOpenings() Openings {
	result := Openings{
		Rows: make([]int, SIZE),
		Cols: make([]int, SIZE),
		Quad: make([][]int, BLOCK_SIZE),
	}

	for br := 0; br < BLOCK_SIZE; br++ {
		result.Quad[br] = make([]int, BLOCK_SIZE)
	}

	for r := 0; r < 16; r++ {
		for c := 0; c < 16; c++ {
			if s.Solution[r][c] == EMPTY {
				rblock := r >> 2
				cblock := c >> 2

				result.Rows[r]++
				result.Cols[c]++
				result.Quad[rblock][cblock]++
				result.Coords = append(result.Coords, Coord{r, c})
			}
		}
	}

	return result
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
	return len(FindNeighbors(n[i])) > len(FindNeighbors(n[j]))
}

func (n FilteredNeighbors) Swap(i, j int) {
	temp := n[i]
	n[i] = n[j]
	n[j] = temp
}

func FindNeighbors(s State) FilteredNeighbors {
	openings := s.FindOpenings()
	sort.Sort(openings)
	for _, coord := range openings.Coords {
		validValues := s.FindValidValues(coord.Row, coord.Col)

		cell := make([]State, 0, len(validValues))

		for _, ch := range validValues {
			ns := s.Clone()
			ns.Solution[coord.Row][coord.Col] = ch
			cell = append(cell, ns)
		}

		return cell
	}

	return nil
}

func FindSolution(initial State) State {
	stack := FilteredNeighbors{initial}
	for len(stack) > 0 {
		last := len(stack) - 1
		element := stack[last]
		stack = stack[:last]

		if element.Done() {
			return element
		}

		neighbors := FindNeighbors(element)
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
