package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTest6(t *testing.T) {
	result := FindCount("ball", 4)
	assert.Equal(t, 199, result, "counts should equal")
}
