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
	INDY_INVALID
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
	return INDY_INVALID
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
	case ROOM_LB:
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
	X         int
	Y         int
	Entrance  int
	Temporary bool
}

type Map struct {
	Height       int
	Width        int
	Rooms        [][]Room
	Exit         int
	IndyPosition ObjectCoord
	Rocks        map[int]ObjectCoord
}

func (m Map) Clone() Map {
	result := Map{
		Height:       m.Height,
		Width:        m.Width,
		Exit:         m.Exit,
		IndyPosition: m.IndyPosition,
		Rooms:        make([][]Room, m.Height),
		Rocks:        make(map[int]ObjectCoord),
	}

	for ri := range m.Rooms {
		result.Rooms[ri] = make([]Room, m.Width)
		copy(result.Rooms[ri], m.Rooms[ri])
	}

	for k, v := range m.Rocks {
		result.Rocks[k] = v
	}

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
	Rocks        map[int]ObjectCoord
	IndySolution []ObjectCoord
	RockSolution []map[int]ObjectCoord
}

type PathNeighbor struct {
	IndyPosition ObjectCoord
	Rocks        map[int]ObjectCoord
}

func FindRockPermutations(m Map, st PathState) []map[int]ObjectCoord {
	var currentPerms []map[int]ObjectCoord

	for ri, rock := range st.Rocks {
		if rock.Temporary {
			// this was an echo of a collision,
			// but Indy will still die if he goes here
			continue
		}

		exits := m.Rooms[rock.Y][rock.X].TheoreticalExits(rock.Entrance)

		prevPerms := make([]map[int]ObjectCoord, len(currentPerms))
		copy(prevPerms, currentPerms)
		currentPerms = nil

		for _, exit := range exits {
			nm := ObjectCoord{
				X:        rock.X,
				Y:        rock.Y,
				Entrance: FindEntranceFromExit(exit),
			}

			switch exit {
			case EXIT_BOTTOM:
				if rock.Y+1 < m.Height {
					nm.Y = rock.Y + 1
				} else {
					continue
				}
			case EXIT_LEFT:
				if rock.X-1 >= 0 {
					nm.X = rock.X - 1
				} else {
					continue
				}
			case EXIT_RIGHT:
				if rock.X+1 < m.Width {
					nm.X = rock.X + 1
				} else {
					continue
				}
			}

			if len(prevPerms) > 0 {
				for _, perm := range prevPerms {
					nperm := map[int]ObjectCoord{}

					found := false
					for k, v := range perm {
						if nm.X == v.X && nm.Y == v.Y {
							// two rocks collided and destroyed each other
							// skip copying this over,
							// the rock was destroyed
							found = true
							continue
						}

						nperm[k] = v
					}

					if found {
						// create a collision echo
						// in case indy tries to go here
						nm.Temporary = true
					}

					nperm[ri] = nm

					currentPerms = append(currentPerms, nperm)
				}
			} else {
				currentPerms = append(currentPerms, map[int]ObjectCoord{
					ri: nm,
				})
			}
		}
	}

	return currentPerms
}

func FindIndyExits(m Map, st PathState) []ObjectCoord {
	var indyExits []ObjectCoord

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
		indyExits = append(indyExits, nm)
	}

	return indyExits
}

func HasIndyRockCollision(indyLoc ObjectCoord, rockLoc map[int]ObjectCoord) bool {
	for _, rock := range rockLoc {
		if indyLoc.X == rock.X && indyLoc.Y == rock.Y {
			return true
		}
	}
	return false
}

func FindPathNeighbors(m Map, st PathState) []PathNeighbor {
	var result []PathNeighbor

	rockPermutations := FindRockPermutations(m, st)
	indyExits := FindIndyExits(m, st)

	if len(rockPermutations) > 0 {
		for _, perm := range rockPermutations {
			for _, exit := range indyExits {
				if !HasIndyRockCollision(exit, perm) {
					result = append(result, PathNeighbor{
						IndyPosition: exit,
						Rocks:        perm,
					})
				}
			}
		}
	} else {
		for _, exit := range indyExits {
			result = append(result, PathNeighbor{
				IndyPosition: exit,
			})
		}
	}

	return result
}

type MapPath struct {
	Indy  []ObjectCoord
	Rocks []map[int]ObjectCoord
}

func FindMapPath(m Map) []MapPath {
	var result []MapPath

	stack := []PathState{
		{
			IndyPosition: m.IndyPosition,
			Rocks:        m.Rocks,
		},
	}

	for len(stack) > 0 {
		st := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		for _, mv := range FindPathNeighbors(m, st) {
			ns := PathState{
				IndyPosition: mv.IndyPosition,
				Rocks:        mv.Rocks,
				IndySolution: append(st.IndySolution, st.IndyPosition),
				RockSolution: append(st.RockSolution, st.Rocks),
			}

			if HasIndyRockCollision(mv.IndyPosition, mv.Rocks) {
				// Indy collided with a rock
				continue
			}

			if mv.IndyPosition.Y == m.Height-1 && mv.IndyPosition.X == m.Exit {
				// found a solution
				result = append(result, MapPath{
					Indy:  append(ns.IndySolution, mv.IndyPosition),
					Rocks: append(ns.RockSolution, mv.Rocks),
				})
			} else {
				stack = append(stack, ns)
			}
		}
	}

	return result
}

func FindValidMapPath(m Map) []MapPath {
	var result []MapPath

	for _, path := range FindMapPath(m) {
		budget := 1
		valid := true

		// we can make at most len(path.Indy) rotations,
		// so we only need to loop that many times
		for tick := 0; tick < len(path.Indy); tick++ {
			rocks := make([][]bool, m.Height)
			for ri := range rocks {
				rocks[ri] = make([]bool, m.Width)
			}

			// Move the rocks
			if tick < len(path.Rocks) {
				for ri, rock := range path.Rocks[tick] {
					next := INDY_INVALID
					if tick < len(path.Rocks)-1 {
						if nrock, ok := path.Rocks[tick+1][ri]; ok {
							next = nrock.Entrance
						}
					}

					rr := m.Rooms[rock.Y][rock.X].Clone()

					rexit := FindEntranceFromExit(rr.Exit(rock.Entrance))
					if rexit != next {
						// we need to rotate

						// if we can't rotate, then this path is invalid
						if !rr.Rotatable {
							valid = false
							break
						}

						// cost 1
						rr.Right()
						rexit = FindEntranceFromExit(rr.Exit(rock.Entrance))
						budget--
						if rexit != next {
							// cost 1
							rr.Left()
							rr.Left()
							rexit = FindEntranceFromExit(rr.Exit(rock.Entrance))
							if rexit != next {
								// cost 2
								budget--
							}
						}
					}

					if budget < 0 {
						valid = false
						break
					}
				}
			}

			// Move Indy
			current := path.Indy[tick]

			next := INDY_TOP
			if tick < len(path.Indy)-1 {
				next = path.Indy[tick+1].Entrance
			}

			pr := m.Rooms[current.Y][current.X].Clone()

			// we already know that this path works,
			// so one of these rotations is going to work
			rexit := FindEntranceFromExit(pr.Exit(current.Entrance))
			if rexit != next {
				// we need to rotate

				// if we can't rotate, then this path is invalid
				if !pr.Rotatable {
					valid = false
					break
				}

				// cost 1
				pr.Right()
				rexit = FindEntranceFromExit(pr.Exit(current.Entrance))
				budget--
				if rexit != next {
					// cost 1
					pr.Left()
					pr.Left()
					rexit = FindEntranceFromExit(pr.Exit(current.Entrance))
					if rexit != next {
						// cost 2
						budget--
					}
				}
			}

			// Indy would be here before we can rotate all of the way
			if budget < 0 {
				valid = false
				break
			}

			if current.X == m.Exit && current.Y == m.Height-1 {
				// we reached the end of Indy's path
				break
			}

			budget++
		}

		if valid {
			result = append(result, path)
		}
	}

	return result
}

func FindNextMove(m Map) string {
	paths := FindValidMapPath(m)
	if len(paths) > 0 {
		path := paths[0]
		for pindex := 0; pindex < len(path.Indy); pindex++ {
			current := path.Indy[pindex]

			next := INDY_TOP
			if pindex < len(path.Indy)-1 {
				next = path.Indy[pindex+1].Entrance
			}

			pr := m.Rooms[current.Y][current.X].Clone()

			// we already know that this path works,
			// so one of these rotations is going to work
			rexit := FindEntranceFromExit(pr.Exit(current.Entrance))
			if rexit != next {
				// we need to rotate

				// cost 1
				pr.Right()
				rexit = FindEntranceFromExit(pr.Exit(current.Entrance))
				if rexit != next {
					// cost 1
					pr.Left()
					pr.Left()
					rexit = FindEntranceFromExit(pr.Exit(current.Entrance))
					if rexit == next {
						m.Rooms[current.Y][current.X].Left()
						return fmt.Sprintf("%d %d LEFT", current.X, current.Y)
					}
					// otherwise we turn right regardless
				}
				m.Rooms[current.Y][current.X].Right()
				return fmt.Sprintf("%d %d RIGHT", current.X, current.Y)
			}
		}
	}
	return "WAIT"
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
		fmt.Fprintln(os.Stderr, LINE)
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

		initial_map.Rocks = make(map[int]ObjectCoord)

		fmt.Fprintln(os.Stderr, "Rocks:")
		for i := 0; i < R; i++ {
			scanner.Scan()
			rtext := scanner.Text()
			fmt.Fprintln(os.Stderr, rtext)
			initial_map.Rocks[i] = ParseObjectCoord(rtext)
		}

		// fmt.Fprintln(os.Stderr, "Debug messages...")

		// One line containing on of three commands: 'X Y LEFT', 'X Y RIGHT' or 'WAIT'
		fmt.Println(FindNextMove(initial_map))
	}
}
