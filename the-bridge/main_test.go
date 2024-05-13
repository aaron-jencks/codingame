package main

import (
	"fmt"
	"os"
	"testing"
)

type bufferWriter struct {
	buff []byte
}

func (bw *bufferWriter) Write(b []byte) (int, error) {
	bw.buff = append(bw.buff, b...)
	return len(b), nil
}

func TestOneHole(t *testing.T) {
	tfp, err := os.OpenFile("./one-hole.txt", os.O_RDONLY, 0777)
	if err != nil {
		fmt.Printf("failed to open file: %s\n", err)
		return
	}
	out := bufferWriter{}
	runMain(tfp, &out)
}
