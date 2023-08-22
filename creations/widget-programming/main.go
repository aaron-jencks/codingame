package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type dfa struct {
	start       int
	transitions map[int]map[rune]int
	accepting   map[int]bool
}

func NewDfa() dfa {
	return dfa{
		transitions: map[int]map[rune]int{},
		accepting:   map[int]bool{},
	}
}

func (d *dfa) addEdge(source int, target int, symbol rune) {
	if v, ok := d.transitions[source]; ok {
		v[symbol] = target
	} else {
		d.transitions[source] = map[rune]int{
			symbol: target,
		}
	}
}

func (d *dfa) setAccepting(state int, accept bool) {
	d.accepting[state] = accept
}

func (d dfa) validate(s string) bool {
	current := d.start
	for _, sc := range s {
		next := d.transitions[current][sc]
		current = next
	}
	return d.accepting[current]
}

type nfa struct {
	start       int
	transitions map[int]map[rune][]int
	accepting   map[int]bool
}

func NewNfa() nfa {
	return nfa{
		transitions: map[int]map[rune][]int{},
		accepting:   map[int]bool{},
	}
}

func (n *nfa) setAccepting(state int, accept bool) {
	n.accepting[state] = accept
}

func (n nfa) containsEdge(source, target int, symbol rune) bool {
	if v, ok := n.transitions[source]; ok {
		if tl, ok := v[symbol]; ok {
			for _, t := range tl {
				if t == target {
					return true
				}
			}
		}
	}
	return false
}

func (n *nfa) addEdge(source int, target int, symbol rune) {
	if v, ok := n.transitions[source]; ok {
		if !n.containsEdge(source, target, symbol) {
			if vl, ok := v[symbol]; ok {
				v[symbol] = append(vl, target)
			} else {
				v[symbol] = []int{target}
			}
		}
	} else {
		n.transitions[source] = map[rune][]int{
			symbol: {target},
		}
	}
}

// TODO need to do DFS to find multiple epsilon transitions
func (n nfa) epsilonExpansion(source int) []int {
	result := []int{source}
	if em, ok := n.transitions[source]; ok {
		if n, ok := em['.']; ok {
			result = append(result, n...)
		}
	}
	return result
}

func (n nfa) convert() dfa {
	result := NewDfa()

	stateNames := map[string]int{}
	stateCount := 0

	// determine the new state names
	hasher := func(states []int) int {
		sort.Slice(states, func(i, j int) bool { return i < j })
		h := ""
		for _, s := range states {
			h += fmt.Sprint(rune(s))
		}
		stateNames[h] = stateCount
		stateCount++
		return stateCount - 1
	}

	return result
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	var N int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &N)

	scanner.Scan()
	scanner.Text()

	n := NewNfa()

	scanner.Scan()
	fmt.Sscan(scanner.Text(), &n.start)

	var T int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &T)

	for i := 0; i < T; i++ {
		var source, target int
		var symbol string
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &source, &target, &symbol)

		n.setAccepting(source, false)
		n.setAccepting(target, false)

		n.addEdge(source, target, rune(symbol[0]))
	}

	scanner.Scan()
	scanner.Text()

	scanner.Scan()
	aaccepting := strings.Split(scanner.Text(), " ")
	for _, acc := range aaccepting {
		var iacc int
		fmt.Sscan(acc, &iacc)
		n.setAccepting(iacc, true)
	}

	d := n.convert()

	var C int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &C)

	for i := 0; i < C; i++ {
		scanner.Scan()
		testcase := scanner.Text()
		accepted := d.validate(testcase)
		if accepted {
			fmt.Println("accept")
		} else {
			fmt.Println("reject")
		}
	}
}
