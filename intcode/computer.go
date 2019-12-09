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
	OffsetRelBase                 = 9
	Halt                          = 99
)

var instructionTypes = [...]InstructionType{Add, Mul, Input, Output, Halt, JumpIfNonZero, JumpIfZero, LessThan, Equals, OffsetRelBase}
var instructionNames = map[InstructionType]string{
	Add:           "add",
	Mul:           "mul",
	Input:         "in",
	Output:        "out",
	JumpIfNonZero: "jinz",
	JumpIfZero:    "jiz",
	LessThan:      "<",
	Equals:        "==",
	OffsetRelBase: "base",
	Halt:          "halt",
}

type ParamMode int

const (
	Address ParamMode = iota
	Immediate
	Relative
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
	OffsetRelBase: []ParamMode{Immediate},
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
	state              map[int]int
	instructionPointer int
	relativeBase       int
	Debug              bool
}

func NewComputer(program []int) *Computer {
	state := make(map[int]int)
	for i, instruction := range program {
		state[i] = instruction
	}
	return &Computer{
		state:              state,
		instructionPointer: 0,
		relativeBase:       0,
		Debug:              false,
	}
}

var halted = fmt.Errorf("Halted")

type IntReader interface {
	NextInt() (int, bool)
}

type IntWriter interface {
	PushInt(i int)
}

func (c *Computer) debugInstructions(numParams int) []int {
	ptr := c.instructionPointer
	rawParams := make([]int, numParams+1)
	rawParams[0] = c.state[ptr]
	for i := 1; i < numParams+1; i++ {
		rawParams[i] = c.state[ptr+i]
	}
	return rawParams
}

func (c *Computer) step(in IntReader, out IntWriter) error {
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
		expectedMode := expectedModes[i]
		switch mode {
		case Address:
			if expectedMode == Immediate {
				value = c.state[value]
			}
		case Relative:
			switch expectedMode {
			case Address:
				value += c.relativeBase
			case Immediate:
				value = c.state[value+c.relativeBase]
			}
		case Immediate:
			if expectedMode == Address {
				return fmt.Errorf("Unexpected immediate mode for param #%d @ %d (%s %v)", i+1, ptr, instructionNames[opcode], c.debugInstructions(numParams))
			}
		}
		params[i] = value
	}
	if c.Debug {
		fmt.Printf("@%d: %s %v\t(%v/%d)\n", ptr, instructionNames[opcode], params, c.debugInstructions(numParams), c.relativeBase)
	}
	switch opcode {
	case Add:
		c.state[params[2]] = params[0] + params[1]
	case Mul:
		c.state[params[2]] = params[0] * params[1]
	case Input:
		v, ok := in.NextInt()
		if !ok {
			return fmt.Errorf("Input exhausted @ %d", ptr)
		}
		c.state[params[0]] = v
	case Output:
		out.PushInt(params[0])
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
	case OffsetRelBase:
		c.relativeBase += params[0]
	case Halt:
		return halted
	}
	c.instructionPointer += numParams + 1
	return nil
}

func (c *Computer) Run(in IntReader, out IntWriter) error {
	for {
		err := c.step(in, out)
		if err == nil {
			continue
		}
		if err == halted {
			return nil
		}
		return err
	}
}

type inout struct {
	input  []int
	output []int
}

func (s *inout) NextInt() (int, bool) {
	if len(s.input) == 0 {
		return 0, false
	}
	v := s.input[0]
	s.input = s.input[1:]
	return v, true
}

func (s *inout) PushInt(i int) {
	s.output = append(s.output, i)
}

func (c *Computer) RunWith(input ...int) ([]int, error) {
	s := &inout{input, nil}
	if err := c.Run(s, s); err != nil {
		return nil, err
	}
	return s.output, nil
}
