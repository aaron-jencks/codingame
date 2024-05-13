package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	data, err := os.ReadFile("./simple.txt")
	assert.NoError(t, err, "unexpected error while reading test file")
	reader := bytes.NewReader(data)
	writer := bytes.NewBuffer(nil)
	main_reader(reader, writer)
	assert.Equal(t, "4\n", string(writer.Bytes()), "expected output to be equal")
}
