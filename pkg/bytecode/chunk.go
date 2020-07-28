/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package bytecode

import (
	"fmt"
	"io"
	"strings"
)

const (
	OpNop uint8 = iota
	OpConstant
	OpConstantLong
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreater
	OpGreaterEqual
	OpLess
	OpLessEqual
	OpAdd
	OpAddBNum
	OpSubtract
	OpSubtractBNum
	OpMultiply
	OpDivide
	OpPower
	OpNot
	OpNegate
	OpBlend
	OpReturn
)

const (
	// MaxConstantIndex is the maximum number of constants we can have on a
	// single chunk. This is equals to 2^24.
	//
	// To establish the relationship with bytecode, we have two opcodes for
	// reading constants: bytecode.OpConstant is the faster one and supports
	// indices between 0 and 255; when this is not enough, we also have
	// bytecode.OpConstantLong, which can deal with the whole range of supported
	// indices between: 0 to 16777215 (=2^24-1).
	MaxConstantsPerChunk = 16777216
)

// A Chunk is a chunk of bytecode.
type Chunk struct {
	// The code itself.
	Code []uint8

	// The constant values used in Code.
	Constants []Value

	// The source code line that generated each instruction. We have one entry
	// for each entry in Code. Very space-inefficient, but very simple.
	Lines []int
}

// Write writes a byte to the chunk. line is the source code line number that
// generated this byte.
func (c *Chunk) Write(b uint8, line int) {
	c.Code = append(c.Code, b)
	c.Lines = append(c.Lines, line)
}

// AddConstant adds a constant to the chunk and returns the index of the new
// constant into c.Constants.
func (c *Chunk) AddConstant(value Value) int {
	c.Constants = append(c.Constants, value)
	return len(c.Constants) - 1
}

// SearchConstant searches the constant pool for a constant with the given
// value. If found, it returns the index of this constant into c.Constants. If
// not found, it returns a negative value.
func (c *Chunk) SearchConstant(value Value) int {
	for i, v := range c.Constants {
		if ValuesEqual(value, v) {
			return i
		}
	}

	return -1
}

// Disassemble disassembles the chunk amd resturns a string representation of
// it. The chunk name (passed as name) is included in the disassembly.
func (c *Chunk) Disassemble(name string) string {
	var out strings.Builder

	fmt.Fprintf(&out, "== %v ==\n", name)

	for offset := 0; offset < len(c.Code); {
		offset = c.DisassembleInstruction(&out, offset)
	}

	return out.String()
}

// DisassembleInstruction disassembles the instruction at a given offset and
// returns the offset of the next instruction to disassemble. Output is written
// to out.
func (c *Chunk) DisassembleInstruction(out io.Writer, offset int) int { // nolint:gocyclo
	fmt.Fprintf(out, "%04v ", offset)

	if offset > 0 && c.Lines[offset] == c.Lines[offset-1] {
		fmt.Fprint(out, "   | ")
	} else {
		fmt.Fprintf(out, "%4d ", c.Lines[offset])
	}

	instruction := c.Code[offset]

	switch instruction {
	case OpNop:
		return c.disassembleConstantInstruction(out, "NOP", offset)

	case OpConstant:
		return c.disassembleConstantInstruction(out, "CONSTANT", offset)

	case OpConstantLong:
		return c.disassembleConstantLongInstruction(out, "CONSTANT_LONG", offset)

	case OpTrue:
		return c.disassembleSimpleInstruction(out, "TRUE", offset)

	case OpFalse:
		return c.disassembleSimpleInstruction(out, "FALSE", offset)

	case OpEqual:
		return c.disassembleSimpleInstruction(out, "EQUAL", offset)

	case OpNotEqual:
		return c.disassembleSimpleInstruction(out, "NOT_EQUAL", offset)

	case OpGreater:
		return c.disassembleSimpleInstruction(out, "GREATER", offset)

	case OpGreaterEqual:
		return c.disassembleSimpleInstruction(out, "GREATER_EQUAL", offset)

	case OpLess:
		return c.disassembleSimpleInstruction(out, "LESS", offset)

	case OpLessEqual:
		return c.disassembleSimpleInstruction(out, "LESS_EQUAL", offset)

	case OpAdd:
		return c.disassembleSimpleInstruction(out, "ADD", offset)

	case OpAddBNum:
		return c.disassembleSimpleInstruction(out, "ADD_BNUM", offset)

	case OpSubtract:
		return c.disassembleSimpleInstruction(out, "SUBTRACT", offset)

	case OpSubtractBNum:
		return c.disassembleSimpleInstruction(out, "SUBTRACT_BNUM", offset)

	case OpMultiply:
		return c.disassembleSimpleInstruction(out, "MULTIPLY", offset)

	case OpDivide:
		return c.disassembleSimpleInstruction(out, "DIVIDE", offset)

	case OpPower:
		return c.disassembleSimpleInstruction(out, "POWER", offset)

	case OpNot:
		return c.disassembleSimpleInstruction(out, "NOT", offset)

	case OpNegate:
		return c.disassembleSimpleInstruction(out, "NEGATE", offset)

	case OpBlend:
		return c.disassembleSimpleInstruction(out, "BLEND", offset)

	case OpReturn:
		return c.disassembleSimpleInstruction(out, "RETURN", offset)

	default:
		fmt.Fprintf(out, "Unknown opcode %d\n", instruction)
		return offset + 1
	}
}

// disassembleSimpleInstruction disassembles a simple instruction at a given
// offset. name is the instruction name, and the output is written to out.
// Returns the offset to the next instruction.
//
// A simple instruction is one composed of a single byte (just the opcode, no
// operands).
func (c *Chunk) disassembleSimpleInstruction(out io.Writer, name string, offset int) int {
	fmt.Fprintf(out, "%v\n", name)
	return offset + 1
}

// disassembleConstantInstruction disassembles a OpConstant instruction at a
// given offset. name is the instruction name, and the output is written to out.
// Returns the offset to the next instruction.
func (c *Chunk) disassembleConstantInstruction(out io.Writer, name string, offset int) int {
	index := c.Code[offset+1]
	fmt.Fprintf(out, "%-16s %4d '%v'\n", name, index, c.Constants[index])

	return offset + 2
}

// disassembleConstantLongInstruction disassembles a OpConstantLong instruction
// at a given offset. name is the instruction name, and the output is written to
// out. Returns the offset to the next instruction.
func (c *Chunk) disassembleConstantLongInstruction(out io.Writer, name string, offset int) int {
	index := ThreeBytesToInt(c.Code[offset+1], c.Code[offset+2], c.Code[offset+3])
	fmt.Fprintf(out, "%-16s %4d '%v'\n", name, index, c.Constants[index])

	return offset + 4
}

// Converts three bytes to a 24-bit unsigned integer. a is the least significant
// byte; c is the most significant byte.
func ThreeBytesToInt(a, b, c byte) int {
	return (int(c) << 16) | (int(b) << 8) | int(a)
}

// Converts a 24-bit unsigned integer to three bytes. The least significant byte
// is returned in a; the most significant byte is returned in c.
func IntToThreeBytes(v int) (a, b, c byte) {
	a = byte(v & 0x000000FF)
	b = byte((v & 0x0000FF00) >> 8)
	c = byte((v & 0x00FF0000) >> 16)
	return
}
