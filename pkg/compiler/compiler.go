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
	PrecPower                        // ^
	PrecUnary                        // not -
	PrecCall                         // . ()
	PrecPrimary
)

// compileFn is a function used to parse and generate code for a certain kind of
// expression.
type compileFn = func(c *Compiler)

// parseRule encodes one rule of our Pratt parser.
type parseRule struct {
	prefix     compileFn  // For expressions using the token as a prefix operator.
	infix      compileFn  // For expressions using the token as an infix operator.
	precedence precedence // When the token is used as a binary operator.
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
func (c *Compiler) Compile(source string) *bytecode.Chunk {
	// TODO: Candidate to change: I'd like to instantiate the scanner on New().
	c.s = scanner.New(source)

	c.advance()
	c.expression()
	c.consume(token.KindEOF, "Expect end of expression.")

	c.endCompiler()

	if c.p.hadError {
		return nil
	} else {
		return c.chunk
	}
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
func (c *Compiler) parsePrecedence(p precedence) {
	c.advance()
	prefixRule := rules[c.p.previous.Kind].prefix
	if prefixRule == nil {
		c.error("Expect expression.")
		return
	}

	prefixRule(c)

	for p <= rules[c.p.current.Kind].precedence {
		c.advance()
		infixRule := rules[c.p.previous.Kind].infix
		infixRule(c)
	}
}

//
// The parsing rules and compilation functions
//

// rules is the table of parsing rules for our Pratt parser.
var rules []parseRule

// expression parses and generates code for an expression.
func (c *Compiler) expression() {
	c.parsePrecedence(PrecAssignment)
}

// floatNumber parses and generates code for a number literal. The float literal
// token is expected to have been just consumed.
func (c *Compiler) floatNumber() {
	value, err := strconv.ParseFloat(c.p.previous.Lexeme, 64)
	if err != nil {
		panic("Compiler got invalid number lexeme: " + c.p.previous.Lexeme)
	}
	c.emitConstant(bytecode.Value(value))
}

// grouping parses and generates code for a parenthesized expression. The left
// paren token is expected to have been just consumed.
func (c *Compiler) grouping() {
	c.expression()
	c.consume(token.KindRightParen, "Expect ')' after expression.")
}

// unary parses and generates code for a unary expression. The operator token is
// expected to have been just consumed.
func (c *Compiler) unary() {
	operatorKind := c.p.previous.Kind

	// Compile the operand.
	c.parsePrecedence(PrecUnary)

	// Emit the operator instruction.
	switch operatorKind {
	case token.KindMinus:
		c.emitBytes(bytecode.OpNegate)
	default:
		panic(fmt.Sprintf("Unexpected operator kind on unary expression: %v", operatorKind))
	}
}

// binary parses and generates code for a binary operator expression. The left
// operand and the operator token are expected to have been just consumed.
func (c *Compiler) binary() {
	// Remember the operator.
	operatorKind := c.p.previous.Kind

	// Compile the right operand.
	rule := rules[operatorKind]
	c.parsePrecedence(rule.precedence + 1)

	// Emit the operator instruction.
	switch operatorKind {
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
}

// initRules initializes the rules array.
//
// Using block comments to convince gofmt to keep things aligned is ugly as
// hell.
func initRules() { // nolint:funlen
	rules = make([]parseRule, token.NumberOfKinds)

	//                                     prefix                                   infix                     precedence
	//                                    -------------------------------------     ---------------------     --------------
	rules[token.KindLeftParen] = /*     */ parseRule{(*Compiler).grouping /*    */, nil /*                */, PrecNone}
	rules[token.KindRightParen] = /*    */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindLeftBrace] = /*     */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindRightBrace] = /*    */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindLeftBracket] = /*   */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindRightBracket] = /*  */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindComma] = /*         */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindDot] = /*           */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindMinus] = /*         */ parseRule{(*Compiler).unary /*       */, (*Compiler).binary /* */, PrecTerm}
	rules[token.KindPlus] = /*          */ parseRule{nil /*                     */, (*Compiler).binary /* */, PrecTerm}
	rules[token.KindSlash] = /*         */ parseRule{nil /*                     */, (*Compiler).binary /* */, PrecFactor}
	rules[token.KindStar] = /*          */ parseRule{nil /*                     */, (*Compiler).binary /* */, PrecFactor}
	rules[token.KindColon] = /*         */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindTilde] = /*         */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindAt] = /*            */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindHat] = /*           */ parseRule{nil /*                     */, (*Compiler).binary /* */, PrecPower}
	rules[token.KindEqual] = /*         */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindEqualEqual] = /*    */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindBangEqual] = /*     */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindGreater] = /*       */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindGreaterEqual] = /*  */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindLess] = /*          */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindLessEqual] = /*     */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindIdentifier] = /*    */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindStringLiteral] = /* */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindNumberLiteral] = /* */ parseRule{(*Compiler).floatNumber /* */, nil /*                */, PrecNone}
	rules[token.KindAlias] = /*         */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindAnd] = /*           */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindBnum] = /*          */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindBool] = /*          */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindBreak] = /*         */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindCase] = /*          */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindClass] = /*         */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindContinue] = /*      */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindDo] = /*            */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindElse] = /*          */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindElseif] = /*        */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindEnd] = /*           */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindEnum] = /*          */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindFalse] = /*         */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindFloat] = /*         */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindFor] = /*           */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindFunction] = /*      */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindGosub] = /*         */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindGoto] = /*          */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindIf] = /*            */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindIn] = /*            */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindInt] = /*           */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindListen] = /*        */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindMap] = /*           */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindMeta] = /*          */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindNil] = /*           */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindNot] = /*           */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindOr] = /*            */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindPassage] = /*       */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindPrint] = /*         */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindReturn] = /*        */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindSay] = /*           */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindString] = /*        */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindStruct] = /*        */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindSuper] = /*         */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindSwitch] = /*        */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindThen] = /*          */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindTrue] = /*          */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindVars] = /*          */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindVoid] = /*          */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindWhile] = /*         */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindError] = /*         */ parseRule{nil /*                     */, nil /*                */, PrecNone}
	rules[token.KindEOF] = /*           */ parseRule{nil /*                     */, nil /*                */, PrecNone}
}
