/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package compiler

import (
	"gitlab.com/stackedboxes/romulang/pkg/bytecode"
	"gitlab.com/stackedboxes/romulang/pkg/scanner"
	"gitlab.com/stackedboxes/romulang/pkg/token"
)

// parser holds some parsing-related data. I'd say it's not really a parser.
type parser struct {
	current   *token.Token // The current token.
	previous  *token.Token // The previous token.
	hadError  bool         // Did we find at least one error?
	panicMode bool         // Are we in panic mode? (Parsing panic, nothing to do with Go panic!)
}

// Compiler is a Romualdo compiler.
type Compiler struct {
	p     *parser
	s     *scanner.Scanner
	chunk *bytecode.Chunk // The chunk the compiler is generating.
}

// New returns a new Compiler.
func New() *Compiler {
	return &Compiler{
		p:     &parser{},
		chunk: &bytecode.Chunk{},
	}
}

// Compile compiles source and returns the chunk with the compiled bytecode. In
// case of errors, returns nil.
func (c *Compiler) Compile(source string) *bytecode.Chunk {
	// TODO: Candidate to change: I'd like to instantiate the scanner on New().
	c.s = scanner.New(source)

	c.advance()
	c.expression()
	c.consume(token.KindEOF, "Expect end of expression.")

	if c.p.hadError {
		return nil
	} else {
		return c.chunk
	}
}
