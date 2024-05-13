package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultipleIslands(t *testing.T) {
	data, err := os.ReadFile("./multiple-islands.txt")
	assert.NoError(t, err, "unexpected error while reading test file")
	reader := bytes.NewReader(data)
	writer := bytes.NewBuffer(nil)
	mainIo(reader, writer)
	assert.Equal(t, "3 12\n", writer.String(), "expected output to be equal")
}
