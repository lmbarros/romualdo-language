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
	chunk.Write(compiler.OpReturn, 171)
	chunk.Write(compiler.OpConstant, 171)
	chunk.Write(uint8(constant), 171)

	fmt.Print(chunk.Disassemble("test chunk"))
}
