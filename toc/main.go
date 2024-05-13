package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type entry struct {
	indent int
	title  string
	page   int
}

const parseRegexp = `(?P<indent>>*)(?P<title>\w+)\s(?P<page>\d+)`

var parseRe = regexp.MustCompile(parseRegexp)

func parseEntry(l string) entry {
	result := entry{}

	matches := parseRe.FindStringSubmatch(l)
	result.indent = len(matches[parseRe.SubexpIndex("indent")])
	result.title = matches[parseRe.SubexpIndex("title")]
	fmt.Sscan(matches[parseRe.SubexpIndex("page")], &result.page)

	// fmt.Fprintln(os.Stderr, result)

	return result
}

type chapterCount struct {
	indent int
	count  int
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	var ll int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &ll)

	var N int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &N)

	entries := make([]entry, N)
	for i := 0; i < N; i++ {
		scanner.Scan()
		el := scanner.Text()
		entries[i] = parseEntry(el)
	}

	chapterStack := make([]chapterCount, 1)
	for _, ent := range entries {
		lastCount := len(chapterStack) - 1
		if ent.indent > chapterStack[lastCount].indent {
			chapterStack = append(chapterStack, chapterCount{
				indent: ent.indent,
			})
			lastCount++
		} else if ent.indent < chapterStack[lastCount].indent {
			for ent.indent < chapterStack[lastCount].indent {
				chapterStack = chapterStack[:lastCount]
				lastCount--
			}
		}

		chapterStack[lastCount].count++

		indent := strings.Repeat("    ", chapterStack[lastCount].indent)
		truncString := fmt.Sprintf("%s%d %s%d",
			indent, chapterStack[lastCount].count, ent.title, ent.page,
		)

		fmt.Printf("%s%d %s%s%d\n",
			indent,
			chapterStack[lastCount].count,
			ent.title,
			strings.Repeat(".", ll-len(truncString)),
			ent.page,
		)
	}
}
