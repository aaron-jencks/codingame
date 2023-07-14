package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	ROOM_EMPTY = iota
	ROOM_CROSS
	ROOM_HORIZONTAL
	ROOM_VERTICAL
	ROOM_TLRB
	ROOM_TRLB
	ROOM_TOP_T
	ROOM_RIGHT_T
	ROOM_BOTTOM_T
	ROOM_LEFT_T
	ROOM_TL
	ROOM_TR
	ROOM_RB
	ROOM_LB
)

const (
	ROTATION_0 = iota
	ROTATION_90
	ROTATION_180
	ROTATION_270
	ROTATION_360
)

const (
	INDY_TOP = iota
	INDY_LEFT
	INDY_RIGHT
)

const (
	MOVE_LEFT = iota
	MOVE_RIGHT
)

const (
	EXIT_LEFT = iota
	EXIT_BOTTOM
	EXIT_RIGHT
	EXIT_INVALID
)

func FindEntranceFromExit(exit int) int {
	switch exit {
	case EXIT_LEFT:
		return INDY_RIGHT
	case EXIT_BOTTOM:
		return INDY_TOP
	case EXIT_RIGHT:
		return INDY_LEFT
	}
	return INDY_TOP
}

type Room struct {
	Rtype     int
	Rotatable bool
	Rotation  int
}

func (r *Room) Left() {
	if r.Rotatable {
		r.Rotation = (r.Rotation + ROTATION_90) % ROTATION_360
	}
}

func (r *Room) Right() {
	if r.Rotatable {
		r.Rotation -= ROTATION_90
		if r.Rotation < ROTATION_0 {
			r.Rotation = ROTATION_270
		}
	}
}

/**
 * returns exits for this room based on the input regardless of the rotation
 * still checks the rotatability of the room though,
 * if a room cannot be rotated and the input side is invalid, then no exits are returned,
 * otherwise all possible exits for the input side are returned.
 */
func (r Room) TheoreticalExits(in int) []int {
	switch r.Rtype {
	case ROOM_EMPTY:
		return nil
	case ROOM_CROSS:
		return []int{
			EXIT_BOTTOM,
		}
	case ROOM_HORIZONTAL:
		if !r.Rotatable && in == INDY_TOP {
			return nil
		} else {
			switch in {
			case INDY_LEFT:
				return []int{
					EXIT_RIGHT,
				}
			case INDY_RIGHT:
				return []int{
					EXIT_LEFT,
				}
			case INDY_TOP:
				return []int{
					EXIT_BOTTOM,
				}
			}
		}
	case ROOM_VERTICAL:
		if !r.Rotatable && in != INDY_TOP {
			return nil
		} else {
			switch in {
			case INDY_LEFT:
				return []int{
					EXIT_RIGHT,
				}
			case INDY_RIGHT:
				return []int{
					EXIT_LEFT,
				}
			case INDY_TOP:
				return []int{
					EXIT_BOTTOM,
				}
			}
		}
	case ROOM_TLRB:
		if !r.Rotatable && in == INDY_LEFT {
			return nil
		} else {
			switch in {
			case INDY_LEFT:
				return []int{
					EXIT_BOTTOM,
				}
			case INDY_RIGHT:
				return []int{
					EXIT_BOTTOM,
				}
			case INDY_TOP:
				if r.Rotatable {
					return []int{
						EXIT_RIGHT,
						EXIT_LEFT,
					}
				} else {
					return []int{
						EXIT_LEFT,
					}
				}
			}
		}
	case ROOM_TRLB:
		if !r.Rotatable && in == INDY_RIGHT {
			return nil
		} else {
			switch in {
			case INDY_LEFT:
				return []int{
					EXIT_BOTTOM,
				}
			case INDY_RIGHT:
				return []int{
					EXIT_BOTTOM,
				}
			case INDY_TOP:
				if r.Rotatable {
					return []int{
						EXIT_RIGHT,
						EXIT_LEFT,
					}
				} else {
					return []int{
						EXIT_RIGHT,
					}
				}
			}
		}
	case ROOM_TOP_T:
		if !r.Rotatable && in == INDY_TOP {
			return nil
		} else {
			switch in {
			case INDY_LEFT:
				if r.Rotatable {
					return []int{
						EXIT_RIGHT,
						EXIT_BOTTOM,
					}
				} else {
					return []int{
						EXIT_RIGHT,
					}
				}
			case INDY_RIGHT:
				if r.Rotatable {
					return []int{
						EXIT_LEFT,
						EXIT_BOTTOM,
					}
				} else {
					return []int{
						EXIT_LEFT,
					}
				}
			case INDY_TOP:
				return []int{
					EXIT_BOTTOM,
				}
			}
		}
	case ROOM_RIGHT_T:
		if !r.Rotatable && in == INDY_LEFT {
			return nil
		} else {
			if !r.Rotatable {
				return []int{
					EXIT_BOTTOM,
				}
			}

			switch in {
			case INDY_LEFT:
				return []int{
					EXIT_RIGHT,
					EXIT_BOTTOM,
				}
			case INDY_RIGHT:
				return []int{
					EXIT_LEFT,
					EXIT_BOTTOM,
				}
			case INDY_TOP:
				return []int{
					EXIT_BOTTOM,
				}
			}
		}
	case ROOM_BOTTOM_T:
		if !r.Rotatable && in == INDY_TOP {
			return nil
		} else {
			if !r.Rotatable {
				return []int{
					EXIT_BOTTOM,
				}
			}

			switch in {
			case INDY_LEFT:
				return []int{
					EXIT_RIGHT,
					EXIT_BOTTOM,
				}
			case INDY_RIGHT:
				return []int{
					EXIT_LEFT,
					EXIT_BOTTOM,
				}
			case INDY_TOP:
				return []int{
					EXIT_BOTTOM,
				}
			}
		}
	case ROOM_LEFT_T:
		if !r.Rotatable && in == INDY_RIGHT {
			return nil
		} else {
			if !r.Rotatable {
				return []int{
					EXIT_BOTTOM,
				}
			}

			switch in {
			case INDY_LEFT:
				return []int{
					EXIT_RIGHT,
					EXIT_BOTTOM,
				}
			case INDY_RIGHT:
				return []int{
					EXIT_LEFT,
					EXIT_BOTTOM,
				}
			case INDY_TOP:
				return []int{
					EXIT_BOTTOM,
				}
			}
		}
	case ROOM_TL:
		if !r.Rotatable && in != INDY_TOP {
			return nil
		} else {
			switch in {
			case INDY_LEFT:
				fallthrough
			case INDY_RIGHT:
				return []int{
					EXIT_BOTTOM,
				}
			case INDY_TOP:
				if !r.Rotatable {
					return []int{
						EXIT_LEFT,
					}
				}

				return []int{
					EXIT_LEFT,
					EXIT_RIGHT,
				}
			}
		}
	case ROOM_TR:
		if !r.Rotatable && in != INDY_TOP {
			return nil
		} else {
			switch in {
			case INDY_LEFT:
				fallthrough
			case INDY_RIGHT:
				return []int{
					EXIT_BOTTOM,
				}
			case INDY_TOP:
				if !r.Rotatable {
					return []int{
						EXIT_RIGHT,
					}
				}

				return []int{
					EXIT_LEFT,
					EXIT_RIGHT,
				}
			}
		}
	case ROOM_RB:
		if !r.Rotatable && in != INDY_RIGHT {
			return nil
		} else {
			switch in {
			case INDY_LEFT:
				fallthrough
			case INDY_RIGHT:
				return []int{
					EXIT_BOTTOM,
				}
			case INDY_TOP:
				return []int{
					EXIT_LEFT,
					EXIT_RIGHT,
				}
			}
		}
	case ROOM_LB:
		if !r.Rotatable && in != INDY_LEFT {
			return nil
		} else {
			switch in {
			case INDY_LEFT:
				fallthrough
			case INDY_RIGHT:
				return []int{
					EXIT_BOTTOM,
				}
			case INDY_TOP:
				return []int{
					EXIT_LEFT,
					EXIT_RIGHT,
				}
			}
		}
	}

	return nil
}

func (r Room) Exit(in int) int {
	switch r.Rtype {
	case ROOM_EMPTY:
		return EXIT_INVALID
	case ROOM_CROSS:
		return EXIT_BOTTOM
	case ROOM_HORIZONTAL:
		switch r.Rotation {
		case ROTATION_0:
			fallthrough
		case ROTATION_180:
			switch in {
			case INDY_TOP:
				return EXIT_INVALID
			case INDY_LEFT:
				return EXIT_RIGHT
			case INDY_RIGHT:
				return EXIT_LEFT
			}
		case ROTATION_90:
			fallthrough
		case ROTATION_270:
			switch in {
			case INDY_TOP:
				return EXIT_BOTTOM
			case INDY_LEFT:
				fallthrough
			case INDY_RIGHT:
				return EXIT_INVALID
			}
		}
	case ROOM_VERTICAL:
		switch r.Rotation {
		case ROTATION_90:
			fallthrough
		case ROTATION_270:
			switch in {
			case INDY_TOP:
				return EXIT_INVALID
			case INDY_LEFT:
				return EXIT_RIGHT
			case INDY_RIGHT:
				return EXIT_LEFT
			}
		case ROTATION_0:
			fallthrough
		case ROTATION_180:
			switch in {
			case INDY_TOP:
				return EXIT_BOTTOM
			case INDY_LEFT:
				fallthrough
			case INDY_RIGHT:
				return EXIT_INVALID
			}
		}
	case ROOM_TLRB:
		switch r.Rotation {
		case ROTATION_0:
			fallthrough
		case ROTATION_180:
			switch in {
			case INDY_TOP:
				return EXIT_LEFT
			case INDY_LEFT:
				return EXIT_INVALID
			case INDY_RIGHT:
				return EXIT_BOTTOM
			}
		case ROTATION_90:
			fallthrough
		case ROTATION_270:
			switch in {
			case INDY_TOP:
				return EXIT_RIGHT
			case INDY_LEFT:
				return EXIT_BOTTOM
			case INDY_RIGHT:
				return EXIT_INVALID
			}
		}
	case ROOM_TRLB:
		switch r.Rotation {
		case ROTATION_0:
			fallthrough
		case ROTATION_180:
			switch in {
			case INDY_TOP:
				return EXIT_RIGHT
			case INDY_LEFT:
				return EXIT_BOTTOM
			case INDY_RIGHT:
				return EXIT_INVALID
			}
		case ROTATION_90:
			fallthrough
		case ROTATION_270:
			switch in {
			case INDY_TOP:
				return EXIT_LEFT
			case INDY_LEFT:
				return EXIT_INVALID
			case INDY_RIGHT:
				return EXIT_BOTTOM
			}
		}
	case ROOM_TOP_T:
		switch r.Rotation {
		case ROTATION_0:
			switch in {
			case INDY_TOP:
				return EXIT_INVALID
			case INDY_LEFT:
				return EXIT_RIGHT
			case INDY_RIGHT:
				return EXIT_LEFT
			}
		case ROTATION_90:
			switch in {
			case INDY_TOP:
				return EXIT_BOTTOM
			case INDY_LEFT:
				return EXIT_BOTTOM
			case INDY_RIGHT:
				return EXIT_INVALID
			}
		case ROTATION_180:
			switch in {
			case INDY_TOP:
				return EXIT_INVALID
			case INDY_LEFT:
				return EXIT_BOTTOM
			case INDY_RIGHT:
				return EXIT_BOTTOM
			}
		case ROTATION_270:
			switch in {
			case INDY_TOP:
				return EXIT_BOTTOM
			case INDY_LEFT:
				return EXIT_INVALID
			case INDY_RIGHT:
				return EXIT_BOTTOM
			}
		}
	case ROOM_RIGHT_T:
		switch r.Rotation {
		case ROTATION_90:
			switch in {
			case INDY_TOP:
				return EXIT_INVALID
			case INDY_LEFT:
				return EXIT_RIGHT
			case INDY_RIGHT:
				return EXIT_LEFT
			}
		case ROTATION_180:
			switch in {
			case INDY_TOP:
				return EXIT_BOTTOM
			case INDY_LEFT:
				return EXIT_BOTTOM
			case INDY_RIGHT:
				return EXIT_INVALID
			}
		case ROTATION_270:
			switch in {
			case INDY_TOP:
				return EXIT_INVALID
			case INDY_LEFT:
				return EXIT_BOTTOM
			case INDY_RIGHT:
				return EXIT_BOTTOM
			}
		case ROTATION_0:
			switch in {
			case INDY_TOP:
				return EXIT_BOTTOM
			case INDY_LEFT:
				return EXIT_INVALID
			case INDY_RIGHT:
				return EXIT_BOTTOM
			}
		}
	case ROOM_BOTTOM_T:
		switch r.Rotation {
		case ROTATION_180:
			switch in {
			case INDY_TOP:
				return EXIT_INVALID
			case INDY_LEFT:
				return EXIT_RIGHT
			case INDY_RIGHT:
				return EXIT_LEFT
			}
		case ROTATION_270:
			switch in {
			case INDY_TOP:
				return EXIT_BOTTOM
			case INDY_LEFT:
				return EXIT_BOTTOM
			case INDY_RIGHT:
				return EXIT_INVALID
			}
		case ROTATION_0:
			switch in {
			case INDY_TOP:
				return EXIT_INVALID
			case INDY_LEFT:
				return EXIT_BOTTOM
			case INDY_RIGHT:
				return EXIT_BOTTOM
			}
		case ROTATION_90:
			switch in {
			case INDY_TOP:
				return EXIT_BOTTOM
			case INDY_LEFT:
				return EXIT_INVALID
			case INDY_RIGHT:
				return EXIT_BOTTOM
			}
		}
	case ROOM_LEFT_T:
		switch r.Rotation {
		case ROTATION_270:
			switch in {
			case INDY_TOP:
				return EXIT_INVALID
			case INDY_LEFT:
				return EXIT_RIGHT
			case INDY_RIGHT:
				return EXIT_LEFT
			}
		case ROTATION_0:
			switch in {
			case INDY_TOP:
				return EXIT_BOTTOM
			case INDY_LEFT:
				return EXIT_BOTTOM
			case INDY_RIGHT:
				return EXIT_INVALID
			}
		case ROTATION_90:
			switch in {
			case INDY_TOP:
				return EXIT_INVALID
			case INDY_LEFT:
				return EXIT_BOTTOM
			case INDY_RIGHT:
				return EXIT_BOTTOM
			}
		case ROTATION_180:
			switch in {
			case INDY_TOP:
				return EXIT_BOTTOM
			case INDY_LEFT:
				return EXIT_INVALID
			case INDY_RIGHT:
				return EXIT_BOTTOM
			}
		}
	case ROOM_TL:
		switch r.Rotation {
		case ROTATION_0:
			switch in {
			case INDY_TOP:
				return EXIT_LEFT
			case INDY_LEFT:
				return EXIT_INVALID
			case INDY_RIGHT:
				return EXIT_INVALID
			}
		case ROTATION_90:
			switch in {
			case INDY_TOP:
				return EXIT_INVALID
			case INDY_LEFT:
				return EXIT_BOTTOM
			case INDY_RIGHT:
				return EXIT_INVALID
			}
		case ROTATION_180:
			switch in {
			case INDY_TOP:
				return EXIT_INVALID
			case INDY_LEFT:
				return EXIT_INVALID
			case INDY_RIGHT:
				return EXIT_BOTTOM
			}
		case ROTATION_270:
			switch in {
			case INDY_TOP:
				return EXIT_RIGHT
			case INDY_LEFT:
				return EXIT_INVALID
			case INDY_RIGHT:
				return EXIT_INVALID
			}
		}
	case ROOM_TR:
		switch r.Rotation {
		case ROTATION_90:
			switch in {
			case INDY_TOP:
				return EXIT_LEFT
			case INDY_LEFT:
				return EXIT_INVALID
			case INDY_RIGHT:
				return EXIT_INVALID
			}
		case ROTATION_180:
			switch in {
			case INDY_TOP:
				return EXIT_INVALID
			case INDY_LEFT:
				return EXIT_BOTTOM
			case INDY_RIGHT:
				return EXIT_INVALID
			}
		case ROTATION_270:
			switch in {
			case INDY_TOP:
				return EXIT_INVALID
			case INDY_LEFT:
				return EXIT_INVALID
			case INDY_RIGHT:
				return EXIT_BOTTOM
			}
		case ROTATION_0:
			switch in {
			case INDY_TOP:
				return EXIT_RIGHT
			case INDY_LEFT:
				return EXIT_INVALID
			case INDY_RIGHT:
				return EXIT_INVALID
			}
		}
	case ROOM_RB:
		switch r.Rotation {
		case ROTATION_270:
			switch in {
			case INDY_TOP:
				return EXIT_LEFT
			case INDY_LEFT:
				return EXIT_INVALID
			case INDY_RIGHT:
				return EXIT_INVALID
			}
		case ROTATION_0:
			switch in {
			case INDY_TOP:
				return EXIT_INVALID
			case INDY_LEFT:
				return EXIT_BOTTOM
			case INDY_RIGHT:
				return EXIT_INVALID
			}
		case ROTATION_90:
			switch in {
			case INDY_TOP:
				return EXIT_INVALID
			case INDY_LEFT:
				return EXIT_INVALID
			case INDY_RIGHT:
				return EXIT_BOTTOM
			}
		case ROTATION_180:
			switch in {
			case INDY_TOP:
				return EXIT_RIGHT
			case INDY_LEFT:
				return EXIT_INVALID
			case INDY_RIGHT:
				return EXIT_INVALID
			}
		}
	case ROOM_LB:
		switch r.Rotation {
		case ROTATION_180:
			switch in {
			case INDY_TOP:
				return EXIT_LEFT
			case INDY_LEFT:
				return EXIT_INVALID
			case INDY_RIGHT:
				return EXIT_INVALID
			}
		case ROTATION_270:
			switch in {
			case INDY_TOP:
				return EXIT_INVALID
			case INDY_LEFT:
				return EXIT_BOTTOM
			case INDY_RIGHT:
				return EXIT_INVALID
			}
		case ROTATION_0:
			switch in {
			case INDY_TOP:
				return EXIT_INVALID
			case INDY_LEFT:
				return EXIT_INVALID
			case INDY_RIGHT:
				return EXIT_BOTTOM
			}
		case ROTATION_90:
			switch in {
			case INDY_TOP:
				return EXIT_RIGHT
			case INDY_LEFT:
				return EXIT_INVALID
			case INDY_RIGHT:
				return EXIT_INVALID
			}
		}
	}

	return EXIT_INVALID
}

func (r Room) Clone() Room {
	return Room{
		Rtype:     r.Rtype,
		Rotatable: r.Rotatable,
		Rotation:  r.Rotation,
	}
}

func ParseRoom(line string) Room {
	var result Room = Room{
		Rotatable: true,
		Rotation:  ROTATION_0,
	}

	var rt int
	fmt.Sscanf(line, "%d", &rt)
	if rt < 0 {
		result.Rotatable = false
		rt *= -1
	}

	result.Rtype = rt
	return result
}

type ObjectCoord struct {
	X        int
	Y        int
	Entrance int
}

type Map struct {
	Height       int
	Width        int
	Rooms        [][]Room
	Exit         int
	IndyPosition ObjectCoord
	Rocks        []ObjectCoord
}

func (m Map) Clone() Map {
	result := Map{
		Height:       m.Height,
		Width:        m.Width,
		Exit:         m.Exit,
		IndyPosition: m.IndyPosition,
		Rooms:        make([][]Room, m.Height),
		Rocks:        make([]ObjectCoord, len(m.Rocks)),
	}

	for ri := range m.Rooms {
		result.Rooms[ri] = make([]Room, m.Width)
		copy(result.Rooms[ri], m.Rooms[ri])
	}

	copy(result.Rocks, m.Rocks)

	return result
}

func ParseObjectCoord(line string) ObjectCoord {
	var result ObjectCoord = ObjectCoord{}
	var entrance string
	fmt.Sscanf(line, "%d %d %s", &result.X, &result.Y, &entrance)
	switch entrance {
	case "TOP":
		result.Entrance = INDY_TOP
	case "RIGHT":
		result.Entrance = INDY_RIGHT
	case "LEFT":
		result.Entrance = INDY_LEFT
	}
	return result
}

func ParseMap(h, w int, lines []string) Map {
	result := Map{
		Height: h,
		Width:  w,
	}
	fmt.Sscanf(lines[len(lines)-1], "%d", &result.Exit)

	result.Rooms = make([][]Room, h)
	for ri, line := range lines[:len(lines)-1] {
		result.Rooms[ri] = make([]Room, w)

		rooms := strings.Split(line, " ")
		for ci := range rooms {
			result.Rooms[ri][ci] = ParseRoom(rooms[ci])
		}
	}

	return result
}

type SolutionMove struct {
	X         int
	Y         int
	Direction int
}

type PathState struct {
	IndyPosition ObjectCoord
	Solution     []ObjectCoord
}

func FindPathNeighbors(m Map, st PathState) []ObjectCoord {
	var result []ObjectCoord

	for _, exit := range m.Rooms[st.IndyPosition.Y][st.IndyPosition.X].TheoreticalExits(st.IndyPosition.Entrance) {
		nm := ObjectCoord{
			X:        st.IndyPosition.X,
			Y:        st.IndyPosition.Y,
			Entrance: FindEntranceFromExit(exit),
		}
		switch exit {
		case EXIT_BOTTOM:
			if st.IndyPosition.Y+1 < m.Height {
				nm.Y = st.IndyPosition.Y + 1
			} else {
				continue
			}
		case EXIT_LEFT:
			if st.IndyPosition.X-1 >= 0 {
				nm.X = st.IndyPosition.X - 1
			} else {
				continue
			}
		case EXIT_RIGHT:
			if st.IndyPosition.X+1 < m.Width {
				nm.X = st.IndyPosition.X + 1
			} else {
				continue
			}
		}
		result = append(result, nm)
	}

	return result
}

func FindMapPath(m Map, start ObjectCoord) [][]ObjectCoord {
	var result [][]ObjectCoord

	stack := []PathState{
		{
			IndyPosition: start,
		},
	}

	for len(stack) > 0 {
		st := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		for _, mv := range FindPathNeighbors(m, st) {
			ns := PathState{
				IndyPosition: mv,
				Solution:     append(st.Solution, st.IndyPosition),
			}

			if mv.Y == m.Height-1 && mv.X == m.Exit {
				// found a solution
				result = append(result, append(ns.Solution, mv))
			} else {
				stack = append(stack, ns)
			}
		}
	}

	return result
}

func FindValidMapPath(m Map, start ObjectCoord) [][]ObjectCoord {
	var result [][]ObjectCoord

	for _, path := range FindMapPath(m, start) {
		budget := 1
		valid := true

		for pindex := 0; pindex < len(path); pindex++ {
			current := path[pindex]

			next := EXIT_BOTTOM
			if pindex < len(path)-1 {
				next = path[pindex+1].Entrance
			}

			pr := m.Rooms[current.Y][current.X].Clone()

			// we already know that this path works,
			// so one of these rotations is going to work
			rexit := pr.Exit(current.Entrance)
			if rexit != next {
				// cost 1
				pr.Right()
				rexit = pr.Exit(current.Entrance)
				budget--
				if rexit != next {
					// cost 1
					pr.Left()
					pr.Left()
					rexit = pr.Exit(current.Entrance)
					if rexit != next {
						// cost 2
						budget--
					}
				}
			} else {
				budget++
			}

			// Indy would be here before we can rotate all of the way
			if budget < 0 {
				valid = false
				break
			}
		}

		if valid {
			result = append(result, path)
		}
	}

	return result
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	// W: number of columns.
	// H: number of rows.
	var W, H int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &W, &H)

	var LINES []string

	// each line represents a series of room types for that row
	for i := 0; i < H+1; i++ {
		scanner.Scan()
		LINE := scanner.Text()
		LINES = append(LINES, LINE)
	}

	var initial_map Map = ParseMap(H, W, LINES)

	for {
		scanner.Scan()
		initial_map.IndyPosition = ParseObjectCoord(scanner.Text())

		// R: the number of rocks currently in the grid.
		var R int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &R)

		initial_map.Rocks = make([]ObjectCoord, R)

		for i := 0; i < R; i++ {
			scanner.Scan()
			initial_map.Rocks[i] = ParseObjectCoord(scanner.Text())
		}

		// fmt.Fprintln(os.Stderr, "Debug messages...")

		// One line containing on of three commands: 'X Y LEFT', 'X Y RIGHT' or 'WAIT'
		fmt.Println("WAIT")
	}
}
