package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
)

// TODO add word probabilities for neighboring letters
// letter neighbor context
// should help quite a bit

type State int

const (
	UNKNOWN State = iota
	INCORRECT
	PARTIAL
	CORRECT
)

type Set[T comparable] map[T]bool

func (s Set[T]) add(e T) {
	s[e] = true
}

func (s Set[T]) contains(e T) bool {
	_, ok := s[e]
	return ok
}

func (s Set[T]) toArr() []T {
	result := make([]T, 0, len(s))
	for k := range s {
		result = append(result, k)
	}
	return result
}

func createLetterSet(w string) Set[rune] {
	result := Set[rune]{}
	for _, l := range w {
		result.add(l)
	}
	return result
}

type RuneWordSetMap map[rune]Set[string]

type LetterFrequencies map[rune]int

func (lf LetterFrequencies) contains(r rune) bool {
	_, ok := lf[r]
	return ok
}

func (lf LetterFrequencies) toArr() []rune {
	result := make([]rune, 0, len(lf))
	for k := range lf {
		result = append(result, k)
	}
	return result
}

func findLetterFrequencies(words []string) []LetterFrequencies {
	result := make([]LetterFrequencies, 6)
	for i := 0; i < 6; i++ {
		result[i] = LetterFrequencies{}
		for r := 'A'; r <= 'Z'; r++ {
			result[i][r] = 0
		}
	}

	for _, w := range words {
		for i := 0; i < 6; i++ {
			result[i][rune(w[i])]++
		}
	}

	return result
}

func condenseFrequences(freqs []LetterFrequencies) LetterFrequencies {
	result := LetterFrequencies{}
	for lfi, lf := range freqs {
		if lfi == 0 {
			for r, c := range lf {
				result[r] = c
			}
		} else {
			for r, c := range lf {
				if result.contains(r) {
					result[r] += c
				} else {
					result[r] = c
				}
			}
		}
	}
	return result
}

func wordSetIntersection(a, b Set[string]) Set[string] {
	result := Set[string]{}
	for w := range a {
		if b.contains(w) {
			result.add(w)
		}
	}
	for w := range b {
		if a.contains(w) {
			result.add(w)
		}
	}
	return result
}

func wordSetFilterBlacklist(s Set[string], bl Set[rune]) Set[string] {
	result := Set[string]{}

FILTEROUTER:
	for w := range s {
		for _, r := range w {
			if bl.contains(r) {
				continue FILTEROUTER
			}
		}
		result.add(w)
	}

	return result
}

type WordleSim struct {
	letterPositions []RuneWordSetMap
	wordContains    RuneWordSetMap
	blacklist       Set[rune]
	prev            []rune
	letterFreqs     []LetterFrequencies
	words           Set[string]
}

func (ws *WordleSim) initializeSim(words []string) {
	ws.letterPositions = []RuneWordSetMap{
		{}, {}, {}, {}, {}, {},
	}
	ws.wordContains = RuneWordSetMap{}
	ws.blacklist = Set[rune]{}
	ws.prev = nil
	ws.words = Set[string]{}

	for r := 'A'; r <= 'Z'; r++ {
		for li := 0; li < 6; li++ {
			ws.letterPositions[li][r] = Set[string]{}
		}
		ws.wordContains[r] = Set[string]{}
	}

	for _, word := range words {
		for wi := 0; wi < 6; wi++ {
			ws.letterPositions[wi][rune(word[wi])].add(word)
		}

		uletters := createLetterSet(word)
		for r := range uletters {
			ws.wordContains[r].add(word)
		}

		ws.words.add(word)
	}

	ws.letterFreqs = findLetterFrequencies(words)
}

func (ws *WordleSim) guess(states []State) string {
	if ws.prev == nil {
		clmap := condenseFrequences(ws.letterFreqs)
		clmapArr := clmap.toArr()
		sort.Slice(clmapArr, func(i, j int) bool {
			// want descending order
			return clmap[clmapArr[i]] > clmap[clmapArr[j]]
		})
		ws.prev = clmapArr[:6]
		return string(clmapArr[:6])
	}

	var wsets []Set[string]
	for si, s := range states {
		switch s {
		case PARTIAL:
			wsets = append(wsets, ws.wordContains[ws.prev[si]])
		case CORRECT:
			wsets = append(wsets, ws.letterPositions[si][ws.prev[si]])
			// fmt.Fprintf(os.Stderr, "words for letter %c at %d\n", ws.prev[si], si)
			// fmt.Fprintln(os.Stderr, ws.letterPositions[si][ws.prev[si]].toArr()[:5])
		case INCORRECT:
			ws.blacklist.add(ws.prev[si])
		}
	}

	var current Set[string]
	if len(wsets) == 0 {
		current = ws.words
	} else {
		current = wsets[0]
		for _, s := range wsets[1:] {
			current = wordSetIntersection(current, s)
		}
	}

	if len(ws.blacklist) > 0 {
		current = wordSetFilterBlacklist(current, ws.blacklist)
	}

	// fmt.Fprintln(os.Stderr, "possible words:")
	// if len(current) > 5 {
	// 	fmt.Fprintln(os.Stderr, current.toArr()[:5])
	// } else {
	// 	fmt.Fprintln(os.Stderr, current.toArr())
	// }

	lmap := findLetterFrequencies(current.toArr())
	guess := make([]rune, 6)
	for li, m := range lmap {
		clmapArr := m.toArr()

		if len(clmapArr) > 1 {
			sort.Slice(clmapArr, func(i, j int) bool {
				// want descending order
				return lmap[li][clmapArr[i]] > lmap[li][clmapArr[j]]
			})
		}

		// fmt.Fprintf(os.Stderr, "letters for index %d\n", li)
		// if len(clmapArr) > 5 {
		// 	fmt.Fprintln(os.Stderr, clmapArr[:5])
		// } else {
		// 	fmt.Fprintln(os.Stderr, clmapArr)
		// }

		guess[li] = clmapArr[0]
	}

	for string(ws.prev) == string(guess) {
		rand.Shuffle(6, func(i, j int) { guess[i], guess[j] = guess[j], guess[i] })
	}

	ws.prev = guess
	return string(guess)
}

func main() {
	sim := WordleSim{}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	scanner.Scan()
	scanner.Text()

	scanner.Scan()
	inputs := strings.Split(scanner.Text(), " ")
	sim.initializeSim(inputs)

	states := make([]State, 6)

	for {
		scanner.Scan()
		inputs = strings.Split(scanner.Text(), " ")
		for i := 0; i < 6; i++ {
			fmt.Sscan(inputs[i], &states[i])
		}
		fmt.Fprintln(os.Stderr, states)

		fmt.Println(sim.guess(states))
	}
}
