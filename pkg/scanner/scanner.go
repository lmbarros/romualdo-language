/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package scanner

import (
	"unicode"
	"unicode/utf8"

	"gitlab.com/stackedboxes/romulang/pkg/token"
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

// New returns a new Scanner that will scan source.
func New(source string) *Scanner {
	return &Scanner{
		source: source,
		line:   1,
	}
}

// Token returns the next token in the source code being scanned.
func (s *Scanner) Token() *token.Token { // nolint:funlen,gocyclo
	s.skipWhitespace()

	s.start = s.current

	if s.isAtEnd() {
		return s.makeToken(token.KindEOF)
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
		return s.makeToken(token.KindLeftParen)
	case ')':
		return s.makeToken(token.KindRightParen)
	case '{':
		return s.makeToken(token.KindLeftBrace)
	case '}':
		return s.makeToken(token.KindRightBrace)
	case '[':
		return s.makeToken(token.KindLeftBracket)
	case ']':
		return s.makeToken(token.KindRightBracket)
	case ',':
		return s.makeToken(token.KindComma)
	case '.':
		return s.makeToken(token.KindDot)
	case '-':
		return s.makeToken(token.KindMinus)
	case '+':
		return s.makeToken(token.KindPlus)
	case '/':
		return s.makeToken(token.KindSlash)
	case '*':
		return s.makeToken(token.KindStar)
	case ':':
		return s.makeToken(token.KindColon)
	case '~':
		return s.makeToken(token.KindTilde)
	case '@':
		return s.makeToken(token.KindAt)
	case '^':
		return s.makeToken(token.KindHat)
	case '!':
		if s.match('=') {
			return s.makeToken(token.KindBangEqual)
		}
		return s.errorToken("'!' not followed by '='")
	case '=':
		if s.match('=') {
			return s.makeToken(token.KindEqualEqual)
		}
		return s.makeToken(token.KindEqual)
	case '<':
		if s.match('=') {
			return s.makeToken(token.KindLessEqual)
		}
		return s.makeToken(token.KindLess)
	case '>':
		if s.match('=') {
			return s.makeToken(token.KindGreaterEqual)
		}
		return s.makeToken(token.KindGreater)

	case '"':
		return s.scanString()
	}

	return s.errorToken("Unexpected character.")
}

// isAtEnd checks if the scanner reached the end of the input. Specifically,
// this means that s.current is pointing beyond the last valid index of
// s.source.
func (s *Scanner) isAtEnd() bool {
	return s.current == len(s.source)
}

// makeToken returns a token of a given kind.
func (s *Scanner) makeToken(kind token.Kind) *token.Token {
	return &token.Token{
		Kind:   kind,
		Lexeme: s.source[s.start:s.current],
		Line:   s.line,
	}
}

// errorToken returns a new token of kind TokenError containing a given error
// message.s
func (s *Scanner) errorToken(message string) *token.Token {
	return &token.Token{
		Kind:   token.KindError,
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
	if s.current >= len(s.source)-1 {
		return 0
	}
	_, width := utf8.DecodeRuneInString(s.source[s.current:])
	r, _ := utf8.DecodeRuneInString(s.source[s.current+width:])
	return r
}

// scanString scans a string token.
func (s *Scanner) scanString() *token.Token {
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
	return s.makeToken(token.KindStringLiteral)
}

// scanNumber scans a number token.
func (s *Scanner) scanNumber() *token.Token {
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

	return s.makeToken(token.KindNumberLiteral)
}

// scanIdentifier scans an identifier token.
func (s *Scanner) scanIdentifier() *token.Token {
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
func (s *Scanner) identifierKind() token.Kind { // nolint:funlen,gocognit,gocyclo
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
				return s.checkKeyword(2, "ias", token.KindAlias)
			case 'n':
				return s.checkKeyword(2, "d", token.KindAnd)
			}
		}
	case 'b':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'n':
				return s.checkKeyword(2, "um", token.KindBnum)
			case 'o':
				return s.checkKeyword(2, "ol", token.KindBool)
			case 'r':
				return s.checkKeyword(2, "eak", token.KindBreak)
			}
		}
	case 'c':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword(2, "se", token.KindCase)
			case 'l':
				return s.checkKeyword(2, "ass", token.KindClass)
			case 'o':
				return s.checkKeyword(2, "ntinue", token.KindContinue)
			}
		}
	case 'd':
		return s.checkKeyword(1, "o", token.KindDo)
	case 'e':
		switch lexeme {
		case "else":
			return token.KindElse
		case "elseif":
			return token.KindElseif
		case "end":
			return token.KindEnd
		case "enum":
			return token.KindEnum
		}
	case 'f':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword(2, "lse", token.KindFalse)
			case 'l':
				return s.checkKeyword(2, "oat", token.KindFloat)
			case 'o':
				return s.checkKeyword(2, "r", token.KindFor)
			case 'u':
				return s.checkKeyword(2, "nction", token.KindFunction)
			}
		}
	case 'g':
		if len(lexeme) > 2 {
			switch lexeme {
			case "gosub":
				return token.KindGosub
			case "goto":
				return token.KindGoto
			}
		}
	case 'i':
		if len(lexeme) > 2 {
			switch s.source[s.start+1] {
			case 'f':
				return s.checkKeyword(2, "", token.KindIf)
			case 'n':
				switch lexeme {
				case "in":
					return token.KindIn
				case "int":
					return token.KindInt
				}
			}
		}
	case 'l':
		return s.checkKeyword(1, "isten", token.KindListen)
	case 'm':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword(2, "p", token.KindMap)
			case 'e':
				return s.checkKeyword(2, "ta", token.KindMeta)
			}
		}
	case 'n':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'i':
				return s.checkKeyword(2, "l", token.KindNil)
			case 'o':
				return s.checkKeyword(2, "t", token.KindNot)
			}
		}
	case 'o':
		return s.checkKeyword(1, "r", token.KindOr)
	case 'p':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'r':
				return s.checkKeyword(2, "int", token.KindPrint)
			case 'a':
				return s.checkKeyword(2, "ssage", token.KindPassage)
			}
		}
	case 'r':
		return s.checkKeyword(1, "eturn", token.KindReturn)
	case 's':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword(2, "y", token.KindSay)
			case 'u':
				return s.checkKeyword(2, "per", token.KindSuper)
			case 'w':
				return s.checkKeyword(2, "itch", token.KindSwitch)
			case 't':
				switch lexeme {
				case "struct":
					return token.KindStruct
				case "string":
					return token.KindString
				}
			}
		}
	case 't':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'h':
				return s.checkKeyword(2, "en", token.KindThen)
			case 'r':
				return s.checkKeyword(2, "ue", token.KindTrue)
			}
		}
	case 'v':
		if len(lexeme) > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword(2, "rs", token.KindVars)
			case 'o':
				return s.checkKeyword(2, "id", token.KindVoid)
			}
		}
	case 'w':
		return s.checkKeyword(1, "hile", token.KindWhile)

	}

	return token.KindIdentifier
}

// checkKeyword checks if the current lexeme is a given keyword. It start
// checking at the start-th character, checking if it matches rest. If there is
// a match, returns kind. Otherwise, returns TokenIdentifier.
func (s *Scanner) checkKeyword(start int, rest string, kind token.Kind) token.Kind {
	restLength := len(rest)
	lexemeLength := s.current - s.start
	keywordLength := start + restLength

	if lexemeLength == keywordLength && s.source[s.start+start:s.current] == rest {
		return kind
	}

	return token.KindIdentifier
}
