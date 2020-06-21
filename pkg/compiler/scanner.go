/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package compiler

import "unicode/utf8"

// A Scanner is used to scan (tokenize) the Romualdo code.
type Scanner struct {
	// source is the source code being scanned.
	source string

	// start points to the start of the token being currently scanned. It points
	// into source.
	start int

	// current points to the charcater we are currently looking at. It points
	// into source.
	current int

	// line holds the line number we are currently looking at.
	line int
}

// NewScanner returns a new Scanner that will scan source.
func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
		line:   1,
	}
}

// Token returns the next token in the source code being scanned.
func (s *Scanner) Token() *Token {
	s.start = s.current

	if s.isAtEnd() {
		return s.makeToken(TokenEOF)
	}

	r := s.advance()

	switch r {
	case '(':
		return s.makeToken(TokenLeftParen)
	case ')':
		return s.makeToken(TokenRightParen)
	case '{':
		return s.makeToken(TokenLeftBrace)
	case '}':
		return s.makeToken(TokenRightBrace)
	case '[':
		return s.makeToken(TokenLeftBracket)
	case ']':
		return s.makeToken(TokenRightBracket)
	case ',':
		return s.makeToken(TokenComma)
	case '.':
		return s.makeToken(TokenDot)
	case '-':
		return s.makeToken(TokenMinus)
	case '+':
		return s.makeToken(TokenPlus)
	case '/':
		return s.makeToken(TokenSlash)
	case '*':
		return s.makeToken(TokenStar)
	case ':':
		return s.makeToken(TokenColon)
	case '~':
		return s.makeToken(TokenTilde)
	case '@':
		return s.makeToken(TokenAt)
	case '^':
		return s.makeToken(TokenHat)
	case '!':
		if s.match('=') {
			return s.makeToken(TokenBangEqual)
		}
		return s.errorToken("'!' not followed by '='")
	case '=':
		if s.match('=') {
			return s.makeToken(TokenEqualEqual)
		}
		return s.makeToken(TokenEqual)
	case '<':
		if s.match('=') {
			return s.makeToken(TokenLessEqual)
		}
		return s.makeToken(TokenLess)
	case '>':
		if s.match('=') {
			return s.makeToken(TokenGreaterEqual)
		}
		return s.makeToken(TokenGreater)
	}

	return s.errorToken("Unexpected character.")
}

// isAtEnd checks if the scanner reached the end of the input.
func (s *Scanner) isAtEnd() bool {
	return s.current == len(s.source)
}

// makeToken returns a token of a given kind.
func (s *Scanner) makeToken(kind TokenKind) *Token {
	return &Token{
		Kind:   kind,
		Lexeme: s.source[s.start:s.current],
		Line:   s.line,
	}
}

// errorToken returns a new token of kind TokenError containing a given error
// message.s
func (s *Scanner) errorToken(message string) *Token {
	return &Token{
		Kind:   TokenError,
		Lexeme: message,
		Line:   s.line,
	}
}

// advance returns the next rune in the input source and advance the s.current
// index so that it points to the start of the next rune.
func (s *Scanner) advance() rune {
	runeValue, width := utf8.DecodeRuneInString(s.source[s.current:])
	s.current += width

	return runeValue
}

// match checks if the next rune matches the expected one. If it does, the
// scanner consumes the rune and returns true. If not, the scanner leaves the
// rune there (not consuming it) and returns false.
func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	currentRune, width := utf8.DecodeRuneInString(s.source[s.current:])

	if currentRune != expected {
		return false
	}

	s.current += width
	return true
}
