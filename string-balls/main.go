package main

import (
	"bufio"
	"fmt"
	"os"
)

type State struct {
	s string
	d int
}

func FindCount(center string, radius int) int {
	var sum int = 1
	visited := map[string]bool{
		center: true,
	}

	stack := make([]State, 1, 1000000)
	stack[0] = State{
		s: center,
		d: 0,
	}
	for len(stack) > 0 {
		// last := len(stack) - 1
		// element := stack[last]
		// stack = stack[:last]
		element := stack[0]
		stack = stack[1:]

		// Find neighbors of the given string
		ncount := 0
		for ri, r := range element.s {
			celement := []rune(element.s)
			var selement string
			wrap := radius-element.d > 25

			increment := 1
			n := r - 1
			if r == 'a' && wrap {
				n = 'z'
				increment = 25
			}
			if n >= 'a' && n <= 'z' {
				celement[ri] = n
				selement = string(celement)
				if v, ok := visited[selement]; !(ok && v) {
					ncount++
					visited[selement] = true
					if element.d+increment < radius {
						stack = append(stack, State{
							s: selement,
							d: element.d + increment,
						})
					}
				}
			}

			increment = 1
			n = r + 1
			if r == 'z' && wrap {
				n = 'a'
				increment = 25
			}
			if n >= 'a' && n <= 'z' {
				celement[ri] = n
				selement = string(celement)
				if v, ok := visited[selement]; !(ok && v) {
					ncount++
					visited[selement] = true
					if element.d+increment < radius {
						stack = append(stack, State{
							s: selement,
							d: element.d + increment,
						})
					}
				}
			}
		}
		// fmt.Fprintf(os.Stderr, "Found %d neighbors for %s at depth %d\n", ncount, element.s, element.d)
		sum += ncount
	}

	return sum
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	var radius int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &radius)

	scanner.Scan()
	center := scanner.Text()

	fmt.Fprintln(os.Stderr, radius)
	fmt.Fprintln(os.Stderr, center)

	// abcdefghijklmnopqrstubwxyz
	// ab
	// bb ac aa
	// distance of 1: numLettersLtZ + numLettersGtA
	// distance of 2: numLettersLtY + numLettersGtB
	// distance of 3:
	// ...
	// distance of radius:

	// A bat visited you
	// fmt.Fprintln(os.Stderr, "Debug messages...")
	fmt.Println(FindCount(center, radius)) // Write answer to stdout
}
