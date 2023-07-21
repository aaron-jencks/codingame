package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTiming(t *testing.T) {
	var avg float64 = 0

	testCases := []struct {
		name             string
		size             int
		cellInitializers []struct {
			row   int
			col   int
			name  string
			sides string
		}
	}{
		{
			name: "empty board",
			size: 7,
		},
		{
			name: "semi-empty board",
			size: 7,
			cellInitializers: []struct {
				row   int
				col   int
				name  string
				sides string
			}{
				{
					row:   0,
					col:   0,
					name:  "A1",
					sides: "LR",
				},
				{
					row:   0,
					col:   1,
					name:  "B1",
					sides: "LR",
				},
				{
					row:   0,
					col:   2,
					name:  "C1",
					sides: "LR",
				},
				{
					row:   0,
					col:   3,
					name:  "D1",
					sides: "LT",
				},
				{
					row:   0,
					col:   4,
					name:  "E1",
					sides: "TR",
				},
				{
					row:   0,
					col:   5,
					name:  "F1",
					sides: "LTR",
				},
				{
					row:   0,
					col:   0,
					name:  "A2",
					sides: "LT",
				},
				{
					row:   1,
					col:   1,
					name:  "B2",
					sides: "R",
				},
				{
					row:   1,
					col:   2,
					name:  "C2",
					sides: "LTR",
				},
				{
					row:   1,
					col:   5,
					name:  "F2",
					sides: "LBR",
				},
				{
					row:   2,
					col:   1,
					name:  "B3",
					sides: "LTR",
				},
				{
					row:   2,
					col:   2,
					name:  "C3",
					sides: "LBR",
				},
				{
					row:   2,
					col:   5,
					name:  "F3",
					sides: "LT",
				},
				{
					row:   2,
					col:   6,
					name:  "G3",
					sides: "TBR",
				},
				{
					row:   3,
					col:   2,
					name:  "C4",
					sides: "LTR",
				},
				{
					row:   5,
					col:   5,
					name:  "F6",
					sides: "LBR",
				},
				{
					row:   6,
					col:   4,
					name:  "E7",
					sides: "LB",
				},
				{
					row:   6,
					col:   5,
					name:  "F7",
					sides: "TR",
				},
			},
		},
	}

	for ti, tt := range testCases {
		t.Run(fmt.Sprintf("%d: %s", ti, tt.name), func(t *testing.T) {
			for i := 0; i < 5; i++ {
				// initialize the test controller
				controller := GameController{
					boardSize: tt.size,
				}
				controller.cmap = make([][]Cell, controller.boardSize)
				for ri := range controller.cmap {
					controller.cmap[ri] = make([]Cell, controller.boardSize)
				}
				for _, ci := range tt.cellInitializers {
					cc := cellFromBoard(ci.name, ci.sides)
					controller.cells = append(controller.cells, cc)
					controller.cmap[ci.row][ci.col] = cc
				}
				for ri := range controller.cmap {
					for ci := range controller.cmap[ri] {
						if controller.cmap[ri][ci].name == "" {
							box := fmt.Sprintf("%c%d", ci+'A', ri+1)
							cc := cellFromBoard(box, "TLRB")
							controller.cells = append(controller.cells, cc)
							controller.cmap[ri][ci] = cc
						}
					}
				}
				controller.ourTurn = true

				// run a single minimax and output time.
				tstart := time.Now() // start time

				moves := getOurValidMoves(controller)

				bm := MinimaxResult{}
				first := true
				for mi, m := range moves {
					ngc := controller.Copy()

					ngc.Move(m.coord, m.side)

					score := iterMinimax(ngc, 4) // , 0, 4, -2147483647, 2147483647)
					fmt.Printf("Completed minmax %d/%d in %d ms\n", mi+1, len(moves), time.Since(tstart).Milliseconds())
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
				avg += float64(time.Since(tstart))
			}

			assert.Less(t, avg/float64(5*time.Millisecond), float64(100), "Should take 100 ms or less for moves")
		})
	}
}

func TestInValidMoves(t *testing.T) {
	testCases := []struct {
		name             string
		size             int
		cellInitializers []struct {
			row   int
			col   int
			name  string
			sides string
		}
		invalidMoves []Move
	}{
		{
			name: "flowright",
			size: 7,
			invalidMoves: []Move{
				{
					coord: Coord{
						row: 0,
						col: 6,
					},
					side: "T",
				},
				{
					coord: Coord{
						row: 0,
						col: 6,
					},
					side: "L",
				},
			},
			cellInitializers: []struct {
				row   int
				col   int
				name  string
				sides string
			}{
				{
					row:   0,
					col:   0,
					name:  "A1",
					sides: "LR",
				},
				{
					row:   0,
					col:   1,
					name:  "B1",
					sides: "LR",
				},
				{
					row:   0,
					col:   2,
					name:  "C1",
					sides: "LR",
				},
				{
					row:   0,
					col:   3,
					name:  "D1",
					sides: "LR",
				},
				{
					row:   0,
					col:   4,
					name:  "E1",
					sides: "LR",
				},
				{
					row:   0,
					col:   5,
					name:  "F1",
					sides: "LR",
				},
				{
					row:   0,
					col:   6,
					name:  "G1",
					sides: "TL",
				},
				{
					row:   1,
					col:   0,
					name:  "A2",
					sides: "LT",
				},
				{
					row:   1,
					col:   1,
					name:  "B2",
					sides: "TR",
				},
				{
					row:   1,
					col:   2,
					name:  "C2",
					sides: "LT",
				},
				{
					row:   1,
					col:   3,
					name:  "D2",
					sides: "TR",
				},
				{
					row:   1,
					col:   4,
					name:  "E2",
					sides: "LT",
				},
				{
					row:   1,
					col:   5,
					name:  "F2",
					sides: "TR",
				},
				{
					row:   1,
					col:   6,
					name:  "G2",
					sides: "LTB",
				},
				{
					row:   2,
					col:   6,
					name:  "G3",
					sides: "LTB",
				},
				{
					row:   3,
					col:   6,
					name:  "G4",
					sides: "LTB",
				},
				{
					row:   4,
					col:   6,
					name:  "G5",
					sides: "LTB",
				},
				{
					row:   5,
					col:   6,
					name:  "G6",
					sides: "LTB",
				},
				{
					row:   6,
					col:   6,
					name:  "G7",
					sides: "LB",
				},
				{
					row:   6,
					col:   5,
					name:  "F7",
					sides: "LBR",
				},
				{
					row:   6,
					col:   4,
					name:  "E7",
					sides: "LBR",
				},
				{
					row:   6,
					col:   3,
					name:  "D7",
					sides: "LBR",
				},
				{
					row:   6,
					col:   2,
					name:  "C7",
					sides: "LBR",
				},
				{
					row:   6,
					col:   1,
					name:  "B7",
					sides: "LBR",
				},
				{
					row:   6,
					col:   0,
					name:  "A7",
					sides: "LBR",
				},
			},
		},
	}

	for ti, tt := range testCases {
		t.Run(fmt.Sprintf("%d: %s", ti, tt.name), func(t *testing.T) {
			// initialize the test controller
			controller := GameController{
				boardSize: tt.size,
			}
			controller.cmap = make([][]Cell, controller.boardSize)
			for ri := range controller.cmap {
				controller.cmap[ri] = make([]Cell, controller.boardSize)
			}
			for _, ci := range tt.cellInitializers {
				cc := cellFromBoard(ci.name, ci.sides)
				controller.cells = append(controller.cells, cc)
				controller.cmap[ci.row][ci.col] = cc
			}
			for ri := range controller.cmap {
				for ci := range controller.cmap[ri] {
					if controller.cmap[ri][ci].name == "" {
						box := fmt.Sprintf("%c%d", ci+'A', ri+1)
						cc := cellFromBoard(box, "TLRB")
						controller.cells = append(controller.cells, cc)
						controller.cmap[ri][ci] = cc
					}
				}
			}
			controller.ourTurn = true

			moves := getOurValidMoves(controller)

			for _, m := range moves {
				for _, nm := range tt.invalidMoves {
					if nm.coord.row == m.coord.row && nm.coord.col == m.coord.col {
						assert.NotEqual(t, nm.side, m.side, "found invalid move (%d, %d) side %s",
							m.coord.col, m.coord.row, m.side)
					}
				}
			}
		})
	}
}
