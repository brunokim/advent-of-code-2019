package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type coord3 struct {
	x, y, z int
}

func (c coord3) String() string {
	return fmt.Sprintf("<x=% 3d, y=% 3d, z=% 3d>", c.x, c.y, c.z)
}

func add3(p1, p2 coord3) coord3 {
	return coord3{p1.x + p2.x, p1.y + p2.y, p1.z + p2.z}
}

type moon struct {
	pos coord3
	vel coord3
}

func (m moon) String() string {
	return fmt.Sprintf("pos=%v, vel=%v", m.pos, m.vel)
}

func (m moon) potentialEnergy() int {
	return abs(m.pos.x) + abs(m.pos.y) + abs(m.pos.z)
}

func (m moon) kineticEnergy() int {
	return abs(m.vel.x) + abs(m.vel.y) + abs(m.vel.z)
}

func (m moon) totalEnergy() int {
	return m.potentialEnergy() * m.kineticEnergy()
}

var coord3RE = regexp.MustCompile(`<x=([0-9-]+), y=([0-9-]+), z=([0-9-]+)>`)

func mustConvert(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err.Error())
	}
	return i
}

func parsePlanetsInput(input string) []moon {
	var moons []moon
	for _, line := range strings.Split(input, "\n") {
		matches := coord3RE.FindStringSubmatch(line)
		if len(matches) != 4 {
			panic(fmt.Sprintf("%q doesn't match regex", line))
		}
		pos := coord3{mustConvert(matches[1]), mustConvert(matches[2]), mustConvert(matches[3])}
		moons = append(moons, moon{pos, coord3{}})
	}
	return moons
}

func sgn(i int) int {
	if i < 0 {
		return -1
	}
	if i > 0 {
		return +1
	}
	return 0
}

func gravity(m1, m2 moon) (f12, f21 coord3) {
	fx := sgn(m2.pos.x - m1.pos.x)
	fy := sgn(m2.pos.y - m1.pos.y)
	fz := sgn(m2.pos.z - m1.pos.z)
	return coord3{fx, fy, fz}, coord3{-fx, -fy, -fz}
}

func step(moons []moon) []moon {
	n := len(moons)
	forces := make([]coord3, n)
	for i, m1 := range moons {
		for j, m2 := range moons {
			if i <= j {
				continue
			}
			f12, f21 := gravity(m1, m2)
			forces[i] = add3(forces[i], f12)
			forces[j] = add3(forces[j], f21)
		}
	}
	result := make([]moon, n)
	for i := 0; i < n; i++ {
		vel := add3(moons[i].vel, forces[i])
		pos := add3(moons[i].pos, vel)
		result[i] = moon{pos, vel}
	}
	return result
}

const numBodies = 4

type dimension [numBodies]int

type bodies1d struct {
	pos dimension
	vel dimension
}

type solver struct {
	cache map[bodies1d]bodies1d
}

func newSolver() *solver {
	return &solver{cache: map[bodies1d]bodies1d{}}
}

func (s *solver) next(bs bodies1d) (bodies1d, bool) {
	if result, ok := s.cache[bs]; ok {
		return result, true
	}
	var forces dimension
	for i, v1 := range bs.pos {
		for j, v2 := range bs.pos {
			if i <= j {
				continue
			}
			f := sgn(v2 - v1)
			forces[i] += f
			forces[j] -= f
		}
	}
	var result bodies1d
	for i := 0; i < numBodies; i++ {
		result.vel[i] = bs.vel[i] + forces[i]
		result.pos[i] = bs.pos[i] + result.vel[i]
	}
	s.cache[bs] = result
	return result, false
}

type system struct {
	x, y, z bodies1d
}

func newSystem(moons []moon) system {
	s := system{}
	for i, m := range moons {
		s.x.pos[i], s.x.vel[i] = m.pos.x, m.vel.x
		s.y.pos[i], s.y.vel[i] = m.pos.y, m.vel.y
		s.z.pos[i], s.z.vel[i] = m.pos.z, m.vel.z
	}
	return s
}

func (s system) toMoons() []moon {
	moons := make([]moon, numBodies)
	for i := 0; i < numBodies; i++ {
		moons[i].pos.x, moons[i].vel.x = s.x.pos[i], s.x.vel[i]
		moons[i].pos.y, moons[i].vel.y = s.y.pos[i], s.y.vel[i]
		moons[i].pos.z, moons[i].vel.z = s.z.pos[i], s.z.vel[i]
	}
	return moons
}

func (s *solver) step(sys system) system {
	x, _ := s.next(sys.x)
	y, _ := s.next(sys.y)
	z, _ := s.next(sys.z)
	return system{x, y, z}
}

func measureRecurrence(bs bodies1d) int {
	s := newSolver()
	cached := false
	i := 0
	for !cached {
		bs, cached = s.next(bs)
		i++
	}
	return i - 1
}

func day12() {
	moons := parsePlanetsInput(day12Input)
	output := func(moons []moon) {
		for _, m := range moons {
			fmt.Println(m)
		}
	}
	for i := 0; i < 1000; i++ {
		moons = step(moons)
	}
	output(moons)
	var energy int
	for _, m := range moons {
		energy += m.totalEnergy()
	}
	fmt.Println("Total energy:", energy)
	// Part 2 validation
	solver := newSolver()
	sys := newSystem(parsePlanetsInput(day12Input))
	for i := 0; i < 1000; i++ {
		sys = solver.step(sys)
	}
	output(sys.toMoons())
	// Recurrence testing
	dim := bodies1d{
		pos: dimension{0, 0, 1, 1},
		vel: dimension{0, 0, 0, 0},
	}
	recur := measureRecurrence(dim)
	fmt.Println("0 -", dim)
	for i := 0; i < recur; i++ {
		dim, _ = solver.next(dim)
		fmt.Println(i+1, "-", dim)
	}
	// Part 2
	sys = newSystem(parsePlanetsInput(day12Input))
	xPeriod := measureRecurrence(sys.x)
	yPeriod := measureRecurrence(sys.y)
	zPeriod := measureRecurrence(sys.z)
	fmt.Println(xPeriod, yPeriod, zPeriod)

	lcm := xPeriod * (yPeriod / gcd(xPeriod, yPeriod))
	lcm = zPeriod * (lcm / gcd(zPeriod, lcm))
	fmt.Println("Part 2:", lcm)
}

const testInput = `<x=-1, y=0, z=2>
<x=2, y=-10, z=-7>
<x=4, y=-8, z=8>
<x=3, y=5, z=-1>`

const testInput2 = `<x=-8, y=-10, z=0>
<x=5, y=5, z=10>
<x=2, y=-7, z=3>
<x=9, y=-8, z=-3>`

const day12Input = `<x=19, y=-10, z=7>
<x=1, y=2, z=-3>
<x=14, y=-4, z=1>
<x=8, y=7, z=-6>`
