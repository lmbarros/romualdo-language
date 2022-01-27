/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2022 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package bytecode

import (
	"encoding/binary"
	"math"
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

// A Chunk is a chunk of bytecode.
type Chunk struct {
	// The code itself.
	Code []uint8
}

// Write writes a byte to the chunk. line is the source code line number that
// generated this byte.
func (c *Chunk) Write(b uint8) {
	c.Code = append(c.Code, b)
}

// Decodes the first four bytes in bytecode into an unsigned 31-bit integer.
// Panics if the value read does not fit into 31 bits.
func DecodeUInt31(bytecode []byte) int {
	v := binary.LittleEndian.Uint32(bytecode)
	if v > math.MaxInt32 {
		panic("Value does not fit into 31 bits")
	}
	return int(v)
}

// Encodes an unsigned 31-bit integer into the four first bytes of bytecode.
// Panics if v does not fit into 31 bits.
func EncodeUInt31(bytecode []byte, v int) {
	if v < 0 || v > math.MaxInt32 {
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
	if v < math.MinInt32 || v > math.MaxInt32 {
		panic("Value does not fit into 32 bits")
	}
	binary.LittleEndian.PutUint32(bytecode, uint32(v))
}
