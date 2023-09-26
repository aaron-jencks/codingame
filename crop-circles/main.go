package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

const (
	MOW = iota
	PLANT
	PLANTMOW
)

func ParseInstruction(ins string) (int, int, int, int) {
	mode := MOW
	if ins[0] == 'P' {
		mode = PLANT
		if len(ins) > 8 && ins[5] == 'M' {
			mode = PLANTMOW
		}
	}

	ilen := len(ins)
	var rin, cin int = ilen - 2, ilen - 3
	diam := 0

	if ins[rin] >= '0' && ins[rin] <= '9' {
		rin--
		cin--
		diam = 10*int(ins[ilen-2]-'0') + int(ins[ilen-1]-'0')
	} else {
		diam = int(ins[ilen-1] - '0')
	}

	col := int(ins[cin] - 'a')
	row := int(ins[rin] - 'a')

	return row, col, diam, mode
}

func GenerateCircle(row, col, diam, mode int, field [][]bool) {
	radius := diam >> 1
	for r := -radius; r < radius+1; r++ {
		for c := -radius; c < radius+1; c++ {
			if int(math.Round(math.Sqrt(float64(r*r+c*c)))) <= radius {
				rr := r + row
				rc := c + col
				if rr >= 0 && rr < 25 && rc >= 0 && rc < 19 {
					switch mode {
					case MOW:
						field[rr][rc] = true
					case PLANT:
						field[rr][rc] = false
					case PLANTMOW:
						field[rr][rc] = !field[rr][rc]
					}
				}
			}
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	field := make([][]bool, 25)
	for i := 0; i < 25; i++ {
		field[i] = make([]bool, 19)
	}

	scanner.Scan()
	instructions := scanner.Text()
	fmt.Fprintln(os.Stderr, instructions)
	circles := strings.Split(instructions, " ")

	for _, circle := range circles {
		r, c, d, m := ParseInstruction(circle)
		GenerateCircle(r, c, d, m, field)
	}

	for ri := range field {
		for _, v := range field[ri] {
			if !v {
				fmt.Print("{}")
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Print("\n")
	}
}
