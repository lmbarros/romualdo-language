/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package parser

import (
	"unicode"
	"unicode/utf8"
)

// A scanner is used to scan (tokenize) the Romualdo code.
type scanner struct {
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

// newScanner returns a new scanner that will scan source.
func newScanner(source string) *scanner {
	return &scanner{
		source: source,
		line:   1,
	}
}

// token returns the next token in the source code being scanned.
func (s *scanner) token() *token { // nolint:funlen,gocyclo
	s.skipWhitespace()

	s.start = s.current

	if s.isAtEnd() {
		return s.makeToken(tokenKindEOF)
	}

	r := s.advance()

	if unicode.IsLetter(r) {
		return s.scanIdentifier()
	}

	if unicode.IsDigit(r) {
		return s.scanNumber()
	}

	switch r {
	case '(':
		return s.makeToken(tokenKindLeftParen)
	case ')':
		return s.makeToken(tokenKindRightParen)
	case '{':
		return s.makeToken(tokenKindLeftBrace)
	case '}':
		return s.makeToken(tokenKindRightBrace)
	case '[':
		return s.makeToken(tokenKindLeftBracket)
	case ']':
		return s.makeToken(tokenKindRightBracket)
	case ',':
		return s.makeToken(tokenKindComma)
	case '.':
		return s.makeToken(tokenKindDot)
	case '-':
		return s.makeToken(tokenKindMinus)
	case '+':
		return s.makeToken(tokenKindPlus)
	case '/':
		return s.makeToken(tokenKindSlash)
	case '*':
		return s.makeToken(tokenKindStar)
	case ':':
		return s.makeToken(tokenKindColon)
	case '~':
		return s.makeToken(tokenKindTilde)
	case '@':
		return s.makeToken(tokenKindAt)
	case '^':
		return s.makeToken(tokenKindHat)
	case '!':
		if s.match('=') {
			return s.makeToken(tokenKindBangEqual)
		}
		return s.errorToken("'!' not followed by '='")
	case '=':
		if s.match('=') {
			return s.makeToken(tokenKindEqualEqual)
		}
		return s.makeToken(tokenKindEqual)
	case '<':
		if s.match('=') {
			return s.makeToken(tokenKindLessEqual)
		}
		return s.makeToken(tokenKindLess)
	case '>':
		if s.match('=') {
			return s.makeToken(tokenKindGreaterEqual)
		}
		return s.makeToken(tokenKindGreater)

	case '"':
		return s.scanString()
	}

	return s.errorToken("Unexpected character.")
}

// isAtEnd checks if the scanner reached the end of the input. Specifically,
// this means that s.current is pointing beyond the last valid index of
// s.source.
func (s *scanner) isAtEnd() bool {
	return s.current == len(s.source)
}

// makeToken returns a token of a given kind.
func (s *scanner) makeToken(kind tokenKind) *token {
	return &token{
		kind:   kind,
		lexeme: s.source[s.start:s.current],
		line:   s.line,
	}
}

// errorToken returns a new token of kind TokenError containing a given error
// message.s
func (s *scanner) errorToken(message string) *token {
	return &token{
		kind:   tokenKindError,
		lexeme: message,
		line:   s.line,
	}
}

// advance returns the next rune in the input source and advance the s.current
// index so that it points to the start of the next rune.
func (s *scanner) advance() rune {
	runeValue, width := utf8.DecodeRuneInString(s.source[s.current:])
	s.current += width

	return runeValue
}

// match checks if the next rune matches the expected one. If it does, the
// scanner consumes the rune and returns true. If not, the scanner leaves the
// rune there (not consuming it) and returns false.
func (s *scanner) match(expected rune) bool {
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

// skipWhitespace skips all whitespace and comments, leaving s.current pointing
// to the start of a non-space, non-comment rune.
func (s *scanner) skipWhitespace() {
	for {
		r, width := utf8.DecodeRuneInString(s.source[s.current:])

		switch {
		case r == '#':
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		case r == '\n':
			s.line++
			s.current += width
		case unicode.IsSpace(r):
			s.current += width
		default:
			return
		}
	}
}

// peek returns the current rune from the input without advancing the s.current
// pointer.
func (s *scanner) peek() rune {
	r, _ := utf8.DecodeRuneInString(s.source[s.current:])
	return r
}

// peekNext returns the next rune from the input (one rune past s.current)
// without the advancing s.current pointer.
func (s *scanner) peekNext() rune {
	if s.current >= len(s.source)-1 {
		return 0
	}
	_, width := utf8.DecodeRuneInString(s.source[s.current:])
	r, _ := utf8.DecodeRuneInString(s.source[s.current+width:])
	return r
}

// scanString scans a string parser.
func (s *scanner) scanString() *token {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return s.errorToken("Unterminated string.")
	}

	// The closing quote.
	s.advance()
	return s.makeToken(tokenKindStringLiteral)
}

// scanNumber scans a number token.
func (s *scanner) scanNumber() *token {
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part.
	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		// Consume the ".".
		s.advance()

		for unicode.IsDigit(s.peek()) {
			s.advance()
		}
	}

	return s.makeToken(tokenKindNumberLiteral)
}

// scanIdentifier scans an identifier token.
func (s *scanner) scanIdentifier() *token {
	for {
		r := s.peek()
		if r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) {
			s.advance()
		} else {
			break
		}
	}

	return s.makeToken(s.identifierKind())
}

// identifierKind returns the token kind corresponding to the current token
// (assumed to be either a reserved word or an identifier).
func (s *scanner) identifierKind() tokenKind { // nolint:funlen,gocognit,gocyclo
	lexeme := s.source[s.start:s.current]

	// It's fine to not decode a UTF-8 character here: we are trying to identify
	// keywords, and keywords are 100% 7-bit clean ASCII. Any non-ASCII
	// character will not match any keyword and therefore will be classified as
	// TokenIdentifier.
	switch s.source[s.start] {
	case 'a':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'l':
				return s.checkKeyword(2, "ias", tokenKindAlias)
			case 'n':
				return s.checkKeyword(2, "d", tokenKindAnd)
			}
		}
	case 'b':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'n':
				return s.checkKeyword(2, "um", tokenKindBnum)
			case 'o':
				return s.checkKeyword(2, "ol", tokenKindBool)
			case 'r':
				return s.checkKeyword(2, "eak", tokenKindBreak)
			}
		}
	case 'c':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword(2, "se", tokenKindCase)
			case 'l':
				return s.checkKeyword(2, "ass", tokenKindClass)
			case 'o':
				return s.checkKeyword(2, "ntinue", tokenKindContinue)
			}
		}
	case 'd':
		return s.checkKeyword(1, "o", tokenKindDo)
	case 'e':
		switch lexeme {
		case "else":
			return tokenKindElse
		case "elseif":
			return tokenKindElseif
		case "end":
			return tokenKindEnd
		case "enum":
			return tokenKindEnum
		}
	case 'f':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword(2, "lse", tokenKindFalse)
			case 'l':
				return s.checkKeyword(2, "oat", tokenKindFloat)
			case 'o':
				return s.checkKeyword(2, "r", tokenKindFor)
			case 'u':
				return s.checkKeyword(2, "nction", tokenKindFunction)
			}
		}
	case 'g':
		switch lexeme {
		case "gosub":
			return tokenKindGosub
		case "goto":
			return tokenKindGoto
		}
	case 'i':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'f':
				return s.checkKeyword(2, "", tokenKindIf)
			case 'n':
				switch lexeme {
				case "in":
					return tokenKindIn
				case "int":
					return tokenKindInt
				}
			}
		}
	case 'l':
		return s.checkKeyword(1, "isten", tokenKindListen)
	case 'm':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword(2, "p", tokenKindMap)
			case 'e':
				return s.checkKeyword(2, "ta", tokenKindMeta)
			}
		}
	case 'n':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'i':
				return s.checkKeyword(2, "l", tokenKindNil)
			case 'o':
				return s.checkKeyword(2, "t", tokenKindNot)
			}
		}
	case 'o':
		return s.checkKeyword(1, "r", tokenKindOr)
	case 'p':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'r':
				return s.checkKeyword(2, "int", tokenKindPrint)
			case 'a':
				return s.checkKeyword(2, "ssage", tokenKindPassage)
			}
		}
	case 'r':
		return s.checkKeyword(1, "eturn", tokenKindReturn)
	case 's':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword(2, "y", tokenKindSay)
			case 'u':
				return s.checkKeyword(2, "per", tokenKindSuper)
			case 'w':
				return s.checkKeyword(2, "itch", tokenKindSwitch)
			case 't':
				switch lexeme {
				case "struct":
					return tokenKindStruct
				case "string":
					return tokenKindString
				}
			}
		}
	case 't':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'h':
				return s.checkKeyword(2, "en", tokenKindThen)
			case 'r':
				return s.checkKeyword(2, "ue", tokenKindTrue)
			}
		}
	case 'v':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword(2, "rs", tokenKindVars)
			case 'o':
				return s.checkKeyword(2, "id", tokenKindVoid)
			}
		}
	case 'w':
		return s.checkKeyword(1, "hile", tokenKindWhile)

	}

	return tokenKindIdentifier
}

// checkKeyword checks if the current lexeme is a given keyword. It start
// checking at the start-th character, checking if it matches rest. If there is
// a match, returns kind. Otherwise, returns TokenIdentifier.
func (s *scanner) checkKeyword(start int, rest string, kind tokenKind) tokenKind {
	restLength := len(rest)
	lexemeLength := s.current - s.start
	keywordLength := start + restLength

	if lexemeLength == keywordLength && s.source[s.start+start:s.current] == rest {
		return kind
	}

	return tokenKindIdentifier
}
