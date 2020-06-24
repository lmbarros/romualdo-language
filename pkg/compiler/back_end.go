/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package compiler

import "gitlab.com/stackedboxes/romulang/pkg/bytecode"

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
