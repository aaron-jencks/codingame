package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//  abcdefghijklmnopqrs
// a
// b
// c
// d
// e
// f
// g
// h
// i
// j
// k
// l
// m
// n
// o
// p
// q
// r
// s
// t
// u
// v
// w
// x
// y

func TestParsing(t *testing.T) {
	tcs := []struct {
		in         string
		r, c, d, m int
	}{
		// {
		// 	"ls11",
		// 	18, 11, 11, MOW,
		// },
		{
			"PLANTjm14",
			12, 9, 14, PLANT,
		},
	}

	for _, tc := range tcs {
		r, c, d, m := ParseInstruction(tc.in)
		assert.Equal(t, tc.r, r, "row mismatch for %s", tc.in)
		assert.Equal(t, tc.c, c, "column mismatch for %s", tc.in)
		assert.Equal(t, tc.d, d, "diameter mismatch for %s", tc.in)
		assert.Equal(t, tc.m, m, "mode mismatch for %s", tc.in)
	}
}

func TestLastCase(t *testing.T) {
	field := make([][]bool, 25)
	for i := 0; i < 25; i++ {
		field[i] = make([]bool, 19)
	}

	expected := make([][]bool, 25)
	for i := 0; i < 25; i++ {
		expected[i] = make([]bool, 19)
	}

	r, c, d, m := ParseInstruction("jm13")
	GenerateCircle(r, c, d, m, field)
	r, c, d, m = ParseInstruction("PLANTjm14")
	GenerateCircle(r, c, d, m, field)

	assert.Equal(t, expected, field, "expected fields to be equal")
}
