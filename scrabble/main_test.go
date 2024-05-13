package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func RunFileTest(t *testing.T, iname, oname string) {
	data, err := os.ReadFile(iname)
	assert.NoError(t, err, "unexpected error while reading input file")
	dataout, err := os.ReadFile(oname)
	assert.NoError(t, err, "unexpected error while reading output file")
	gout := bytes.NewBuffer(nil)
	main_io(bytes.NewReader(data), gout)
	assert.Equal(t, string(dataout), gout.String(), "expected outputs to match")
}

func TestEasy(t *testing.T) {
	RunFileTest(t, "./easy.txt", "./easy.out.txt")
}

func TestDouble(t *testing.T) {
	RunFileTest(t, "./double.txt", "./double.out.txt")
}
