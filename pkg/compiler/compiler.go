/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package compiler

import (
	"fmt"
	"strconv"

	"gitlab.com/stackedboxes/romulang/pkg/ast"
	"gitlab.com/stackedboxes/romulang/pkg/bytecode"
	"gitlab.com/stackedboxes/romulang/pkg/scanner"
	"gitlab.com/stackedboxes/romulang/pkg/token"
)

func init() {
	initRules()
}

// precedence is the precedence of expressions.
type precedence int

const (
	PrecNone       precedence = iota // Means: cannot be the "center" of an expression.
	PrecAssignment                   // =
	PrecOr                           // or
	PrecAnd                          // and
	PrecEquality                     // == !=
	PrecComparison                   // < > <= >=
	PrecTerm                         // + -
	PrecFactor                       // * /
	PrecUnary                        // not -
	PrecPower                        // ^
	PrecCall                         // . ()
	PrecPrimary
)

// compileFn is a function used to parse and generate code for a certain kind of
// prefix expression.
type compileFn = func(c *Compiler) ast.Node

// infixCompileFn is a function used to parse and generate code for a certain
// kind of infix expression. lhs is the left-hand side expression previously
// parsed.
type infixCompileFn = func(c *Compiler, lhs ast.Node) ast.Node

// parseRule encodes one rule of our Pratt parser.
type parseRule struct {
	prefix     compileFn      // For expressions using the token as a prefix operator.
	infix      infixCompileFn // For expressions using the token as an infix operator.
	precedence precedence     // When the token is used as a binary operator.
}

// parser holds some parsing-related data. I'd say it's not really a parser.
type parser struct {
	current   *token.Token // The current token.
	previous  *token.Token // The previous token.
	hadError  bool         // Did we find at least one error?
	panicMode bool         // Are we in panic mode? (Parsing panic, nothing to do with Go panic!)
}

// Compiler is a Romualdo compiler.
type Compiler struct {
	// Set DebugPrintCode to true to make the compiler print a disassembly of
	// the generated code when it finishes compiling it.
	DebugPrintCode bool

	p     *parser
	s     *scanner.Scanner
	chunk *bytecode.Chunk // The chunk the compiler is generating.
}

// New returns a new Compiler.
func New() *Compiler {
	return &Compiler{
		p:     &parser{},
		chunk: &bytecode.Chunk{},
	}
}

// Compile compiles source and returns the chunk with the compiled bytecode. In
// case of errors, returns nil.
func (c *Compiler) Compile(source string) (*bytecode.Chunk, ast.Node) {
	// TODO: Candidate to change: I'd like to instantiate the scanner on New().
	c.s = scanner.New(source)

	c.advance()
	node := c.expression()
	c.consume(token.KindEOF, "Expect end of expression.")

	c.endCompiler()

	if c.p.hadError {
		return nil, nil
	}

	return c.chunk, node
}

// endCompiler wraps up the compilation.
func (c *Compiler) endCompiler() {
	c.emitReturn()

	if c.DebugPrintCode && !c.p.hadError {
		c.chunk.Disassemble("code")
	}
}

// parsePrecedence parses and generates code for expressions with a precedence
// level equal to or greater than p.
func (c *Compiler) parsePrecedence(p precedence) ast.Node {
	c.advance()
	prefixRule := rules[c.p.previous.Kind].prefix
	if prefixRule == nil {
		c.error("Expect expression.")
		return nil
	}

	node := prefixRule(c)

	for p <= rules[c.p.current.Kind].precedence {
		c.advance()
		infixRule := rules[c.p.previous.Kind].infix
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
func (c *Compiler) expression() ast.Node {
	return c.parsePrecedence(PrecAssignment)
}

// floatLiteral parses and generates code for a number literal. The float
// literal token is expected to have been just consumed.
func (c *Compiler) floatLiteral() ast.Node {
	value, err := strconv.ParseFloat(c.p.previous.Lexeme, 64)
	if err != nil {
		panic("Compiler got invalid number lexeme: " + c.p.previous.Lexeme)
	}
	c.emitConstant(bytecode.NewValueFloat(value))

	return &ast.FloatLiteral{
		Value: value,
	}
}

// stringLiteral parses and generates code for a string literal. The string
// literal token is expected to have been just consumed.
func (c *Compiler) stringLiteral() ast.Node {
	value := c.p.previous.Lexeme[1 : len(c.p.previous.Lexeme)-1] // remove the quotes
	c.emitConstant(bytecode.NewValueString(value))

	return &ast.StringLiteral{
		Value: value,
	}
}

// grouping parses and generates code for a parenthesized expression. The left
// paren token is expected to have been just consumed.
func (c *Compiler) grouping() ast.Node {
	expr := c.expression()
	c.consume(token.KindRightParen, "Expect ')' after expression.")
	return &ast.Grouping{Expr: expr}
}

// unary parses and generates code for a unary expression. The operator token is
// expected to have been just consumed.
func (c *Compiler) unary() ast.Node {
	operatorKind := c.p.previous.Kind
	operatorLexeme := c.p.previous.Lexeme

	// Compile the operand.
	operand := c.parsePrecedence(PrecUnary)

	// Emit the operator instruction.
	switch operatorKind {
	case token.KindNot:
		c.emitBytes(bytecode.OpNot)
		return &ast.Unary{Operator: operatorLexeme, Operand: operand}
	case token.KindMinus:
		c.emitBytes(bytecode.OpNegate)
		return &ast.Unary{Operator: operatorLexeme, Operand: operand}
	case token.KindPlus:
		// Unary plus is a no-op.
		return operand
	default:
		panic(fmt.Sprintf("Unexpected operator kind on unary expression: %v", operatorKind))
	}
}

// binary parses and generates code for a binary operator expression. The left
// operand and the operator token are expected to have been just consumed.
func (c *Compiler) binary(lhs ast.Node) ast.Node {
	// Remember the operator.
	operatorKind := c.p.previous.Kind
	operatorLexeme := c.p.previous.Lexeme

	// Compile the right operand.
	var rhs ast.Node
	rule := rules[operatorKind]
	if operatorKind == token.KindHat {
		rhs = c.parsePrecedence(rule.precedence)
	} else {
		rhs = c.parsePrecedence(rule.precedence + 1)
	}

	// Emit the operator instruction.
	switch operatorKind {
	case token.KindBangEqual:
		c.emitBytes(bytecode.OpNotEqual)
	case token.KindEqualEqual:
		c.emitBytes(bytecode.OpEqual)
	case token.KindGreater:
		c.emitBytes(bytecode.OpGreater)
	case token.KindGreaterEqual:
		c.emitBytes(bytecode.OpGreaterEqual)
	case token.KindLess:
		c.emitBytes(bytecode.OpLess)
	case token.KindLessEqual:
		c.emitBytes(bytecode.OpLessEqual)
	case token.KindPlus:
		c.emitBytes(bytecode.OpAdd)
	case token.KindMinus:
		c.emitBytes(bytecode.OpSubtract)
	case token.KindStar:
		c.emitBytes(bytecode.OpMultiply)
	case token.KindSlash:
		c.emitBytes(bytecode.OpDivide)
	case token.KindHat:
		c.emitBytes(bytecode.OpPower)
	default:
		panic(fmt.Sprintf("Unexpected token type on binary operator: %v", operatorKind))
	}

	return &ast.Binary{
		Operator: operatorLexeme,
		LHS:      lhs,
		RHS:      rhs,
	}
}

// boolLiteral parses and generates code for a literal Boolean value. The
// corresponding keyword is expected to have been just consumed.
func (c *Compiler) boolLiteral() ast.Node {
	switch c.p.previous.Kind {
	case token.KindTrue:
		c.emitBytes(bytecode.OpTrue)
		return &ast.BoolLiteral{Value: true}
	case token.KindFalse:
		c.emitBytes(bytecode.OpFalse)
		return &ast.BoolLiteral{Value: false}
	default:
		panic(fmt.Sprintf("Unexpected token type on boolLiteral: %v", c.p.previous.Kind))
	}
}

// initRules initializes the rules array.
//
// Using block comments to convince gofmt to keep things aligned is ugly as
// hell.
func initRules() { // nolint:funlen
	rules = make([]parseRule, token.NumberOfKinds)

	//                                     prefix                                      infix                          precedence
	//                                    ---------------------------------------     --------------------------     --------------
	rules[token.KindLeftParen] = /*     */ parseRule{(*Compiler).grouping /*       */, nil /*                     */, PrecNone}
	rules[token.KindRightParen] = /*    */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindLeftBrace] = /*     */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindRightBrace] = /*    */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindLeftBracket] = /*   */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindRightBracket] = /*  */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindComma] = /*         */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindDot] = /*           */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindMinus] = /*         */ parseRule{(*Compiler).unary /*          */, (*Compiler).binary /*      */, PrecTerm}
	rules[token.KindPlus] = /*          */ parseRule{(*Compiler).unary /*          */, (*Compiler).binary /*      */, PrecTerm}
	rules[token.KindSlash] = /*         */ parseRule{nil /*                        */, (*Compiler).binary /*      */, PrecFactor}
	rules[token.KindStar] = /*          */ parseRule{nil /*                        */, (*Compiler).binary /*      */, PrecFactor}
	rules[token.KindColon] = /*         */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindTilde] = /*         */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindAt] = /*            */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindHat] = /*           */ parseRule{nil /*                        */, (*Compiler).binary /*      */, PrecPower}
	rules[token.KindEqual] = /*         */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindEqualEqual] = /*    */ parseRule{nil /*                        */, (*Compiler).binary /*      */, PrecEquality}
	rules[token.KindBangEqual] = /*     */ parseRule{nil /*                        */, (*Compiler).binary /*      */, PrecEquality}
	rules[token.KindGreater] = /*       */ parseRule{nil /*                        */, (*Compiler).binary /*      */, PrecComparison}
	rules[token.KindGreaterEqual] = /*  */ parseRule{nil /*                        */, (*Compiler).binary /*      */, PrecComparison}
	rules[token.KindLess] = /*          */ parseRule{nil /*                        */, (*Compiler).binary /*      */, PrecComparison}
	rules[token.KindLessEqual] = /*     */ parseRule{nil /*                        */, (*Compiler).binary /*      */, PrecComparison}
	rules[token.KindIdentifier] = /*    */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindStringLiteral] = /* */ parseRule{(*Compiler).stringLiteral /*  */, nil /*                     */, PrecNone}
	rules[token.KindNumberLiteral] = /* */ parseRule{(*Compiler).floatLiteral /*   */, nil /*                     */, PrecNone}
	rules[token.KindAlias] = /*         */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindAnd] = /*           */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindBnum] = /*          */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindBool] = /*          */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindBreak] = /*         */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindCase] = /*          */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindClass] = /*         */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindContinue] = /*      */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindDo] = /*            */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindElse] = /*          */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindElseif] = /*        */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindEnd] = /*           */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindEnum] = /*          */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindFalse] = /*         */ parseRule{(*Compiler).boolLiteral /*    */, nil /*                     */, PrecNone}
	rules[token.KindFloat] = /*         */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindFor] = /*           */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindFunction] = /*      */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindGosub] = /*         */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindGoto] = /*          */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindIf] = /*            */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindIn] = /*            */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindInt] = /*           */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindListen] = /*        */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindMap] = /*           */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindMeta] = /*          */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindNil] = /*           */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindNot] = /*           */ parseRule{(*Compiler).unary /*          */, nil /*                     */, PrecNone}
	rules[token.KindOr] = /*            */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindPassage] = /*       */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindPrint] = /*         */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindReturn] = /*        */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindSay] = /*           */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindString] = /*        */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindStruct] = /*        */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindSuper] = /*         */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindSwitch] = /*        */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindThen] = /*          */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindTrue] = /*          */ parseRule{(*Compiler).boolLiteral /*    */, nil /*                     */, PrecNone}
	rules[token.KindVars] = /*          */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindVoid] = /*          */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindWhile] = /*         */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindError] = /*         */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
	rules[token.KindEOF] = /*           */ parseRule{nil /*                        */, nil /*                     */, PrecNone}
}
