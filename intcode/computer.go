package intcode

import (
	"fmt"
)

type InstructionType int

const (
	Add           InstructionType = 1
	Mul                           = 2
	Input                         = 3
	Output                        = 4
	JumpIfNonZero                 = 5
	JumpIfZero                    = 6
	LessThan                      = 7
	Equals                        = 8
	Halt                          = 99
)

var instructionTypes = [...]InstructionType{Add, Mul, Input, Output, Halt, JumpIfNonZero, JumpIfZero, LessThan, Equals}
var instructionNames = map[InstructionType]string{
	Add:           "add",
	Mul:           "mul",
	Input:         "in",
	Output:        "out",
	JumpIfNonZero: "jinz",
	JumpIfZero:    "jiz",
	LessThan:      "<",
	Equals:        "==",
	Halt:          "halt",
}

type ParamMode int

const (
	Address ParamMode = iota
	Immediate
)

var expectedParamModes = map[InstructionType][]ParamMode{
	Add:           []ParamMode{Immediate, Immediate, Address},
	Mul:           []ParamMode{Immediate, Immediate, Address},
	Input:         []ParamMode{Address},
	Output:        []ParamMode{Immediate},
	JumpIfNonZero: []ParamMode{Immediate, Immediate},
	JumpIfZero:    []ParamMode{Immediate, Immediate},
	LessThan:      []ParamMode{Immediate, Immediate, Address},
	Equals:        []ParamMode{Immediate, Immediate, Address},
	Halt:          []ParamMode{},
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

type Computer struct {
	state              []int
	instructionPointer int
	inputs             []int
	Debug              bool
}

func NewComputer(state []int, inputs ...int) *Computer {
	return &Computer{state, 0, inputs, false}
}

func (c *Computer) consumeInput() int {
	result := c.inputs[0]
	c.inputs = c.inputs[1:]
	return result
}

var halted = fmt.Errorf("Halted")

func (c *Computer) step() error {
	ptr := c.instructionPointer
	opcode, modes, err := decodeInstruction(c.state[ptr])
	if err != nil {
		return err
	}
	numParams := len(modes)
	expectedModes := expectedParamModes[opcode]
	params := make([]int, numParams)
	for i, mode := range modes {
		value := c.state[ptr+1+i]
		if mode == Address && expectedModes[i] == Immediate {
			value = c.state[value]
		}
		if mode == Immediate && expectedModes[i] == Address {
			rawParams := c.state[ptr+1 : ptr+1+numParams]
			return fmt.Errorf("Unexpected immediate mode for param #%d @ %d (%s %v)", i+1, ptr, instructionNames[opcode], rawParams)
		}
		params[i] = value
	}
	if c.Debug {
		fmt.Printf("@%d: %s %v\n", ptr, instructionNames[opcode], params)
	}
	switch opcode {
	case Add:
		c.state[params[2]] = params[0] + params[1]
	case Mul:
		c.state[params[2]] = params[0] * params[1]
	case Input:
		c.state[params[0]] = c.consumeInput()
	case Output:
		fmt.Println(params[0])
	case JumpIfNonZero:
		if params[0] != 0 {
			c.instructionPointer = params[1]
			return nil
		}
	case JumpIfZero:
		if params[0] == 0 {
			c.instructionPointer = params[1]
			return nil
		}
	case LessThan:
		value := 0
		if params[0] < params[1] {
			value = 1
		}
		c.state[params[2]] = value
	case Equals:
		value := 0
		if params[0] == params[1] {
			value = 1
		}
		c.state[params[2]] = value
	case Halt:
		return halted
	}
	c.instructionPointer += numParams + 1
	return nil
}

func (c *Computer) Run() error {
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

