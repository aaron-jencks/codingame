package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEpsilonExpansion(t *testing.T) {
	n := NewNfa()
	n.transitions = map[int]map[rune][]int{
		0: {
			'.': {
				1, 2, 3, 4,
			},
		},
		1: {
			'.': {
				3, 5,
			},
			'a': {
				6,
			},
		},
		4: {
			'.': {
				7, 8,
			},
		},
	}

	neighbors := n.epsilonExpansion([]int{0})

	assert.Equal(t, []int{0, 1, 2, 3, 4, 7, 8, 5}, neighbors, "neighbor array mismatch")
}

func TestHasher(t *testing.T) {
	sn := map[string]int{}
	sc := 0
	assert.Equal(t, 0, NodeHasher([]int{0, 3, 2, 1}, sn, &sc), "should create new node id using state count")
	h, ok := sn["0123"]
	assert.True(t, ok, "node hasher mismatch %v", sn)
	assert.Equal(t, 0, h, "node assignment mismatch")
	assert.Equal(t, 0, NodeHasher([]int{0, 3, 2, 1}, sn, &sc), "expected to reuse old value")
}

// see https://youtu.be/jMxuL4Xzi_A?si=BnrXf9SGN_gNYM5A
func TestConversion1(t *testing.T) {
	n := NewNfa()
	n.alphabet = []rune{'a', 'b'}
	n.addEdge(0, 1, '.')
	n.addEdge(1, 1, 'a')
	n.addEdge(1, 2, 'a')
	n.addEdge(1, 2, 'b')
	n.addEdge(2, 0, 'a')
	n.addEdge(2, 2, 'a')
	n.addEdge(2, 3, 'b')
	n.addEdge(3, 1, 'b')
	n.setAccepting(0, true)

	assert.True(t, n.containsEdge(0, 1, '.'), "nfa should contain generated edges: 0->1 '.'")
	assert.True(t, n.containsEdge(1, 1, 'a'), "nfa should contain generated edges: 1->1 'a'")
	assert.True(t, n.containsEdge(1, 2, 'a'), "nfa should contain generated edges: 1->2 'a'")
	assert.True(t, n.containsEdge(1, 2, 'b'), "nfa should contain generated edges: 1->2 'b'")
	assert.True(t, n.containsEdge(2, 0, 'a'), "nfa should contain generated edges: 2->0 'a'")
	assert.True(t, n.containsEdge(2, 2, 'a'), "nfa should contain generated edges: 2->2 'a'")
	assert.True(t, n.containsEdge(2, 3, 'b'), "nfa should contain generated edges: 2->3 'b'")
	assert.True(t, n.containsEdge(3, 1, 'b'), "nfa should contain generated edges: 3->1 'b'")

	nodes := n.getNodes()
	sort.Slice(nodes, func(i, j int) bool { return nodes[i] < nodes[j] })
	assert.Equal(t, []int{0, 1, 2, 3}, nodes, "nfa should contain all nodes generated")

	dfa := n.convert()

	tcs := map[string]bool{
		"":                         true,
		"ba":                       true,
		"aa":                       true,
		"abbaaaaaaa":               true,
		"abbba":                    true,
		"bbaababababbbbaaaabbabab": false,
		"bbbbaaaaaa":               true,
		"abaaaa":                   true,
		"bb":                       false,
		"bbb":                      false,
		"abb":                      false,
		"ab":                       false,
		"a":                        false,
		"b":                        false,
		"bbaaaaa":                  false,
		"bababaaaa":                true,
		"bbbababaaa":               true,
	}

	for k, v := range tcs {
		assert.Equal(t, v, dfa.validate(k), "invalid validation for string %s", k)
	}
}

func TestTestCases(t *testing.T) {
	tcs := []string{
		"even_numbers",
		"odd_numbers",
		"wide_example",
		"wide_example_2",
	}

	for _, tc := range tcs {
		in, err := os.ReadFile(fmt.Sprintf("%s.in.txt", tc))
		assert.NoError(t, err, "unexpected error reading test input")

		out, err := os.ReadFile(fmt.Sprintf("%s.out.txt", tc))
		assert.NoError(t, err, "unexpected error reading test output")
		outlines := strings.Split(string(out), "\n")
		for oi := range outlines {
			outlines[oi] = strings.Trim(outlines[oi], "\r")
		}

		t.Run(tc, func(tt *testing.T) {
			assert.Equal(tt, outlines, ParseInput(strings.NewReader(string(in))), "output mismatch")
		})
	}
}
