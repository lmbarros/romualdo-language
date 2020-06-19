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
	chunk.Write(compiler.OpConstant, 171)
	chunk.Write(uint8(constant), 171)

	constant = chunk.AddConstant(3.4)
	chunk.Write(compiler.OpConstant, 171)
	chunk.Write(uint8(constant), 171)

	chunk.Write(compiler.OpAdd, 171)

	constant = chunk.AddConstant(5.6)
	chunk.Write(compiler.OpConstant, 171)
	chunk.Write(uint8(constant), 171)

	chunk.Write(compiler.OpDivide, 171)
	chunk.Write(compiler.OpNegate, 171)

	constant = chunk.AddConstant(2.0)
	chunk.Write(compiler.OpConstant, 171)
	chunk.Write(uint8(constant), 171)

	chunk.Write(compiler.OpPower, 171)

	chunk.Write(compiler.OpReturn, 171)

	fmt.Print(chunk.Disassemble("test chunk"))
	fmt.Print("\n\n")

	vm := compiler.NewVM()
	vm.Interpret(&chunk)
}
