package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

const (
	FOREST_WIDTH = 30
)

type runes []rune

func (r runes) bumpUp(i int) {
	if r[i] == 'Z' {
		r[i] = ' '
		return
	} else if r[i] == ' ' {
		r[i] = 'A'
		return
	}
	r[i]++
}

func (r runes) bumpDown(i int) {
	if r[i] == 'A' {
		r[i] = ' '
		return
	} else if r[i] == ' ' {
		r[i] = 'Z'
		return
	}
	r[i]--
}

func (r runes) letterForwardDistance(i int, target rune) int {
	if target == r[i] {
		return 0
	}

	var fwd rune = 0
	if r[i] == ' ' {
		fwd = target - 'A' + 1
	} else if target == ' ' {
		fwd = 26 - (r[i] - 'A')
	} else if target > r[i] {
		fwd = target - r[i]
	} else {
		fwd = (26 - (r[i] - 'A')) + (target - 'A' + 1) // account for the space
	}

	return int(fwd)
}

func (r runes) letterReverseDistance(i int, target rune) int {
	if target == r[i] {
		return 0
	}

	var rvs rune = 0
	if target == ' ' {
		rvs = r[i] - 'A' + 1
	} else if r[i] == ' ' {
		rvs = 26 - (target - 'A')
	} else if target < r[i] {
		rvs = r[i] - target
	} else {
		rvs = (25 - (target - 'A')) + (r[i] - 'A' + 1) + 1 // account for the space
	}

	return int(rvs)
}

func (r runes) letterDistance(i int, target rune) int {
	fwd := r.letterForwardDistance(i, target)

	if fwd == 0 {
		return fwd
	}

	rvs := r.letterReverseDistance(i, target)

	if fwd < rvs {
		return fwd
	}
	return -rvs
}

type blub int

func (b *blub) left() {
	*b = (*b - 1) % FOREST_WIDTH
}

func (b *blub) right() {
	*b = (*b + 1) % FOREST_WIDTH
}

func (b blub) binForwardDistance(target int) int {
	if int(b) == target {
		return 0
	}

	if target > int(b) {
		return target - int(b)
	} else {
		return 30 - int(b) + target
	}
}

func (b blub) binReverseDistance(target int) int {
	if int(b) == target {
		return 0
	}

	if target < int(b) {
		return int(b) - target
	} else {
		return 30 - target + int(b)
	}
}

func (b blub) binDistance(target int) int {
	fwd := b.binForwardDistance(target)

	if fwd == 0 {
		return fwd
	}

	rvs := b.binReverseDistance(target)

	if fwd < rvs {
		return fwd
	}
	return -rvs
}

var forest = make(runes, FOREST_WIDTH)
var position blub

func initializeForest() {
	for i := 0; i < FOREST_WIDTH; i++ {
		forest[i] = ' '
	}
}

type path struct {
	bin            int
	binDistance    int
	letterDistance int
}

func (p path) total() int {
	return int(math.Abs(float64(p.binDistance))) + int(math.Abs(float64(p.letterDistance)))
}

func findShortestPath(target rune, position blub, forest runes) path {
	fmt.Fprintln(os.Stderr, "finding shortest path")

	result := path{
		bin:            int(position),
		letterDistance: forest.letterDistance(int(position), target),
	}

	fmt.Fprintf(os.Stderr, "initial path: %v\n", result)

	for bi := 1; bi < FOREST_WIDTH; bi++ {
		b := (int(position) + bi) % FOREST_WIDTH
		temp := path{
			bin:            b,
			binDistance:    position.binDistance(b),
			letterDistance: forest.letterDistance(b, target),
		}
		if temp.total() < result.total() {
			fmt.Fprintf(os.Stderr, "updating path to bin %d with value %v\n", b, temp)
			result = temp
		}
	}

	return result
}

func generateLoopString(target rune) string {
	ldist := runes{' '}.letterDistance(0, target)
	var lstring string
	if ldist < 0 {
		lstring = strings.Repeat("-", int(math.Abs(float64(ldist))))
	} else if ldist > 0 {
		lstring = strings.Repeat("+", int(math.Abs(float64(ldist))))
	}
	return "[+]" + lstring
}

func generateString(target rune, position *blub, forest runes) string {
	fmt.Fprintf(os.Stderr, "finding path to: %c\n", target)
	p := findShortestPath(target, *position, forest)
	fmt.Fprintln(os.Stderr, p)
	var bstring, lstring string = "", ""
	if p.binDistance < 0 {
		if forest[p.bin] == ' ' && p.binDistance < -3 {
			bstring = "[<]"
		} else {
			bstring = strings.Repeat("<", int(math.Abs(float64(p.binDistance))))
		}
	} else if p.binDistance > 0 {
		if forest[p.bin] == ' ' && p.binDistance > 3 {
			bstring = "[>]"
		} else {
			bstring = strings.Repeat(">", int(math.Abs(float64(p.binDistance))))
		}
	}
	if p.letterDistance < 0 {
		lstring = strings.Repeat("-", int(math.Abs(float64(p.letterDistance))))
	} else if p.letterDistance > 0 {
		lstring = strings.Repeat("+", int(math.Abs(float64(p.letterDistance))))
	}

	semifinal := bstring + lstring

	loopstring := generateLoopString(target)
	if len(loopstring) < len(semifinal) {
		semifinal = loopstring
		forest[int(*position)] = target
	} else {
		*position = blub(p.bin)
		forest[p.bin] = target
	}

	return semifinal + "."
}

func alphabetCipher(phrase string, position blub) string {
	result := "<-[>[+>]+[<]>-]"
	position = blub(29)

	for _, r := range phrase {
		if r == ' ' {
			if int(position) > 23 {
				result += strings.Repeat(">", 26-int(position)) + "."
				position = blub(26)
				continue
			}
			if int(position) < 2 {
				result += strings.Repeat("<", int(position)+1) + "."
				position = blub(29)
				continue
			}
			result += "[>]."
			position = blub(26)
			continue
		}

		tb := 25 - int(r-'A')
		bdist := position.binDistance(tb)

		if int(position) < 26 {
			if int(math.Abs(float64(bdist))) > (tb + 3) {
				result += "[<]" + strings.Repeat(">", tb+1) + "."
				position = blub(tb)
				continue
			}

			if math.Abs(float64(bdist)) > math.Abs(float64(26-tb)) {
				result += "[>]" + strings.Repeat("<", 26-tb) + "."
				position = blub(tb)
				continue
			}
		}

		if bdist < 0 {
			result += strings.Repeat("<", int(math.Abs(float64(bdist))))
		} else if bdist > 0 {
			result += strings.Repeat(">", int(math.Abs(float64(bdist))))
		}
		result += "."

		position = blub(tb)
	}

	return result
}

func uniqueLetterCount(s string) int {
	l := map[rune]bool{}
	for _, r := range s {
		if _, ok := l[r]; !ok {
			l[r] = true
		}
	}
	return len(l)
}

func main() {
	initializeForest()

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	scanner.Scan()
	magicPhrase := scanner.Text()
	fmt.Fprintln(os.Stderr, magicPhrase)

	result := ""
	lcount := uniqueLetterCount(magicPhrase)
	if lcount > 10 && len(magicPhrase) < 30 {
		result = alphabetCipher(magicPhrase, position)
	} else {
		for _, r := range magicPhrase {
			result += generateString(r, &position, forest)
		}
	}

	fmt.Println(result)
}
