package main

import (
	"bufio"
	"fmt"
	"io"
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

func (n nfa) epsilonExpansion(sources []int) []int {
	result := make([]int, len(sources))
	visited := map[int]bool{}
	stack := make([]int, len(sources))

	for si, s := range sources {
		result[si] = s
		stack[si] = s
		visited[s] = true
	}

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

// determine the new state names
func NodeHasher(states []int, stateNames map[string]int, stateCount *int) int {
	sort.Slice(states, func(i, j int) bool { return states[i] < states[j] }) // must be a slice
	h := ""
	for _, s := range states {
		h += fmt.Sprint(rune(s))
	}
	if si, ok := stateNames[h]; ok {
		return si
	}
	stateNames[h] = *stateCount
	*stateCount++
	return *stateCount - 1
}

func (n nfa) convert() dfa {
	result := NewDfa()

	visited := map[int]bool{n.start: true}
	stateNames := map[string]int{}
	stateCount := 0

	for _, node := range n.getNodes() {
		stateNames[fmt.Sprint(node)] = node
		if node >= stateCount {
			stateCount = node + 1
		}
	}

	// setup trap state
	stateNames[fmt.Sprint(stateCount)] = stateCount

	trapId := stateCount

	for _, l := range n.alphabet {
		result.addEdge(stateCount, stateCount, l)
	}

	stateCount++

	new_start := n.createExpandedNode([]int{n.start}, stateNames, &stateCount)

	result.start = new_start.nodeId
	if new_start.accepting {
		result.setAccepting(new_start.nodeId, true)
	}

	stack := []expandedNode{new_start}
	for len(stack) > 0 {
		last := len(stack) - 1
		element := stack[last]
		stack = stack[:last]

		for _, letter := range n.alphabet {
			connected_nodes := map[int]bool{}

			for _, cn := range element.nodes {
				if trans, ok := n.transitions[cn]; ok {
					if neighbors, ok := trans[letter]; ok {
						for _, neigh := range neighbors {
							connected_nodes[neigh] = true
						}
					}
				}
			}

			if len(connected_nodes) == 0 {
				connected_nodes[trapId] = true
			}

			new_node := make([]int, 0, len(connected_nodes))
			for k := range connected_nodes {
				new_node = append(new_node, k)
			}

			enn := n.createExpandedNode(new_node, stateNames, &stateCount)

			result.addEdge(element.nodeId, enn.nodeId, letter)

			if _, ok := visited[enn.nodeId]; !ok {
				visited[enn.nodeId] = true

				if enn.accepting {
					result.accepting[enn.nodeId] = true
				}

				stack = append(stack, enn)
			}
		}
	}

	return result
}

type expandedNode struct {
	nodes     []int
	accepting bool
	nodeId    int
}

func (n nfa) createExpandedNode(nodes []int, stateNames map[string]int, stateCount *int) expandedNode {
	result := expandedNode{}

	result.nodes = n.epsilonExpansion(nodes)
	result.nodeId = NodeHasher(result.nodes, stateNames, stateCount)

	for _, node := range result.nodes {
		if acc, ok := n.accepting[node]; ok && acc {
			result.accepting = true
			break
		}
	}

	return result
}

func ParseInput(r io.Reader) []string {
	scanner := bufio.NewScanner(r)
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

	var result []string
	for i := 0; i < C; i++ {
		scanner.Scan()
		testcase := scanner.Text()
		accepted := d.validate(testcase)
		if accepted {
			result = append(result, "accept")
		} else {
			result = append(result, "reject")
		}
	}
	return result
}

func main() {
	for _, o := range ParseInput(os.Stdin) {
		fmt.Println(o)
	}
}
