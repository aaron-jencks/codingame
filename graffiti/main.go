package main

import (
	"fmt"
	"sort"
)

type interval struct {
	start int
	end   int
}

func main() {
	var L int
	fmt.Scan(&L)

	var N int
	fmt.Scan(&N)

	intervals := make([]interval, N)

	for i := 0; i < N; i++ {
		fmt.Scan(&intervals[i].start, &intervals[i].end)
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].start < intervals[j].start
	})

	var filtered []interval
	// get rid of overlapped intervals
	current := intervals[0]
	for i := 1; i < N; i++ {
		inter := intervals[i]
		if inter.end <= current.end {
			// interval is completely overlapped
			// eliminate it
			continue
		}
		if inter.start <= current.end {
			// interval is an extension of the previous
			current.end = inter.end
		} else {
			// interval is disjoint from the previous
			// add to the filtered
			filtered = append(filtered, current)
			current = inter
		}
	}

	// append the last interval
	filtered = append(filtered, current)

	if filtered[0].start == 0 && filtered[0].end == L {
		fmt.Println("All painted")
		return
	}

	if filtered[0].start > 0 {
		fmt.Printf("0 %d\n", filtered[0].start)
	}

	// iterate over the list of now disjoint intervals and find the gaps
	previous := filtered[0]
	for i := 1; i < len(filtered); i++ {
		fmt.Printf("%d %d\n", previous.end, filtered[i].start)
		previous = filtered[i]
	}

	if filtered[len(filtered)-1].end < L {
		fmt.Printf("%d %d\n", filtered[len(filtered)-1].end, L)
	}
}
