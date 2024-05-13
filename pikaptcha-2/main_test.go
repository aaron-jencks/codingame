package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSparseValidator(t *testing.T) {
	data, err := os.ReadFile("./sparse.validator.txt")
	assert.NoError(t, err, "unexpected error while reading input file")
	dataout, err := os.ReadFile("./sparse.validator.out.txt")
	assert.NoError(t, err, "unexpected error while reading output file")
	gout := bytes.NewBuffer(nil)
	main_io(bytes.NewReader(data), gout)
	assert.Equal(t, string(dataout), gout.String(), "expected outputs to match")
}

func TestCrossesValidator(t *testing.T) {
	data, err := os.ReadFile("./crosses.validator.txt")
	assert.NoError(t, err, "unexpected error while reading input file")
	gout := bytes.NewBuffer(nil)
	main_io(bytes.NewReader(data), gout)
}
