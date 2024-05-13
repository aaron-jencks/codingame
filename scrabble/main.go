package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
)

const (
	DOWN = iota
	RIGHT
)

type coord struct {
	x int
	y int
}

type wordOrigin struct {
	origin    coord
	direction int
}

type wordScore struct {
	word  string
	score int
}

func main_io(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	// nbTiles: Number of tiles in the tile set
	var nbTiles int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &nbTiles)
	fmt.Fprintln(os.Stderr, nbTiles)

	tiles := map[rune]int{}

	for i := 0; i < nbTiles; i++ {
		var character string
		var score int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &character, &score)
		fmt.Fprintf(os.Stderr, "%s %d\n", character, score)
		tiles[rune(character[0])] = score
	}

	var width, height int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &width, &height)
	fmt.Fprintln(os.Stderr, width, height)

	scoreTiles := make([]string, height)
	for i := 0; i < height; i++ {
		scanner.Scan()
		scoreTiles[i] = scanner.Text()
		fmt.Fprintln(os.Stderr, scoreTiles[i])
	}

	grid := make([][]rune, height)
	for i := 0; i < height; i++ {
		scanner.Scan()
		grid[i] = []rune(scanner.Text())
		fmt.Fprintln(os.Stderr, string(grid[i]))
	}

	hvisited := make([][]bool, height)
	vvisited := make([][]bool, height)
	newLetters := make([]coord, 0, width*height)
	newGrid := make([][]bool, height)
	for i := 0; i < height; i++ {
		hvisited[i] = make([]bool, width)
		vvisited[i] = make([]bool, width)
		newGrid[i] = make([]bool, width)

		scanner.Scan()
		line := scanner.Text()
		fmt.Fprintln(os.Stderr, line)
		for li, c := range line {
			if grid[i][li] == '.' && c != '.' {
				fmt.Fprintf(os.Stderr, "possible word at (%d, %d)\n", li, i)
				newLetters = append(newLetters, coord{
					li, i,
				})
				grid[i][li] = rune(c)
				newGrid[i][li] = true
			}
		}
	}

	bonus := len(newLetters) >= 7

	words := make([]wordOrigin, 0, len(newLetters))
WORDSEARCH:
	for len(newLetters) > 0 {
		l := newLetters[0]
		newLetters = newLetters[1:]

		if hvisited[l.y][l.x] && vvisited[l.y][l.x] {
			continue
		}

		// find total word horizontal
		count := 0

		var cx int
		invalid := false
		for cx = l.x; cx >= 0 && grid[l.y][cx] != '.'; cx-- {
			if hvisited[l.y][cx] {
				invalid = true
				break
			}
			hvisited[l.y][cx] = true
			count++
		}

		if !invalid && (count > 1 || (l.x < width-1 && grid[l.y][l.x+1] != '.')) {
			fmt.Fprintf(os.Stderr, "found a word starting at (%d, %d)\n", cx+1, l.y)
			words = append(words, wordOrigin{
				origin: coord{
					cx + 1, l.y,
				},
				direction: RIGHT,
			})
		}

		// find total word vertical
		count = 0

		var cy int
		for cy = l.y; cy >= 0 && grid[cy][l.x] != '.'; cy-- {
			if vvisited[cy][l.x] {
				continue WORDSEARCH
			}
			vvisited[cy][l.x] = true
			count++
		}

		if count > 1 || (l.y < height-1 && grid[l.y+1][l.x] != '.') {
			fmt.Fprintf(os.Stderr, "found a word starting at (%d, %d)\n", l.x, cy+1)
			words = append(words, wordOrigin{
				origin: coord{
					l.x, cy + 1,
				},
				direction: DOWN,
			})
		}
	}

	total := 0
	scores := make([]wordScore, 0, len(words)<<1)

	for _, w := range words {
		doubleCount := 0
		tripleCount := 0
		sum := 0

		// we know that this coordinate is the leftmost or topmost letter

		var rw string

		// horizontal
		if w.direction == RIGHT {
			for cx := w.origin.x; cx < width && grid[w.origin.y][cx] != '.'; cx++ {
				rw += string(grid[w.origin.y][cx])
				multiplier := 1
				if newGrid[w.origin.y][cx] {
					switch scoreTiles[w.origin.y][cx] {
					case 'l':
						multiplier = 2
					case 'L':
						multiplier = 3
					case 'w':
						doubleCount++
					case 'W':
						tripleCount++
					}
				}
				sum += tiles[grid[w.origin.y][cx]] * multiplier
			}

			sum <<= doubleCount
			for tripleCount > 0 {
				sum *= 3
				tripleCount--
			}

			if len(rw) > 1 {
				total += sum
				scores = append(scores, wordScore{
					word:  rw,
					score: sum,
				})
			}
		} else {
			// vertical
			for cy := w.origin.y; cy < height && grid[cy][w.origin.x] != '.'; cy++ {
				rw += string(grid[cy][w.origin.x])
				multiplier := 1
				if newGrid[cy][w.origin.x] {
					switch scoreTiles[cy][w.origin.x] {
					case 'l':
						multiplier = 2
					case 'L':
						multiplier = 3
					case 'w':
						doubleCount++
					case 'W':
						tripleCount++
					}
				}
				sum += tiles[grid[cy][w.origin.x]] * multiplier
			}

			sum <<= doubleCount
			for tripleCount > 0 {
				sum *= 3
				tripleCount--
			}

			if len(rw) > 1 {
				total += sum
				scores = append(scores, wordScore{
					word:  rw,
					score: sum,
				})
			}
		}
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].word < scores[j].word
	})

	for _, w := range scores {
		fmt.Fprintf(out, "%s %d\n", w.word, w.score)

	}
	if bonus {
		fmt.Fprintln(out, "Bonus 50")
		total += 50
	}
	fmt.Fprintf(out, "Total %d\n", total)
}

func main() {
	main_io(os.Stdin, os.Stdout)
}
