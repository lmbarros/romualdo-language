/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package frontend

import (
	"fmt"
	"os"
	"strconv"

	"gitlab.com/stackedboxes/romulang/pkg/ast"
)

// precedence is the precedence of expressions.
type precedence int

const (
	precNone       precedence = iota // Means: cannot be the "center" of an expression.
	precAssignment                   // =
	precOr                           // or
	precAnd                          // and
	precEquality                     // == !=
	precComparison                   // < > <= >=
	precTerm                         // + -
	precFactor                       // * /
	PrecUnary                        // not -
	precPower                        // ^
	precCall                         // . ()
	precPrimary
)

// prefixParseFn is a function used to parse code for a certain kind of prefix
// expression.
type prefixParseFn = func(c *parser) ast.Node

// infixParseFn is a function used to parse code for a certain kind of infix
// expression. lhs is the left-hand side expression previously parsed.
type infixParseFn = func(c *parser, lhs ast.Node) ast.Node

// parseRule encodes one rule of our Pratt parser.
type parseRule struct {
	prefix     prefixParseFn // For expressions using the token as a prefix operator.
	infix      infixParseFn  // For expressions using the token as an infix operator.
	precedence precedence    // When the token is used as a binary operator.
}

// parser is a parser for the Romualdo language. It converts source code into an AST.
type parser struct {
	// currentToken is the current token we are parsing.
	currentToken *token

	// previousToken is the previous token we have parsed.
	previousToken *token

	// hadError indicates whether we found at least one syntax error.
	hadError bool

	// panicMode indicates whether we are in panic mode. This has nothing to do
	// with Go panics. Right after finding a syntax error it is hard to generate
	// good error messages because the parser is "out of sync" with the code, so
	// we enter panic mode (during which we don't report any errors). Once we
	// find a "synchronization point", we leave panicmode.
	panicMode bool

	// The scanner from where we get our tokens.
	scanner *scanner
}

// newParser returns a new parser that will parse source.
func newParser(source string) *parser {
	return &parser{
		scanner: newScanner(source),
	}
}

// parse parses source and returns the root of the resulting AST. Returns nil in
// case of error.
func (p *parser) parse() ast.Node {
	p.advance()
	node := p.expression()
	p.consume(tokenKindEOF, "Expect end of expression.")
	if p.hadError {
		return nil
	}

	return node
}

// parsePrecedence parses and generates code for expressions with a precedence
// level equal to or greater than p.
func (p *parser) parsePrecedence(prec precedence) ast.Node {
	p.advance()
	prefixRule := rules[p.previousToken.kind].prefix
	if prefixRule == nil {
		p.error("Expect expression.")
		return nil
	}

	node := prefixRule(p)

	for prec <= rules[p.currentToken.kind].precedence {
		p.advance()
		infixRule := rules[p.previousToken.kind].infix
		node = infixRule(p, node)
	}

	return node
}

//
// The parsing rules and compilation functions
//

// rules is the table of parsing rules for our Pratt parser.
var rules []parseRule

// expression parses an expression.
func (p *parser) expression() ast.Node {
	return p.parsePrecedence(precAssignment)
}

// floatLiteral parses a floting-point number literal. The float literal token
// is expected to have been just consumed.
func (p *parser) floatLiteral() ast.Node {
	value, err := strconv.ParseFloat(p.previousToken.lexeme, 64)
	if err != nil {
		panic("Compiler got invalid number lexeme: " + p.previousToken.lexeme)
	}

	return &ast.FloatLiteral{
		BaseNode: ast.BaseNode{
			LineNumber:   p.previousToken.line,
			SourceLexeme: p.previousToken.lexeme,
		},
		Value: value,
	}
}

// stringLiteral parses a string literal. The string literal token is expected
// to have been just consumed.
func (p *parser) stringLiteral() ast.Node {
	value := p.previousToken.lexeme[1 : len(p.previousToken.lexeme)-1] // remove the quotes

	return &ast.StringLiteral{
		BaseNode: ast.BaseNode{
			LineNumber:   p.previousToken.line,
			SourceLexeme: p.previousToken.lexeme,
		},
		Value: value,
	}
}

// grouping parses a parenthesized expression. The left paren token is expected
// to have been just consumed.
func (p *parser) grouping() ast.Node {
	expr := p.expression()
	p.consume(tokenKindRightParen, "Expect ')' after expression.")
	return expr
}

// unary parses a unary expression. The operator token is expected to have been
// just consumed.
func (p *parser) unary() ast.Node {
	operatorKind := p.previousToken.kind
	operatorLexeme := p.previousToken.lexeme

	// Parse the operand.
	operand := p.parsePrecedence(PrecUnary)

	// Return the node.
	switch operatorKind {
	case tokenKindNot:
		return &ast.Unary{Operator: operatorLexeme, Operand: operand}
	case tokenKindMinus:
		return &ast.Unary{Operator: operatorLexeme, Operand: operand}
	case tokenKindPlus:
		// Unary plus is a no-op.
		return operand
	default:
		panic(fmt.Sprintf("Unexpected operator kind on unary expression: %v", operatorKind))
	}
}

// binary parses a binary operator expression. The left operand and the operator
// token are expected to have been just consumed.
func (p *parser) binary(lhs ast.Node) ast.Node {
	// Remember the operator.
	operatorKind := p.previousToken.kind
	operatorLexeme := p.previousToken.lexeme

	// Parse the right operand.
	var rhs ast.Node
	rule := rules[operatorKind]
	if operatorKind == tokenKindHat {
		rhs = p.parsePrecedence(rule.precedence)
	} else {
		rhs = p.parsePrecedence(rule.precedence + 1)
	}

	return &ast.Binary{
		BaseNode: ast.BaseNode{
			LineNumber:   p.previousToken.line,
			SourceLexeme: p.previousToken.lexeme,
		},
		Operator: operatorLexeme,
		LHS:      lhs,
		RHS:      rhs,
	}
}

// boolLiteral parses a literal Boolean value. The corresponding keyword is
// expected to have been just consumed.
func (p *parser) boolLiteral() ast.Node {
	switch p.previousToken.kind {
	case tokenKindTrue:
		return &ast.BoolLiteral{
			BaseNode: ast.BaseNode{
				LineNumber:   p.previousToken.line,
				SourceLexeme: p.previousToken.lexeme,
			},
			Value: true,
		}
	case tokenKindFalse:
		return &ast.BoolLiteral{
			BaseNode: ast.BaseNode{
				LineNumber:   p.previousToken.line,
				SourceLexeme: p.previousToken.lexeme,
			},
			Value: false,
		}
	default:
		panic(fmt.Sprintf("Unexpected token type on boolLiteral: %v", p.previousToken.kind))
	}
}

// advance advances the parser by one token. This will report errors for each
// error token found; callers will only see the non-error tokens.
func (p *parser) advance() {
	p.previousToken = p.currentToken

	for {
		p.currentToken = p.scanner.token()
		if p.currentToken.kind != tokenKindError {
			break
		}

		p.errorAtCurrent(p.currentToken.lexeme)
	}
}

// consume consumes the current token (and advances the parser), assuming it is
// of a given kind. If it is not of this kind, reports this is an error with a
// given error message.
func (p *parser) consume(kind tokenKind, message string) {
	if p.currentToken.kind == kind {
		p.advance()
		return
	}

	p.errorAtCurrent(message)
}

// errorAtCurrent reports an error at the current (c.currentToken) token.
func (p *parser) errorAtCurrent(message string) {
	p.errorAt(p.currentToken, message)
}

// error reports an error at the token we just consumed (c.previousToken).
func (p *parser) error(message string) {
	p.errorAt(p.previousToken, message)
}

// errorAt reports an error at a given token, with a given error message.
func (p *parser) errorAt(tok *token, message string) {
	if p.panicMode {
		return
	}

	p.panicMode = true

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
	p.hadError = true
}

func init() {
	initRules()
}

// initRules initializes the rules array.
//
// Using block comments to convince gofmt to keep things aligned is ugly as
// hell.
func initRules() { // nolint:funlen
	rules = make([]parseRule, numberOfTokenKinds)

	//                                     prefix                                      infix                          precedence
	//                                    ---------------------------------------     --------------------------     --------------
	rules[tokenKindLeftParen] = /*     */ parseRule{(*parser).grouping /*       */, nil /*                     */, precNone}
	rules[tokenKindRightParen] = /*    */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindLeftBrace] = /*     */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindRightBrace] = /*    */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindLeftBracket] = /*   */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindRightBracket] = /*  */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindComma] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindDot] = /*           */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindMinus] = /*         */ parseRule{(*parser).unary /*          */, (*parser).binary /*      */, precTerm}
	rules[tokenKindPlus] = /*          */ parseRule{(*parser).unary /*          */, (*parser).binary /*      */, precTerm}
	rules[tokenKindSlash] = /*         */ parseRule{nil /*                        */, (*parser).binary /*      */, precFactor}
	rules[tokenKindStar] = /*          */ parseRule{nil /*                        */, (*parser).binary /*      */, precFactor}
	rules[tokenKindColon] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindTilde] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindAt] = /*            */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindHat] = /*           */ parseRule{nil /*                        */, (*parser).binary /*      */, precPower}
	rules[tokenKindEqual] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindEqualEqual] = /*    */ parseRule{nil /*                        */, (*parser).binary /*      */, precEquality}
	rules[tokenKindBangEqual] = /*     */ parseRule{nil /*                        */, (*parser).binary /*      */, precEquality}
	rules[tokenKindGreater] = /*       */ parseRule{nil /*                        */, (*parser).binary /*      */, precComparison}
	rules[tokenKindGreaterEqual] = /*  */ parseRule{nil /*                        */, (*parser).binary /*      */, precComparison}
	rules[tokenKindLess] = /*          */ parseRule{nil /*                        */, (*parser).binary /*      */, precComparison}
	rules[tokenKindLessEqual] = /*     */ parseRule{nil /*                        */, (*parser).binary /*      */, precComparison}
	rules[tokenKindIdentifier] = /*    */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindStringLiteral] = /* */ parseRule{(*parser).stringLiteral /*  */, nil /*                     */, precNone}
	rules[tokenKindNumberLiteral] = /* */ parseRule{(*parser).floatLiteral /*   */, nil /*                     */, precNone}
	rules[tokenKindAlias] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindAnd] = /*           */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindBnum] = /*          */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindBool] = /*          */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindBreak] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindCase] = /*          */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindClass] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindContinue] = /*      */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindDo] = /*            */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindElse] = /*          */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindElseif] = /*        */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindEnd] = /*           */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindEnum] = /*          */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindFalse] = /*         */ parseRule{(*parser).boolLiteral /*    */, nil /*                     */, precNone}
	rules[tokenKindFloat] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindFor] = /*           */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindFunction] = /*      */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindGosub] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindGoto] = /*          */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindIf] = /*            */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindIn] = /*            */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindInt] = /*           */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindListen] = /*        */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindMap] = /*           */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindMeta] = /*          */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindNil] = /*           */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindNot] = /*           */ parseRule{(*parser).unary /*          */, nil /*                     */, precNone}
	rules[tokenKindOr] = /*            */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindPassage] = /*       */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindPrint] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindReturn] = /*        */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindSay] = /*           */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindString] = /*        */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindStruct] = /*        */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindSuper] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindSwitch] = /*        */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindThen] = /*          */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindTrue] = /*          */ parseRule{(*parser).boolLiteral /*    */, nil /*                     */, precNone}
	rules[tokenKindVars] = /*          */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindVoid] = /*          */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindWhile] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindError] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindEOF] = /*           */ parseRule{nil /*                        */, nil /*                     */, precNone}
}
