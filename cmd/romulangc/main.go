/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package main

import (
	"fmt"

	"gitlab.com/stackedboxes/romulang/pkg/compiler"
)

func main() {
	var chunk compiler.Chunk

	constant := chunk.AddConstant(1.2)
	chunk.EmitOp(compiler.OpReturn)
	chunk.EmitOp(compiler.OpConstant)
	chunk.EmitByte(uint8(constant))

	fmt.Print(chunk.Disassemble("test chunk"))
}
