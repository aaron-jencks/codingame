package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
)

type particle struct {
	q, m   float64
	symbol string
}

func (p particle) err(g float64) float64 {
	gp := math.Abs(p.q) / p.m
	return math.Abs(gp-g) / gp
}

type coord struct {
	x float64
	y float64
}

type icoord struct {
	x int
	y int
}

func (c icoord) toCoord() coord {
	return coord{
		float64(c.x), float64(c.y),
	}
}

func (c coord) mid(o coord) coord {
	return coord{
		x: (c.x + o.x) / 2,
		y: (c.y + o.y) / 2,
	}
}

func (c coord) gradient(o coord) float64 {
	return (o.y - c.y) / (o.x - c.x)
}

func (c coord) distance(o coord) float64 {
	a := c.x - o.x
	b := c.y - o.y
	return math.Sqrt(a*a + b*b)
}

type line struct {
	m, b float64
}

func (l line) solve(x float64) float64 {
	return l.m*x + l.b
}

func (l line) intersection(o line) coord {
	x := (o.b - l.b) / (l.m - o.m)
	return coord{
		x: x,
		y: l.solve(x),
	}
}

func findRadius(p, q, r coord) float64 {
	mpq := p.mid(q)
	mqr := q.mid(r)
	glpq := -1. / p.gradient(q)
	glqr := -1. / q.gradient(r)
	center := line{
		m: glpq,
		b: glpq*-mpq.x + mpq.y,
	}.intersection(line{
		m: glqr,
		b: glqr*-mqr.x + mqr.y,
	})
	return p.distance(center)
}

func findPoints(h, w int, grid [][]rune) []icoord {
	var result []icoord

	for i := 0; i < h; i++ {
		found := false
		for j := 0; j < w; j++ {
			if grid[i][j] == ' ' {
				found = true
				result = append(result, icoord{
					x: j,
					y: i,
				})
				break
			}
		}
		if !found {
			break
		}
	}

	return result
}

var particles = []particle{
	{-1, 0.511, "e-"},
	{0, 938, "n0"},
	{1, 938, "p+"},
	{2, 3727, "alpha"},
	{1, 140, "pi+"},
}

func main_io(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	// w: width of ASCII-art picture (one meter per column)
	var w int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &w)

	// h: height of ASCII-art picture (one meter per line)
	var h int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &h)

	// B: strengh of magnetic field (tesla)
	var B float64
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &B)

	// V: speed of the particle (speed-of-light unit)
	var V float64
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &V)

	grid := make([][]rune, h)
	for i := 0; i < h; i++ {
		scanner.Scan()
		line := scanner.Text()
		fmt.Fprintln(os.Stderr, line)
		grid[i] = []rune(line)
	}

	coords := findPoints(h, w, grid)

	var a, b, c coord
	a = coords[0].toCoord()
	b = coords[len(coords)>>1].toCoord()
	c = coords[len(coords)-1].toCoord()
	if a.gradient(b) == b.gradient(c) {
		fmt.Fprintln(os.Stderr, "found a line")
		fmt.Fprintln(out, "n0 inf")
		return
	}

	fmt.Fprintln(os.Stderr, a, b, c)

	R := math.Round(findRadius(a, b, c))
	Rrem := int(R) % 10
	if Rrem >= 5 {
		R += float64(10 - Rrem)
	} else {
		R = R - float64(Rrem) // round to multiple of 10
	}

	g := 1e6 * (1 / math.Sqrt(1-V*V)) * V / (B * R * 299792458)

	fmt.Fprintln(os.Stderr, R, g)

	index := -1
	err := 0.6
	for pi, p := range particles {
		perr := p.err(g)
		fmt.Fprintf(os.Stderr, "%s: %.2f\n", p.symbol, perr)
		if perr < 0.5 && perr < err {
			index = pi
			err = perr
		}
	}

	if index < 0 {
		fmt.Fprintln(out, "I just won the Nobel prize in physics !")
		return
	} else if particles[index].q == 0 {
		fmt.Fprintln(out, "n0 inf")
		return
	}

	// fmt.Fprintln(os.Stderr, "Debug messages...")

	fmt.Fprintf(out, "%s %d\n", particles[index].symbol, int(R))
}

func main() {
	main_io(os.Stdin, os.Stdout)
}
