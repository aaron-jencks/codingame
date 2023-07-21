package main

import (
	"fmt"
	"strings"
)

func lettersToNumber(s string) int {
	var result int = 0
	for _, c := range s {
		result *= 26
		result += int(c - 'A' + 1)
	}
	return result
}

type Coord struct {
	row int
	col int
}

func coordFromBoard(bc string) Coord {
	c := Coord{}

	var split int
	for ci, c := range bc {
		if strings.ContainsRune("0123456789", c) {
			split = ci
			break
		}
	}

	c.col = lettersToNumber(bc[:split]) - 1

	fmt.Sscanf(bc[split:], "%d", &c.row)
	c.row--

	return c
}

func (c Coord) ToBoardStyle() string {
	return fmt.Sprintf("%c%d", 'A'+rune(c.row), c.col)
}

type Cell struct {
	name   string
	coord  Coord
	left   bool
	right  bool
	top    bool
	bottom bool
}

func cellFromBoard(coord string, sides string) Cell {
	c := Cell{
		name: coord,
	}
	c.coord = coordFromBoard(coord)
	for _, s := range sides {
		switch s {
		case 'L':
			c.left = true
		case 'R':
			c.right = true
		case 'T':
			c.top = true
		case 'B':
			c.bottom = true
		}
	}

	return c
}

func (c Cell) availableSideCount() int {
	var result int = 0

	if c.left {
		result++
	}
	if c.right {
		result++
	}
	if c.bottom {
		result++
	}
	if c.top {
		result++
	}

	return result
}

func (c Cell) Copy() Cell {
	var result Cell = c
	return result
}

type GameController struct {
	ourTurn       bool
	boardSize     int
	ourScore      int
	opponentScore int
	cells         []Cell
	cmap          [][]Cell
}

func (gc GameController) Display() string {
	var result string = "Current Board:\n"
	player := "Opponent"
	if gc.ourTurn {
		player = "Us"
	}
	result += fmt.Sprintf("Player: %s\nBoard Size: %d\nPlayer Score: %d\nOpponent Score: %d\nCells:\n",
		player, gc.boardSize, gc.ourScore, gc.opponentScore)
	for ci, c := range gc.cells {
		result += fmt.Sprintf("%d: Location=(%d, %d) Left=%v Right=%v Bottom=%v Top=%v\n",
			ci, c.coord.row, c.coord.col, c.left, c.right, c.bottom, c.top)
	}
	return result
}

func (gc *GameController) Reset() {
	gc.cells = nil
	for ri := range gc.cmap {
		for ci := range gc.cmap[ri] {
			gc.cmap[ri][ci] = Cell{}
		}
	}
	gc.ourTurn = true
}

func (gc GameController) Copy() GameController {
	var result GameController = gc
	gc.cmap = make([][]Cell, gc.boardSize)
	for ri := range gc.cmap {
		gc.cmap[ri] = make([]Cell, gc.boardSize)
	}
	result.cells = nil
	for _, c := range gc.cells {
		cc := c.Copy()
		result.cells = append(result.cells, cc)
		gc.cmap[cc.coord.row][cc.coord.col] = cc
	}
	return result
}

func (gc *GameController) Move(loc Coord, side string) {
	c := gc.cmap[loc.row][loc.col]

	switch side {
	case "R":
		c.right = false
	case "L":
		c.left = false
	case "T":
		c.top = false
	case "B":
		c.bottom = false
	}

	if c.availableSideCount() == 0 {
		if gc.ourTurn {
			gc.ourScore++
		} else {
			gc.opponentScore++
		}
	} else {
		gc.ourTurn = !gc.ourTurn
	}
}

func (gc GameController) ScoreState() int {
	return gc.ourScore - gc.opponentScore
}

type Move struct {
	coord Coord
	side  string
}

func getOurValidMoves(gc GameController) []Move {
	var result []Move = make([]Move, 0, gc.boardSize*gc.boardSize*4)

	var leftEdge, rightEdge, topEdge, bottomEdge, middle bool
	var side string
	var ri, ci int
	var c Cell

	parseMove := func(c Cell, ci, ri int) {
		if c.left && (ci > 0 && gc.cmap[ri][ci-1].availableSideCount() > 2 || ci == 0) {
			result = append(result, Move{
				coord: Coord{
					row: ri,
					col: ci,
				},
				side: "L",
			})
		}
		if c.right && (ci < gc.boardSize-1 && gc.cmap[ri][ci+1].availableSideCount() > 2 || ci == gc.boardSize-1) {
			result = append(result, Move{
				coord: Coord{
					row: ri,
					col: ci,
				},
				side: "R",
			})
		}
		if c.top && (ri > 0 && gc.cmap[ri-1][ci].availableSideCount() > 2 || ri == 0) {
			result = append(result, Move{
				coord: Coord{
					row: ri,
					col: ci,
				},
				side: "T",
			})
		}
		if c.bottom && (ri < gc.boardSize-1 && gc.cmap[ri+1][ci].availableSideCount() > 2 || ri == gc.boardSize-1) {
			result = append(result, Move{
				coord: Coord{
					row: ri,
					col: ci,
				},
				side: "B",
			})
		}
	}

	for ri = range gc.cmap {
		for ci = range gc.cmap[ri] {
			c = gc.cmap[ri][ci]

			if c.name != "" {
				switch c.availableSideCount() {
				case 1:
					// we can get this cell
					if c.left {
						side = "L"
					} else if c.right {
						side = "R"
					} else if c.bottom {
						side = "B"
					} else if c.top {
						side = "T"
					}

					result = append(result, Move{
						coord: Coord{
							row: ri,
							col: ci,
						},
						side: side,
					})
				case 3:
					parseMove(c, ci, ri)
				case 4:
					if ri == 0 && ci > 0 && ci < gc.boardSize-1 {
						if topEdge {
							continue
						}
						topEdge = true
					}
					if ri == gc.boardSize-1 && ci > 0 && ci < gc.boardSize-1 {
						if bottomEdge {
							continue
						}
						bottomEdge = true
					}
					if ci == 0 && ri > 0 && ri < gc.boardSize-1 {
						if leftEdge {
							continue
						}
						leftEdge = true
					}
					if ci == gc.boardSize-1 && ri > 0 && ri < gc.boardSize-1 {
						if rightEdge {
							continue
						}
						rightEdge = true
					}
					if ri > 0 && ri < gc.boardSize-1 && ci > 0 && ci < gc.boardSize-1 {
						if middle {
							continue
						}
						middle = true
					}

					parseMove(c, ci, ri)
				}
			}
		}
	}

	if len(result) == 0 {
		// fmt.Fprintf(os.Stderr, "We failed to find a cell, looking for any position\n")

		// Just move somewhere we already lost
		for ri = range gc.cmap {
			for ci = range gc.cmap[ri] {
				c = gc.cmap[ri][ci]
				if c.name != "" && c.availableSideCount() > 0 {
					if c.left {
						side = "L"
					} else if c.right {
						side = "R"
					} else if c.bottom {
						side = "B"
					} else if c.top {
						side = "T"
					}

					result = append(result, Move{
						coord: Coord{
							row: ri,
							col: ci,
						},
						side: side,
					})
				}
			}
		}
	}

	// fmt.Fprintf(os.Stderr, "Found %d moves for us\n", len(result))

	return result
}

type MinimaxResult struct {
	m     Move
	score int
}

type MinimaxMoveState struct {
	gc           GameController
	alpha        int
	beta         int
	currentDepth int
	bestScore    int
	bestMove     Move
}

// store alpha and beta for each parent move
// store the parent move in the state
// for each child move update alpha and beta

type MinimaxState struct {
	parentMove  Move
	currentMove Move
}

func iterMinimax(gc GameController, depth int) int {
	var nd int
	var ngc GameController
	var m Move
	var moves []Move
	var omoves []Move

	maxStackSize := 9
	for i := 1; i < depth; i++ {
		maxStackSize += 9 + i
	}

	moveStates := map[Move]*MinimaxMoveState{
		{}: {
			gc:        gc,
			alpha:     -2147483647,
			beta:      2147483647,
			bestScore: -2147483647,
		},
	}

	stack := make([]MinimaxState, 0, maxStackSize)

	// initialize the stack
	omoves = getOurValidMoves(gc)

	for _, m = range omoves {
		ngc = gc.Copy()
		nd = 1

		ngc.Move(m.coord, m.side)

		mms := MinimaxMoveState{
			gc:           ngc,
			alpha:        -2147483647,
			beta:         2147483647,
			currentDepth: nd,
		}

		mms.bestScore = -2147483647
		if ngc.ourTurn {
			mms.bestScore = 2147483647
		}

		moveStates[m] = &mms

		stack = append(stack, MinimaxState{
			parentMove:  Move{},
			currentMove: m,
		})
	}

	// perform dfs
	for len(stack) > 0 {
		// pop from the stack
		lastElem := len(stack) - 1
		elem := stack[lastElem]
		stack = stack[:lastElem]

		cm := moveStates[elem.currentMove]
		pm := moveStates[elem.parentMove]

		// prune leftover children as we encounter them
		if pm.beta <= pm.alpha {
			continue
		}

		// calculate the score of the move and update
		// the alpha and beta
		if cm.currentDepth == depth {
			nc := cm.gc.ScoreState()
			if pm.gc.ourTurn && nc > pm.bestScore {
				pm.bestMove = elem.currentMove
				pm.bestScore = nc
				if pm.bestScore >= pm.alpha {
					pm.alpha = pm.bestScore
				}
			} else if !pm.gc.ourTurn && nc < pm.bestScore {
				pm.bestMove = elem.currentMove
				pm.bestScore = nc
				if pm.bestScore <= pm.beta {
					pm.beta = pm.bestScore
				}
			}
			continue
		}

		// discover new children
		moves = getOurValidMoves(gc)

		// push the children onto the stack
		for _, m = range moves {
			ngc = cm.gc.Copy()
			nd = cm.currentDepth + 1

			ngc.Move(m.coord, m.side)

			mms := MinimaxMoveState{
				gc:           ngc,
				alpha:        cm.alpha,
				beta:         cm.beta,
				currentDepth: nd,
			}

			mms.bestScore = 2147483647
			if gc.ourTurn {
				mms.bestScore = -2147483647
			}

			moveStates[m] = &mms

			stack = append(stack, MinimaxState{
				parentMove:  elem.currentMove,
				currentMove: m,
			})
		}
	}

	return moveStates[Move{}].bestScore
}

func main() {

	controller := GameController{}

	// boardSize: The size of the board.
	fmt.Scan(&controller.boardSize)

	controller.cmap = make([][]Cell, controller.boardSize)
	for ri := range controller.cmap {
		controller.cmap[ri] = make([]Cell, controller.boardSize)
	}

	// playerId: The ID of the player. 'A'=first player, 'B'=second player.
	var playerId string
	fmt.Scan(&playerId)
	controller.ourTurn = true

	for {
		// Reset cell array
		controller.Reset()

		// playerScore: The player's score.
		// opponentScore: The opponent's score.
		fmt.Scan(&controller.ourScore, &controller.opponentScore)

		// numBoxes: The number of playable boxes.
		var numBoxes int
		fmt.Scan(&numBoxes)

		for i := 0; i < numBoxes; i++ {
			// box: The ID of the playable box.
			// sides: Playable sides of the box.
			var box, sides string
			fmt.Scan(&box, &sides)
			cc := cellFromBoard(box, sides)
			controller.cells = append(controller.cells, cc)
			controller.cmap[cc.coord.row][cc.coord.col] = cc
		}

		// fmt.Fprintln(os.Stderr, controller.Display())

		moves := getOurValidMoves(controller)

		bm := MinimaxResult{}
		first := true
		for _, m := range moves {
			ngc := controller.Copy()

			ngc.Move(m.coord, m.side)

			score := iterMinimax(ngc, 4) // , 0, 4, -2147483647, 2147483647)
			if first {
				bm.m = m
				bm.score = score
				first = false
				continue
			}

			if score > bm.score {
				bm.m = m
				bm.score = score
			}
		}

		fmt.Printf("%s %s\n", controller.cmap[bm.m.coord.row][bm.m.coord.col].name, bm.m.side)
	}
}
