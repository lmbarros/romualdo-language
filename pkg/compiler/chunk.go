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

// An Opcode is a code representing one of the instructions (or operations) in
// the RomualdoVM.
type Opcode uint8

const (
	// The OpReturn instruction is used to return values from functions.
	OpReturn Opcode = iota
)

// A Chunk is a chunk of bytecode.
type Chunk struct {
	// The code itself.
	Code []Opcode
}

// Write writes an instructions to the chunk.
func (c *Chunk) Write(opcode Opcode) {
	c.Code = append(c.Code, opcode)
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

	instruction := c.Code[offset]
	switch instruction {
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
// parameters).
func (c *Chunk) disassembleSimpleInstruction(out io.Writer, name string, offset int) int {
	fmt.Fprintf(out, "%v\n", name)
	return offset + 1
}
