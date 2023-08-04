package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	START = -1
	END   = -2
)

func GeneratePicture(data [][]int) []string {
	result := []string{
		"+---[CODINGAME]---+",
	}

	for _, row := range data {
		line := "|"
		for _, c := range row {
			switch c % 15 { // apply wrap-around logic
			case 0:
				line += " "
			case 1:
				line += "."
			case 2:
				line += "o"
			case 3:
				line += "+"
			case 4:
				line += "="
			case 5:
				line += "*"
			case 6:
				line += "B"
			case 7:
				line += "O"
			case 8:
				line += "X"
			case 9:
				line += "@"
			case 10:
				line += "%"
			case 11:
				line += "&"
			case 12:
				line += "#"
			case 13:
				line += "/"
			case 14:
				line += "^"
			case START:
				line += "S"
			case END:
				line += "E"
			}
		}
		result = append(result, line+"|")
	}

	result = append(result, "+-----------------+")

	return result
}

func ParseInts(fingerprint string) []uint8 {
	codes := strings.Split(fingerprint, ":")

	result := make([]uint8, len(codes))
	for i, code := range codes {
		fmt.Sscanf(code, "%x", &result[i])
	}

	return result
}

func PerformWalk(codes []uint8) [][]int {
	m := make([][]int, 9)
	for i := 0; i < 9; i++ {
		m[i] = make([]int, 17)
	}

	x := 8
	y := 4

	for _, code := range codes {
		for c := 0; c < 4; c++ {
			// look at the 2 bits for the direction
			dir := code & 0x03

			switch dir {
			case 0:
				// up-left
				if y > 0 {
					y--
				}
				if x > 0 {
					x--
				}
			case 1:
				// up-right
				if y > 0 {
					y--
				}
				if x < 16 {
					x++
				}
			case 2:
				// down-left
				if y < 8 {
					y++
				}
				if x > 0 {
					x--
				}
			case 3:
				// down-right
				if y < 8 {
					y++
				}
				if x < 16 {
					x++
				}
			}

			m[y][x]++ // place the coin

			// move the next 2 bits into position
			code >>= 2
		}
	}

	m[4][8] = START
	m[y][x] = END

	return m
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	scanner.Scan()
	fingerprint := scanner.Text()
	codes := ParseInts(fingerprint)
	counts := PerformWalk(codes)
	lines := GeneratePicture(counts)

	// fmt.Fprintln(os.Stderr, "Debug messages...")
	fmt.Println(strings.Join(lines, "\n")) // Write answer to stdout
}
