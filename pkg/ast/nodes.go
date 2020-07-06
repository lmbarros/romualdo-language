/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package ast

// BaseNode contains the functionality common to all AST nodes.
type BaseNode struct {
	// LineNumber stores the line number from where this node comes.
	LineNumber int

	// SourceLexeme stores the lexeme from where this node comes.
	SourceLexeme string
}

func (n *BaseNode) Line() int {
	return n.LineNumber
}

func (n *BaseNode) Lexeme() string {
	return n.SourceLexeme
}

// FloatLiteral is an AST node representing a floating point number literal.
type FloatLiteral struct {
	BaseNode

	// Value is the float literal's value.
	Value float64
}

func (n *FloatLiteral) Type() Type {
	return Type{TypeFloat}
}

func (n *FloatLiteral) Walk(v Visitor) {
	v.Enter(n)
	v.Leave(n)
}

// IntLiteral is an AST node representing an integer number literal.
type IntLiteral struct {
	BaseNode

	// Value is the int literal's value.
	Value int64
}

func (n *IntLiteral) Type() Type {
	return Type{TypeInt}
}

func (n *IntLiteral) Walk(v Visitor) {
	v.Enter(n)
	v.Leave(n)
}

// BoolLiteral is an AST node representing a Boolean value literal.
type BoolLiteral struct {
	BaseNode

	// Value is the bool literal's value.
	Value bool
}

func (n *BoolLiteral) Type() Type {
	return Type{TypeBool}
}

func (n *BoolLiteral) Walk(v Visitor) {
	v.Enter(n)
	v.Leave(n)
}

// StringLiteral is an AST node representing a string value literal.
type StringLiteral struct {
	BaseNode

	// Value is the string literal's value.
	Value string
}

func (n *StringLiteral) Type() Type {
	return Type{TypeString}
}

func (n *StringLiteral) Walk(v Visitor) {
	v.Enter(n)
	v.Leave(n)
}

// Unary is an AST node representing a unary operator.
type Unary struct {
	BaseNode

	// Operator contains the lexeme used as the unary operator.
	Operator string

	// Operand is the expression on which the operator is applied.
	Operand Node
}

func (n *Unary) Type() Type {
	return n.Operand.Type()
}

func (n *Unary) Walk(v Visitor) {
	v.Enter(n)
	n.Operand.Walk(v)
	v.Leave(n)
}

// Binary is an AST node representing a binary operator.
type Binary struct {
	BaseNode

	// Operator contains the lexeme used as the binary operator.
	Operator string

	// LHS is the expression on the left-hand side of the operator.
	LHS Node

	// RHS is the expression on the right-hand side of the operator.
	RHS Node
}

func (n *Binary) Type() Type {
	switch n.Operator {
	case "==", "!=", "<", "<=", ">", ">=":
		return Type{TypeBool}
	case "+":
		return n.LHS.Type()
	default:
		return Type{TypeFloat}
	}
}

func (n *Binary) Walk(v Visitor) {
	v.Enter(n)
	n.LHS.Walk(v)
	n.RHS.Walk(v)
	v.Leave(n)
}
