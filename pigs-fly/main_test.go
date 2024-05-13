package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTraits(t *testing.T) {
	tokens := [][]string{
		{"PIGS", "have", "FEET"},
		{"FEET", "are", "LIMBS"},
	}

	var traits set = set{}
	for _, line := range tokens {
		findTraits(line, traits)
	}

	assert.True(t, traits.contains("FEET"), "traits should contain feet")
	assert.True(t, traits.contains("LIMBS"), "traits should contain limbs")
}
