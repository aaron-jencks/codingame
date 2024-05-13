package main

import (
	"fmt"
	"io"
	"os"
)

type action int8

const (
	A_SPEED action = iota
	A_SLOW
	A_JUMP
	A_WAIT
	A_UP
	A_DOWN
)

const (
	S_SAFE byte = '.'
	S_HOLE byte = '0'
)

type bike struct {
	speed    int8
	alive    bool
	jumping  int8
	lane     int8
	position int8
}

func (b *bike) clone() *bike {
	return &*b
}

func (b *bike) tick(lanes []string, a action) {
	if b.alive {
		switch a {
		case A_SPEED:
			if b.speed < 10 {
				b.speed++
			}
		case A_SLOW:
			if b.speed > 0 {
				b.speed--
			}
		case A_JUMP:
			if b.jumping == 0 {
				b.jumping = b.speed
			}
		case A_UP:
			if b.lane > 0 {
				b.lane--
				for pi := int8(0); pi < b.speed-1; pi++ {
					if lanes[b.lane][b.position+1+pi] == S_HOLE || lanes[b.lane+1][b.position+1+pi] == S_HOLE {
						b.alive = false
						return
					}
				}
			}
		case A_DOWN:
			if b.lane < 3 {
				b.lane++
				for pi := int8(0); pi < b.speed-1; pi++ {
					if lanes[b.lane][b.position+1+pi] == S_HOLE || lanes[b.lane-1][b.position+1+pi] == S_HOLE {
						b.alive = false
						return
					}
				}
			}
		default:
		}
	}
}

type game struct {
	bikes []*bike
	lanes []string
	turn  uint8
}

func (g *game) getActions() []action {
	var result []action

	cup := true
	csup := true
	cdown := true
	csdown := true
	airborne := false
	for _, b := range g.bikes {
		if b.jumping > 0 {
			airborne = true
			break
		}
		if b.lane == 0 {
			cup = false
		} else if b.lane == 3 {
			cdown = false
		}
		if b.speed == 10 {
			csup = false
		} else if b.speed == 1 {
			csdown = false
		}
	}

	if !airborne {
		if cup {
			result = append(result, A_UP)
		}
		if cdown {
			result = append(result, A_DOWN)
		}

		if csup {
			result = append(result, A_SPEED)
		}
		if csdown {
			result = append(result, A_SLOW)
		}
	} else {
		result = append(result, A_WAIT)
	}

	return result
}

func (g *game) tick(a action) *game {
	gc := &game{
		make([]*bike, len(g.bikes)),
		make([]string, 4),
		g.turn,
	}
	copy(gc.lanes, g.lanes)
	for bi, b := range g.bikes {
		gc.bikes[bi] = b.clone()
		gc.bikes[bi].tick(gc.lanes, a)
	}
	gc.turn++
	return gc
}

func (g game) score() int {
	result := 0
	for _, b := range g.bikes {
		if b.alive && b.speed > 0 {
			result += int(b.position)
		} else {
			result -= 100
		}
	}
	return result
}

func (g *game) clone() *game {
	gc := &game{
		make([]*bike, len(g.bikes)),
		make([]string, 4),
		g.turn,
	}
	copy(gc.lanes, g.lanes)
	for bi, b := range g.bikes {
		gc.bikes[bi] = b.clone()
	}
	return gc
}

type stackGame struct {
	g     *game
	depth int
	move  action
}

func simulateGame(g game, depth int) action {
	if depth == 0 {
		return A_WAIT
	}

	stack := []stackGame{}

	for _, a := range g.getActions() {
		stack = append(stack, stackGame{
			g:     g.tick(a),
			depth: 1,
			move:  a,
		})
	}

	bestScore := -1000
	bestMove := A_WAIT

	for len(stack) > 0 {
		last := len(stack) - 1
		element := stack[last]
		stack = stack[:last]

		if element.depth == depth || element.g.turn == 49 {
			sc := element.g.score()
			if sc > bestScore {
				bestScore = sc
				bestMove = element.move
			}
			continue
		}

		for _, a := range g.getActions() {
			stack = append(stack, stackGame{
				g:     g.tick(a),
				depth: element.depth + 1,
				move:  element.move,
			})
		}
	}

	return bestMove
}

func (g *game) initTurn(r io.Reader) {
	var speed int8
	fmt.Fscan(r, &speed)
	fmt.Fprintln(os.Stderr, speed)
	for bi := 0; bi < len(g.bikes); bi++ {
		var sa int8
		fmt.Fscan(r, &g.bikes[bi].position, &g.bikes[bi].lane, &sa)
		g.bikes[bi].alive = sa == 1
		g.bikes[bi].speed = speed
		fmt.Fprintln(os.Stderr, g.bikes[bi].position, g.bikes[bi].lane, sa)
	}
}

func initGame(r io.Reader) *game {
	var m, v int8
	var l0, l1, l2, l3 string
	fmt.Fscan(r, &m, &v, &l0, &l1, &l2, &l3)
	fmt.Fprintf(os.Stderr, "%d\n%d\n%s\n%s\n%s\n%s\n", m, v, l0, l1, l2, l3)

	g := &game{
		make([]*bike, m),
		[]string{
			l0, l1, l2, l3,
		},
		0,
	}

	for bi := int8(0); bi < m; bi++ {
		g.bikes[bi] = &bike{
			alive: true,
		}
	}

	return g
}

func runMain(r io.Reader, w io.Writer) {
	og := initGame(r)
	first := true

	for {
		og.initTurn(r)

		if first {
			fmt.Fprintln(w, "SPEED")
			first = false
			continue
		}
		nm := simulateGame(*og, 5)

		switch nm {
		case A_SPEED:
			fmt.Fprintln(w, "SPEED")
		case A_SLOW:
			fmt.Fprintln(w, "SLOW")
		case A_WAIT:
			fmt.Fprintln(w, "WAIT")
		case A_DOWN:
			fmt.Fprintln(w, "DOWN")
		case A_UP:
			fmt.Fprintln(w, "UP")
		case A_JUMP:
			fmt.Fprintln(w, "JUMP")
		}
	}
}

func main() {
	runMain(os.Stdin, os.Stdout)
}
