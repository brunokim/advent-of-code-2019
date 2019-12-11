package main

import (
	"fmt"
	"math"
	"sort"
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

func difference(s1, s2 coordSet) coordSet {
	s := make(coordSet)
	for p := range s1 {
		if _, ok := s2[p]; ok {
			continue
		}
		s[p] = struct{}{}
	}
	return s
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

func blockedFrom(station coord, asteroids coordSet, bbox boundingBox) coordSet {
	s := newCoordSet()
	for asteroid := range asteroids {
		if asteroid == station {
			continue
		}
		v := vector(station, asteroid)
		p := add(asteroid, v)
		for withinBoundingBox(p, bbox) {
			if _, ok := asteroids[p]; ok {
				s[p] = struct{}{}
			}
			p = add(p, v)
		}
	}
	return s
}

func day10Instance(input string) []coord {
	asteroids := newCoordSet(parseAsteroidMap(input)...)
	asteroidList := asteroids.toList()
	bbox := calcBoundingBox(asteroidList)
	// fmt.Println(bbox, asteroidList)
	blocked := make(map[coord]coordSet)
	for a1 := range asteroids {
		blocked[a1] = blockedFrom(a1, asteroids, bbox)
	}
	station := asteroidList[0]
	for asteroid, invisible := range blocked {
		//fmt.Println(asteroid, blocked.toList())
		if len(invisible) < len(blocked[station]) {
			station = asteroid
		}
	}
	fmt.Println("Part 1:", station, len(asteroids)-len(blocked[station])-1)
	angle := func(p coord) float64 {
		rad := math.Atan2(float64(p.y-station.y), float64(p.x-station.x))
		deg := rad / math.Pi * 180
		ang := deg + 90
		if ang < 0 {
			ang += 360
		}
		return ang
	}
	var destroyed []coord
	remaining := make(coordSet)
	for p := range asteroids {
		if p == station {
			continue
		}
		remaining[p] = struct{}{}
	}
	for i := 0; len(remaining) > 0 && i < 100; i++ {
		blocked := blockedFrom(station, remaining, bbox)
		visible := difference(remaining, blocked)
		visibleList := visible.toList()
		sort.Slice(visibleList, func(i, j int) bool {
			a1, a2 := visibleList[i], visibleList[j]
			return angle(a1) < angle(a2)
		})
		fmt.Printf("Iteration #%d: %d of %d destroyed\n", i+1, len(visible), len(remaining))
		destroyed = append(destroyed, visibleList...)
		remaining = blocked
	}
	return destroyed
}

func day10() {
	for i, test := range day10Tests {
		fmt.Println("Test", i)
		day10Instance(test)
	}
	destroyed := day10Instance(day10Input)
	fmt.Println(destroyed[:15])
	fmt.Println(destroyed[198:201])
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
