package main

import (
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
