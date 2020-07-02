/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package compiler

import (
	"math"

	"gitlab.com/stackedboxes/romulang/pkg/bytecode"
)

// currentChunk returns the current chunk we are compiling into.
func (c *Compiler) currentChunk() *bytecode.Chunk {
	return c.chunk
}

// emitBytes writes one or more bytes to the bytecode chunk being generated.
func (c *Compiler) emitBytes(bytes ...byte) {
	for _, b := range bytes {
		c.currentChunk().Write(b, c.p.previous.Line)
	}
}

// emitReturn wraps up the compilation.
func (c *Compiler) emitReturn() {
	c.emitBytes(bytecode.OpReturn)
}

// emitConstant emits the bytecode for a constant having a given value.
func (c *Compiler) emitConstant(value bytecode.Value) {
	c.emitBytes(bytecode.OpConstant, c.makeConstant(value))
}

// makeConstant adds value to the pool of constants and returns the index in
// which it was added. If there is already a constant with this value, its index
// is returned (hey, we don't need duplicate constants, right? They are
// constant, after all!)
func (c *Compiler) makeConstant(value bytecode.Value) byte {
	// TODO: Support a more reasonable number of constants.
	if i := c.currentChunk().SearchConstant(value); i >= 0 {
		return byte(i)
	}

	constant := c.currentChunk().AddConstant(value)
	if constant > math.MaxUint8 {
		c.error("Too many constants in one chunk.")
		return 0
	}

	return byte(constant)
}
