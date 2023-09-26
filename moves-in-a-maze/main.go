package main

import (
	"bufio"
	"fmt"
	"os"
)

type State struct {
	row      int
	col      int
	distance int
}

func GenerateIntMapFromString(h, w int, ss []string) [][]int {
	result := make([][]int, h)
	for i := 0; i < h; i++ {
		result[i] = make([]int, w)
		for j := 0; j < w; j++ {
			switch ss[i][j] {
			case 'S':
				result[i][j] = -1
			case '#':
				result[i][j] = -2
			case '.':
				result[i][j] = -1
			}
		}
	}
	return result
}

func CheckNeighbor(q []State, ca, cb bool, ar, ac, br, bc, d int, m [][]int) []State {
	if ca && (m[ar][ac] == -1 || m[ar][ac] > d) {
		q = append(q, State{
			row:      ar,
			col:      ac,
			distance: d + 1,
		})
	} else if cb && (m[br][bc] == -1 || m[br][bc] > d) {
		q = append(q, State{
			row:      br,
			col:      bc,
			distance: d + 1,
		})
	}
	return q
}

func FindDistances(sx, sy, h, w int, m []string) [][]int {
	md := GenerateIntMapFromString(h, w, m)

	q := make([]State, 0, 900)
	q = append(q, State{
		row: sy,
		col: sx,
	})
	for len(q) > 0 {
		e := q[0]
		q = q[1:]

		if md[e.row][e.col] != -1 && md[e.row][e.col] <= e.distance {
			continue
		}

		md[e.row][e.col] = e.distance

		// find neighbors

		// up
		q = CheckNeighbor(q, e.row > 0, e.row == 0, e.row-1, e.col, h-1, e.col, e.distance, md)

		// down
		q = CheckNeighbor(q, e.row < h-1, e.row == h-1, e.row+1, e.col, 0, e.col, e.distance, md)

		// left
		q = CheckNeighbor(q, e.col > 0, e.col == 0, e.row, e.col-1, e.row, w-1, e.distance, md)

		// right
		q = CheckNeighbor(q, e.col < w-1, e.col == w-1, e.row, e.col+1, e.row, 0, e.distance, md)
	}

	return md
}

func GenerateOutputFromArray(h, w int, md [][]int) []string {
	var result []string
	for i := 0; i < h; i++ {
		// fmt.Fprintln(os.Stderr, "Debug messages...")
		rout := ""
		for j := 0; j < w; j++ {
			switch md[i][j] {
			case -2:
				rout += "#"
			case -1:
				rout += "."
			default:
				start := '0'
				d := md[i][j]
				if d > 9 {
					start = 'A'
					d -= 10
				}
				rout += string(start + rune(d))
			}
		}
		result = append(result, rout)
	}
	return result
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	var w, h int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &w, &h)

	fmt.Fprintf(os.Stderr, "%d %d\n", w, h)

	m := make([]string, h)
	var sx, sy int
	var foundS bool

	for i := 0; i < h; i++ {
		scanner.Scan()
		ROW := scanner.Text()
		fmt.Fprintln(os.Stderr, ROW)
		m[i] = ROW
		if !foundS {
			for j, rc := range ROW {
				if rc == 'S' {
					sx = j
					sy = i
					foundS = true
					break
				}
			}
		}
	}

	md := FindDistances(sx, sy, h, w, m)

	for _, l := range GenerateOutputFromArray(h, w, md) {
		// fmt.Fprintln(os.Stderr, "Debug messages...")
		fmt.Println(l) // Write answer to stdout
	}
}
