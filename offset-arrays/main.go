package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func debugf(f string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, f, args...)
}

type ArraySpec struct {
	name     string
	start    int
	stop     int
	elements []int
}

func (as ArraySpec) index(i int) int {
	// diff := as.stop - as.start
	return as.elements[(i - as.start)]
}

var assign_re *regexp.Regexp = regexp.MustCompile(`(?P<identifier>[A-Z]+)\[(?P<start>-?\d+)\.+(?P<stop>-?\d+)\]\s*=(?P<contents>(\s-?\d+)+)`)
var ident_re *regexp.Regexp = regexp.MustCompile(`[A-Z]+`)
var index_re *regexp.Regexp = regexp.MustCompile(`-?\d+`)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	var n int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &n)

	var arrs map[string]ArraySpec = map[string]ArraySpec{}

	for i := 0; i < n; i++ {
		scanner.Scan()
		assignment := scanner.Text()
		match := assign_re.FindStringSubmatch(assignment)

		var arr_spec ArraySpec = ArraySpec{}

		for m, name := range assign_re.SubexpNames() {
			switch name {
			case "identifier":
				arr_spec.name = match[m]
			case "start":
				fmt.Sscanf(match[m], "%d", &arr_spec.start)
			case "stop":
				fmt.Sscanf(match[m], "%d", &arr_spec.stop)
			case "contents":
				bits := strings.Split(match[m][1:], " ")
				debugf("Found elements for array: %v\n", bits)
				arr_spec.elements = make([]int, len(bits))
				for n := range bits {
					fmt.Sscanf(bits[n], "%d", &arr_spec.elements[n])
				}
			}
		}

		debugf("array def: %v\n", arr_spec)

		arrs[arr_spec.name] = arr_spec
	}
	scanner.Scan()
	x := scanner.Text()
	debugf("Parsing index: %s\n", x)

	identifiers := ident_re.FindAllString(x, -1)

	var num int
	snum := index_re.FindString(x)
	fmt.Sscanf(snum, "%d", &num)

	for fi := range identifiers {
		i := len(identifiers) - (fi + 1)
		debugf("Parsing index for %s[%d]\n", identifiers[i], num)
		num = arrs[identifiers[i]].index(num)
	}

	// fmt.Fprintln(os.Stderr, "Debug messages...")
	fmt.Println(num) // Write answer to stdout
}
