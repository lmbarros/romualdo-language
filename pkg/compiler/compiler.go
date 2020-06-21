/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package compiler

import "fmt"

// Compiler is a Romualdo compiler.
type Compiler struct {
}

// Compile compiles source.
func (c *Compiler) Compile(source string) {
	s := NewScanner(source)

	line := -1

	for {
		token := s.Token()

		if token.Line != line {
			fmt.Printf("%4d ", token.Line)
			line = token.Line
		} else {
			fmt.Print("   | ")
		}

		fmt.Printf("%18v '%v'\n", token.Kind, token.Lexeme)

		if token.Kind == TokenEOF {
			break
		}
	}

}
