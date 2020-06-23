/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/stackedboxes/romulang/pkg/token"
)

// Tests Scanner.Token() with simple cases (zero or one-token only).
func TestScannerTokenSimpleCases(t *testing.T) { // nolint: funlen
	tokens := tokenizeString("")
	assert.Equal(t, []token.Kind{token.KindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1}, tokenLines(tokens))

	tokens = tokenizeString("foo")
	assert.Equal(t, []token.Kind{token.KindIdentifier, token.KindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{"foo", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString("listen")
	assert.Equal(t, []token.Kind{token.KindListen, token.KindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{"listen", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString("struct")
	assert.Equal(t, []token.Kind{token.KindStruct, token.KindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{"struct", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString("lístên")
	assert.Equal(t, []token.Kind{token.KindIdentifier, token.KindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{"lístên", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString("# foo")
	assert.Equal(t, []token.Kind{token.KindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1}, tokenLines(tokens))

	tokens = tokenizeString("123.456")
	assert.Equal(t, []token.Kind{token.KindNumberLiteral, token.KindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{"123.456", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString(">=")
	assert.Equal(t, []token.Kind{token.KindGreaterEqual, token.KindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{">=", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString("=")
	assert.Equal(t, []token.Kind{token.KindEqual, token.KindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{"=", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString("==")
	assert.Equal(t, []token.Kind{token.KindEqualEqual, token.KindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{"==", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString(`"turtles"`)
	assert.Equal(t, []token.Kind{token.KindStringLiteral, token.KindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{`"turtles"`, ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString("!=")
	assert.Equal(t, []token.Kind{token.KindBangEqual, token.KindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{"!=", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString("!")
	assert.Equal(t, []token.Kind{token.KindError}, tokenKinds(tokens))

	tokens = tokenizeString("⟨")
	assert.Equal(t, []token.Kind{token.KindError}, tokenKinds(tokens))
}

// Tests Scanner.Token() token sequences longer than one token.
func TestScannerTokenTokenSequences(t *testing.T) { // nolint: funlen
	tokens := tokenizeString("while true do")
	assert.Equal(t, []token.Kind{
		token.KindWhile, token.KindTrue, token.KindDo, token.KindEOF},
		tokenKinds(tokens))
	assert.Equal(t, []string{"while", "true", "do", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1, 1, 1}, tokenLines(tokens))

	tokens = tokenizeString("vars x: int = 1 + 2 end")
	assert.Equal(t, []token.Kind{
		token.KindVars, token.KindIdentifier, token.KindColon, token.KindInt,
		token.KindEqual, token.KindNumberLiteral, token.KindPlus,
		token.KindNumberLiteral, token.KindEnd, token.KindEOF},
		tokenKinds(tokens))
	assert.Equal(t, []string{"vars", "x", ":", "int", "=", "1", "+", "2", "end", ""},
		tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, tokenLines(tokens))

	tokens = tokenizeString(`struct,物語 "string"+~`)
	assert.Equal(t, []token.Kind{
		token.KindStruct, token.KindComma, token.KindIdentifier,
		token.KindStringLiteral, token.KindPlus, token.KindTilde, token.KindEOF},
		tokenKinds(tokens))
	assert.Equal(t, []string{"struct", ",", "物語", `"string"`, "+", "~", ""},
		tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1, 1, 1, 1, 1, 1}, tokenLines(tokens))

	tokens = tokenizeString("(alias{and / or}    super) # ⟨123.2⟩")
	assert.Equal(t, []token.Kind{
		token.KindLeftParen, token.KindAlias, token.KindLeftBrace, token.KindAnd,
		token.KindSlash, token.KindOr, token.KindRightBrace, token.KindSuper,
		token.KindRightParen, token.KindEOF},
		tokenKinds(tokens))
	assert.Equal(t, []string{"(", "alias", "{", "and", "/", "or", "}", "super", ")", ""},
		tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, tokenLines(tokens))

}

// tokenKinds extract the token kinds from a slice of tokens.
func tokenKinds(tokens []*token.Token) []token.Kind {
	result := make([]token.Kind, 0, len(tokens))

	for _, tok := range tokens {
		result = append(result, tok.Kind)
	}

	return result
}

// tokenLexemes extract the lexemes from a slice of tokens.
func tokenLexemes(tokens []*token.Token) []string {
	result := make([]string, 0, len(tokens))

	for _, tok := range tokens {
		result = append(result, tok.Lexeme)
	}

	return result
}

// tokenLexemes extract the line numbers from a slice of tokens.
func tokenLines(tokens []*token.Token) []int {
	result := make([]int, 0, len(tokens))

	for _, tok := range tokens {
		result = append(result, tok.Line)
	}

	return result
}

// tokenizeString creates a Scanner and calls Token() on it until getting an
// EOF or error. Then it returns a slice with the resulting Tokens.
func tokenizeString(source string) []*token.Token {
	s := New(source)
	result := make([]*token.Token, 0, 16)

	tok := s.Token()
	result = append(result, tok)
	for tok.Kind != token.KindEOF && tok.Kind != token.KindError {
		tok = s.Token()
		result = append(result, tok)
	}

	return result
}
