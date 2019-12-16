package main

import (
	"fmt"
	"strings"

	"brunokim.xyz/advent-of-code-2019/intcode"
)

type robotState int

const (
	painting robotState = iota
	walking
)

type direction int

const (
	north direction = iota
	west
	south
	east
)

type rotationDirection int

const (
	counterClockwise rotationDirection = 0
	clockwise                          = 1
)

var nextClockwise = map[direction]direction{
	north: east,
	west:  north,
	south: west,
	east:  south,
}
var nextCounterClockwise = map[direction]direction{
	north: west,
	west:  south,
	south: east,
	east:  north,
}

type robot struct {
	pos    coord
	dir    direction
	panels map[coord]int
	state  robotState
}

func newRobot() *robot {
	return &robot{
		pos:    coord{0, 0},
		dir:    north,
		panels: map[coord]int{},
		state:  painting,
	}
}

var colorChar = [2]rune{'.', '#'}
var robotChar = map[direction]rune{
	north: '^',
	west:  '<',
	south: 'v',
	east:  '>',
}

func positionFromIndex(i, j int, bbox boundingBox) coord {
	return coord{
		x: bbox.topLeft.x + j,
		y: bbox.topLeft.y + i,
	}
}

func indicesFromPosition(pos coord, bbox boundingBox) (i, j int) {
	return pos.y - bbox.topLeft.y, pos.x - bbox.topLeft.x
}

func (r *robot) String() string {
	coords := make([]coord, len(r.panels))
	var i int
	for panel := range r.panels {
		coords[i] = panel
		i++
	}
	bbox := calcBoundingBox(coords)
	width := bbox.bottomRight.x - bbox.topLeft.x + 1
	height := bbox.bottomRight.y - bbox.topLeft.y + 1
	panels := make([][]rune, height)
	for i := 0; i < height; i++ {
		panels[i] = make([]rune, width)
		for j := 0; j < width; j++ {
			pos := positionFromIndex(i, j, bbox)
			if color, ok := r.panels[pos]; ok {
				panels[i][j] = colorChar[color]
			} else {
				panels[i][j] = ' '
			}
		}
	}
	i, j := indicesFromPosition(r.pos, bbox)
	panels[i][j] = robotChar[r.dir]

	b := new(strings.Builder)
	for _, row := range panels {
		b.WriteString(string(row))
		b.WriteString("\n")
	}
	return b.String()
}

func (r *robot) newPosition() coord {
	switch r.dir {
	case north:
		return coord{r.pos.x, r.pos.y - 1}
	case west:
		return coord{r.pos.x - 1, r.pos.y}
	case south:
		return coord{r.pos.x, r.pos.y + 1}
	case east:
		return coord{r.pos.x + 1, r.pos.y}
	default:
		panic("Invalid direction: " + string(r.dir))
	}
}

func (r *robot) newDirection(i int) direction {
	switch rotationDirection(i) {
	case counterClockwise:
		return nextCounterClockwise[r.dir]
	case clockwise:
		return nextClockwise[r.dir]
	default:
		panic("Invalid rotation direction: " + string(i))
	}
}

func (r *robot) NextInt() (int, bool) {
	return r.panels[r.pos], true
}

func (r *robot) PushInt(i int) {
	switch r.state {
	case painting:
		r.panels[r.pos] = i
		r.state = walking
	case walking:
		r.dir = r.newDirection(i)
		r.pos = r.newPosition()
		r.state = painting
	default:
		panic("Invalid state: " + string(r.state))
	}
}

func day11() {
	r := newRobot()
	r.panels[r.pos] = 1
	c := intcode.NewComputer(intcode.ParseProgram(day11Input))
	c.Run(r, r)
	fmt.Println(len(r.panels))
	fmt.Println(r)
}

const day11Input = `3,8,1005,8,314,1106,0,11,0,0,0,104,1,104,0,3,8,1002,8,-1,10,1001,10,1,10,4,10,108,1,8,10,4,10,1002,8,1,28,2,2,16,10,1,1108,7,10,1006,0,10,1,5,14,10,3,8,102,-1,8,10,101,1,10,10,4,10,108,1,8,10,4,10,102,1,8,65,1006,0,59,2,109,1,10,1006,0,51,2,1003,12,10,3,8,102,-1,8,10,1001,10,1,10,4,10,108,1,8,10,4,10,1001,8,0,101,1006,0,34,1,1106,0,10,1,1101,17,10,3,8,102,-1,8,10,101,1,10,10,4,10,1008,8,0,10,4,10,1001,8,0,135,3,8,1002,8,-1,10,101,1,10,10,4,10,108,0,8,10,4,10,1001,8,0,156,3,8,1002,8,-1,10,101,1,10,10,4,10,108,0,8,10,4,10,1001,8,0,178,1,108,19,10,3,8,102,-1,8,10,101,1,10,10,4,10,108,0,8,10,4,10,1002,8,1,204,1,1006,17,10,3,8,102,-1,8,10,101,1,10,10,4,10,108,1,8,10,4,10,102,1,8,230,1006,0,67,1,103,11,10,1,1009,19,10,1,109,10,10,3,8,102,-1,8,10,101,1,10,10,4,10,1008,8,0,10,4,10,101,0,8,268,3,8,102,-1,8,10,101,1,10,10,4,10,1008,8,1,10,4,10,1002,8,1,290,2,108,13,10,101,1,9,9,1007,9,989,10,1005,10,15,99,109,636,104,0,104,1,21101,48210224024,0,1,21101,0,331,0,1105,1,435,21101,0,937264165644,1,21101,0,342,0,1105,1,435,3,10,104,0,104,1,3,10,104,0,104,0,3,10,104,0,104,1,3,10,104,0,104,1,3,10,104,0,104,0,3,10,104,0,104,1,21101,235354025051,0,1,21101,389,0,0,1105,1,435,21102,29166169280,1,1,21102,400,1,0,1105,1,435,3,10,104,0,104,0,3,10,104,0,104,0,21102,709475849060,1,1,21102,1,423,0,1106,0,435,21102,868498428684,1,1,21101,434,0,0,1105,1,435,99,109,2,21201,-1,0,1,21101,0,40,2,21102,1,466,3,21101,456,0,0,1105,1,499,109,-2,2105,1,0,0,1,0,0,1,109,2,3,10,204,-1,1001,461,462,477,4,0,1001,461,1,461,108,4,461,10,1006,10,493,1101,0,0,461,109,-2,2106,0,0,0,109,4,2102,1,-1,498,1207,-3,0,10,1006,10,516,21102,1,0,-3,21201,-3,0,1,21201,-2,0,2,21102,1,1,3,21102,535,1,0,1106,0,540,109,-4,2106,0,0,109,5,1207,-3,1,10,1006,10,563,2207,-4,-2,10,1006,10,563,21202,-4,1,-4,1106,0,631,21201,-4,0,1,21201,-3,-1,2,21202,-2,2,3,21101,582,0,0,1105,1,540,22102,1,1,-4,21102,1,1,-1,2207,-4,-2,10,1006,10,601,21101,0,0,-1,22202,-2,-1,-2,2107,0,-3,10,1006,10,623,22102,1,-1,1,21101,623,0,0,105,1,498,21202,-2,-1,-2,22201,-4,-2,-4,109,-5,2105,1,0`
