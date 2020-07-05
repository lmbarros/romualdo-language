/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package parser

import (
	"fmt"
	"os"
)

// advance advances the parser by one token. This will report errors for each
// error token found; callers will only see the non-error tokens.
func (c *compiler) advance() {
	c.p.previous = c.p.current

	for {
		c.p.current = c.s.token()
		if c.p.current.kind != tokenKindError {
			break
		}

		c.errorAtCurrent(c.p.current.lexeme)
	}
}

// consume consumes the current token (and advances the parser), assuming it is
// of a given kind. If it is not of this kind, reports this is an error with a
// given error message.
func (c *compiler) consume(kind tokenKind, message string) {
	if c.p.current.kind == kind {
		c.advance()
		return
	}

	c.errorAtCurrent(message)
}

// errorAtCurrent reports an error at the current (c.p.current) token.
func (c *compiler) errorAtCurrent(message string) {
	c.errorAt(c.p.current, message)
}

// error reports an error at the token we just consumed (c.p.previous).
func (c *compiler) error(message string) {
	c.errorAt(c.p.previous, message)
}

// errorAt reports an error at a given token, with a given error message.
func (c *compiler) errorAt(tok *token, message string) {
	if c.p.panicMode {
		return
	}

	c.p.panicMode = true

	fmt.Fprintf(os.Stderr, "[line %v] Error", tok.line)

	switch tok.kind {
	case tokenKindEOF:
		fmt.Fprintf(os.Stderr, " at end")
	case tokenKindError:
		// Nothing.
	default:
		fmt.Fprintf(os.Stderr, " at '%v'", tok.lexeme)
	}

	fmt.Fprintf(os.Stderr, ": %v\n", message)
	c.p.hadError = true
}
