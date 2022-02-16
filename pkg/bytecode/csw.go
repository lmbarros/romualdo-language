/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2022 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package bytecode

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// CSWMagic is the "magic number" identifying a Romualdo Compiled Storyworld. It
// is comprised of the "RmldCSW" string followed by a SUB character (which in
// times long gone used to represent a "soft end-of-file").
var CSWMagic = []byte{0x52, 0x6D, 0x6C, 0x64, 0x43, 0x53, 0x57, 0x1A}

// CSWVersion is the current version of a Romualdo Compiled Storyworld.
const CSWVersion byte = 0

// GlobalVar represents a global variable.
type GlobalVar struct {
	// Name is the name of the global variable.
	Name string

	// Value is the value of the global variable.
	Value Value
}

// CompiledStoryworld is a compiled, binary version of a Romualdo Language
// "program". It can be serialized and deserialized. All serialized data is
// little endian.
//
// You probably don't create one manually. Instead, either run the compiler to
// generate one from source code, or read one from a file.
//
// Overall file format:
//
// - Magic
//
// - 32-bit version (currently 0)
//
// - 32-bit size (binary data size in bytes)
//
// - 32-bit CRC32 of the binary data (using the IEEE polynomial)
//
// - Binary data
type CompiledStoryworld struct {
	// Chunks is a slide with all Chunks of bytecode containing the compiled
	// data. There is one Chunk for each
	Chunks []*Chunk

	// FirstChunk indexes the element in Chunks from where the Storyworld
	// execution starts. In other words, it points to the "main" chunk.
	FirstChunk int

	// Globals contains all the global variables.
	Globals []GlobalVar

	// The constant values used in all Chunks.
	Constants []Value

	// Strings contains all the strings used in all Chunks.
	Strings *StringInterner
}

// NewCompiledStoryworld creates a new CompiledStoryworld. Goes without saying.
func NewCompiledStoryworld() *CompiledStoryworld {
	return &CompiledStoryworld{
		Strings: NewStringInterner(),
	}
}

// AddConstant adds a constant to the CompiledStoryworld and returns the index
// of the new constant into csw.Constants.
func (csw *CompiledStoryworld) AddConstant(value Value) int {
	csw.Constants = append(csw.Constants, value)
	return len(csw.Constants) - 1
}

// SearchConstant searches the constant pool for a constant with the given
// value. If found, it returns the index of this constant into csw.Constants. If
// not found, it returns a negative value.
func (csw *CompiledStoryworld) SearchConstant(value Value) int {
	for i, v := range csw.Constants {
		if ValuesEqual(value, v) {
			return i
		}
	}

	return -1
}

// ReadCompiledStoryworld deserializes a CompiledStoryworld, reading the binary
// data from r.
func ReadCompiledStoryworld(r io.Reader) (*CompiledStoryworld, error) {
	// TODO: not implemented yet. This should probably use a string interner
	// passed as parameter to allow all strings to be interned (both constants
	// and the ones created in run-time).
	return nil, errors.New("not implemented yet")
}

// WriteTo serializes a CompiledStoryworld, writing the binary data to w.
func (csw *CompiledStoryworld) WriteTo(w io.Writer) (n int64, err error) {
	// TODO: not implemented yet.
	return 0, errors.New("not implemented yet")
}

// GetGlobalIndex returns the index into csw.Globals where the global variable
// named name is stored. If no such variable exists, returns a negative value.
func (csw *CompiledStoryworld) GetGlobalIndex(name string) int {
	for i, v := range csw.Globals {
		if v.Name == name {
			return i
		}
	}

	return -1
}

// SetGlobal sets the global variable name to value, creating it if it doesn't
// exist yet. Returns a value telling if the variable was created on this call
// or not.
func (csw *CompiledStoryworld) SetGlobal(name string, value Value) bool {
	i := csw.GetGlobalIndex(name)
	if i < 0 {
		csw.Globals = append(csw.Globals, GlobalVar{Name: name, Value: value})
		return true
	}

	csw.Globals[i].Value = value
	return false
}

// Disassemble disassembles the compiled storyworld and returns a string
// representation of it. The fi argument can be nil, but in this case the
// disassembling will be less friendly.
// TODO: Make di really optional.
func (csw *CompiledStoryworld) Disassemble(di *DebugInfo) string {
	var out strings.Builder

	fmt.Fprint(&out, "== Globals ==\n")

	for _, global := range csw.Globals {
		name := global.Name
		value := global.Value.Value
		// TODO: This is showing the Go type. OK for now, but should be the Romualdo type.
		fmt.Fprintf(&out, "Global  %v '%v' (%T)\n", name, value, value)
	}

	fmt.Fprint(&out, "\n\n")

	for i, chunk := range csw.Chunks {
		fmt.Fprintf(&out, "== %v ==\n", di.ChunksNames[i])

		for offset := 0; offset < len(chunk.Code); {
			offset = csw.DisassembleInstruction(chunk, &out, offset, di.ChunksLines[i])
		}
	}

	return out.String()
}

// DisassembleInstruction disassembles the instruction at a given offset and
// returns the offset of the next instruction to disassemble. Output is written
// to out.
func (csw *CompiledStoryworld) DisassembleInstruction(chunk *Chunk, out io.Writer, offset int, lines []int) int { // nolint: gocyclo, funlen
	fmt.Fprintf(out, "%04v ", offset)

	if offset > 0 && lines[offset] == lines[offset-1] {
		fmt.Fprint(out, "   | ")
	} else {
		fmt.Fprintf(out, "%4d ", lines[offset])
	}

	instruction := chunk.Code[offset]

	switch instruction {
	case OpNop:
		return csw.disassembleSimpleInstruction(out, "NOP", offset)

	case OpConstant:
		return csw.disassembleConstantInstruction(chunk, out, "CONSTANT", offset)

	case OpConstantLong:
		return csw.disassembleConstantLongInstruction(chunk, out, "CONSTANT_LONG", offset)

	case OpTrue:
		return csw.disassembleSimpleInstruction(out, "TRUE", offset)

	case OpFalse:
		return csw.disassembleSimpleInstruction(out, "FALSE", offset)

	case OpPop:
		return csw.disassembleSimpleInstruction(out, "POP", offset)

	case OpEqual:
		return csw.disassembleSimpleInstruction(out, "EQUAL", offset)

	case OpNotEqual:
		return csw.disassembleSimpleInstruction(out, "NOT_EQUAL", offset)

	case OpGreater:
		return csw.disassembleSimpleInstruction(out, "GREATER", offset)

	case OpGreaterEqual:
		return csw.disassembleSimpleInstruction(out, "GREATER_EQUAL", offset)

	case OpLess:
		return csw.disassembleSimpleInstruction(out, "LESS", offset)

	case OpLessEqual:
		return csw.disassembleSimpleInstruction(out, "LESS_EQUAL", offset)

	case OpAdd:
		return csw.disassembleSimpleInstruction(out, "ADD", offset)

	case OpAddBNum:
		return csw.disassembleSimpleInstruction(out, "ADD_BNUM", offset)

	case OpSubtract:
		return csw.disassembleSimpleInstruction(out, "SUBTRACT", offset)

	case OpSubtractBNum:
		return csw.disassembleSimpleInstruction(out, "SUBTRACT_BNUM", offset)

	case OpMultiply:
		return csw.disassembleSimpleInstruction(out, "MULTIPLY", offset)

	case OpDivide:
		return csw.disassembleSimpleInstruction(out, "DIVIDE", offset)

	case OpPower:
		return csw.disassembleSimpleInstruction(out, "POWER", offset)

	case OpNot:
		return csw.disassembleSimpleInstruction(out, "NOT", offset)

	case OpNegate:
		return csw.disassembleSimpleInstruction(out, "NEGATE", offset)

	case OpBlend:
		return csw.disassembleSimpleInstruction(out, "BLEND", offset)

	case OpJump:
		return csw.disassembleSByteInstruction(chunk, out, "JUMP", offset)

	case OpJumpLong:
		return csw.disassembleSIntInstruction(chunk, out, "JUMP_LONG", offset)

	case OpJumpIfFalse:
		return csw.disassembleSByteInstruction(chunk, out, "JUMP_IF_FALSE", offset)

	case OpJumpIfFalseLong:
		return csw.disassembleSIntInstruction(chunk, out, "JUMP_IF_FALSE_LONG", offset)

	case OpCall:
		return csw.disassembleUByteInstruction(chunk, out, "CALL", offset)

	case OpReturnValue:
		return csw.disassembleSimpleInstruction(out, "RETURN_VALUE", offset)

	case OpReturnVoid:
		return csw.disassembleSimpleInstruction(out, "RETURN_VOID", offset)

	case OpToInt:
		return csw.disassembleSimpleInstruction(out, "TO_INT", offset)

	case OpToFloat:
		return csw.disassembleSimpleInstruction(out, "TO_FLOAT", offset)

	case OpToBNum:
		return csw.disassembleSimpleInstruction(out, "TO_BNUM", offset)

	case OpToString:
		return csw.disassembleSimpleInstruction(out, "TO_STRING", offset)

	case OpPrint:
		return csw.disassembleSimpleInstruction(out, "PRINT", offset)

	case OpReadGlobal:
		return csw.disassembleGlobalInstruction(chunk, out, "READ_GLOBAL", offset)

	case OpWriteGlobal:
		return csw.disassembleGlobalInstruction(chunk, out, "WRITE_GLOBAL", offset)

	case OpReadLocal:
		return csw.disassembleUByteInstruction(chunk, out, "READ_LOCAL", offset)

	case OpWriteLocal:
		return csw.disassembleUByteInstruction(chunk, out, "WRITE_LOCAL", offset)

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
func (csw *CompiledStoryworld) disassembleSimpleInstruction(out io.Writer, name string, offset int) int {
	fmt.Fprintf(out, "%v\n", name)
	return offset + 1
}

// disassembleConstantInstruction disassembles a OpConstant instruction at a
// given offset. name is the instruction name, and the output is written to out.
// Returns the offset to the next instruction.
func (csw *CompiledStoryworld) disassembleConstantInstruction(chunk *Chunk, out io.Writer, name string, offset int) int {
	index := chunk.Code[offset+1]
	fmt.Fprintf(out, "%-16s %4d '%v'\n", name, index, csw.Constants[index])

	return offset + 2
}

// disassembleConstantLongInstruction disassembles a OpConstantLong instruction
// at a given offset. name is the instruction name, and the output is written to
// out. Returns the offset to the next instruction.
func (csw *CompiledStoryworld) disassembleConstantLongInstruction(chunk *Chunk, out io.Writer, name string, offset int) int {
	index := DecodeUInt31(chunk.Code[offset+1:])
	fmt.Fprintf(out, "%-16s %4d '%v'\n", name, index, csw.Constants[index])
	return offset + 5
}

// disassembleGlobalInstruction disassembles an OpReadGlobal or opWriteGlobal
// instruction at a given offset. name is the instruction name, and the output
// is written to out. Returns the offset to the next instruction.
func (csw *CompiledStoryworld) disassembleGlobalInstruction(chunk *Chunk, out io.Writer, name string, offset int) int {
	index := chunk.Code[offset+1]
	fmt.Fprintf(out, "%-16s %4d '%v'\n", name, index, csw.Globals[index].Name)

	return offset + 2
}

// disassembleSByteInstruction disassembles an instruction that has a signed
// byte immediate argument at a given offset. name is the instruction name, and
// the output is written to out. Returns the offset to the next instruction.
func (csw *CompiledStoryworld) disassembleSByteInstruction(chunk *Chunk, out io.Writer, name string, offset int) int {
	arg := int8(chunk.Code[offset+1])
	fmt.Fprintf(out, "%-16s %4d\n", name, arg)

	return offset + 1
}

// disassembleUByteInstruction disassembles an instruction that has an unsigned
// byte immediate argument instruction at a given offset. name is the
// instruction name, and the output is written to out. Returns the offset to the
// next instruction.
func (csw *CompiledStoryworld) disassembleUByteInstruction(chunk *Chunk, out io.Writer, name string, offset int) int {
	arg := chunk.Code[offset+1]
	fmt.Fprintf(out, "%-16s %4d\n", name, arg)

	return offset + 1
}

// disassembleSIntInstruction disassembles an instruction that has a 32-bit
// signed integer immediate argument at a given offset. name is the
// instruction name, and the output is written to out. Returns the offset to the
// next instruction.
func (csw *CompiledStoryworld) disassembleSIntInstruction(chunk *Chunk, out io.Writer, name string, offset int) int {
	arg := DecodeSInt32(chunk.Code[offset+1:])
	fmt.Fprintf(out, "%-16s %4d\n", name, arg)

	return offset + 4
}
