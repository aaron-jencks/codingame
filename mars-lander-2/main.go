package main

import (
	"fmt"
	"math"
)

const (
	Y_MARGIN     = 20
	SPEED_MARGIN = 5
	MAX_DY       = 40
	MAX_DX       = 20
	GRAVITY      = 3.711
)

type lander struct {
	x       int
	y       int
	dx      int
	dy      int
	fuel    int
	angle   int
	power   int
	targetL int
	targetR int
	targetY int
}

func (l lander) isOverTarget() bool {
	return l.targetL <= l.x && l.x <= l.targetR
}

func (l lander) isFinishing() bool {
	return l.y < l.targetY+Y_MARGIN
}

func (l lander) hasSafeSpeed() bool {
	return math.Abs(float64(l.dx)) <= MAX_DX-SPEED_MARGIN &&
		math.Abs(float64(l.dy)) <= MAX_DY-SPEED_MARGIN
}

func (l lander) goesInWrongDirection() bool {
	return l.x < l.targetL && l.dx < 0 || l.targetR < l.x && l.dx > 0
}

func (l lander) goesTooFastHorizontally() bool {
	return math.Abs(float64(l.dx)) > 4*MAX_DX
}

func (l lander) goesTooSlowHorizontally() bool {
	return math.Abs(float64(l.dx)) < 2*MAX_DX
}

func (l lander) angleToSlow() int {
	speed := math.Sqrt(float64(l.dx*l.dx + l.dy*l.dy))
	return int(math.Asin(float64(l.dx)/speed) * 180 / math.Pi)
}

func (l lander) angleToAimTarget() int {
	angle := int(math.Acos(GRAVITY/4.) * 180 / math.Pi)
	if l.x < l.targetL {
		return -angle
	} else if l.targetR < l.x {
		return angle
	}
	return 0
}

func (l lander) powerToHover() int {
	if l.dy >= 0 {
		return 3
	}
	return 4
}

func main() {
	ship := lander{}

	// surfaceN: the number of points used to draw the surface of Mars.
	var surfaceN int
	fmt.Scan(&surfaceN)

	var prevX, prevY int = -1, -1
	for i := 0; i < surfaceN; i++ {
		// landX: X coordinate of a surface point. (0 to 6999)
		// landY: Y coordinate of a surface point. By linking all the points together in a sequential fashion, you form the surface of Mars.
		var landX, landY int
		fmt.Scan(&landX, &landY)
		if landY == prevY {
			ship.targetL = prevX
			ship.targetR = landX
			ship.targetY = landY
		} else {
			prevX = landX
			prevY = landY
		}
	}

	for {
		fmt.Scan(&ship.x, &ship.x, &ship.dx, &ship.dy, &ship.fuel, &ship.angle, &ship.power)

		if !ship.isOverTarget() {
			if ship.goesInWrongDirection() || ship.goesTooFastHorizontally() {
				fmt.Printf("%d 4\n", ship.angleToSlow())
			} else if ship.goesTooSlowHorizontally() {
				fmt.Printf("%d 4\n", ship.angleToAimTarget())
			} else {
				fmt.Printf("0 %d\n", ship.powerToHover())
			}
		} else {
			if ship.isFinishing() {
				fmt.Println("0 3")
			} else if ship.hasSafeSpeed() {
				fmt.Println("0 2")
			} else {
				fmt.Printf("%d 4\n", ship.angleToSlow())
			}
		}
	}
}
