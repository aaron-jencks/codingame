package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type coin struct {
	value int
	count int
}

func (c coin) total() int {
	return c.value * c.count
}

type coins []coin

func (c coins) total() int {
	total := 0
	for _, cc := range c {
		total += cc.total()
	}
	return total
}

func (c coins) count() int {
	total := 0
	for _, cc := range c {
		total += cc.count
	}
	return total
}

func mainio(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 1000000), 1000000)
	var inputs []string

	var valueToReach int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &valueToReach)

	var N int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &N)

	coins := make(coins, N)

	scanner.Scan()
	inputs = strings.Split(scanner.Text(), " ")
	for i := 0; i < N; i++ {
		fmt.Sscan(inputs[i], &coins[i].count)
	}

	scanner.Scan()
	inputs = strings.Split(scanner.Text(), " ")
	for i := 0; i < N; i++ {
		fmt.Sscan(inputs[i], &coins[i].value)
	}

	sort.Slice(coins, func(i, j int) bool {
		return coins[i].value < coins[j].value
	})

	if coins.total() < valueToReach {
		fmt.Fprintln(out, "-1")
	} else if coins.total() == valueToReach {
		fmt.Fprintln(out, coins.count())
	} else {
		total := 0
		count := 0
		for _, c := range coins {
			if total+c.total() < valueToReach {
				total += c.total()
				count += c.count
				continue
			} else if total+c.total() == valueToReach {
				total += c.total()
				count += c.count
				break
			}
			diff := valueToReach - total
			rem := diff % c.value
			ccount := diff / c.value
			if rem != 0 {
				ccount++
			}
			count += ccount
			break
		}
		fmt.Fprintln(out, count)
	}
}

func main() {
	mainio(os.Stdin, os.Stdout)
}
