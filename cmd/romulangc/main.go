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

	chunk.Write(compiler.OpReturn)

	fmt.Printf(chunk.Disassemble("test chunk"))
}
