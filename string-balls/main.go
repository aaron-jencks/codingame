package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	NONE = iota
	DOWN
	UP
)

type State struct {
	s            string
	d            int
	modIndex     int
	modDirection int
}

var dupCount int = 0

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
			var increment int
			var n rune

			if ri != element.modIndex || element.modDirection != UP {
				increment = 1
				n = r - 1
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
								s:            selement,
								d:            element.d + increment,
								modIndex:     ri,
								modDirection: DOWN,
							})
						}
					} else {
						dupCount++
					}
				}
			}

			if r == 'z' && !wrap {
				continue
			}

			if ri != element.modIndex || element.modDirection != DOWN {
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
								s:            selement,
								d:            element.d + increment,
								modIndex:     ri,
								modDirection: UP,
							})
						}
					} else {
						dupCount++
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
	fmt.Println(FindCount(center, radius)) // Write answer to stdout
	fmt.Fprintf(os.Stderr, "Duplicates found: %d\n", dupCount)
}

// TODO
// if we can restrict the direction that each letter can go, then we can eliminate the possibilities of duplicates
// we begin by populating the stack with words where each individual letter is assigned every possible direction
// this takes care of radius 1
// from there we can for each element of the stack, create a copy where each other letter is assigned every possible direction
// checking to see if that combination of directions has been covered yet.
// after each round of checks we move the letters in the given direction and push them back onto the stack
// we repeat this until we reach the desired depth
// this should avoid duplicate checks because we will have an entity on the stack for each unique possible combinations of directions
// and we can complete this in exactly radius ticks of the stack

// keep the old version for testing since we know it works and compare it with the new version.
