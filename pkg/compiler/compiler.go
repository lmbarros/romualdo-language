/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package compiler

import (
	"fmt"

	"gitlab.com/stackedboxes/romulang/pkg/bytecode"
	"gitlab.com/stackedboxes/romulang/pkg/scanner"
	"gitlab.com/stackedboxes/romulang/pkg/token"
)

// Compiler is a Romualdo compiler.
type Compiler struct {
}

// New returns a new Compiler.
func New() *Compiler {
	return &Compiler{}
}

// Compile compiles source.
func (c *Compiler) Compile(source string) *bytecode.Chunk {
	s := scanner.New(source)

	line := -1

	for {
		tok := s.Token()

		if tok.Line != line {
			fmt.Printf("%4d ", tok.Line)
			line = tok.Line
		} else {
			fmt.Print("   | ")
		}

		fmt.Printf("%18v '%v'\n", tok.Kind, tok.Lexeme)

		if tok.Kind == token.KindEOF {
			break
		}
	}

	// TODO: Dummy return
	return &bytecode.Chunk{}
}
