/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package parser

import (
	"fmt"
	"strconv"

	"gitlab.com/stackedboxes/romulang/pkg/ast"
	"gitlab.com/stackedboxes/romulang/pkg/bytecode"
)

func init() {
	initRules()
}

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

// prefixCompileFn is a function used to parse and generate code for a certain
// kind of prefix expression.
type prefixCompileFn = func(c *compiler) ast.Node

// infixCompileFn is a function used to parse and generate code for a certain
// kind of infix expression. lhs is the left-hand side expression previously
// parsed.
type infixCompileFn = func(c *compiler, lhs ast.Node) ast.Node

// parseRule encodes one rule of our Pratt parser.
type parseRule struct {
	prefix     prefixCompileFn // For expressions using the token as a prefix operator.
	infix      infixCompileFn  // For expressions using the token as an infix operator.
	precedence precedence      // When the token is used as a binary operator.
}

// parser holds some parsing-related data. I'd say it's not really a parser.
type parser struct {
	current   *token // The current token.
	previous  *token // The previous token.
	hadError  bool   // Did we find at least one error?
	panicMode bool   // Are we in panic mode? (Parsing panic, nothing to do with Go panic!)
}

// compiler is a Romualdo compiler.
type compiler struct {
	// Set DebugPrintCode to true to make the compiler print a disassembly of
	// the generated code when it finishes compiling it.
	DebugPrintCode bool

	p     *parser
	s     *scanner
	chunk *bytecode.Chunk // The chunk the compiler is generating.
}

// NewCompiler returns a new Compiler.
func NewCompiler() *compiler {
	return &compiler{
		p:     &parser{},
		chunk: &bytecode.Chunk{},
	}
}

// Compile compiles source and returns the chunk with the compiled bytecode. In
// case of errors, returns nil.
func (c *compiler) Compile(source string) (*bytecode.Chunk, ast.Node) {
	// TODO: Candidate to change: I'd like to instantiate the scanner on New().
	c.s = newScanner(source)

	c.advance()
	node := c.expression()
	c.consume(tokenKindEOF, "Expect end of expression.")

	var err error
	c.chunk, err = GenerateCode(node)
	if err != nil {
		return nil, nil
	}

	c.endCompiler()

	if c.p.hadError {
		return nil, nil
	}

	return c.chunk, node
}

// endCompiler wraps up the compilation.
func (c *compiler) endCompiler() {
	if c.DebugPrintCode && !c.p.hadError {
		c.chunk.Disassemble("code")
	}
}

// parsePrecedence parses and generates code for expressions with a precedence
// level equal to or greater than p.
func (c *compiler) parsePrecedence(p precedence) ast.Node {
	c.advance()
	prefixRule := rules[c.p.previous.kind].prefix
	if prefixRule == nil {
		c.error("Expect expression.")
		return nil
	}

	node := prefixRule(c)

	for p <= rules[c.p.current.kind].precedence {
		c.advance()
		infixRule := rules[c.p.previous.kind].infix
		node = infixRule(c, node)
	}

	return node
}

//
// The parsing rules and compilation functions
//

// rules is the table of parsing rules for our Pratt parser.
var rules []parseRule

// expression parses and generates code for an expression.
func (c *compiler) expression() ast.Node {
	return c.parsePrecedence(precAssignment)
}

// floatLiteral parses and generates code for a number literal. The float
// literal token is expected to have been just consumed.
func (c *compiler) floatLiteral() ast.Node {
	value, err := strconv.ParseFloat(c.p.previous.lexeme, 64)
	if err != nil {
		panic("Compiler got invalid number lexeme: " + c.p.previous.lexeme)
	}

	return &ast.FloatLiteral{
		Value: value,
	}
}

// stringLiteral parses and generates code for a string literal. The string
// literal token is expected to have been just consumed.
func (c *compiler) stringLiteral() ast.Node {
	value := c.p.previous.lexeme[1 : len(c.p.previous.lexeme)-1] // remove the quotes

	return &ast.StringLiteral{
		Value: value,
	}
}

// grouping parses and generates code for a parenthesized expression. The left
// paren token is expected to have been just consumed.
func (c *compiler) grouping() ast.Node {
	expr := c.expression()
	c.consume(tokenKindRightParen, "Expect ')' after expression.")
	return expr
}

// unary parses and generates code for a unary expression. The operator token is
// expected to have been just consumed.
func (c *compiler) unary() ast.Node {
	operatorKind := c.p.previous.kind
	operatorLexeme := c.p.previous.lexeme

	// Compile the operand.
	operand := c.parsePrecedence(PrecUnary)

	// Emit the operator instruction.
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

// binary parses and generates code for a binary operator expression. The left
// operand and the operator token are expected to have been just consumed.
func (c *compiler) binary(lhs ast.Node) ast.Node {
	// Remember the operator.
	operatorKind := c.p.previous.kind
	operatorLexeme := c.p.previous.lexeme

	// Compile the right operand.
	var rhs ast.Node
	rule := rules[operatorKind]
	if operatorKind == tokenKindHat {
		rhs = c.parsePrecedence(rule.precedence)
	} else {
		rhs = c.parsePrecedence(rule.precedence + 1)
	}

	return &ast.Binary{
		Operator: operatorLexeme,
		LHS:      lhs,
		RHS:      rhs,
	}
}

// boolLiteral parses and generates code for a literal Boolean value. The
// corresponding keyword is expected to have been just consumed.
func (c *compiler) boolLiteral() ast.Node {
	switch c.p.previous.kind {
	case tokenKindTrue:
		return &ast.BoolLiteral{Value: true}
	case tokenKindFalse:
		return &ast.BoolLiteral{Value: false}
	default:
		panic(fmt.Sprintf("Unexpected token type on boolLiteral: %v", c.p.previous.kind))
	}
}

// initRules initializes the rules array.
//
// Using block comments to convince gofmt to keep things aligned is ugly as
// hell.
func initRules() { // nolint:funlen
	rules = make([]parseRule, numberOfTokenKinds)

	//                                     prefix                                      infix                          precedence
	//                                    ---------------------------------------     --------------------------     --------------
	rules[tokenKindLeftParen] = /*     */ parseRule{(*compiler).grouping /*       */, nil /*                     */, precNone}
	rules[tokenKindRightParen] = /*    */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindLeftBrace] = /*     */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindRightBrace] = /*    */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindLeftBracket] = /*   */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindRightBracket] = /*  */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindComma] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindDot] = /*           */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindMinus] = /*         */ parseRule{(*compiler).unary /*          */, (*compiler).binary /*      */, precTerm}
	rules[tokenKindPlus] = /*          */ parseRule{(*compiler).unary /*          */, (*compiler).binary /*      */, precTerm}
	rules[tokenKindSlash] = /*         */ parseRule{nil /*                        */, (*compiler).binary /*      */, precFactor}
	rules[tokenKindStar] = /*          */ parseRule{nil /*                        */, (*compiler).binary /*      */, precFactor}
	rules[tokenKindColon] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindTilde] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindAt] = /*            */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindHat] = /*           */ parseRule{nil /*                        */, (*compiler).binary /*      */, precPower}
	rules[tokenKindEqual] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindEqualEqual] = /*    */ parseRule{nil /*                        */, (*compiler).binary /*      */, precEquality}
	rules[tokenKindBangEqual] = /*     */ parseRule{nil /*                        */, (*compiler).binary /*      */, precEquality}
	rules[tokenKindGreater] = /*       */ parseRule{nil /*                        */, (*compiler).binary /*      */, precComparison}
	rules[tokenKindGreaterEqual] = /*  */ parseRule{nil /*                        */, (*compiler).binary /*      */, precComparison}
	rules[tokenKindLess] = /*          */ parseRule{nil /*                        */, (*compiler).binary /*      */, precComparison}
	rules[tokenKindLessEqual] = /*     */ parseRule{nil /*                        */, (*compiler).binary /*      */, precComparison}
	rules[tokenKindIdentifier] = /*    */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindStringLiteral] = /* */ parseRule{(*compiler).stringLiteral /*  */, nil /*                     */, precNone}
	rules[tokenKindNumberLiteral] = /* */ parseRule{(*compiler).floatLiteral /*   */, nil /*                     */, precNone}
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
	rules[tokenKindFalse] = /*         */ parseRule{(*compiler).boolLiteral /*    */, nil /*                     */, precNone}
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
	rules[tokenKindNot] = /*           */ parseRule{(*compiler).unary /*          */, nil /*                     */, precNone}
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
	rules[tokenKindTrue] = /*          */ parseRule{(*compiler).boolLiteral /*    */, nil /*                     */, precNone}
	rules[tokenKindVars] = /*          */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindVoid] = /*          */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindWhile] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindError] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindEOF] = /*           */ parseRule{nil /*                        */, nil /*                     */, precNone}
}
