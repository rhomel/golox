package vm

import "fmt"

const (
	OP_CONSTANT uint8 = iota
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_NEGATE
	OP_RETURN
)

// the chunk implementation in C needs its own capacity management. But Go
// already implements this as part of the append function. So we just need to
// use Go's built-in machinary to do the same thing. In order to follow the
// book, I will add similarly named methods.
type Chunk struct {
	Code      []uint8
	Lines     []int
	Constants *ValueArray
}

func InitChunk() *Chunk {
	return &Chunk{
		Constants: InitValueArray(),
	}
}

func (c *Chunk) Write(byt uint8, line int) {
	c.Code = append(c.Code, byt)
	c.Lines = append(c.Lines, line)
}

func (c *Chunk) AddConstant(value Value) int {
	c.Constants.Write(value)
	return c.Constants.Count() - 1
}

func (c *Chunk) Count() int {
	return len(c.Code)
}

func (c *Chunk) Disassemble(name string) {
	fmt.Printf("== %s ==\n", name)
	for offset := 0; offset < c.Count(); {
		offset = c.DisassembleInstruction(offset)
	}
}

func (c *Chunk) DisassembleInstruction(offset int) int {
	fmt.Printf("%04d ", offset)
	if offset > 0 && c.Lines[offset] == c.Lines[offset-1] {
		fmt.Printf("   | ")
	} else {
		fmt.Printf("%4d ", c.Lines[offset])
	}
	instruction := c.Code[offset]
	switch instruction {
	case OP_CONSTANT:
		return constantInstruction("OP_CONSTANT", c, offset)
	case OP_ADD:
		return simpleInstruction("OP_ADD", offset)
	case OP_SUBTRACT:
		return simpleInstruction("OP_SUBTRACT", offset)
	case OP_MULTIPLY:
		return simpleInstruction("OP_MULTIPLY", offset)
	case OP_DIVIDE:
		return simpleInstruction("OP_DIVIDE", offset)
	case OP_NEGATE:
		return simpleInstruction("OP_NEGATE", offset)
	case OP_RETURN:
		return simpleInstruction("OP_RETURN", offset)
	default:
		fmt.Printf("Unknown opcode %04d\n", instruction)
		return offset + 1
	}
}

func constantInstruction(name string, chunk *Chunk, offset int) int {
	constant := chunk.Code[offset+1]
	fmt.Printf("%-16s %4d '", name, constant)
	printValue(chunk.Constants.values[constant])
	fmt.Printf("'\n")
	return offset + 2
}

func simpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1
}

func printValue(value Value) {
	fmt.Printf("%g", value.AsNumber())
}
