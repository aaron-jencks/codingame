package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readFile(t *testing.T, f string) []byte {
	data, err := os.ReadFile(f)
	assert.NoError(t, err, "unexpected error while reading input file")
	return data
}

func TestT1(t *testing.T) {
	data := readFile(t, "./t1.txt")
	gout := bytes.NewBuffer(nil)
	main_io(bytes.NewReader(data), gout)
	assert.Equal(t, "2\n", gout.String(), "expected outputs to match")
}
