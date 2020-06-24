/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package compiler

import (
	"fmt"
	"os"

	"gitlab.com/stackedboxes/romulang/pkg/token"
)

// advance advances the parser by one token. This will report errors for each
// error token found; callers will only see the non-error tokens.
func (c *Compiler) advance() {
	c.p.previous = c.p.current

	for {
		c.p.current = c.s.Token()
		if c.p.current.Kind != token.KindError {
			break
		}

		c.errorAtCurrent(c.p.current.Lexeme)
	}
}

// consume consumes the current token (and advances the parser), assuming it is
// of a given kind. If it is not of this kind, reports this is an error with a
// given error message.
func (c *Compiler) consume(kind token.Kind, message string) {
	if c.p.current.Kind == kind {
		c.advance()
		return
	}

	c.errorAtCurrent(message)
}

// errorAtCurrent reports an error at the current (c.p.current) token.
func (c *Compiler) errorAtCurrent(message string) {
	c.errorAt(c.p.current, message)
}

// error reports an error at the token we just consumed (c.p.previous).
func (c *Compiler) error(message string) {
	c.errorAt(c.p.previous, message)
}

func (c *Compiler) errorAt(tok *token.Token, message string) {
	if c.p.panicMode {
		return
	}

	c.p.panicMode = true

	fmt.Fprintf(os.Stderr, "[line %v] Error", tok.Line)

	switch tok.Kind {
	case token.KindEOF:
		fmt.Fprintf(os.Stderr, " at end")
	case token.KindError:
		// Nothing.
	default:
		fmt.Fprintf(os.Stderr, " at '%v'", tok.Lexeme)
	}

	fmt.Fprintf(os.Stderr, ": %v\n", message)
	c.p.hadError = true
}
