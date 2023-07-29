package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLetterNeighbors(t *testing.T) {
	for _, l := range "abcdefghijklmnopqrstuvwxyz" {
		neighbors := FindLetterNeighbor(l, false)
		if l > 'a' {
			assert.Contains(t, neighbors, l-1, "missing left neighbor for letter %c", l)
		}
		if l < 'z' {
			assert.Contains(t, neighbors, l+1, "missing right neighbor for letter %c", l)
		}
	}
	ns := FindLetterNeighbor('a', true)
	assert.Contains(t, ns, 'z', "wrapping should include wrapped neighbors")
	ns = FindLetterNeighbor('z', true)
	assert.Contains(t, ns, 'z', "wrapping should include wrapped neighbors")
}

func TestTestdepth(t *testing.T) {
	result := FindCount("ab", 4)
	assert.Equal(t, 19, result, "counts should equal")
}

func TestTest6(t *testing.T) {
	result := FindCount("ball", 4)
	assert.Equal(t, 199, result, "counts should equal")
}

func TestTest8(t *testing.T) {
	result := FindCount("greece", 11)
	assert.Equal(t, 181477, result, "counts should equal")
}

func TestTest10(t *testing.T) {
	FindCount("portugal", 10)
	// assert.Equal(t, 181477, result, "counts should equal")
}

func TestTest11(t *testing.T) {
	FindCount("codingame", 9)
	// assert.Equal(t, 181477, result, "counts should equal")
}
