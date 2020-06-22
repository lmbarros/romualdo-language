/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package compiler

import (
	"unicode"
	"unicode/utf8"
)

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
	s.skipWhitespace()

	s.start = s.current

	if s.isAtEnd() {
		return s.makeToken(TokenEOF)
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

	case '"':
		return s.scanString()
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

// skipWhitespace skips all whitespace and comments, leaving s.current pointing
// to the start of a non-space, non-comment rune.
func (s *Scanner) skipWhitespace() {
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
func (s *Scanner) peek() rune {
	r, _ := utf8.DecodeRuneInString(s.source[s.current:])
	return r
}

// peekNext returns the next rune from the input (one rune past s.current)
// without the advancing s.current pointer.
func (s *Scanner) peekNext() rune {
	if s.isAtEnd() {
		return 0
	}
	_, width := utf8.DecodeRuneInString(s.source[s.current:])
	r, _ := utf8.DecodeRuneInString(s.source[s.current+width:])
	return r
}

// scanString scans a string token.
func (s *Scanner) scanString() *Token {
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
	return s.makeToken(TokenString)
}

// scanNumber scans a number token.
func (s *Scanner) scanNumber() *Token {
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

	return s.makeToken(TokenNumberLiteral)
}

// scanIdentifier scans an identifier token.
func (s *Scanner) scanIdentifier() *Token {
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
func (s *Scanner) identifierKind() TokenKind {
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
				return s.checkKeyword(2, "ias", TokenAlias)
			case 'n':
				return s.checkKeyword(2, "d", TokenAnd)
			}
		}
	case 'b':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'n':
				return s.checkKeyword(2, "um", TokenBnum)
			case 'o':
				return s.checkKeyword(2, "ol", TokenBool)
			case 'r':
				return s.checkKeyword(2, "eak", TokenBreak)
			}
		}
	case 'c':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword(2, "se", TokenCase)
			case 'l':
				return s.checkKeyword(2, "ass", TokenClass)
			}
		}
	case 'd':
		return s.checkKeyword(1, "o", TokenDo)
	case 'e':
		switch lexeme {
		case "else":
			return TokenElse
		case "elseif":
			return TokenElseif
		case "end":
			return TokenEnd
		case "enum":
			return TokenEnum
		}
	case 'f':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword(2, "lse", TokenFalse)
			case 'l':
				return s.checkKeyword(2, "oat", TokenFloat)
			case 'o':
				return s.checkKeyword(2, "r", TokenFor)
			case 'u':
				return s.checkKeyword(2, "nction", TokenFunction)
			}
		}
	case 'g':
		if len(lexeme) > 2 {
			switch lexeme {
			case "gosub":
				return TokenGosub
			case "goto":
				return TokenGoto
			}
		}
	case 'i':
		if len(lexeme) > 2 {
			switch s.source[s.start+1] {
			case 'f':
				return s.checkKeyword(2, "", TokenIf)
			case 'n':
				switch lexeme {
				case "in":
					return TokenIn
				case "int":
					return TokenInt
				}
			}
		}
	case 'l':
		return s.checkKeyword(1, "isten", TokenListen)
	case 'm':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword(2, "p", TokenMap)
			case 'e':
				return s.checkKeyword(2, "ta", TokenMeta)
			}
		}
	case 'n':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'i':
				return s.checkKeyword(2, "l", TokenNil)
			case 'o':
				return s.checkKeyword(2, "t", TokenNot)
			}
		}
	case 'o':
		return s.checkKeyword(1, "r", TokenOr)
	case 'p':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'r':
				return s.checkKeyword(2, "int", TokenPrint)
			case 'a':
				return s.checkKeyword(2, "ssage", TokenPassage)
			}
		}
	case 'r':
		return s.checkKeyword(1, "eturn", TokenReturn)
	case 's':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword(2, "y", TokenSay)
			case 'u':
				return s.checkKeyword(2, "per", TokenSuper)
			case 'w':
				return s.checkKeyword(2, "itch", TokenSwitch)
			case 't':
				switch lexeme {
				case "struct":
					return TokenStruct
				case "string":
					return TokenString
				}
			}
		}
	case 't':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'h':
				return s.checkKeyword(2, "en", TokenThen)
			case 'r':
				return s.checkKeyword(2, "ue", TokenTrue)
			}
		}
	case 'v':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword(2, "rs", TokenVars)
			case 'o':
				return s.checkKeyword(2, "id", TokenVoid)
			}
		}
	case 'w':
		return s.checkKeyword(1, "hile", TokenWhile)

	}

	return TokenIdentifier
}

// checkKeyword checks if the current lexeme is a given keyword. It start
// checking at the start-th character, checking if it matches rest. If there is
// a match, returns kind. Otherwise, returns TokenIdentifier.
func (s *Scanner) checkKeyword(start int, rest string, kind TokenKind) TokenKind {
	restLength := len(rest)
	lexemeLength := s.current - s.start
	keywordLength := s.start + restLength

	if lexemeLength == keywordLength && s.source[s.start+start:s.current] == rest {
		return kind
	}

	return TokenIdentifier
}
