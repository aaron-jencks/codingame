package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Validate(s string, transitions map[int]map[rune]int, start int, accepting map[int]bool) bool {
	current := start
	for _, sc := range s {
		next := transitions[current][sc]
		// fmt.Fprintf(os.Stderr, "%d -> %d for %c\n", current, next, sc)
		current = next
	}
	// fmt.Fprintf(os.Stderr, "case %s ended up in state %d which is %v\n", s, current, accepting[current])
	return accepting[current]
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	var N int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &N)

	scanner.Scan()
	scanner.Text()

	var start int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &start)

	var T int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &T)

	transitions := map[int]map[rune]int{}
	accepting := map[int]bool{}

	for i := 0; i < T; i++ {
		var source, target int
		var symbol string
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &source, &target, &symbol)

		accepting[source] = false
		accepting[target] = false

		if _, ok := transitions[source]; !ok {
			transitions[source] = map[rune]int{
				rune(symbol[0]): target,
			}
		} else {
			transitions[source][rune(symbol[0])] = target
		}
	}

	scanner.Scan()
	scanner.Text()

	scanner.Scan()
	aaccepting := strings.Split(scanner.Text(), " ")
	for _, acc := range aaccepting {
		var iacc int
		fmt.Sscan(acc, &iacc)
		accepting[iacc] = true
	}

	var C int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &C)

	for i := 0; i < C; i++ {
		scanner.Scan()
		testcase := scanner.Text()
		accepted := Validate(testcase, transitions, start, accepting)
		if accepted {
			fmt.Println("accept")
		} else {
			fmt.Println("reject")
		}
	}
}
