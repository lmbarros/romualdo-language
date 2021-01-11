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
}

func (n *BaseNode) Line() int {
	return n.LineNumber
}

// Storyworld is an AST node representing the whole storyworld. It is the root
// of the AST.
type Storyworld struct {
	BaseNode

	// Declarations stores all the declarations that make up the Storyworld.
	Declarations []Node
}

func (n *Storyworld) Type() Type {
	return Type{TypeVoid}
}

func (n *Storyworld) Walk(v Visitor) {
	v.Enter(n)
	for _, decl := range n.Declarations {
		decl.Walk(v)
	}
	v.Leave(n)
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

// BNumLiteral is an AST node representing a bounded number (bnum) literal.
type BNumLiteral struct {
	BaseNode

	// Value is the bnum literal's value.
	Value float64
}

func (n *BNumLiteral) Type() Type {
	return Type{TypeBNum}
}

func (n *BNumLiteral) Walk(v Visitor) {
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

func (n *Binary) Type() Type { // nolint: gocognit
	switch n.Operator {
	case "==", "!=", "<", "<=", ">", ">=":
		return Type{TypeBool}
	case "+", "-", "*":
		if n.LHS.Type().Tag == TypeString || n.LHS.Type().Tag == TypeBNum {
			return n.LHS.Type()
		}
		if n.LHS.Type().Tag == TypeInt && n.RHS.Type().Tag == TypeInt {
			return n.LHS.Type()
		}
		return Type{TypeFloat}

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

// Blend is an AST node representing a blend operator.
type Blend struct {
	BaseNode

	// X is the first of the BNumbers to be blended.
	X Node

	// Y is the second of the BNumbers to be blended.
	Y Node

	// Weight is the BNumber to be used as the blend weighting factor.
	Weight Node
}

func (n *Blend) Type() Type {
	return Type{TypeBNum}
}

func (n *Blend) Walk(v Visitor) {
	v.Enter(n)
	n.X.Walk(v)
	n.Y.Walk(v)
	n.Weight.Walk(v)
	v.Leave(n)
}

// TypeConversion is an AST node representing a type conversion expression.
type TypeConversion struct {
	BaseNode

	// Operator contains the lexeme used as the conversion operator.
	Operator string

	// Value is the value to be converted.
	Value Node

	// Default is the default value to return if the conversion fails. This
	// can't be nil, the parser must provide one even if the code itself
	// doesn't.
	Default Node
}

func (n *TypeConversion) Type() Type {
	switch n.Operator {
	case "int":
		return Type{TypeInt}
	case "float":
		return Type{TypeFloat}
	case "bnum":
		return Type{TypeBNum}
	case "string":
		return Type{TypeString}
	default:
		return Type{TypeInvalid}
	}
}

func (n *TypeConversion) Walk(v Visitor) {
	v.Enter(n)
	n.Value.Walk(v)
	if n.Operator != "string" {
		n.Default.Walk(v)
	}
	v.Leave(n)
}

// BuiltInFunction is an AST node representing a Romualdo built-in function.
type BuiltInFunction struct {
	BaseNode

	// Function contains the name of the built-in function used here.
	Function string

	// Args contains the arguments passed to the built-in function.
	Args []Node
}

func (n *BuiltInFunction) Type() Type {
	switch n.Function {
	case "print":
		return Type{TypeVoid}
	default:
		return Type{TypeInvalid}
	}
}

func (n *BuiltInFunction) Walk(v Visitor) {
	v.Enter(n)
	for _, arg := range n.Args {
		arg.Walk(v)
	}
	v.Leave(n)
}
