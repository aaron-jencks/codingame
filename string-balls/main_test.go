package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestdepth(t *testing.T) {
	result := FindCount("ab", 4)
	assert.Equal(t, 19, result, "counts should equal")
}

func TestTest6(t *testing.T) {
	result := FindCount("ball", 4)
	assert.Equal(t, 199, result, "counts should equal")
}

func TestTest7(t *testing.T) {
	result := FindCount("bye", 100)
	assert.Equal(t, 17576, result, "counts should equal")
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
