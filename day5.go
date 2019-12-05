package main

import (
	"fmt"
	"strconv"
	"strings"
)

type InstructionType int

const (
	Add    InstructionType = 1
	Mul                    = 2
	Input                  = 3
	Output                 = 4
	Halt                   = 99
)

var instructionTypes = [...]InstructionType{Add, Mul, Input, Output, Halt}

type ParamMode int

const (
	Address ParamMode = iota
	Immediate
)

var expectedParamModes = map[InstructionType][]ParamMode{
	Add:    []ParamMode{Immediate, Immediate, Address},
	Mul:    []ParamMode{Immediate, Immediate, Address},
	Input:  []ParamMode{Address},
	Output: []ParamMode{Immediate},
	Halt:   []ParamMode{},
}

func decodeOpcode(i int) (InstructionType, error) {
	for _, instr := range instructionTypes {
		if int(instr) == i {
			return instr, nil
		}
	}
	return -1, fmt.Errorf("Unknown opcode: %d", i)
}

func decodeMode(i int) (ParamMode, error) {
	if i < 0 || i > 2 {
		return -1, fmt.Errorf("Unknown param mode: %d", i)
	}
	return ParamMode(i), nil
}

func decodeInstruction(instr int) (InstructionType, []ParamMode, error) {
	opcode, err := decodeOpcode(instr % 100)
	if err != nil {
		return -1, nil, err
	}
	modeMask := instr / 100
	numParams := len(expectedParamModes[opcode])
	modes := make([]ParamMode, numParams)
	for i := 0; i < numParams; i++ {
		mode, err := decodeMode(modeMask % 10)
		if err != nil {
			return opcode, nil, fmt.Errorf("Invalid mode at position %d: %v", i, err)
		}
		modes[i] = mode
		modeMask /= 10
	}
	return opcode, modes, nil
}

type computer struct {
	state              []int
	instructionPointer int
	inputs             []int
}

func NewComputer(state []int, inputs ...int) *computer {
	return &computer{state, 0, inputs}
}

func (c *computer) consumeInput() int {
	result := c.inputs[0]
	c.inputs = c.inputs[1:]
	return result
}

var halted = fmt.Errorf("Halted")

func (c *computer) step() error {
	opcode, modes, err := decodeInstruction(c.state[c.instructionPointer])
	if err != nil {
		return err
	}
	numParams := len(modes)
	expectedModes := expectedParamModes[opcode]
	params := make([]int, numParams)
	for i, mode := range modes {
		value := c.state[c.instructionPointer+1+i]
		if mode == Address && expectedModes[i] == Immediate {
			value = c.state[value]
		}
		if mode == Immediate && expectedModes[i] == Address {
			return fmt.Errorf("Unexpected immediate mode for param #%d @ %d", i+1, c.instructionPointer)
		}
		params[i] = value
	}
	//fmt.Println(c.instructionPointer, opcode, params)
	switch opcode {
	case Add:
		c.state[params[2]] = params[0] + params[1]
	case Mul:
		c.state[params[2]] = params[0] * params[1]
	case Input:
		c.state[params[0]] = c.consumeInput()
	case Output:
		fmt.Println(params[0])
	case Halt:
		return halted
	}
	c.instructionPointer += numParams+1
	return nil
}

func (c *computer) run() error {
	for {
		err := c.step()
		if err == nil {
			continue
		}
		if err == halted {
			return nil
		}
		return err
	}
}

func parseInput(s string) []int {
	strs := strings.Split(s, ",")
	ints := make([]int, len(strs))
	for i, s := range strs {
		v, err := strconv.Atoi(s)
		if err != nil {
			panic(err.Error())
		}
		ints[i] = v
	}
	return ints
}

func main() {
	c := NewComputer(parseInput(day5Input), 1)
	if err := c.run(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(c)
}

const day2Input = `1,95,7,3,1,1,2,3,1,3,4,3,1,5,0,3,2,1,6,19,1,19,5,23,2,13,23,27,1,10,27,31,2,6,31,35,1,9,35,39,2,10,39,43,1,43,9,47,1,47,9,51,2,10,51,55,1,55,9,59,1,59,5,63,1,63,6,67,2,6,67,71,2,10,71,75,1,75,5,79,1,9,79,83,2,83,10,87,1,87,6,91,1,13,91,95,2,10,95,99,1,99,6,103,2,13,103,107,1,107,2,111,1,111,9,0,99,2,14,0,0`
const day5Input = `3,225,1,225,6,6,1100,1,238,225,104,0,1102,79,14,225,1101,17,42,225,2,74,69,224,1001,224,-5733,224,4,224,1002,223,8,223,101,4,224,224,1,223,224,223,1002,191,83,224,1001,224,-2407,224,4,224,102,8,223,223,101,2,224,224,1,223,224,223,1101,18,64,225,1102,63,22,225,1101,31,91,225,1001,65,26,224,101,-44,224,224,4,224,102,8,223,223,101,3,224,224,1,224,223,223,101,78,13,224,101,-157,224,224,4,224,1002,223,8,223,1001,224,3,224,1,224,223,223,102,87,187,224,101,-4698,224,224,4,224,102,8,223,223,1001,224,4,224,1,223,224,223,1102,79,85,224,101,-6715,224,224,4,224,1002,223,8,223,1001,224,2,224,1,224,223,223,1101,43,46,224,101,-89,224,224,4,224,1002,223,8,223,101,1,224,224,1,223,224,223,1101,54,12,225,1102,29,54,225,1,17,217,224,101,-37,224,224,4,224,102,8,223,223,1001,224,3,224,1,223,224,223,1102,20,53,225,4,223,99,0,0,0,677,0,0,0,0,0,0,0,0,0,0,0,1105,0,99999,1105,227,247,1105,1,99999,1005,227,99999,1005,0,256,1105,1,99999,1106,227,99999,1106,0,265,1105,1,99999,1006,0,99999,1006,227,274,1105,1,99999,1105,1,280,1105,1,99999,1,225,225,225,1101,294,0,0,105,1,0,1105,1,99999,1106,0,300,1105,1,99999,1,225,225,225,1101,314,0,0,106,0,0,1105,1,99999,107,226,226,224,1002,223,2,223,1006,224,329,101,1,223,223,1108,677,226,224,1002,223,2,223,1006,224,344,101,1,223,223,7,677,226,224,102,2,223,223,1006,224,359,101,1,223,223,108,226,226,224,1002,223,2,223,1005,224,374,101,1,223,223,8,226,677,224,1002,223,2,223,1006,224,389,101,1,223,223,1108,226,226,224,102,2,223,223,1006,224,404,101,1,223,223,1007,677,677,224,1002,223,2,223,1006,224,419,101,1,223,223,8,677,677,224,1002,223,2,223,1005,224,434,1001,223,1,223,1008,226,226,224,102,2,223,223,1005,224,449,1001,223,1,223,1008,226,677,224,102,2,223,223,1006,224,464,101,1,223,223,1107,677,677,224,102,2,223,223,1006,224,479,101,1,223,223,107,677,677,224,1002,223,2,223,1005,224,494,1001,223,1,223,1107,226,677,224,1002,223,2,223,1005,224,509,101,1,223,223,1108,226,677,224,102,2,223,223,1006,224,524,101,1,223,223,7,226,226,224,1002,223,2,223,1005,224,539,101,1,223,223,108,677,677,224,1002,223,2,223,1005,224,554,101,1,223,223,8,677,226,224,1002,223,2,223,1005,224,569,1001,223,1,223,1008,677,677,224,102,2,223,223,1006,224,584,101,1,223,223,107,226,677,224,102,2,223,223,1005,224,599,1001,223,1,223,7,226,677,224,102,2,223,223,1005,224,614,101,1,223,223,1007,226,226,224,1002,223,2,223,1005,224,629,101,1,223,223,1107,677,226,224,1002,223,2,223,1006,224,644,101,1,223,223,108,226,677,224,102,2,223,223,1006,224,659,101,1,223,223,1007,677,226,224,102,2,223,223,1006,224,674,101,1,223,223,4,223,99,226`

