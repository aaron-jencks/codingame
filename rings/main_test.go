package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBumps(t *testing.T) {
	tape := make(runes, 1)
	tcs := []struct {
		name     string
		initial  rune
		ups      int
		downs    int
		expected rune
	}{
		{
			name:     "wrap forward",
			initial:  ' ',
			ups:      28,
			expected: 'A',
		},
		{
			name:     "wrap backward",
			initial:  ' ',
			downs:    28,
			expected: 'Z',
		},
		{
			name:     "wrap backward one",
			initial:  ' ',
			downs:    1,
			expected: 'Z',
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			tape[0] = tc.initial
			for i := 0; i < tc.ups; i++ {
				tape.bumpUp(0)
			}
			for i := 0; i < tc.downs; i++ {
				tape.bumpDown(0)
			}
			assert.Equal(tt, tc.expected, tape[0],
				"expected final rune to be equal, expected %c, got %c",
				tc.expected, tape[0],
			)
		})
	}
}

func TestMovement(t *testing.T) {
	var b blub

	tcs := []struct {
		name     string
		initial  blub
		lefts    int
		rights   int
		expected blub
	}{
		{
			name:     "wrap forward",
			initial:  0,
			rights:   FOREST_WIDTH,
			expected: 0,
		},
		{
			name:     "wrap backward",
			initial:  0,
			lefts:    FOREST_WIDTH,
			expected: 0,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			b = tc.initial
			for i := 0; i < tc.lefts; i++ {
				b.left()
			}
			for i := 0; i < tc.rights; i++ {
				b.right()
			}
			assert.Equal(tt, tc.expected, b,
				"expected final rune to be equal",
			)
		})
	}
}

func TestForwardDistance(t *testing.T) {
	tcs := []struct {
		name     string
		initial  rune
		target   rune
		distance int
	}{
		{
			name:     "wrap around",
			initial:  'C',
			target:   'A',
			distance: 25,
		},
		{
			name:     "same",
			initial:  'A',
			target:   'A',
			distance: 0,
		},
		{
			name:     "no wrap",
			initial:  'A',
			target:   'Z',
			distance: 25,
		},
		{
			name:     "no wrap target space",
			initial:  'A',
			target:   ' ',
			distance: 26,
		},
		{
			name:     "no wrap origin space",
			initial:  ' ',
			target:   'Z',
			distance: 26,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			tape := make(runes, 1)
			tape[0] = tc.initial
			result := tape.letterForwardDistance(0, tc.target)
			assert.Equal(tt, tc.distance, result, "expected distance to be correct")
		})
	}
}

func TestReverseDistance(t *testing.T) {
	tcs := []struct {
		name     string
		initial  rune
		target   rune
		distance int
	}{
		{
			name:     "wrap around",
			initial:  'A',
			target:   'C',
			distance: 25,
		},
		{
			name:     "wrap around",
			initial:  'A',
			target:   'Z',
			distance: 2,
		},
		{
			name:     "same",
			initial:  'A',
			target:   'A',
			distance: 0,
		},
		{
			name:     "no wrap",
			initial:  'Z',
			target:   'A',
			distance: 25,
		},
		{
			name:     "no wrap target space",
			initial:  'Z',
			target:   ' ',
			distance: 26,
		},
		{
			name:     "no wrap origin space",
			initial:  ' ',
			target:   'A',
			distance: 26,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			tape := make(runes, 1)
			tape[0] = tc.initial
			result := tape.letterReverseDistance(0, tc.target)
			assert.Equal(tt, tc.distance, result, "expected distance to be correct")
		})
	}
}

func TestBinForwardDistance(t *testing.T) {
	tcs := []struct {
		name     string
		initial  int
		target   int
		distance int
	}{
		{
			name:     "wrap around",
			initial:  3,
			target:   0,
			distance: 27,
		},
		{
			name:     "same",
			initial:  16,
			target:   16,
			distance: 0,
		},
		{
			name:     "no wrap",
			initial:  0,
			target:   29,
			distance: 29,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			bin := blub(tc.initial)
			result := bin.binForwardDistance(tc.target)
			assert.Equal(tt, tc.distance, result, "expected distance to be correct")
		})
	}
}

func TestBinReverseDistance(t *testing.T) {
	tcs := []struct {
		name     string
		initial  int
		target   int
		distance int
	}{
		{
			name:     "wrap around",
			initial:  0,
			target:   3,
			distance: 27,
		},
		{
			name:     "same",
			initial:  16,
			target:   16,
			distance: 0,
		},
		{
			name:     "no wrap",
			initial:  29,
			target:   0,
			distance: 29,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			bin := blub(tc.initial)
			result := bin.binReverseDistance(tc.target)
			assert.Equal(tt, tc.distance, result, "expected distance to be correct")
		})
	}
}
