package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSample(t *testing.T) {
	testCases := []struct {
		name          string
		numbers       []string
		expectedCount int
	}{
		{
			name: "single",
			numbers: []string{
				"0467123456",
			},
			expectedCount: 10,
		},
	}

	for ti, tt := range testCases {
		t.Run(fmt.Sprintf("%d: %s", ti, tt.name), func(t *testing.T) {
			tree := &n{map[byte]n{}}
			for _, number := range tt.numbers {
				tree.i(number)
			}
			assert.Equal(t, tt.expectedCount, tree.l(), "Tree size should match")
		})
	}
}
