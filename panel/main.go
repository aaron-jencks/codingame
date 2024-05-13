package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	var P int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &P)

	properties := make([]string, P)

	for i := 0; i < P; i++ {
		scanner.Scan()
		properties[i] = scanner.Text()
	}

	var N int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &N)

	people := make([]map[string]string, N)

	for i := 0; i < N; i++ {
		scanner.Scan()

		person := scanner.Text()
		pprops := strings.Split(person, " ")

		people[i] = map[string]string{}
		for j := 0; j < P; j++ {
			people[i][properties[j]] = pprops[j+1]
		}
	}

	var F int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &F)
	for i := 0; i < F; i++ {
		scanner.Scan()

		formula := scanner.Text()
		fprops := strings.Split(formula, " AND ")
		fmap := map[string]string{}
		for _, fprop := range fprops {
			chunks := strings.Split(fprop, "=")
			fmap[chunks[0]] = chunks[1]
		}

		count := 0

	PERSONLOOP:
		for _, person := range people {
			for prop, val := range fmap {
				if person[prop] != val {
					continue PERSONLOOP
				}
			}
			count++
		}

		fmt.Println(count)
	}
}
