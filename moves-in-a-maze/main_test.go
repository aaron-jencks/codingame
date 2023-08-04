package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateMapFromString(s string) []string {
	lines := strings.Split(s, "\n")
	return lines
}

func TestEasy(t *testing.T) {
	ms := "##########\n#........#\n##.#####.#\n##.#.....#\n##########"
	es := "##########\n#01234567#\n##2#####8#\n##3#DCBA9#\n##########"
	h := 5
	w := 10
	m := generateMapFromString(ms)
	e := generateMapFromString(es)
	result := GenerateOutputFromArray(h, w, FindDistances(1, 1, h, w, m))
	assert.Equal(t, e, result, "output map mismatch")
}

func TestUnreachable(t *testing.T) {
	ms := "....#...#..........\n....#####..........\n....#...#..........\n....#.#.#..........\n......#............\n...S#####..........\n....#...#.........."
	es := "##########\n#01234567#\n##2#####8#\n##3#DCBA9#\n##########" // TODO fix
	h := 7
	w := 19
	m := generateMapFromString(ms)
	e := generateMapFromString(es)
	result := GenerateOutputFromArray(h, w, FindDistances(1, 1, h, w, m))
	assert.Equal(t, e, result, "output map mismatch")
}
