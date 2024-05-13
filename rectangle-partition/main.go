package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func find_lesser_values(values []int, limit int) []int {
	result := make([]int, 0, len(values))
	for _, v := range values {
		if v >= limit {
			break
		}
		result = append(result, v)
	}
	return result
}

func find_coords_northwest(row, column int, verts, horts []int) ([]int, []int) {
	return find_lesser_values(verts, column), find_lesser_values(horts, row)
}

func find_squares(row, column int, verts, horts []int, grid [][]bool) int {
	count := 0
	fverts, fhorts := find_coords_northwest(row, column, verts, horts)
	for _, v := range fverts {
		for _, h := range fhorts {
			fmt.Fprintf(os.Stderr, "checking (%d, %d) compared to (%d, %d)\n", v, h, column, row)
			if column-v == row-h && grid[h][v] {
				fmt.Fprintf(os.Stderr, "found a square of size %dx%d\n", column-v, row-h)
				count++
			}
		}
	}
	return count
}

func main_reader(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 1000000), 1000000)
	var inputs []string

	var w, h, countX, countY int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &w, &h, &countX, &countY)
	fmt.Fprintf(os.Stderr, "%d %d %d %d\n", w, h, countX, countY)

	grid := make([][]bool, h+1)
	for row := 0; row < h+1; row++ {
		grid[row] = make([]bool, w+1)
	}

	vs := make([]int, countX+1)
	hs := make([]int, countY+1)

	vs[0] = 0
	hs[0] = 0

	scanner.Scan()
	inputs = strings.Split(scanner.Text(), " ")
	for xi, xs := range inputs {
		fmt.Sscan(xs, &vs[xi+1])
		fmt.Fprint(os.Stderr, xs, " ")
	}
	fmt.Fprintln(os.Stderr)

	scanner.Scan()
	inputs = strings.Split(scanner.Text(), " ")
	for yi, ys := range inputs {
		fmt.Sscan(ys, &hs[yi+1])
		fmt.Fprint(os.Stderr, ys, " ")
	}
	fmt.Fprintln(os.Stderr)

	vs = append(vs, w)
	hs = append(hs, h)

	// four corners
	grid[0][0] = true
	grid[h][0] = true
	grid[0][w] = true
	grid[h][w] = true

	for _, col := range vs {
		grid[0][col] = true
		grid[h][col] = true

		for _, row := range hs {
			grid[row][col] = true
		}
	}

	for _, row := range hs {
		grid[row][0] = true
		grid[row][w] = true

		for _, col := range vs {
			grid[row][col] = true
		}
	}

	count := 0

	for _, row := range hs {
		for _, col := range vs {
			count += find_squares(row, col, vs, hs, grid)
		}
	}

	// fmt.Fprintln(os.Stderr, "Debug messages...")
	fmt.Fprintln(out, count) // Write answer to stdout
}

func main() {
	main_reader(os.Stdin, os.Stdout)
}
