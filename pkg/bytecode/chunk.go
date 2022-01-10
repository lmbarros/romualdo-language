/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2021 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package bytecode

import (
	"encoding/binary"
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
	OpPop
	OpPower
	OpNot
	OpNegate
	OpBlend
	OpJump
	OpJumpLong // Must be right after OpJump
	OpJumpIfFalse
	OpJumpIfFalseLong // Must be right after OpJumpIfFalse
	OpJumpIfFalseNoPop
	OpJumpIfFalseNoPopLong // Must be right after OpJumpIfFalseNoPop
	OpJumpIfTrueNoPop
	OpJumpIfTrueNoPopLong // Must be right after OpJumpIfTrueNoPop
	OpReturn
	OpToInt
	OpToFloat
	OpToBNum
	OpToString
	OpPrint
	OpReadGlobal
	OpWriteGlobal
	OpReadLocal
	OpWriteLocal
)

const (
	// MaxConstantsPerChunk is the maximum number of constants we can have on a
	// single chunk. This is equal to 2^31, so that it fits on an int even on
	// platforms that use 32-bit ints.
	//
	// To establish the relationship of this constant with our bytecode: we have
	// two opcodes for reading constants. bytecode.OpConstant is the faster one
	// and supports indices between 0 and 255. When this is not enough, we also
	// have bytecode.OpConstantLong, which can deal with the whole range of
	// supported indices between: 0 to 2_147_483_647 (=2^31-1).
	MaxConstantsPerChunk = 2_147_483_648
)

// GlobalVar represents a global variable.
type GlobalVar struct {
	// Name is the name of the global variable.
	Name string

	// Value is the value of the global variable.
	Value Value
}

// A Chunk is a chunk of bytecode.
type Chunk struct {
	// The code itself.
	Code []uint8

	// The constant values used in Code.
	Constants []Value

	// The source code line that generated each instruction. We have one entry
	// for each entry in Code. Very space-inefficient, but very simple.
	Lines []int

	// Strings contains all the strings used in this
	Strings *StringInterner

	// Globals contains all the global variables.
	Globals []GlobalVar
}

// NewChunk creates and returns a new Chunk.
func NewChunk() *Chunk {
	return &Chunk{
		Strings: NewStringInterner(),
		Globals: []GlobalVar{},
	}
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

// GetGlobalIndex returns the index into c.Globals where the global variable
// named name is stored. If no such variable exists, returns a negative value.
func (c *Chunk) GetGlobalIndex(name string) int {
	for i, v := range c.Globals {
		if v.Name == name {
			return i
		}
	}

	return -1
}

// SetGlobal sets the global variable name to value, creating it if it doesn't
// exist yet. Returns a value telling if the variable was created on this call
// or not.
func (c *Chunk) SetGlobal(name string, value Value) bool {
	i := c.GetGlobalIndex(name)
	if i < 0 {
		c.Globals = append(c.Globals, GlobalVar{Name: name, Value: value})
		return true
	}

	c.Globals[i].Value = value
	return false
}

// Disassemble disassembles the chunk and returns a string representation of
// it. The chunk name (passed as name) is included in the disassembly.
func (c *Chunk) Disassemble(name string) string {
	var out strings.Builder

	fmt.Fprintf(&out, "== %v ==\n", name)

	for _, global := range c.Globals {
		name := global.Name
		value := global.Value.Value
		// TODO: This is showing the Go type. OK for now, but should be the Romualdo type.
		fmt.Fprintf(&out, "Global  %v '%v' (%T)\n", name, value, value)
	}

	fmt.Fprint(&out, "\n")

	for offset := 0; offset < len(c.Code); {
		offset = c.DisassembleInstruction(&out, offset)
	}

	return out.String()
}

// DisassembleInstruction disassembles the instruction at a given offset and
// returns the offset of the next instruction to disassemble. Output is written
// to out.
func (c *Chunk) DisassembleInstruction(out io.Writer, offset int) int { // nolint: gocyclo, funlen
	fmt.Fprintf(out, "%04v ", offset)

	if offset > 0 && c.Lines[offset] == c.Lines[offset-1] {
		fmt.Fprint(out, "   | ")
	} else {
		fmt.Fprintf(out, "%4d ", c.Lines[offset])
	}

	instruction := c.Code[offset]

	switch instruction {
	case OpNop:
		return c.disassembleSimpleInstruction(out, "NOP", offset)

	case OpConstant:
		return c.disassembleConstantInstruction(out, "CONSTANT", offset)

	case OpConstantLong:
		return c.disassembleConstantLongInstruction(out, "CONSTANT_LONG", offset)

	case OpTrue:
		return c.disassembleSimpleInstruction(out, "TRUE", offset)

	case OpFalse:
		return c.disassembleSimpleInstruction(out, "FALSE", offset)

	case OpPop:
		return c.disassembleSimpleInstruction(out, "POP", offset)

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

	case OpJump:
		return c.disassembleSByteInstruction(out, "JUMP", offset)

	case OpJumpLong:
		return c.disassembleSIntInstruction(out, "JUMP_LONG", offset)

	case OpJumpIfFalse:
		return c.disassembleSByteInstruction(out, "JUMP_IF_FALSE", offset)

	case OpJumpIfFalseLong:
		return c.disassembleSIntInstruction(out, "JUMP_IF_FALSE_LONG", offset)

	case OpReturn:
		return c.disassembleSimpleInstruction(out, "RETURN", offset)

	case OpToInt:
		return c.disassembleSimpleInstruction(out, "TO_INT", offset)

	case OpToFloat:
		return c.disassembleSimpleInstruction(out, "TO_FLOAT", offset)

	case OpToBNum:
		return c.disassembleSimpleInstruction(out, "TO_BNUM", offset)

	case OpToString:
		return c.disassembleSimpleInstruction(out, "TO_STRING", offset)

	case OpPrint:
		return c.disassembleSimpleInstruction(out, "PRINT", offset)

	case OpReadGlobal:
		return c.disassembleGlobalInstruction(out, "READ_GLOBAL", offset)

	case OpWriteGlobal:
		return c.disassembleGlobalInstruction(out, "WRITE_GLOBAL", offset)

	case OpReadLocal:
		return c.disassembleUByteInstruction(out, "READ_LOCAL", offset)

	case OpWriteLocal:
		return c.disassembleUByteInstruction(out, "WRITE_LOCAL", offset)

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
	index := DecodeUInt31(c.Code[offset+1:])
	fmt.Fprintf(out, "%-16s %4d '%v'\n", name, index, c.Constants[index])
	return offset + 5
}

// disassembleGlobalInstruction disassembles an OpReadGlobal or opWriteGlobal
// instruction at a given offset. name is the instruction name, and the output
// is written to out. Returns the offset to the next instruction.
func (c *Chunk) disassembleGlobalInstruction(out io.Writer, name string, offset int) int {
	index := c.Code[offset+1]
	fmt.Fprintf(out, "%-16s %4d '%v'\n", name, index, c.Globals[index].Name)

	return offset + 2
}

// disassembleSByteInstruction disassembles an instruction that has a signed
// byte immediate argument at a given offset. name is the instruction name, and
// the output is written to out. Returns the offset to the next instruction.
func (c *Chunk) disassembleSByteInstruction(out io.Writer, name string, offset int) int {
	arg := int8(c.Code[offset+1])
	fmt.Fprintf(out, "%-16s %4d\n", name, arg)

	return offset + 1
}

// disassembleUByteInstruction disassembles an instruction that has an unsigned
// byte immediate argument instruction at a given offset. name is the
// instruction name, and the output is written to out. Returns the offset to the
// next instruction.
func (c *Chunk) disassembleUByteInstruction(out io.Writer, name string, offset int) int {
	arg := c.Code[offset+1]
	fmt.Fprintf(out, "%-16s %4d\n", name, arg)

	return offset + 1
}

// disassembleSIntInstruction disassembles an instruction that has a 32-bit
// signed integer immediate argument at a given offset. name is the
// instruction name, and the output is written to out. Returns the offset to the
// next instruction.
func (c *Chunk) disassembleSIntInstruction(out io.Writer, name string, offset int) int {
	arg := DecodeSInt32(c.Code[offset+1:])
	fmt.Fprintf(out, "%-16s %4d\n", name, arg)

	return offset + 4
}

// Decodes the first four bytes in bytecode into an unsigned 31-bit integer.
// Panics if the value read does not fit into 31 bits.
func DecodeUInt31(bytecode []byte) int {
	v := binary.LittleEndian.Uint32(bytecode)
	if v > 2_147_483_647 {
		panic("Value does not fit into 31 bits")
	}
	return int(v)
}

// Encodes an unsigned 31-bit integer into the four first bytes of bytecode.
// Panics if v does not fit into 31 bits.
func EncodeUInt31(bytecode []byte, v int) {
	if v < 0 || v > 2_147_483_647 {
		panic("Value does not fit into 31 bits")
	}
	binary.LittleEndian.PutUint32(bytecode, uint32(v))
}

// Decodes the first four bytes in bytecode into a signed 32-bit integer.
func DecodeSInt32(bytecode []byte) int {
	v := binary.LittleEndian.Uint32(bytecode)
	return int(v)
}

// Encodes an signed 32-bit integer into the four first bytes of bytecode.
// Panics if v does not fit into 32 bits.
func EncodeSInt32(bytecode []byte, v int) {
	if v < -2_147_483_648 || v > 2_147_483_647 {
		panic("Value does not fit into 32 bits")
	}
	binary.LittleEndian.PutUint32(bytecode, uint32(v))
}
