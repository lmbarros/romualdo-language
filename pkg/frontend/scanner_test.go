/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package frontend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests Scanner.Token() with simple cases (zero or one-token only).
func TestScannerTokenSimpleCases(t *testing.T) { // nolint: funlen
	tokens := tokenizeString("")
	assert.Equal(t, []tokenKind{tokenKindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1}, tokenLines(tokens))

	tokens = tokenizeString("foo")
	assert.Equal(t, []tokenKind{tokenKindIdentifier, tokenKindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{"foo", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString("listen")
	assert.Equal(t, []tokenKind{tokenKindListen, tokenKindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{"listen", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString("struct")
	assert.Equal(t, []tokenKind{tokenKindStruct, tokenKindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{"struct", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString("lístên")
	assert.Equal(t, []tokenKind{tokenKindIdentifier, tokenKindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{"lístên", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString("# foo")
	assert.Equal(t, []tokenKind{tokenKindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1}, tokenLines(tokens))

	tokens = tokenizeString("123.456")
	assert.Equal(t, []tokenKind{tokenKindFloatLiteral, tokenKindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{"123.456", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString(">=")
	assert.Equal(t, []tokenKind{tokenKindGreaterEqual, tokenKindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{">=", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString("=")
	assert.Equal(t, []tokenKind{tokenKindEqual, tokenKindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{"=", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString("==")
	assert.Equal(t, []tokenKind{tokenKindEqualEqual, tokenKindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{"==", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString(`"turtles"`)
	assert.Equal(t, []tokenKind{tokenKindStringLiteral, tokenKindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{`"turtles"`, ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString("!=")
	assert.Equal(t, []tokenKind{tokenKindBangEqual, tokenKindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{"!=", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString("!")
	assert.Equal(t, []tokenKind{tokenKindError}, tokenKinds(tokens))

	tokens = tokenizeString("⟨")
	assert.Equal(t, []tokenKind{tokenKindError}, tokenKinds(tokens))
}

// Tests Scanner.Token() with token sequences longer than one token.
func TestScannerTokenSequences(t *testing.T) { // nolint: funlen
	tokens := tokenizeString("while true do")
	assert.Equal(t, []tokenKind{
		tokenKindWhile, tokenKindTrue, tokenKindDo, tokenKindEOF},
		tokenKinds(tokens))
	assert.Equal(t, []string{"while", "true", "do", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1, 1, 1}, tokenLines(tokens))

	tokens = tokenizeString("vars x: int = 1 + 2.0 end")
	assert.Equal(t, []tokenKind{
		tokenKindVars, tokenKindIdentifier, tokenKindColon, tokenKindInt,
		tokenKindEqual, tokenKindIntLiteral, tokenKindPlus,
		tokenKindFloatLiteral, tokenKindEnd, tokenKindEOF},
		tokenKinds(tokens))
	assert.Equal(t, []string{"vars", "x", ":", "int", "=", "1", "+", "2.0", "end", ""},
		tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, tokenLines(tokens))

	tokens = tokenizeString(`struct,物語 "string"+~`)
	assert.Equal(t, []tokenKind{
		tokenKindStruct, tokenKindComma, tokenKindIdentifier,
		tokenKindStringLiteral, tokenKindPlus, tokenKindTilde, tokenKindEOF},
		tokenKinds(tokens))
	assert.Equal(t, []string{"struct", ",", "物語", `"string"`, "+", "~", ""},
		tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1, 1, 1, 1, 1, 1}, tokenLines(tokens))

	tokens = tokenizeString("(alias{and / or}    super) # ⟨123.2⟩")
	assert.Equal(t, []tokenKind{
		tokenKindLeftParen, tokenKindAlias, tokenKindLeftBrace, tokenKindAnd,
		tokenKindSlash, tokenKindOr, tokenKindRightBrace, tokenKindSuper,
		tokenKindRightParen, tokenKindEOF},
		tokenKinds(tokens))
	assert.Equal(t, []string{"(", "alias", "{", "and", "/", "or", "}", "super", ")", ""},
		tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, tokenLines(tokens))

}

// Tests Scanner.Token() with multiline input.
func TestScannerTokenMultiline(t *testing.T) {
	tokens := tokenizeString(
		`break "starts # goes on
and continues" [ elseif # now this is a comment
       int inti^`)
	assert.Equal(t, []tokenKind{
		tokenKindBreak, tokenKindStringLiteral, tokenKindLeftBracket,
		tokenKindElseif, tokenKindInt, tokenKindIdentifier, tokenKindHat,
		tokenKindEOF},
		tokenKinds(tokens))
	assert.Equal(t, []string{"break", "\"starts # goes on\nand continues\"",
		"[", "elseif", "int", "inti", "^", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 2, 2, 2, 3, 3, 3, 3}, tokenLines(tokens))

	tokens = tokenizeString(
		`goto continue conti
nue @


.]



meta		map
`)
	assert.Equal(t, []tokenKind{
		tokenKindGoto, tokenKindContinue, tokenKindIdentifier,
		tokenKindIdentifier, tokenKindAt, tokenKindDot,
		tokenKindRightBracket, tokenKindMeta, tokenKindMap, tokenKindEOF},
		tokenKinds(tokens))
	assert.Equal(t, []string{"goto", "continue", "conti", "nue", "@", ".", "]",
		"meta", "map", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1, 1, 2, 2, 5, 5, 9, 9, 10}, tokenLines(tokens))

	tokens = tokenizeString(`
# foobar

# nothing here
# nor here

`)
	assert.Equal(t, []tokenKind{tokenKindEOF}, tokenKinds(tokens))
	assert.Equal(t, []string{""}, tokenLexemes(tokens))
	assert.Equal(t, []int{7}, tokenLines(tokens))

	tokens = tokenizeString(
		`goto continue conti
nue @


.]



meta		map
`)
	assert.Equal(t, []tokenKind{
		tokenKindGoto, tokenKindContinue, tokenKindIdentifier,
		tokenKindIdentifier, tokenKindAt, tokenKindDot,
		tokenKindRightBracket, tokenKindMeta, tokenKindMap, tokenKindEOF},
		tokenKinds(tokens))
	assert.Equal(t, []string{"goto", "continue", "conti", "nue", "@", ".", "]",
		"meta", "map", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1, 1, 2, 2, 5, 5, 9, 9, 10}, tokenLines(tokens))

	tokens = tokenizeString(
		`1<bnum
		function<=0.55b
		void> string`)
	assert.Equal(t, []tokenKind{
		tokenKindIntLiteral, tokenKindLess, tokenKindBNum,
		tokenKindFunction, tokenKindLessEqual, tokenKindBNumLiteral,
		tokenKindVoid, tokenKindGreater, tokenKindString, tokenKindEOF},
		tokenKinds(tokens))
	assert.Equal(t, []string{"1", "<", "bnum", "function", "<=", "0.55b",
		"void", ">", "string", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1, 1, 2, 2, 2, 3, 3, 3, 3}, tokenLines(tokens))

	tokens = tokenizeString(
		`ifin
		inif
		if in
		in if`)
	assert.Equal(t, []tokenKind{
		tokenKindIdentifier, tokenKindIdentifier, tokenKindIf, tokenKindIn,
		tokenKindIn, tokenKindIf, tokenKindEOF},
		tokenKinds(tokens))
	assert.Equal(t, []string{"ifin", "inif", "if", "in", "in", "if", ""},
		tokenLexemes(tokens))
	assert.Equal(t, []int{1, 2, 3, 3, 4, 4, 4}, tokenLines(tokens))

	tokens = tokenizeString(`# boring test case, just to get lots of coverage
		- * bool case class
		else enum false # a comment...
		float function  # to make it...
		gosub nil not   # less boring
		print passage say switch string then void return`)
	assert.Equal(t, []tokenKind{
		tokenKindMinus, tokenKindStar, tokenKindBool, tokenKindCase,
		tokenKindClass, tokenKindElse, tokenKindEnum, tokenKindFalse,
		tokenKindFloat, tokenKindFunction, tokenKindGosub, tokenKindNil,
		tokenKindNot, tokenKindPrint, tokenKindPassage, tokenKindSay,
		tokenKindSwitch, tokenKindString, tokenKindThen, tokenKindVoid,
		tokenKindReturn, tokenKindEOF},
		tokenKinds(tokens))
	assert.Equal(t, []string{"-", "*", "bool", "case", "class", "else", "enum",
		"false", "float", "function", "gosub", "nil", "not", "print", "passage",
		"say", "switch", "string", "then", "void", "return", ""},
		tokenLexemes(tokens))
	assert.Equal(t, []int{2, 2, 2, 2, 2, 3, 3, 3, 4, 4, 5, 5, 5,
		6, 6, 6, 6, 6, 6, 6, 6, 6}, tokenLines(tokens))
}

// Tests Scanner.Token() with numbers.
func TestScannerTokenNumbers(t *testing.T) {
	tokens := tokenizeString("9876")
	assert.Equal(t, []tokenKind{tokenKindIntLiteral, tokenKindEOF},
		tokenKinds(tokens))
	assert.Equal(t, []string{"9876", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString("9876.54")
	assert.Equal(t, []tokenKind{tokenKindFloatLiteral, tokenKindEOF},
		tokenKinds(tokens))
	assert.Equal(t, []string{"9876.54", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	tokens = tokenizeString("0.123b")
	assert.Equal(t, []tokenKind{tokenKindBNumLiteral, tokenKindEOF},
		tokenKinds(tokens))
	assert.Equal(t, []string{"0.123b", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1}, tokenLines(tokens))

	// We like floats looking like the complete thing, so we don't recognize
	// ".111" or "111." as floats.
	tokens = tokenizeString("9876.")
	assert.Equal(t, []tokenKind{tokenKindIntLiteral, tokenKindDot, tokenKindEOF},
		tokenKinds(tokens))
	assert.Equal(t, []string{"9876", ".", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1, 1}, tokenLines(tokens))

	tokens = tokenizeString(".7890")
	assert.Equal(t, []tokenKind{tokenKindDot, tokenKindIntLiteral, tokenKindEOF},
		tokenKinds(tokens))
	assert.Equal(t, []string{".", "7890", ""}, tokenLexemes(tokens))
	assert.Equal(t, []int{1, 1, 1}, tokenLines(tokens))
}

// Tests Scanner.Token() with strings.
func TestScannerTokenStrings(t *testing.T) {
	tokens := tokenizeString("\"the neverending string")
	assert.Equal(t, []tokenKind{tokenKindError}, tokenKinds(tokens))
}

// tokenKinds extract the token kinds from a slice of tokens.
func tokenKinds(tokens []*token) []tokenKind {
	result := make([]tokenKind, 0, len(tokens))

	for _, tok := range tokens {
		result = append(result, tok.kind)
	}

	return result
}

// tokenLexemes extract the lexemes from a slice of tokens.
func tokenLexemes(tokens []*token) []string {
	result := make([]string, 0, len(tokens))

	for _, tok := range tokens {
		result = append(result, tok.lexeme)
	}

	return result
}

// tokenLexemes extract the line numbers from a slice of tokens.
func tokenLines(tokens []*token) []int {
	result := make([]int, 0, len(tokens))

	for _, tok := range tokens {
		result = append(result, tok.line)
	}

	return result
}

// tokenizeString creates a Scanner and calls Token() on it until getting an
// EOF or error. Then it returns a slice with the resulting Tokens.
func tokenizeString(source string) []*token {
	s := newScanner(source)
	result := make([]*token, 0, 16)

	tok := s.token()
	result = append(result, tok)
	for tok.kind != tokenKindEOF && tok.kind != tokenKindError {
		tok = s.token()
		result = append(result, tok)
	}

	return result
}
