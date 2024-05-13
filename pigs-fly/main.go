package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type set map[string]bool

func (s set) contains(e string) bool {
	_, ok := s[e]
	return ok
}

func (s set) add(e string) {
	s[e] = true
}

type stats struct {
	traits    []string
	abilities []string
}

type transition struct {
	requirements stats
	target       string
}

var objects = map[string][]transition{}
var globalStats = map[string]stats{}

func findTraits(tokens []string, traits set) {
	if len(tokens) == 3 && tokens[1] == "are" && (traits.contains(tokens[0]) || traits.contains(tokens[2])) {
		traits.add(tokens[0])
		traits.add(tokens[2])
		return
	}

	parseTraitList := func(start int) int {
		if start >= len(tokens) || tokens[start] != "and" {
			return start
		}
		for start < len(tokens) && tokens[start] == "and" {
			traits.add(tokens[start+1])
			start += 2
		}
		return start
	}
	parseTraits := func(start int) int {
		if start >= len(tokens) || tokens[start] != "with" {
			return start
		}
		traits.add(tokens[start+1])
		return parseTraitList(start + 2)
	}
	parseAbilityList := func(start int) int {
		if start >= len(tokens) || tokens[start] != "and" {
			return start
		}
		for start < len(tokens) && tokens[start] == "and" {
			start += 2
		}
		return start
	}
	parseAbilities := func(start int) int {
		if start >= len(tokens) || tokens[start] != "that" {
			return start
		}
		return parseAbilityList(start + 3)
	}
	parseObject := func(start int) int {
		return parseAbilities(parseTraits(start + 1))
	}

	vindex := parseObject(0)
	switch tokens[vindex] {
	case "are":
		parseObject(vindex + 1)
	case "have":
		traits.add(tokens[vindex+1])
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	var N int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &N)

	var tokenList [][]string

	for i := 0; i < N; i++ {
		scanner.Scan()
		S := scanner.Text()
		fmt.Fprintln(os.Stderr, S)
		tokens := strings.Split(S, " ")
		tokenList = append(tokenList, tokens)
	}

	// first pass to determine traits vs objects
	var traits set = set{}
	for _, tok := range tokenList {
		findTraits(tok, traits)
	}

	fmt.Fprintln(os.Stderr, traits)

	// fmt.Fprintln(os.Stderr, "Debug messages...")
	fmt.Printf("All|Some|No  pigs can fly\n")
}
