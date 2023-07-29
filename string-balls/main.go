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

func FindLetterNeighbor(l rune, wrap bool) []rune {
	var result []rune
	if l > 'a' {
		result = append(result, l-1)
	}
	if l < 'z' {
		result = append(result, l+1)
	}
	if wrap {
		if l == 'a' {
			result = append(result, 'z')
		} else if l == 'z' {
			result = append(result, 'a')
		}
	}
	return result
}

func FindCount(center string, radius int) int {
	var sum int = 1
	visited := map[string]bool{
		center: true,
	}

	stack := []State{
		{
			s: center,
			d: 0,
		},
	} // set capacity to 1e6
	for len(stack) > 0 {
		// last := len(stack) - 1
		// element := stack[last]
		// stack = stack[:last]
		element := stack[0]
		stack = stack[1:]

		fmt.Fprintln(os.Stderr, element)

		if element.d >= radius {
			continue
		}

		// Find neighbors of the given string
		for ri, r := range element.s {
			wrap := radius-element.d > 25
			for _, n := range FindLetterNeighbor(r, wrap) {
				increment := 1
				if (r == 'a' && n == 'z') || (r == 'z' && n == 'a') {
					increment = 25
				}
				celement := []rune(element.s)
				celement[ri] = n
				if v, ok := visited[string(celement)]; !(ok && v) {
					sum++
					visited[string(celement)] = true
					stack = append(stack, State{
						s: string(celement),
						d: element.d + increment,
					})
				}
			}
		}
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
