package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

const EPSILON = '.'

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
	alphabet    []rune
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

func (n nfa) getNodes() []int {
	result := []int{n.start}
	visited := map[int]bool{n.start: true}
	for k, v := range n.transitions {
		if _, ok := visited[k]; !ok {
			result = append(result, k)
			visited[k] = true
		}
		for _, neighbors := range v {
			for _, neigh := range neighbors {
				if _, ok := visited[neigh]; !ok {
					result = append(result, neigh)
					visited[neigh] = true
				}
			}
		}
	}
	return result
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

// TODO convert to use an array so that
// we can find all of the neighbors of several nodes at the same time
// since during exploration we'll find nodes that are mixed with other nodes
func (n nfa) epsilonExpansion(source int) []int {
	result := []int{source}
	visited := map[int]bool{source: true}
	stack := []int{source}
	for len(stack) > 0 {
		last := len(stack) - 1
		element := stack[last]
		stack = stack[:last]

		if em, ok := n.transitions[element]; ok {
			if n, ok := em[EPSILON]; ok {
				for _, neigh := range n {
					if _, ok := visited[neigh]; !ok {
						visited[neigh] = true
						result = append(result, neigh)
						stack = append(stack, neigh)
					}
				}
			}
		}
	}
	return result
}

func (n nfa) convert() dfa {
	result := NewDfa()

	stateNames := map[string]int{}
	stateCount := 0

	for _, node := range n.getNodes() {
		stateNames[fmt.Sprint(node)] = node
	}

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

	stack := [][]int{{n.start}}
	for len(stack) > 0 {
		last := len(stack) - 1
		element := stack[last]
		stack = stack[:last]

		node := hasher(element)

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
	alpha := strings.Split(scanner.Text(), " ")

	n := NewNfa()
	for _, a := range alpha {
		n.alphabet = append(n.alphabet, rune(a[0]))
	}

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
