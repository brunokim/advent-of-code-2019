package main

import (
	"fmt"
	"sync"

	"brunokim.xyz/advent-of-code-2019/intcode"
)

func day7Part1Instance(phases ...int) (int, error) {
	input := 0
	for i, phase := range phases {
		amp := intcode.NewComputer(parseInput(day7Input))
		outputs, err := amp.RunWith(phase, input)
		if err != nil {
			return 0, fmt.Errorf("Amp %d: %v", i+1, err)
		}
		if len(outputs) != 1 {
			return 0, fmt.Errorf("More than 1 output generated: %v", outputs)
		}
		input = outputs[0]
	}
	return input, nil
}

func permutations(xs []int) [][]int {
	if len(xs) == 1 {
		return [][]int{[]int{xs[0]}}
	}
	var cs [][]int
	for i, x := range xs {
		rest := make([]int, len(xs)-1)
		copy(rest, xs[:i])
		copy(rest[i:], xs[i+1:])
		for _, comb := range permutations(rest) {
			c := make([]int, len(comb)+1)
			copy(c, comb)
			c[len(comb)] = x
			cs = append(cs, c)
		}
	}
	return cs
}

type pipe struct {
	id        int
	lastInput int
	ch        chan int
}

func newPipe(id int) *pipe {
	return &pipe{
		id:        id,
		lastInput: -1,
		ch:        make(chan int, 1),
	}
}

func (p *pipe) NextInt() (int, bool) {
	i, ok := <-p.ch
	if !ok {
		return 0, false
	}
	return i, true
}

func (p *pipe) PushInt(i int) {
	p.lastInput = i
	p.ch <- i
}

func day7Part2Instance(phases ...int) (int, error) {
	n := len(phases)
	pipes := make([]*pipe, n)
	amps := make([]*intcode.Computer, n)
	errs := make([]error, n)
	for i := 0; i < n; i++ {
		pipes[i] = newPipe(i + 1)
		amps[i] = intcode.NewComputer(parseInput(day7Input))
	}
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		pipes[i].PushInt(phases[i])
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			errs[i] = amps[i].Run(pipes[i], pipes[(i+1)%n])
		}(i)
	}
	pipes[0].PushInt(0)
	wg.Wait()
	for i, err := range errs {
		if err != error(nil) {
			return 0, fmt.Errorf("Amp #%d: %v", i+1, err)
		}
	}
	for _, pipe := range pipes {
		close(pipe.ch)
	}
	output := pipes[0].lastInput
	return output, nil
}

func day7() {
	fmt.Println(day7Part1Instance(4, 3, 2, 1, 0))
	fmt.Println(permutations([]int{1, 2, 3, 4}))
	maxOutput := -1
	for _, comb := range permutations([]int{0, 1, 2, 3, 4}) {
		output, err := day7Part1Instance(comb...)
		if err != nil {
			panic(err.Error())
		}
		if output > maxOutput {
			maxOutput = output
			fmt.Println("Part 1:", output, comb)
		}
	}
	fmt.Println(day7Part2Instance(9, 8, 7, 6, 5))
	for _, comb := range permutations([]int{5, 6, 7, 8, 9}) {
		output, err := day7Part2Instance(comb...)
		if err != nil {
			panic(err.Error())
		}
		if output > maxOutput {
			maxOutput = output
			fmt.Println("Part 2:", output, comb)
		}
	}
}

const day7Input = `3,8,1001,8,10,8,105,1,0,0,21,42,67,84,97,118,199,280,361,442,99999,3,9,101,4,9,9,102,5,9,9,101,2,9,9,1002,9,2,9,4,9,99,3,9,101,5,9,9,102,5,9,9,1001,9,5,9,102,3,9,9,1001,9,2,9,4,9,99,3,9,1001,9,5,9,1002,9,2,9,1001,9,5,9,4,9,99,3,9,1001,9,5,9,1002,9,3,9,4,9,99,3,9,102,4,9,9,101,4,9,9,102,2,9,9,101,3,9,9,4,9,99,3,9,102,2,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,1001,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,1002,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,101,2,9,9,4,9,99,3,9,1001,9,1,9,4,9,3,9,101,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,1002,9,2,9,4,9,3,9,101,2,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,101,1,9,9,4,9,3,9,101,2,9,9,4,9,99,3,9,101,1,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,1002,9,2,9,4,9,3,9,1002,9,2,9,4,9,3,9,1002,9,2,9,4,9,3,9,1001,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,102,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,1001,9,2,9,4,9,99,3,9,102,2,9,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,102,2,9,9,4,9,3,9,101,1,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,1002,9,2,9,4,9,99,3,9,101,1,9,9,4,9,3,9,101,1,9,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,1001,9,2,9,4,9,3,9,1002,9,2,9,4,9,3,9,101,1,9,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,1001,9,2,9,4,9,99`
