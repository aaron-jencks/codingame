package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	var W int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &W)

	var H int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &H)

	scanner.Scan()
	T := scanner.Text()
	segs := strings.Split(T, " ")

	// Generate the unbroken data stream
	var content string = ""
	for si, seg := range segs {
		var count int
		fmt.Sscan(seg, &count)

		var c string = "*"
		if si%2 == 1 {
			c = " "
		}
		content += strings.Repeat(c, count)
	}

	// Insert the page borders
	var result string = ""
	current := 0
	for current < len(content) && len(content) > W {
		nc := current + W
		result += "|" + content[current:nc] + "|\n"
		current = nc
	}
	if current < len(content) {
		result += "|" + content[current:] + "|\n"
	}

	// fmt.Fprintln(os.Stderr, "Debug messages...")
	fmt.Println(result[:len(result)-1])
}
