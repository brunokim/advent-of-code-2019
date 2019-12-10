package main

import (
	"fmt"
	"strings"
)

type coord struct {
	x, y int
}

type boundingBox struct {
	topLeft, bottomRight coord
}

func parseAsteroidMap(input string) []coord {
	var coords []coord
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		for j, ch := range line {
			if ch != '#' {
				continue
			}
			coords = append(coords, coord{j, i})
		}
	}
	return coords
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func add(p, v coord) coord {
	return coord{p.x + v.x, p.y + v.y}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func vector(from, to coord) coord {
	dx := to.x - from.x
	dy := to.y - from.y
	div := abs(gcd(dx, dy))
	return coord{dx / div, dy / div}
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

type coordSet map[coord]struct{}

func newCoordSet(coords ...coord) coordSet {
	m := make(coordSet)
	for _, p := range coords {
		m[p] = struct{}{}
	}
	return m
}

func (s coordSet) toList() []coord {
	coords := make([]coord, len(s))
	i := 0
	for p := range s {
		coords[i] = p
		i++
	}
	return coords
}

func calcBoundingBox(coords []coord) boundingBox {
	bbox := boundingBox{coords[0], coords[0]}
	for _, p := range coords[1:] {
		bbox.topLeft.x = min(bbox.topLeft.x, p.x)
		bbox.topLeft.y = min(bbox.topLeft.y, p.y)
		bbox.bottomRight.x = max(bbox.bottomRight.x, p.x)
		bbox.bottomRight.y = max(bbox.bottomRight.y, p.y)
	}
	return bbox
}

func withinBoundingBox(p coord, bbox boundingBox) bool {
	return p.x >= bbox.topLeft.x &&
		p.x <= bbox.bottomRight.x &&
		p.y >= bbox.topLeft.y &&
		p.y <= bbox.bottomRight.y
}

func day10Instance(input string) {
	asteroids := newCoordSet(parseAsteroidMap(input)...)
	asteroidList := asteroids.toList()
	bbox := calcBoundingBox(asteroidList)
	// fmt.Println(bbox, asteroidList)
	blockedFrom := make(map[coord]coordSet)
	for a1 := range asteroids {
		blockedFrom[a1] = newCoordSet()
		for a2 := range asteroids {
			if a1 == a2 {
				continue
			}
			v := vector(a1, a2)
			p := add(a2, v)
			for withinBoundingBox(p, bbox) {
				if _, ok := asteroids[p]; ok {
					//	fmt.Printf("%v --> %v --> %v (vector: %v)\n", a1, a2, p, v)
					blockedFrom[a1][p] = struct{}{}
				}
				p = add(p, v)
			}
		}
	}
	minBlockedPos := asteroidList[0]
	for asteroid, blocked := range blockedFrom {
		//fmt.Println(asteroid, blocked.toList())
		if len(blocked) < len(blockedFrom[minBlockedPos]) {
			minBlockedPos = asteroid
		}
	}
	fmt.Println("Part 1:", minBlockedPos, len(asteroids)-len(blockedFrom[minBlockedPos])-1)
}

func day10() {
	for i, test := range day10Tests {
		fmt.Println("Test", i)
		day10Instance(test)
	}
	day10Instance(day10Input)
}

var day10Tests = []string{
	`.#..#
.....
#####
....#
...##`,

	`......#.#.
#..#.#....
..#######.
.#.#.###..
.#..#.....
..#....#.#
#..#....#.
.##.#..###
##...#..#.
.#....####`,
	`#.#...#.#.
.###....#.
.#....#...
##.#.#.#.#
....#.#.#.
.##..###.#
..#...##..
..##....##
......#...
.####.###.`,
	`.#..#..###
####.###.#
....###.#.
..###.##.#
##.##.#.#.
....###..#
..#.#..#.#
#..#.#.###
.##...##.#
.....#.#..`,
	`.#..##.###...#######
##.############..##.
.#.######.########.#
.###.#######.####.#.
#####.##.#.##.###.##
..#####..#.#########
####################
#.####....###.#.#.##
##.#################
#####.##.###..####..
..######..##.#######
####.##.####...##..#
.#####..#.######.###
##...#.##########...
#.##########.#######
.####.#.###.###.#.##
....##.##.###..#####
.#.#.###########.###
#.#.#.#####.####.###
###.##.####.##.#..##`,
}

const day10Input = `##.###.#.......#.#....#....#..........#.
....#..#..#.....#.##.............#......
...#.#..###..#..#.....#........#......#.
#......#.....#.##.#.##.##...#...#......#
.............#....#.....#.#......#.#....
..##.....#..#..#.#.#....##.......#.....#
.#........#...#...#.#.....#.....#.#..#.#
...#...........#....#..#.#..#...##.#.#..
#.##.#.#...#..#...........#..........#..
........#.#..#..##.#.##......##.........
................#.##.#....##.......#....
#............#.........###...#...#.....#
#....#..#....##.#....#...#.....#......#.
.........#...#.#....#.#.....#...#...#...
.............###.....#.#...##...........
...#...#.......#....#.#...#....#...#....
.....#..#...#.#.........##....#...#.....
....##.........#......#...#...#....#..#.
#...#..#..#.#...##.#..#.............#.##
.....#...##..#....#.#.##..##.....#....#.
..#....#..#........#.#.......#.##..###..
...#....#..#.#.#........##..#..#..##....
.......#.##.....#.#.....#...#...........
........#.......#.#...........#..###..##
...#.....#..#.#.......##.###.###...#....
...............#..#....#.#....#....#.#..
#......#...#.....#.#........##.##.#.....
###.......#............#....#..#.#......
..###.#.#....##..#.......#.............#
##.#.#...#.#..........##.#..#...##......
..#......#..........#.#..#....##........
......##.##.#....#....#..........#...#..
#.#..#..#.#...........#..#.......#..#.#.
#.....#.#.........#............#.#..##.#
.....##....#.##....#.....#..##....#..#..
.#.......#......#.......#....#....#..#..
...#........#.#.##..#.#..#..#........#..
#........#.#......#..###....##..#......#
...#....#...#.....#.....#.##.#..#...#...
#.#.....##....#...........#.....#...#...`
