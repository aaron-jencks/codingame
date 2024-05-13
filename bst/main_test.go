package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExample(t *testing.T) {
	gout := bytes.NewBuffer(nil)
	main_io(bytes.NewReader([]byte("5\n8 6 13 10 5\n")), gout)
	assert.Equal(t, "8 6 5 13 10\n5 6 8 10 13\n5 6 10 13 8\n8 6 13 5 10\n", gout.String(), "expected outputs to match")
}
