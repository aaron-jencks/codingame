package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoomParsing(t *testing.T) {
	tcs := []struct {
		line   string
		result Room
	}{
		{
			line: "4",
			result: Room{
				Rtype:     4,
				Rotatable: true,
				Rotation:  ROTATION_0,
			},
		},
		{
			line: "-4",
			result: Room{
				Rtype:     4,
				Rotatable: false,
				Rotation:  ROTATION_0,
			},
		},
	}

	for ti, tc := range tcs {
		t.Run(fmt.Sprintf("TestRoomParsing(%d)", ti), func(tt *testing.T) {
			result := ParseRoom(tc.line)
			assert.Equal(tt, tc.result, result, "parsed rooms should be equal")
		})
	}
}

func TestRoomRotation(t *testing.T) {
	// Rotate right (clockwise)
	tcs := []struct {
		initial Room
		result  Room
	}{
		{
			initial: Room{
				Rtype:     4,
				Rotatable: true,
				Rotation:  ROTATION_0,
			},
			result: Room{
				Rtype:     4,
				Rotatable: true,
				Rotation:  ROTATION_270,
			},
		},
		{
			initial: Room{
				Rtype:     4,
				Rotatable: true,
				Rotation:  ROTATION_90,
			},
			result: Room{
				Rtype:     4,
				Rotatable: true,
				Rotation:  ROTATION_0,
			},
		},
		{
			initial: Room{
				Rtype:     4,
				Rotatable: false,
				Rotation:  ROTATION_0,
			},
			result: Room{
				Rtype:     4,
				Rotatable: false,
				Rotation:  ROTATION_0,
			},
		},
	}

	for ti, tc := range tcs {
		t.Run(fmt.Sprintf("TestRoomRotating(Right: %d)", ti), func(tt *testing.T) {
			tc.initial.Right()
			assert.Equal(tt, tc.result, tc.initial, "rotated rooms should be equal")
		})
	}

	// Rotate left (counter-clockwise)
	tcs = []struct {
		initial Room
		result  Room
	}{
		{
			initial: Room{
				Rtype:     4,
				Rotatable: true,
				Rotation:  ROTATION_0,
			},
			result: Room{
				Rtype:     4,
				Rotatable: true,
				Rotation:  ROTATION_90,
			},
		},
		{
			initial: Room{
				Rtype:     4,
				Rotatable: true,
				Rotation:  ROTATION_270,
			},
			result: Room{
				Rtype:     4,
				Rotatable: true,
				Rotation:  ROTATION_0,
			},
		},
		{
			initial: Room{
				Rtype:     4,
				Rotatable: false,
				Rotation:  ROTATION_0,
			},
			result: Room{
				Rtype:     4,
				Rotatable: false,
				Rotation:  ROTATION_0,
			},
		},
	}

	for ti, tc := range tcs {
		t.Run(fmt.Sprintf("TestRoomRotating(Right: %d)", ti), func(tt *testing.T) {
			tc.initial.Left()
			assert.Equal(tt, tc.result, tc.initial, "rotated rooms should be equal")
		})
	}
}

func TestCoordParsing(t *testing.T) {
	tcs := []struct {
		line   string
		result ObjectCoord
	}{
		{
			line: "4 3 TOP",
			result: ObjectCoord{
				X:        4,
				Y:        3,
				Entrance: INDY_TOP,
			},
		},
		{
			line: "3 4 RIGHT",
			result: ObjectCoord{
				X:        3,
				Y:        4,
				Entrance: INDY_RIGHT,
			},
		},
		{
			line: "10 1 LEFT",
			result: ObjectCoord{
				X:        10,
				Y:        1,
				Entrance: INDY_LEFT,
			},
		},
	}

	for ti, tc := range tcs {
		t.Run(fmt.Sprintf("TestCoordParsing(%d)", ti), func(tt *testing.T) {
			result := ParseObjectCoord(tc.line)
			assert.Equal(tt, tc.result, result, "parsed coords should be equal")
		})
	}
}

func TestMapParsing(t *testing.T) {
	tcs := []struct {
		h, w   int
		lines  []string
		result Map
	}{
		{
			w: 2, h: 4,
			lines: []string{
				"4 -3",
				"11 -10",
				"11 5",
				"2 3",
				"1",
			},
			result: Map{
				Width:  2,
				Height: 4,
				Exit:   1,
				Rooms: [][]Room{
					{
						ParseRoom("4"),
						ParseRoom("-3"),
					},
					{
						ParseRoom("11"),
						ParseRoom("-10"),
					},
					{
						ParseRoom("11"),
						ParseRoom("5"),
					},
					{
						ParseRoom("2"),
						ParseRoom("3"),
					},
				},
			},
		},
	}

	for ti, tc := range tcs {
		t.Run(fmt.Sprintf("TestMapParsing(%d)", ti), func(tt *testing.T) {
			result := ParseMap(tc.h, tc.w, tc.lines)
			assert.Equal(tt, tc.result, result, "parsed maps should be equal")
		})
	}
}

func TestMapCloning(t *testing.T) {
	initial := Map{
		Width:  2,
		Height: 4,
		Exit:   1,
		Rooms: [][]Room{
			{
				ParseRoom("4"),
				ParseRoom("-3"),
			},
			{
				ParseRoom("11"),
				ParseRoom("-10"),
			},
			{
				ParseRoom("11"),
				ParseRoom("5"),
			},
			{
				ParseRoom("2"),
				ParseRoom("3"),
			},
		},
		Rocks: map[int]ObjectCoord{
			0: ParseObjectCoord("5 5 TOP"),
		},
	}

	output := initial.Clone()

	assert.Equal(t, initial, output, "cloned map should match the initial")

	initial.Rooms[0][0].Rotatable = false

	assert.NotEqual(t, initial, output, "modifying the initial should not change the clone")
}

func TestTheoreticalExits(t *testing.T) {
	tcs := []struct {
		name string
		r    Room
		out  [][]int
	}{
		{
			name: "cross",
			r:    ParseRoom("1"),
			out: [][]int{
				{EXIT_BOTTOM},
				{EXIT_BOTTOM},
				{EXIT_BOTTOM},
			},
		},
		{
			name: "horizontal/vertical",
			r:    ParseRoom("2"),
			out: [][]int{
				{EXIT_BOTTOM},
				{EXIT_RIGHT},
				{EXIT_LEFT},
			},
		},
		{
			name: "horizontal locked",
			r:    ParseRoom("-2"),
			out: [][]int{
				nil,
				{EXIT_RIGHT},
				{EXIT_LEFT},
			},
		},
		{
			name: "vertical locked",
			r:    ParseRoom("-3"),
			out: [][]int{
				{EXIT_BOTTOM},
				nil, nil,
			},
		},
		{
			name: "tlrb/trlb",
			r:    ParseRoom("4"),
			out: [][]int{
				{EXIT_RIGHT, EXIT_LEFT},
				{EXIT_BOTTOM},
				{EXIT_BOTTOM},
			},
		},
		{
			name: "tlrb locked",
			r:    ParseRoom("-4"),
			out: [][]int{
				{EXIT_LEFT},
				nil,
				{EXIT_BOTTOM},
			},
		},
		{
			name: "trlb locked",
			r:    ParseRoom("-5"),
			out: [][]int{
				{EXIT_RIGHT},
				{EXIT_BOTTOM},
				nil,
			},
		},
		{
			name: "T",
			r:    ParseRoom("6"),
			out: [][]int{
				{EXIT_BOTTOM},
				{EXIT_RIGHT, EXIT_BOTTOM},
				{EXIT_LEFT, EXIT_BOTTOM},
			},
		},
		{
			name: "top T locked",
			r:    ParseRoom("-6"),
			out: [][]int{
				nil,
				{EXIT_RIGHT},
				{EXIT_LEFT},
			},
		},
		{
			name: "right T locked",
			r:    ParseRoom("-7"),
			out: [][]int{
				{EXIT_BOTTOM},
				nil,
				{EXIT_BOTTOM},
			},
		},
		{
			name: "bottom T locked",
			r:    ParseRoom("-8"),
			out: [][]int{
				nil,
				{EXIT_BOTTOM},
				{EXIT_BOTTOM},
			},
		},
		{
			name: "left T locked",
			r:    ParseRoom("-9"),
			out: [][]int{
				{EXIT_BOTTOM},
				{EXIT_BOTTOM},
				nil,
			},
		},
		{
			name: "tl/tr/rb/lb",
			r:    ParseRoom("10"),
			out: [][]int{
				{EXIT_LEFT, EXIT_RIGHT},
				{EXIT_BOTTOM},
				{EXIT_BOTTOM},
			},
		},
		{
			name: "tl locked",
			r:    ParseRoom("-10"),
			out: [][]int{
				{EXIT_LEFT},
				nil, nil,
			},
		},
		{
			name: "tr locked",
			r:    ParseRoom("-11"),
			out: [][]int{
				{EXIT_RIGHT},
				nil, nil,
			},
		},
		{
			name: "rb locked",
			r:    ParseRoom("-12"),
			out: [][]int{
				nil, nil,
				{EXIT_BOTTOM},
			},
		},
		{
			name: "lb locked",
			r:    ParseRoom("-13"),
			out: [][]int{
				nil,
				{EXIT_BOTTOM},
				nil,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("TestTheoreticalExits(%s)", tc.name), func(tt *testing.T) {
			for ent, out := range tc.out {
				result := tc.r.TheoreticalExits(ent)
				assert.Equal(tt, out, result, "Expected outputs to match for entrance %d", ent)
			}
		})
	}
}

func TestPathFinding(t *testing.T) {
	tcs := []struct {
		m     Map
		start ObjectCoord
		out   [][]ObjectCoord
	}{
		{
			m: ParseMap(9, 6, []string{
				"0 0 0 0 0 -3",
				"8 3 3 2 2 10",
				"2 0 0 0 10 13",
				"11 3 -2 3 1 13",
				"-3 10 0 0 2 0",
				"0 6 3 3 4 13",
				"0 3 0 13 -4 10",
				"0 13 2 4 10 0",
				"0 0 0 -3 0 0",
				"3",
			}),
			start: ParseObjectCoord("5 0 TOP"),
			out: [][]ObjectCoord{
				{
					ParseObjectCoord("5 0 TOP"),
					ParseObjectCoord("5 1 TOP"),
					ParseObjectCoord("4 1 RIGHT"),
					ParseObjectCoord("3 1 RIGHT"),
					ParseObjectCoord("2 1 RIGHT"),
					ParseObjectCoord("1 1 RIGHT"),
					ParseObjectCoord("0 1 RIGHT"),
					ParseObjectCoord("0 2 TOP"),
					ParseObjectCoord("0 3 TOP"),
					ParseObjectCoord("1 3 LEFT"),
					ParseObjectCoord("2 3 LEFT"),
					ParseObjectCoord("3 3 LEFT"),
					ParseObjectCoord("4 3 LEFT"),
					ParseObjectCoord("4 4 TOP"),
					ParseObjectCoord("4 5 TOP"),
					ParseObjectCoord("3 5 RIGHT"),
					ParseObjectCoord("2 5 RIGHT"),
					ParseObjectCoord("1 5 RIGHT"),
					ParseObjectCoord("1 6 TOP"),
					ParseObjectCoord("1 7 TOP"),
					ParseObjectCoord("2 7 LEFT"),
					ParseObjectCoord("3 7 LEFT"),
					ParseObjectCoord("3 8 TOP"),
				},
				{
					ParseObjectCoord("5 0 TOP"),
					ParseObjectCoord("5 1 TOP"),
					ParseObjectCoord("4 1 RIGHT"),
					ParseObjectCoord("3 1 RIGHT"),
					ParseObjectCoord("2 1 RIGHT"),
					ParseObjectCoord("1 1 RIGHT"),
					ParseObjectCoord("0 1 RIGHT"),
					ParseObjectCoord("0 2 TOP"),
					ParseObjectCoord("0 3 TOP"),
					ParseObjectCoord("1 3 LEFT"),
					ParseObjectCoord("2 3 LEFT"),
					ParseObjectCoord("3 3 LEFT"),
					ParseObjectCoord("4 3 LEFT"),
					ParseObjectCoord("4 4 TOP"),
					ParseObjectCoord("4 5 TOP"),
					ParseObjectCoord("5 5 LEFT"),
					ParseObjectCoord("5 6 TOP"),
					ParseObjectCoord("4 6 RIGHT"),
					ParseObjectCoord("4 7 TOP"),
					ParseObjectCoord("3 7 RIGHT"),
					ParseObjectCoord("3 8 TOP"),
				},
			},
		},
		{
			m: ParseMap(4, 8, []string{
				"0 -3 0 0 0 0 0 0",
				"0 12 3 3 2 3 12 0",
				"0 0 0 0 0 0 2 0",
				"0 -12 3 2 2 3 13 0",
				"1",
			}),
			start: ObjectCoord{
				X:        1,
				Y:        0,
				Entrance: INDY_TOP,
			},
			out: [][]ObjectCoord{
				{
					ParseObjectCoord("1 0 TOP"),
					ParseObjectCoord("1 1 TOP"),
					ParseObjectCoord("2 1 LEFT"),
					ParseObjectCoord("3 1 LEFT"),
					ParseObjectCoord("4 1 LEFT"),
					ParseObjectCoord("5 1 LEFT"),
					ParseObjectCoord("6 1 LEFT"),
					ParseObjectCoord("6 2 TOP"),
					ParseObjectCoord("6 3 TOP"),
					ParseObjectCoord("5 3 RIGHT"),
					ParseObjectCoord("4 3 RIGHT"),
					ParseObjectCoord("3 3 RIGHT"),
					ParseObjectCoord("2 3 RIGHT"),
					ParseObjectCoord("1 3 RIGHT"),
				},
			},
		},
	}

	for ti, tc := range tcs {
		t.Run(fmt.Sprintf("TestPathfinding(%d)", ti), func(tt *testing.T) {
			tc.m.IndyPosition = tc.start
			result := FindMapPath(tc.m)
			assert.Equal(tt, len(tc.out), len(result), "number of found paths should be equal, expected %d, found %d", len(tc.out), len(result))

			for pi, path := range tc.out {
				found := false
				for _, tpath := range result {
					if reflect.DeepEqual(path, tpath.Indy) {
						found = true
					}
				}
				if !found {
					assert.True(tt, false, "all known paths should be found in the result, path %d was not found", pi)
				}
			}
		})
	}
}

func TestPathValidating(t *testing.T) {
	tcs := []struct {
		m     Map
		start ObjectCoord
		out   [][]ObjectCoord
	}{
		{
			m: ParseMap(9, 6, []string{
				"0 0 0 0 0 -3",
				"8 3 3 2 2 10",
				"2 0 0 0 10 13",
				"11 3 -2 3 1 13",
				"-3 10 0 0 2 0",
				"0 6 3 3 4 13",
				"0 3 0 13 -4 10",
				"0 13 2 4 10 0",
				"0 0 0 -3 0 0",
				"3",
			}),
			start: ParseObjectCoord("5 0 TOP"),
			out: [][]ObjectCoord{
				{
					ParseObjectCoord("5 0 TOP"),
					ParseObjectCoord("5 1 TOP"),
					ParseObjectCoord("4 1 RIGHT"),
					ParseObjectCoord("3 1 RIGHT"),
					ParseObjectCoord("2 1 RIGHT"),
					ParseObjectCoord("1 1 RIGHT"),
					ParseObjectCoord("0 1 RIGHT"),
					ParseObjectCoord("0 2 TOP"),
					ParseObjectCoord("0 3 TOP"),
					ParseObjectCoord("1 3 LEFT"),
					ParseObjectCoord("2 3 LEFT"),
					ParseObjectCoord("3 3 LEFT"),
					ParseObjectCoord("4 3 LEFT"),
					ParseObjectCoord("4 4 TOP"),
					ParseObjectCoord("4 5 TOP"),
					ParseObjectCoord("3 5 RIGHT"),
					ParseObjectCoord("2 5 RIGHT"),
					ParseObjectCoord("1 5 RIGHT"),
					ParseObjectCoord("1 6 TOP"),
					ParseObjectCoord("1 7 TOP"),
					ParseObjectCoord("2 7 LEFT"),
					ParseObjectCoord("3 7 LEFT"),
					ParseObjectCoord("3 8 TOP"),
				},
				{
					ParseObjectCoord("5 0 TOP"),
					ParseObjectCoord("5 1 TOP"),
					ParseObjectCoord("4 1 RIGHT"),
					ParseObjectCoord("3 1 RIGHT"),
					ParseObjectCoord("2 1 RIGHT"),
					ParseObjectCoord("1 1 RIGHT"),
					ParseObjectCoord("0 1 RIGHT"),
					ParseObjectCoord("0 2 TOP"),
					ParseObjectCoord("0 3 TOP"),
					ParseObjectCoord("1 3 LEFT"),
					ParseObjectCoord("2 3 LEFT"),
					ParseObjectCoord("3 3 LEFT"),
					ParseObjectCoord("4 3 LEFT"),
					ParseObjectCoord("4 4 TOP"),
					ParseObjectCoord("4 5 TOP"),
					ParseObjectCoord("5 5 LEFT"),
					ParseObjectCoord("5 6 TOP"),
					ParseObjectCoord("4 6 RIGHT"),
					ParseObjectCoord("4 7 TOP"),
					ParseObjectCoord("3 7 RIGHT"),
					ParseObjectCoord("3 8 TOP"),
				},
			},
		},
		{
			m: ParseMap(4, 8, []string{
				"0 -3 0 0 0 0 0 0",
				"0 12 3 3 2 3 12 0",
				"0 0 0 0 0 0 2 0",
				"0 -12 3 2 2 3 13 0",
				"1",
			}),
			start: ObjectCoord{
				X:        1,
				Y:        0,
				Entrance: INDY_TOP,
			},
			out: [][]ObjectCoord{
				{
					ParseObjectCoord("1 0 TOP"),
					ParseObjectCoord("1 1 TOP"),
					ParseObjectCoord("2 1 LEFT"),
					ParseObjectCoord("3 1 LEFT"),
					ParseObjectCoord("4 1 LEFT"),
					ParseObjectCoord("5 1 LEFT"),
					ParseObjectCoord("6 1 LEFT"),
					ParseObjectCoord("6 2 TOP"),
					ParseObjectCoord("6 3 TOP"),
					ParseObjectCoord("5 3 RIGHT"),
					ParseObjectCoord("4 3 RIGHT"),
					ParseObjectCoord("3 3 RIGHT"),
					ParseObjectCoord("2 3 RIGHT"),
					ParseObjectCoord("1 3 RIGHT"),
				},
			},
		},
	}

	for ti, tc := range tcs {
		t.Run(fmt.Sprintf("TestPathfinding(%d)", ti), func(tt *testing.T) {
			tc.m.IndyPosition = tc.start
			result := FindValidMapPath(tc.m)
			assert.Equal(tt, len(tc.out), len(result), "number of found paths should be equal, expected %d, found %d", len(tc.out), len(result))

			for pi, path := range tc.out {
				found := false
				for _, tpath := range result {
					if reflect.DeepEqual(path, tpath.Indy) {
						found = true
					}
				}
				if !found {
					assert.True(tt, false, "all known paths should be found in the result, path %d was not found", pi)
				}
			}
		})
	}
}

func TestNextMove(t *testing.T) {
	tcs := []struct {
		m     Map
		start ObjectCoord
		out   []string
	}{
		{
			m: ParseMap(9, 6, []string{
				"0 0 0 0 0 -3",
				"8 3 3 2 2 10",
				"2 0 0 0 10 13",
				"11 3 -2 3 1 13",
				"-3 10 0 0 2 0",
				"0 6 3 3 4 13",
				"0 3 0 13 -4 10",
				"0 13 2 4 10 0",
				"0 0 0 -3 0 0",
				"3",
			}),
			start: ParseObjectCoord("5 0 TOP"),
			out: []string{
				"2 1 RIGHT",
				"1 1 RIGHT",
				"0 2 RIGHT",
				"1 3 RIGHT",
				"3 3 RIGHT",
				"4 4 RIGHT",
				"3 5 RIGHT",
				"2 5 RIGHT",
				"1 5 RIGHT",
				"1 7 RIGHT",
				"1 7 RIGHT",
				"3 7 RIGHT",
				"WAIT",
			},
		},
	}

	for ti, tc := range tcs {
		t.Run(fmt.Sprintf("TestNextMove(%d)", ti), func(tt *testing.T) {
			for mi, mv := range tc.out {
				tc.m.IndyPosition = tc.start
				result := FindNextMove(tc.m)
				assert.Equal(t, mv, result, "Function should generate a valid rotated path, failed at move %d", mi)
			}
		})
	}
}

func TestRockPermutations(t *testing.T) {
	tcs := []struct {
		m   Map
		st  PathState
		out []map[int]ObjectCoord
	}{
		{
			m: ParseMap(6, 6, []string{
				"0 -3 -3 -5 1 -1",
				"-13 11 4 0 -3 1",
				"11 6 1 0 10 1",
				"0 0 10 3 3 1",
				"0 0 0 0 0 3",
				"0 0 0 0 0 -3",
				"5",
			}),
			st: PathState{
				IndyPosition: ParseObjectCoord("1 0 TOP"),
				Rocks: map[int]ObjectCoord{
					0: ParseObjectCoord("2 0 TOP"),
					1: ParseObjectCoord("0 1 LEFT"),
					2: ParseObjectCoord("4 2 TOP"),
					3: ParseObjectCoord("5 1 TOP"),
				},
			},
			out: []map[int]ObjectCoord{
				{
					0: ParseObjectCoord("2 1 TOP"),
					1: ParseObjectCoord("0 2 TOP"),
					2: ParseObjectCoord("3 2 RIGHT"),
					3: ParseObjectCoord("5 2 TOP"),
				},
				{
					// should eliminate the two collisions
					// but leave an echo
					0: ParseObjectCoord("2 1 TOP"),
					1: ParseObjectCoord("0 2 TOP"),
					3: ObjectCoord{
						X:         5,
						Y:         2,
						Entrance:  INDY_TOP,
						Temporary: true,
					},
				},
			},
		},
	}

	for ti, tc := range tcs {
		t.Run(fmt.Sprintf("TestRockPermutations(%d)", ti), func(tt *testing.T) {
			result := FindRockPermutations(tc.m, tc.st)

			assert.Equal(tt, len(tc.out), len(result), "result and testcase should have same number of permutations")

			for permi, perm := range tc.out {
				found := false
				for _, rperm := range result {
					if reflect.DeepEqual(perm, rperm) {
						found = true
						break
					}
				}
				if !found {
					assert.True(tt, false, "Function should all valid permutations, permutation %d not found", permi)
				}
			}
		})
	}
}

func TestPathFindingWRocks(t *testing.T) {
	tcs := []struct {
		m     Map
		start ObjectCoord
		rocks map[int]ObjectCoord
		out   [][]ObjectCoord
	}{
		{
			m: ParseMap(8, 10, []string{
				"0 0 0 0 0 0 0 0 -3 0",
				"0 7 -2 3 -2 3 -2 3 11 0",
				"0 -7 -2 2 2 2 2 2 2 -2",
				"0 6 -2 2 2 2 2 2 2 -2",
				"0 -7 -2 2 2 2 2 2 2 -2",
				"0 8 -2 2 2 2 2 2 2 -2",
				"0 -7 -2 2 2 2 2 2 2 -2",
				"0 -3 0 0 0 0 0 0 0 0",
				"1",
			}),
			start: ParseObjectCoord("8 0 TOP"),
			out: [][]ObjectCoord{
				{
					ParseObjectCoord("8 0 TOP"),
					ParseObjectCoord("8 1 TOP"),
					ParseObjectCoord("7 1 RIGHT"),
					ParseObjectCoord("6 1 RIGHT"),
					ParseObjectCoord("5 1 RIGHT"),
					ParseObjectCoord("4 1 RIGHT"),
					ParseObjectCoord("3 1 RIGHT"),
					ParseObjectCoord("2 1 RIGHT"),
					ParseObjectCoord("1 1 RIGHT"),
					ParseObjectCoord("1 2 TOP"),
					ParseObjectCoord("1 3 TOP"),
					ParseObjectCoord("1 4 TOP"),
					ParseObjectCoord("1 5 TOP"),
					ParseObjectCoord("1 6 TOP"),
					ParseObjectCoord("1 7 TOP"),
				},
			},
			rocks: map[int]ObjectCoord{
				0: ParseObjectCoord("9 2 RIGHT"),
			},
		},
	}

	for ti, tc := range tcs {
		t.Run(fmt.Sprintf("TestPathfindingWRocks(%d)", ti), func(tt *testing.T) {
			tc.m.IndyPosition = tc.start
			tc.m.Rocks = tc.rocks

			result := FindMapPath(tc.m)
			assert.Equal(tt, len(tc.out), len(result), "number of found paths should be equal, expected %d, found %d", len(tc.out), len(result))

			for pi, path := range tc.out {
				found := false
				for _, tpath := range result {
					if reflect.DeepEqual(path, tpath.Indy) {
						found = true
					}
				}
				if !found {
					assert.True(tt, false, "all known paths should be found in the result, path %d was not found", pi)
				}
			}
		})
	}
}
