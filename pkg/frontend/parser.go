/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2022 Leandro Motta Barros                                     *
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
	precBlend                        // ~
	PrecUnary                        // not -
	precPower                        // ^
	precCall                         // . ()
	precPrimary
)

// prefixParseFn is a function used to parse code for a certain kind of prefix
// expression. canAssign tells if the expression we parsing accepts to be the
// target of an assignment.
type prefixParseFn = func(c *parser, canAssign bool) ast.Node

// infixParseFn is a function used to parse code for a certain kind of infix
// expression. lhs is the left-hand side expression previously parsed. canAssign
// tells if the expression we parsing accepts to be the target of an assignment.
type infixParseFn = func(c *parser, lhs ast.Node, canAssign bool) ast.Node

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
func (p *parser) parse() *ast.Storyworld {

	sw := ast.Storyworld{}

	p.advance()

	for !p.match(tokenKindEOF) {
		node := p.declaration()
		if p.hadError {
			return nil
		}
		sw.Declarations = append(sw.Declarations, node)
	}

	return &sw
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

	canAssign := prec <= precAssignment
	node := prefixRule(p, canAssign)

	for prec <= rules[p.currentToken.kind].precedence {
		p.advance()
		infixRule := rules[p.previousToken.kind].infix
		node = infixRule(p, node, canAssign)
	}

	if canAssign && p.match(tokenKindEqual) {
		p.error("Invalid assignment target.")
	}

	return node
}

//
// The parsing rules and compilation functions
//

// rules is the table of parsing rules for our Pratt parser.
var rules []parseRule

// declaration parses a declaration.
func (p *parser) declaration() ast.Node {
	var n ast.Node

	switch {
	case p.match(tokenKindGlobals):
		n = p.globalsDeclaration()
	case p.match(tokenKindFunction):
		n = p.functionDeclaration()
	default:
		p.errorAtCurrent("Expect a declaration")
	}

	if p.panicMode {
		p.synchronize()
	}

	return n
}

// statement parses a statement.
func (p *parser) statement() ast.Node {
	switch {
	// For now at least, built-in functions are called with a leading dot, like
	// this: .funcName(args).
	case p.match(tokenKindDot):
		if p.match(tokenKindIdentifier) {
			funcName := p.previousToken.lexeme
			return p.builtInFunction(funcName)
		}
		p.errorAtCurrent("Expect built-in function name.")

	case p.match(tokenKindIf):
		return p.ifStatement()

	case p.match(tokenKindWhile):
		return p.whileStatement()

	case p.match(tokenKindDo):
		return p.block()

	case p.match(tokenKindVar):
		return p.varDeclaration()

	default:
		expr := p.expression()
		return &ast.ExpressionStmt{
			BaseNode: ast.BaseNode{
				LineNumber: p.previousToken.line,
			},
			Expr: expr,
		}
	}

	panic("Can't happen")
}

// builtInFunction parses a built-in function named funcName. The current token
// should be the opening parenthesis after the function name.
func (p *parser) builtInFunction(funcName string) ast.Node {
	if funcName != "print" {
		p.errorAt(p.previousToken, "Unknown built-in function.")
	}

	p.consume(tokenKindLeftParen, "Expect '(' after function name.")

	bif := &ast.BuiltInFunction{
		BaseNode: ast.BaseNode{
			LineNumber: p.previousToken.line,
		},
		Function: funcName,
		Args:     []ast.Node{p.expression()}, // For now, support only funcs with arity 1
	}

	p.consume(tokenKindRightParen, "Expect ')' after expression.")

	return bif
}

// synchronize skips tokens until we find something that looks like a statement
// boundary. This is used to recover from panic mode.
func (p *parser) synchronize() {
	p.panicMode = false

	// TODO: This is basically the same as in Lox. Must adapt to Romualdo
	// more properly.
	for p.currentToken.kind != tokenKindEOF {
		switch p.currentToken.kind {
		case tokenKindClass:
		case tokenKindFunction:
		case tokenKindGlobals:
		case tokenKindPassage:
		case tokenKindFor:
		case tokenKindIf:
		case tokenKindWhile:
		case tokenKindReturn:
			return

		default:
			// Do nothing.
		}

		p.advance()
	}
}

// expression parses an expression.
func (p *parser) expression() ast.Node {
	return p.parsePrecedence(precAssignment)
}

// block parses a block, as in "block of code". What other block could it be?!
// This is a compiler! (My comments are usually more polite than this.)
func (p *parser) block() *ast.Block {
	block := &ast.Block{
		BaseNode: ast.BaseNode{
			LineNumber: p.previousToken.line,
		},
	}

	for !p.check(tokenKindEnd) && !p.check(tokenKindEOF) {
		stmt := p.statement()
		block.Statements = append(block.Statements, stmt)
	}
	p.consume(tokenKindEnd, fmt.Sprintf("Expect 'end' to block started at line %v.", block.LineNumber))
	return block
}

// globalsDeclaration parses a globals block.
func (p *parser) globalsDeclaration() ast.Node {
	globals := &ast.GlobalsBlock{
		BaseNode: ast.BaseNode{
			LineNumber: p.previousToken.line,
		},
	}

	for p.currentToken.kind != tokenKindEnd {
		v := p.varDeclaration()
		globals.Vars = append(globals.Vars, v)
	}

	// TODO: A kind compiler would tell the line where the block started.
	p.consume(tokenKindEnd, "Expect 'end' to close 'globals' block")

	return globals
}

// parseType parses a type. The first token of the type is supposed to have been
// just consumed.
func (p *parser) parseType() *ast.Type {
	switch p.previousToken.kind {
	case tokenKindInt:
		return ast.TheTypeInt
	case tokenKindFloat:
		return ast.TheTypeFloat
	case tokenKindBNum:
		return ast.TheTypeBNum
	case tokenKindString:
		return ast.TheTypeString
	case tokenKindBool:
		return ast.TheTypeBool
	case tokenKindVoid:
		return ast.TheTypeVoid
	case tokenKindFunction:
		p.consume(tokenKindLeftParen, "Expect '('")
		paramTypes := p.parseTypeList()
		p.consume(tokenKindColon, "Expect ':'")
		p.advance()
		retType := p.parseType()
		return &ast.Type{
			Tag:            ast.TypeFunction,
			ParameterTypes: paramTypes,
			ReturnType:     retType,
		}
	default:
		p.errorAtCurrent("Expect type.")
	}
	return ast.TheTypeInvalid
}

// parseParameterList parses a list of parameters. The left paretheses is
// supposed to have just been consumed.
func (p *parser) parseParameterList() []ast.Parameter {
	params := []ast.Parameter{}
	if p.check(tokenKindRightParen) {
		p.advance()
		return params
	}

	for ok := true; ok; ok = p.match(tokenKindComma) {
		p.consume(tokenKindIdentifier, "Expect identifier (the parameter name) or a right parenthesis")
		n := p.previousToken.lexeme
		p.consume(tokenKindColon, "Expect ':' after parameter name")
		p.advance()
		t := p.parseType()
		if t.Tag == ast.TypeVoid {
			p.errorAt(p.previousToken, "Cannot use 'void' as a parameter type")
		}
		params = append(params, ast.Parameter{Name: n, Type: t})
	}

	p.consume(tokenKindRightParen, "Expect ')' to close parameter list")

	return params
}

// parseTypeList parses a list of types. The left paretheses is supposed to have
// just been consumed.
func (p *parser) parseTypeList() []*ast.Type {
	types := []*ast.Type{}
	if p.check(tokenKindRightParen) {
		p.advance()
		return types
	}

	for ok := true; ok; ok = p.match(tokenKindComma) {
		p.advance()
		t := p.parseType()
		if t.Tag == ast.TypeVoid {
			p.errorAt(p.previousToken, "Cannot use 'void' as a parameter type")
		}
		types = append(types, t)
	}

	p.consume(tokenKindRightParen, "Expect ')' to close parameter list")

	return types
}

// varDeclaration parses a variable declaration. The next token is supposed to
// be the variable name.
func (p *parser) varDeclaration() *ast.VarDecl {
	p.consume(tokenKindIdentifier, "Expect identifier (the variable name).")
	name := p.previousToken.lexeme

	baseNode := ast.BaseNode{
		LineNumber: p.previousToken.line,
	}

	// TODO: Make type optional if initializer is present.
	p.consume(tokenKindColon, "Expect ':' after variable name.")

	p.advance()
	varType := p.parseType()

	// TODO: Make initializer optional (use default value if not provided).
	p.consume(tokenKindEqual, "Expect '=' after variable type.")

	initializer := p.expression()

	v := ast.NewVarDecl(baseNode, name, varType, initializer)
	v.BaseNode = ast.BaseNode{
		LineNumber: p.previousToken.line,
	}

	return v
}

// functionDeclaration parses a function declaration. The function keyword is
// expected to have just been consumed.
func (p *parser) functionDeclaration() *ast.FunctionDecl {
	f := &ast.FunctionDecl{
		BaseNode: ast.BaseNode{
			LineNumber: p.previousToken.line,
		},
	}

	p.consume(tokenKindIdentifier, "Expect identifier (the function name).")
	f.Name = p.previousToken.lexeme

	p.consume(tokenKindLeftParen, "Expect '(' after function name.")
	f.Parameters = p.parseParameterList()

	p.consume(tokenKindColon, "Expect ':' after parameter list.")
	p.advance()
	f.ReturnType = p.parseType()

	f.Body = p.block()

	return f
}

// ifStatement parses an if statement. The if keyword is expected to have just
// been consumed.
func (p *parser) ifStatement() ast.Node {
	n := &ast.IfStmt{
		BaseNode: ast.BaseNode{
			LineNumber: p.previousToken.line,
		},
	}

	n.Condition = p.expression()
	p.consume(tokenKindThen, "Expect 'then' after condition.")

	thenBlock := &ast.Block{
		BaseNode: ast.BaseNode{
			LineNumber: p.previousToken.line,
		},
	}

	for !(p.check(tokenKindEnd) || p.check(tokenKindElse) || p.check(tokenKindElseif)) && !p.check(tokenKindEOF) {
		stmt := p.statement()
		thenBlock.Statements = append(thenBlock.Statements, stmt)
	}
	n.Then = thenBlock

	switch {
	case p.match(tokenKindEnd):
		n.Else = nil

	case p.match(tokenKindElse):
		elseBlock := &ast.Block{
			BaseNode: ast.BaseNode{
				LineNumber: p.previousToken.line,
			},
		}
		for !p.check(tokenKindEnd) && !p.check(tokenKindEOF) {
			stmt := p.statement()
			elseBlock.Statements = append(elseBlock.Statements, stmt)
		}
		p.consume(tokenKindEnd, fmt.Sprintf("Expect: 'end' to close 'if' statement started at line %v'.", n.LineNumber))
		n.Else = elseBlock

	case p.match(tokenKindElseif):
		n.Else = p.ifStatement()

	default:
		p.error(fmt.Sprintf("Unterminated 'if' statement at line %v.", n.LineNumber))
	}
	return n
}

// whileStatement parses a while statement. The while keyword is expected to
// have just been consumed.
func (p *parser) whileStatement() ast.Node {
	n := &ast.WhileStmt{
		BaseNode: ast.BaseNode{
			LineNumber: p.previousToken.line,
		},
	}

	n.Condition = p.expression()
	p.consume(tokenKindDo, fmt.Sprintf("Expect: 'do' after 'while' condition ar line %v'.", n.LineNumber))

	n.Body = p.block()

	return n
}

// numberLiteral parses a number literal (int, float, or bnum). The number
// literal token is expected to have been just consumed.
func (p *parser) numberLiteral(canAssign bool) ast.Node {
	baseNode := ast.BaseNode{
		LineNumber: p.previousToken.line,
	}

	switch p.previousToken.kind {
	case tokenKindFloatLiteral:
		value, err := strconv.ParseFloat(p.previousToken.lexeme, 64)
		if err != nil {
			panic("Compiler got invalid float lexeme: " + p.previousToken.lexeme)
		}
		return &ast.FloatLiteral{
			BaseNode: baseNode,
			Value:    value,
		}

	case tokenKindIntLiteral:
		value, err := strconv.ParseInt(p.previousToken.lexeme, 10, 64)
		if err != nil {
			panic("Compiler got invalid int lexeme: " + p.previousToken.lexeme)
		}
		return &ast.IntLiteral{
			BaseNode: baseNode,
			Value:    value,
		}

	case tokenKindBNumLiteral:
		s := p.previousToken.lexeme
		s = s[0 : len(s)-1] // remove trailing "b"
		value, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic("Compiler got invalid bnum lexeme: " + p.previousToken.lexeme)
		}
		if value <= 0.0 || value >= 1.0 {
			p.error(fmt.Sprintf(
				"BNum must be greater than 0.0 and less than 1.0; got %v", value))
		}
		return &ast.BNumLiteral{
			BaseNode: baseNode,
			Value:    value,
		}

	default:
		panic(fmt.Sprintf("Got unexpected token kind when parsing number literal: %v (%v)",
			p.previousToken.kind, int64(p.previousToken.kind)))
	}
}

// stringLiteral parses a string literal. The string literal token is expected
// to have been just consumed.
func (p *parser) stringLiteral(canAssign bool) ast.Node {
	value := p.previousToken.lexeme[1 : len(p.previousToken.lexeme)-1] // remove the quotes

	return &ast.StringLiteral{
		BaseNode: ast.BaseNode{
			LineNumber: p.previousToken.line,
		},
		Value: value,
	}
}

// variable parses a variable reference or (if canAssign == true) assignment in
// the source code. The identifier token with the variable name is expected to
// have been just consumed.
//
// TODO: In the book, this just calls a function namedVariable, that does the
// actual work. I skipped this intermediary call that is useless for now, but I
// might need to add it later on, when it will start to be useful somehow.
func (p *parser) variable(canAssign bool) ast.Node {
	varName := p.previousToken.lexeme

	if canAssign && p.match(tokenKindEqual) {
		rhs := p.expression()
		return &ast.Assignment{
			BaseNode: ast.BaseNode{
				LineNumber: p.previousToken.line,
			},
			VarName: varName,
			Value:   rhs,
		}
	}

	return &ast.VarRef{
		BaseNode: ast.BaseNode{
			LineNumber: p.previousToken.line,
		},
		Name:    varName,
		VarType: ast.TheTypeInvalid, // Filled in a later pass
	}
}

// grouping parses a parenthesized expression. The left paren token is expected
// to have been just consumed.
func (p *parser) grouping(canAssign bool) ast.Node {
	expr := p.expression()
	p.consume(tokenKindRightParen, "Expect ')' after expression.")
	return expr
}

// unary parses a unary expression. The operator token is expected to have been
// just consumed.
func (p *parser) unary(canAssign bool) ast.Node {
	operatorKind := p.previousToken.kind
	operatorLexeme := p.previousToken.lexeme
	operatorLine := p.previousToken.line

	// Parse the operand.
	operand := p.parsePrecedence(PrecUnary)

	// Return the node.
	switch operatorKind {
	case tokenKindNot, tokenKindMinus, tokenKindPlus:
		return &ast.Unary{
			BaseNode: ast.BaseNode{
				LineNumber: operatorLine,
			},
			Operator: operatorLexeme,
			Operand:  operand,
		}
	default:
		panic(fmt.Sprintf("Unexpected operator kind on unary expression: %v", operatorKind))
	}
}

// binary parses a binary operator expression. The left operand and the operator
// token are expected to have been just consumed.
func (p *parser) binary(lhs ast.Node, canAssign bool) ast.Node {
	// Remember the operator.
	operatorKind := p.previousToken.kind
	operatorLexeme := p.previousToken.lexeme
	operatorLine := p.previousToken.line

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
			LineNumber: operatorLine,
		},
		Operator: operatorLexeme,
		LHS:      lhs,
		RHS:      rhs,
	}
}

// and parses an "and" expression. The left-hand-side argument and the "and"
// operator are expected to have been just consumed.
func (p *parser) and(lhs ast.Node, canAssign bool) ast.Node {
	rhs := p.parsePrecedence(precAnd)
	return &ast.And{
		BaseNode: ast.BaseNode{
			LineNumber: p.previousToken.line,
		},
		LHS: lhs,
		RHS: rhs,
	}
}

// or parses an "or" expression. The left-hand-side argument and the "or"
// operator are expected to have been just consumed.
func (p *parser) or(lhs ast.Node, canAssign bool) ast.Node {
	rhs := p.parsePrecedence(precAnd)
	return &ast.Or{
		BaseNode: ast.BaseNode{
			LineNumber: p.previousToken.line,
		},
		LHS: lhs,
		RHS: rhs,
	}
}

// blend parses a blend operator expression. The first operand (parameter x)
// and the first tilde token are expected to have been just consumed.
func (p *parser) blend(x ast.Node, canAssign bool) ast.Node {
	// Remember the operator
	operatorKind := p.previousToken.kind // always a ast.tokenKindTilde
	operatorLine := p.previousToken.line
	rule := rules[operatorKind]

	// Parse the second operand (y), the second tilde, and the third operand (weight)
	y := p.parsePrecedence(rule.precedence + 1)
	p.consume(tokenKindTilde, "Expect '~' after expression.")
	weight := p.parsePrecedence(rule.precedence + 1)

	return &ast.Blend{
		BaseNode: ast.BaseNode{
			LineNumber: operatorLine,
		},
		X:      x,
		Y:      y,
		Weight: weight,
	}
}

// boolLiteral parses a literal Boolean value. The corresponding keyword is
// expected to have been just consumed.
func (p *parser) boolLiteral(canAssign bool) ast.Node {
	switch p.previousToken.kind {
	case tokenKindTrue:
		return &ast.BoolLiteral{
			BaseNode: ast.BaseNode{
				LineNumber: p.previousToken.line,
			},
			Value: true,
		}
	case tokenKindFalse:
		return &ast.BoolLiteral{
			BaseNode: ast.BaseNode{
				LineNumber: p.previousToken.line,
			},
			Value: false,
		}
	default:
		panic(fmt.Sprintf("Unexpected token type on boolLiteral: %v", p.previousToken.kind))
	}
}

// typeConversion parses a type conversion expression. The corresponding keyword
// is expected to have been just consumed.
func (p *parser) typeConversion(canAssign bool) ast.Node {
	conversionLexeme := p.previousToken.lexeme

	// Consume the open paren and parse the expression to be converted
	p.consume(tokenKindLeftParen, "Expect '(' after conversion operator.")
	v := p.parsePrecedence(precAssignment)

	var d ast.Node
	bn := ast.BaseNode{
		LineNumber: p.previousToken.line,
	}

	// If we have a comma, consume it and parse the expression with the default
	// value; otherwise, use a, er, default default value.
	switch {
	// Check for string first, as it can't have a default (but the parser is
	// supposed to always generate a non-nil default, so we use an empty
	// string).
	case conversionLexeme == "string":
		d = &ast.StringLiteral{BaseNode: bn, Value: ""}

	// If not string, we may have an explicit default value
	case p.currentToken.kind == tokenKindComma:
		p.consume(tokenKindComma, "Expect ',' in conversion expresion.")
		d = p.parsePrecedence(precAssignment)

	// If no default is provided, use default default
	case conversionLexeme == "int":
		d = &ast.IntLiteral{BaseNode: bn, Value: 0}
	case conversionLexeme == "float":
		d = &ast.FloatLiteral{BaseNode: bn, Value: 0.0}
	case conversionLexeme == "bnum":
		d = &ast.BNumLiteral{BaseNode: bn, Value: 0.0}
	}

	// Consume the closing paren
	p.consume(tokenKindRightParen, "Expect ')' after conversion expresion.")

	// Voilà, return the node
	return &ast.TypeConversion{
		BaseNode: ast.BaseNode{
			LineNumber: p.previousToken.line,
		},
		Operator: conversionLexeme,
		Value:    v,
		Default:  d,
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

// check checks if the current token is of a given kind.
func (p *parser) check(kind tokenKind) bool {
	return p.currentToken.kind == kind
}

// match consumes the current token if it is of a given type and returns true;
// otherwise, it simply returns false without consuming any token.
func (p *parser) match(kind tokenKind) bool {
	if !p.check(kind) {
		return false
	}
	p.advance()
	return true
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
	rules[tokenKindLeftParen] = /*     */ parseRule{(*parser).grouping /*         */, nil /*                     */, precNone}
	rules[tokenKindRightParen] = /*    */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindLeftBrace] = /*     */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindRightBrace] = /*    */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindLeftBracket] = /*   */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindRightBracket] = /*  */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindComma] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindDot] = /*           */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindMinus] = /*         */ parseRule{(*parser).unary /*            */, (*parser).binary /*        */, precTerm}
	rules[tokenKindPlus] = /*          */ parseRule{(*parser).unary /*            */, (*parser).binary /*        */, precTerm}
	rules[tokenKindSlash] = /*         */ parseRule{nil /*                        */, (*parser).binary /*        */, precFactor}
	rules[tokenKindStar] = /*          */ parseRule{nil /*                        */, (*parser).binary /*        */, precFactor}
	rules[tokenKindColon] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindTilde] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindAt] = /*            */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindHat] = /*           */ parseRule{nil /*                        */, (*parser).binary /*        */, precPower}
	rules[tokenKindEqual] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindEqualEqual] = /*    */ parseRule{nil /*                        */, (*parser).binary /*        */, precEquality}
	rules[tokenKindBangEqual] = /*     */ parseRule{nil /*                        */, (*parser).binary /*        */, precEquality}
	rules[tokenKindGreater] = /*       */ parseRule{nil /*                        */, (*parser).binary /*        */, precComparison}
	rules[tokenKindGreaterEqual] = /*  */ parseRule{nil /*                        */, (*parser).binary /*        */, precComparison}
	rules[tokenKindLess] = /*          */ parseRule{nil /*                        */, (*parser).binary /*        */, precComparison}
	rules[tokenKindLessEqual] = /*     */ parseRule{nil /*                        */, (*parser).binary /*        */, precComparison}
	rules[tokenKindIdentifier] = /*    */ parseRule{(*parser).variable /*         */, nil /*                     */, precNone}
	rules[tokenKindStringLiteral] = /* */ parseRule{(*parser).stringLiteral /*    */, nil /*                     */, precNone}
	rules[tokenKindFloatLiteral] = /*  */ parseRule{(*parser).numberLiteral /*    */, nil /*                     */, precNone}
	rules[tokenKindAlias] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindAnd] = /*           */ parseRule{nil /*                        */, (*parser).and /*           */, precAnd}
	rules[tokenKindBNum] = /*          */ parseRule{(*parser).typeConversion /*   */, nil /*                     */, precNone}
	rules[tokenKindBNumLiteral] = /*   */ parseRule{(*parser).numberLiteral /*    */, nil /*                     */, precNone}
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
	rules[tokenKindFalse] = /*         */ parseRule{(*parser).boolLiteral /*      */, nil /*                     */, precNone}
	rules[tokenKindFloat] = /*         */ parseRule{(*parser).typeConversion /*   */, nil /*                     */, precNone}
	rules[tokenKindFor] = /*           */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindFunction] = /*      */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindGlobals] = /*       */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindGosub] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindGoto] = /*          */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindIf] = /*            */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindIn] = /*            */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindInt] = /*           */ parseRule{(*parser).typeConversion /*   */, nil /*                     */, precNone}
	rules[tokenKindIntLiteral] = /*    */ parseRule{(*parser).numberLiteral /*    */, nil /*                     */, precNone}
	rules[tokenKindListen] = /*        */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindMap] = /*           */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindMeta] = /*          */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindNil] = /*           */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindNot] = /*           */ parseRule{(*parser).unary /*            */, nil /*                     */, precNone}
	rules[tokenKindOr] = /*            */ parseRule{nil /*                        */, (*parser).or /*            */, precOr}
	rules[tokenKindPassage] = /*       */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindReturn] = /*        */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindSay] = /*           */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindString] = /*        */ parseRule{(*parser).typeConversion /*   */, nil /*                     */, precNone}
	rules[tokenKindStruct] = /*        */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindSuper] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindSwitch] = /*        */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindThen] = /*          */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindTilde] = /*         */ parseRule{nil /*                        */, (*parser).blend /*         */, precBlend}
	rules[tokenKindTrue] = /*          */ parseRule{(*parser).boolLiteral /*      */, nil /*                     */, precNone}
	rules[tokenKindVar] = /*           */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindVoid] = /*          */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindWhile] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindError] = /*         */ parseRule{nil /*                        */, nil /*                     */, precNone}
	rules[tokenKindEOF] = /*           */ parseRule{nil /*                        */, nil /*                     */, precNone}
}
