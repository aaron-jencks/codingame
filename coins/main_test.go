package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEasy(t *testing.T) {
	buff := strings.NewReader("8\n1\n5\n6\n")
	buffout := bytes.NewBuffer(nil)
	mainio(buff, buffout)
	assert.Equal(t, "2\n", buffout.String())
}
