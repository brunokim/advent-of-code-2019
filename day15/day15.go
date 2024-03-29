package day15

import (
	"fmt"
	"strings"

	"brunokim.xyz/advent-of-code-2019/intcode"
)

type coord struct{ x, y int }

type boundingBox struct{ topLeft, bottomRight coord }

type tile int

const (
	unknown tile = iota
	empty
	wall
	target
	self
	step
)

var tileRune = map[tile]rune{
	unknown: '#',
	empty:   ' ',
	wall:    '0',
	target:  '*',
	self:    'D',
	step:    '.',
}

type direction int

const (
	north direction = iota
	south
	west
	east
)

var directionCode = map[direction]int{
	north: 1,
	south: 2,
	west:  3,
	east:  4,
}

func walk(pos coord, dir direction) coord {
	switch dir {
	case north:
		return coord{pos.x, pos.y - 1}
	case south:
		return coord{pos.x, pos.y + 1}
	case west:
		return coord{pos.x - 1, pos.y}
	case east:
		return coord{pos.x + 1, pos.y}
	default:
		panic(fmt.Sprintf("Unknown direction: %d", dir))
	}
}

type robot struct {
	pos           coord
	place         map[coord]tile
	visited       map[coord]bool
	lastDirection direction
}

func newRobot() *robot {
	origin := coord{0, 0}
	return &robot{
		pos:     origin,
		place:   map[coord]tile{origin: empty},
		visited: map[coord]bool{},
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func calcBoundingBox(coords []coord) boundingBox {
	minX, maxX, minY, maxY := coords[0].x, coords[0].x, coords[0].y, coords[0].y
	for _, pos := range coords[1:] {
		minX = min(minX, pos.x)
		maxX = max(maxX, pos.x)
		minY = min(minY, pos.y)
		maxY = max(maxY, pos.y)
	}
	return boundingBox{coord{minX, minY}, coord{maxX, maxY}}
}

func (bbox boundingBox) height() int {
	return bbox.bottomRight.y - bbox.topLeft.y + 1
}

func (bbox boundingBox) width() int {
	return bbox.bottomRight.x - bbox.topLeft.x + 1
}

func (bbox boundingBox) coordToIndices(pos coord) (i, j int) {
	return pos.y - bbox.topLeft.y, pos.x - bbox.topLeft.x
}

func (bbox boundingBox) indicesToCoord(i, j int) coord {
	return coord{bbox.topLeft.x + j, bbox.topLeft.y + i}
}

func coordsToGrid(coords map[coord]tile) ([][]tile, boundingBox) {
	var coordList []coord
	for pos := range coords {
		coordList = append(coordList, pos)
	}
	bbox := calcBoundingBox(coordList)
	h, w := bbox.height(), bbox.width()
	grid := make([][]tile, h)
	for i := 0; i < h; i++ {
		grid[i] = make([]tile, w)
	}
	for pos, tile := range coords {
		i, j := bbox.coordToIndices(pos)
		grid[i][j] = tile
	}
	return grid, bbox
}

func printGrid(grid [][]tile) string {
	b := new(strings.Builder)
	for _, row := range grid {
		for _, ch := range row {
			b.WriteRune(tileRune[ch])
		}
		b.WriteRune('\n')
	}
	return b.String()
}

func (r *robot) String() string {
	grid, bbox := coordsToGrid(r.place)
	i, j := bbox.coordToIndices(r.pos)
	grid[i][j] = self
	return printGrid(grid)
}

func (r *robot) NextInt() (int, bool) {
	for dir, code := range directionCode {
		nextPos := walk(r.pos, dir)
		if r.place[nextPos] == unknown {
			r.lastDirection = dir
			return code, true
		}
	}
	for dir, code := range directionCode {
		nextPos := walk(r.pos, dir)
		t := r.place[nextPos]
		if (t == empty || t == target) && !r.visited[nextPos] {
			r.lastDirection = dir
			r.visited[r.pos] = true
			return code, true
		}
	}
	return 0, false
}

func (r *robot) PushInt(response int) {
	nextPos := walk(r.pos, r.lastDirection)
	switch response {
	case 0:
		r.place[nextPos] = wall
	case 1:
		r.pos = nextPos
		r.place[nextPos] = empty
	case 2:
		r.pos = nextPos
		r.place[nextPos] = target
	default:
		panic(fmt.Sprintf("Unknown response: %d", response))
	}
	fmt.Println(r)
}

func bfs(origin coord, m map[coord]tile) []coord {
	queue := []coord{origin}
	visited := map[coord]bool{}
	parent := map[coord]coord{}
	var targetPos coord
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]
		visited[pos] = true
		if m[pos] == target {
			targetPos = pos
			break
		}
		for dir := range directionCode {
			neighbor := walk(pos, dir)
			if m[neighbor] == wall {
				continue
			}
			if visited[neighbor] {
				continue
			}
			parent[neighbor] = pos
			queue = append(queue, neighbor)
		}
	}
	var path []coord
	pos := targetPos
	for pos != origin {
		path = append(path, pos)
		pos = parent[pos]
	}
	path = append(path, origin)
	return path
}

func printWithPath(m map[coord]tile, path []coord) {
	grid, bbox := coordsToGrid(m)
	for _, pos := range path {
		i, j := bbox.coordToIndices(pos)
		grid[i][j] = step
	}
	fmt.Println(printGrid(grid))
}

func Run() {
	c := intcode.NewComputer(intcode.ParseProgram(day15Input))
	r := newRobot()
	c.Run(r, r)
	path := bfs(coord{0, 0}, r.place)
	printWithPath(r.place, path)
	fmt.Println("Part 1:", len(path)-1)
	maxPath := path
	for pos, t := range r.place {
		if t != empty {
			continue
		}
		subPath := bfs(pos, r.place)
		if len(subPath) > len(maxPath) {
			maxPath = subPath
		}
	}
	printWithPath(r.place, maxPath)
	fmt.Println("Part 2:", len(maxPath)-1)
}

const day15Input = `3,1033,1008,1033,1,1032,1005,1032,31,1008,1033,2,1032,1005,1032,58,1008,1033,3,1032,1005,1032,81,1008,1033,4,1032,1005,1032,104,99,1002,1034,1,1039,102,1,1036,1041,1001,1035,-1,1040,1008,1038,0,1043,102,-1,1043,1032,1,1037,1032,1042,1105,1,124,101,0,1034,1039,102,1,1036,1041,1001,1035,1,1040,1008,1038,0,1043,1,1037,1038,1042,1106,0,124,1001,1034,-1,1039,1008,1036,0,1041,1002,1035,1,1040,1001,1038,0,1043,101,0,1037,1042,1106,0,124,1001,1034,1,1039,1008,1036,0,1041,101,0,1035,1040,102,1,1038,1043,1002,1037,1,1042,1006,1039,217,1006,1040,217,1008,1039,40,1032,1005,1032,217,1008,1040,40,1032,1005,1032,217,1008,1039,35,1032,1006,1032,165,1008,1040,9,1032,1006,1032,165,1101,0,2,1044,1105,1,224,2,1041,1043,1032,1006,1032,179,1102,1,1,1044,1105,1,224,1,1041,1043,1032,1006,1032,217,1,1042,1043,1032,1001,1032,-1,1032,1002,1032,39,1032,1,1032,1039,1032,101,-1,1032,1032,101,252,1032,211,1007,0,26,1044,1105,1,224,1101,0,0,1044,1106,0,224,1006,1044,247,102,1,1039,1034,101,0,1040,1035,102,1,1041,1036,1002,1043,1,1038,1001,1042,0,1037,4,1044,1106,0,0,22,11,19,72,14,9,6,73,82,17,41,18,83,18,49,19,12,14,39,17,20,69,20,12,48,8,8,59,36,7,33,1,15,13,10,46,96,15,2,22,80,99,12,68,99,79,22,84,16,45,25,51,4,20,95,4,51,43,13,89,2,91,48,2,46,55,24,84,8,88,10,98,46,57,15,27,7,1,19,20,63,24,50,13,63,13,59,19,13,53,75,8,20,8,44,44,21,5,11,76,9,21,2,11,27,61,6,12,72,22,40,11,9,50,18,2,38,21,78,18,13,99,9,74,5,22,30,35,5,16,34,91,55,4,19,28,42,21,62,12,74,94,16,40,2,95,54,21,2,23,56,34,9,49,47,14,39,9,65,35,53,23,25,68,15,95,25,70,27,3,33,2,31,17,40,60,24,94,34,6,99,9,92,1,92,7,49,32,8,46,47,13,37,15,11,2,15,24,8,73,8,21,64,19,74,24,5,60,9,21,47,12,12,72,18,39,90,16,6,85,13,71,19,14,24,2,65,11,51,9,19,23,34,12,9,88,77,17,6,72,19,79,39,19,21,95,87,24,91,53,7,29,20,25,11,39,38,24,72,6,1,97,15,87,11,77,64,17,57,95,9,85,19,77,8,18,97,8,39,49,4,16,81,12,36,7,7,81,22,52,56,22,47,42,4,46,75,21,19,85,37,22,90,20,10,56,24,85,55,4,91,7,22,86,1,89,13,68,35,14,27,35,9,44,79,12,42,20,16,28,89,11,57,10,60,15,13,95,3,48,24,90,86,51,18,8,71,11,80,91,5,4,93,9,80,94,9,31,7,6,90,6,57,18,19,41,69,57,8,3,42,21,16,5,79,9,13,56,99,98,19,22,85,14,35,12,21,69,16,23,3,5,78,68,2,24,12,35,36,24,93,72,12,16,7,7,19,56,8,69,45,94,18,49,44,61,21,25,19,96,7,13,27,50,76,14,5,60,4,11,90,60,9,31,85,17,11,18,74,37,20,53,53,1,42,93,66,24,10,10,73,36,19,84,14,87,71,18,64,58,3,9,70,14,10,62,81,25,19,52,5,3,78,10,66,84,84,14,66,9,19,81,8,56,11,7,39,84,31,98,22,25,56,4,12,43,78,20,19,43,88,23,10,62,90,22,38,29,5,29,32,20,14,1,3,44,13,92,79,11,59,22,77,38,3,83,18,22,37,24,32,8,19,47,20,23,32,14,72,80,24,37,33,20,8,12,17,31,20,13,51,68,65,19,31,1,1,47,88,15,31,25,94,4,11,95,87,16,77,86,92,3,2,48,39,52,62,22,63,1,70,18,61,78,14,12,50,75,10,30,2,10,96,13,58,87,9,90,3,83,5,13,28,3,67,66,21,46,10,1,70,64,8,10,50,13,22,93,3,58,13,58,2,69,1,44,2,18,22,61,61,25,36,20,7,31,6,2,7,29,2,27,22,93,16,25,8,79,93,22,2,29,27,12,56,48,34,6,40,14,13,8,14,2,8,64,32,19,18,99,22,83,83,79,16,84,58,22,88,19,31,18,35,18,31,85,20,30,16,75,16,46,16,65,16,3,44,6,2,65,97,24,40,20,25,31,88,14,66,20,13,11,76,18,43,67,13,92,47,9,81,78,20,51,12,7,43,17,24,99,14,4,89,13,84,48,13,60,13,51,23,66,7,61,19,91,17,72,64,48,10,74,13,85,8,76,11,72,3,32,22,37,80,44,18,86,50,71,5,36,21,76,23,64,23,61,40,62,24,61,0,0,21,21,1,10,1,0,0,0,0,0,0`
