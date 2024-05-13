package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func simulateWordle(t *testing.T, answer string, words []string) int {
	sim := WordleSim{}
	sim.initializeSim(words)
	var states = make([]State, 6)
	gcount := 0
	for {
		g := sim.guess(states)

		if g == answer {
			return gcount + 1
		}
		gcount++
		assert.Equal(t, 6, len(g), "expected length of guess to be 6")

		if gcount > 27 {
			assert.Fail(t, "too many guesses",
				"expected correct answer to be guess within 27 tries for word '%s'", answer)
			return gcount
		}

		for li := 0; li < 6; li++ {
			states[li] = INCORRECT
			for _, r := range answer {
				if rune(g[li]) == r {
					states[li] = PARTIAL
					break
				}
			}
			if g[li] == answer[li] {
				states[li] = CORRECT
			}
		}
	}
}

func readWords(t *testing.T) []string {
	data, err := os.ReadFile("./words.txt")
	assert.NoError(t, err, "unexpected error while reading test words")
	words := strings.Split(string(data), "\n")
	return words
}

func TestFull(t *testing.T) {
	words := readWords(t)
	fmt.Printf("testing %d words\n", len(words))
	for _, w := range words {
		t.Run(w, func(tt *testing.T) {
			simulateWordle(tt, w, words)
		})
	}
}

func TestArcade(t *testing.T) {
	words := readWords(t)
	simulateWordle(t, "ARCADE", words)
}

func TestBohunk(t *testing.T) {
	words := readWords(t)
	simulateWordle(t, "BOHUNK", words)
}

func TestHelper(t *testing.T) {
	words := readWords(t)
	simulateWordle(t, "HELPER", words)
}
