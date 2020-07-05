/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package ast

// FloatLiteral is an AST node representing a floating point number literal.
type FloatLiteral struct {
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
	// Operator contains the lexeme used as the binary operator.
	Operator string

	// LHS is the expression on the left-hand side of the operator.
	LHS Node

	// RHS is the expression on the right-hand side of the operator.
	RHS Node
}

func (n *Binary) Type() Type {
	return n.LHS.Type()
}

func (n *Binary) Walk(v Visitor) {
	v.Enter(n)
	n.LHS.Walk(v)
	n.RHS.Walk(v)
	v.Leave(n)
}

// Grouping is an AST node representing a parenthesized expression.
//
// TODO: This is silly, isn't it? I don't need a grouping node, because the AST
// structure itself represents the grouping. Gotta get rid of this!
type Grouping struct {
	// Expr is the the expression in parentheses.
	Expr Node
}

func (n *Grouping) Type() Type {
	return n.Expr.Type()
}

func (n *Grouping) Walk(v Visitor) {
	v.Enter(n)
	n.Expr.Walk(v)
	v.Leave(n)
}
