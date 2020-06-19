/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package compiler

import (
	"fmt"
	"io"
	"strings"
)

const (
	// OpConstant loads a constant from the constants pool.
	OpConstant uint8 = iota

	// OpNegate is the "unary minus" operator, as in -3.14.
	OpNegate

	// OpReturn is used to return values from functions.
	OpReturn
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

// Disassemble disassembles the chunk amd resturns a string representation of
// it. The chunk name (passed as name) is included in the disassembly.
func (c *Chunk) Disassemble(name string) string {
	var out strings.Builder

	fmt.Fprintf(&out, "== %v ==\n", name)

	for offset := 0; offset < len(c.Code); {
		offset = c.disassembleInstruction(&out, offset)
	}

	return out.String()
}

// disassembleInstruction disassembles the instruction at a given offset and
// returns the offset of the next instruction to disassemble. Output is written
// to out.
func (c *Chunk) disassembleInstruction(out io.Writer, offset int) int {
	fmt.Fprintf(out, "%04v ", offset)

	if offset > 0 && c.Lines[offset] == c.Lines[offset-1] {
		fmt.Fprint(out, "   | ")
	} else {
		fmt.Fprintf(out, "%4d ", c.Lines[offset])
	}

	instruction := c.Code[offset]

	switch instruction {
	case OpConstant:
		return c.disassembleConstantInstruction(out, "CONSTANT", offset)

	case OpNegate:
		return c.disassembleSimpleInstruction(out, "NEGATE", offset)

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
