package main

import "fmt"
import "strings"

func parseInput(s string) []int {
	list := make([]int, len(s))
	for i, ch := range s {
		list[i] = int(ch - '0')
	}
	return list
}

type patternParams struct{ n, size int }

var cache = map[patternParams][]int{}

func pattern(n, size int) []int {
	if result, ok := cache[patternParams{n, size}]; ok {
		return result
	}
	xs := make([]int, size+1)
	state := 0
	stateValue := [4]int{0, 1, 0, -1}
	for i := 0; i < size+1; i += n {
		for j := i; j < i+n && j < size+1; j++ {
			xs[j] = stateValue[state]
		}
		state = (state + 1) % 4
	}
	result := xs[1:]
	cache[patternParams{n, size}] = result
	return result
}

func dot(xs, ys []int) int {
	var s int
	for i := 0; i < len(xs); i++ {
		s += xs[i] * ys[i]
	}
	return s
}

func phase(xs []int) []int {
	size := len(xs)
	result := make([]int, size)
	for n := 1; n < size+1; n++ {
		pat := pattern(n, size)
		result[n-1] = abs(dot(xs, pat)) % 10
	}
	return result
}

func day16Part1(input string) []int {
	xs := parseInput(input)
	for i := 0; i < 100; i++ {
		xs = phase(xs)
	}
	return xs[:8]
}

func times(n int, s string) string {
	b := new(strings.Builder)
	for i := 0; i < n; i++ {
		b.WriteString(s)
	}
	return b.String()
}

func toNum(xs []int) int {
	s := 0
	for _, x := range xs {
		s = 10*s + x
	}
	return s
}

func day16Part2(input string) []int {
	s := times(10000, input)
	xs := parseInput(s)
	for i := 0; i < 100; i++ {
		xs = phase(xs)
	}
	offset := toNum(xs[:8])
	return xs[offset : offset+8]
}

func day16() {
	fmt.Println(day16Part1("12345678"))
	fmt.Println(day16Part1("80871224585914546619083218645595"))
	fmt.Println(day16Part1("19617804207202209144916044189917"))
	fmt.Println("Part 1:", day16Part1(day16Input))
	fmt.Println(day16Part2("03036732577212944063491565474664"))
}

const day16Input = `59756772370948995765943195844952640015210703313486295362653878290009098923609769261473534009395188480864325959786470084762607666312503091505466258796062230652769633818282653497853018108281567627899722548602257463608530331299936274116326038606007040084159138769832784921878333830514041948066594667152593945159170816779820264758715101494739244533095696039336070510975612190417391067896410262310835830006544632083421447385542256916141256383813360662952845638955872442636455511906111157861890394133454959320174572270568292972621253460895625862616228998147301670850340831993043617316938748361984714845874270986989103792418940945322846146634931990046966552`
