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
