/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2021 Leandro Motta Barros                                     *
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

	// cachedType caches the type of this node. Used to memoize Type().
	cachedType *Type
}

func (n *Unary) Type() Type {
	if n.cachedType == nil {
		t := n.Operand.Type()
		n.cachedType = &t
	}
	return *n.cachedType
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

	// cachedType caches the type of this node. Used to memoize Type().
	cachedType *Type
}

func (n *Binary) Type() Type { // nolint: gocognit

	if n.cachedType == nil {
		switch n.Operator {
		case "==", "!=", "<", "<=", ">", ">=":
			n.cachedType = &Type{TypeBool}
		case "+", "-", "*":
			if n.LHS.Type().Tag == TypeString || n.LHS.Type().Tag == TypeBNum {
				t := n.LHS.Type()
				n.cachedType = &t
			} else if n.LHS.Type().Tag == TypeInt && n.RHS.Type().Tag == TypeInt {
				t := n.LHS.Type()
				n.cachedType = &t
			} else {
				n.cachedType = &Type{TypeFloat}
			}
		default:
			n.cachedType = &Type{TypeFloat}
		}
	}

	return *n.cachedType
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

// GlobalsBlock is an AST node representing a globals block.
type GlobalsBlock struct {
	BaseNode

	// Vars contains the variables defined in this block.
	Vars []*VarDecl
}

func (n *GlobalsBlock) Type() Type {
	return Type{TypeVoid}
}

func (n *GlobalsBlock) Walk(v Visitor) {
	v.Enter(n)
	for _, varDecl := range n.Vars {
		varDecl.Walk(v)
	}
	v.Leave(n)
}

// VarDecl is an AST node representing a single variable declaration.
type VarDecl struct {
	BaseNode

	// Name is teh variable name.
	Name string

	// Initializer is the expression used to initialize thr variable.
	Initializer Node

	// varType is the variable type. Use Type() to get it.
	varType Type
}

// NewVarDecl creates a new VarDecl, with the given name, type and initializer.
func NewVarDecl(baseNode BaseNode, name string, varType Type, initializer Node) *VarDecl {
	return &VarDecl{
		BaseNode:    baseNode,
		Name:        name,
		varType:     varType,
		Initializer: initializer,
	}
}

func (n *VarDecl) Type() Type {
	return n.varType
}

func (n *VarDecl) Walk(v Visitor) {
	v.Enter(n)
	n.Initializer.Walk(v)
	v.Leave(n)
}

// VarRef is an AST node representing a reference to a variable. (I mean, a
// variable being used in the code.)
type VarRef struct {
	BaseNode

	// Name is the variable name.
	Name string

	// VarType is the variable type.
	VarType Type
}

func (n *VarRef) Type() Type {
	return n.VarType
}

func (n *VarRef) Walk(v Visitor) {
	v.Enter(n)
	v.Leave(n)
}

// Assignment is an AST node representing an assignment.
type Assignment struct {
	BaseNode

	// VarName is the left-hand side of the assignment. Contains the name of the
	// variable we are assigning to.
	VarName string

	// Value is the right-hand side of the assignment. Contains the value we are
	// assigning to the variable.
	Value Node
}

func (n *Assignment) Type() Type {
	return n.Value.Type()
}

func (n *Assignment) Walk(v Visitor) {
	v.Enter(n)
	n.Value.Walk(v)
	v.Leave(n)
}

// ExpressionStmt is an AST node representing an expression when used as a
// statement. The important point about it is that the expression value is
// discarded.
type ExpressionStmt struct {
	BaseNode

	// Expr is the expression used as a statement.
	Expr Node
}

func (n *ExpressionStmt) Type() Type {
	return Type{TypeVoid}
}

func (n *ExpressionStmt) Walk(v Visitor) {
	v.Enter(n)
	n.Expr.Walk(v)
	v.Leave(n)
}

// Block is an AST node representing a block (specificilly, a block of code).
type Block struct {
	BaseNode

	// The statements that make up this block.
	Statements []Node
}

func (n *Block) Type() Type {
	return Type{TypeVoid}
}

func (n *Block) Walk(v Visitor) {
	v.Enter(n)
	for _, stmt := range n.Statements {
		stmt.Walk(v)
	}
	v.Leave(n)
}

// IfStmt is an AST node representing an if statement.
type IfStmt struct {
	BaseNode

	// Condition is the if condition.
	Condition Node

	// Then is the block of code executed if the condition is true.
	Then Block

	// Else is the code executed if the condition is false. Can be either a
	// proper block or an `if` statement (in the case of an `elseif`). Might
	// also be nil (when no `else` block is present).
	Else Node

	//
	// Fields used for code generation
	//

	// IfJumpAddress is the address of the jump instruction used for the "if".
	IfJumpAddress int

	// ElseJumpAddress is the address of the jump instruction emitted right
	// before the "else" block.
	ElseJumpAddress int
}

func (n *IfStmt) Type() Type {
	return Type{TypeVoid}
}

func (n *IfStmt) Walk(v Visitor) {
	v.Enter(n)
	n.Condition.Walk(v)
	v.Event(n, EventAfterIfCondition)
	n.Then.Walk(v)
	v.Event(n, EventAfterThenBlock)
	if n.Else != nil {
		v.Event(n, EventBeforeElse)
		n.Else.Walk(v)
		v.Event(n, EventAfterElse)
	}
	v.Leave(n)
}

// And is an AST node representing an "and" expression.
type And struct {
	BaseNode

	// LHS is the expression on the left-hand-side of the expression.
	LHS Node

	// RHS is the expression on the right-hand-side of the expression.
	RHS Node

	//
	// Fields used for code generation
	//

	// JumpAddress is the address of the jump instruction used short-circuiting
	// the execution of the "and".
	JumpAddress int
}

func (n *And) Type() Type {
	return Type{TypeBool}
}

func (n *And) Walk(v Visitor) {
	v.Enter(n)
	n.LHS.Walk(v)
	v.Event(n, EventAfterLogicalBinaryOp)
	n.RHS.Walk(v)
	v.Leave(n)
}

// Or is an AST node representing an "or" expression.
type Or struct {
	BaseNode

	// LHS is the expression on the left-hand-side of the expression.
	LHS Node

	// RHS is the expression on the right-hand-side of the expression.
	RHS Node

	//
	// Fields used for code generation
	//

	// JumpAddress is the address of the jump instruction used short-circuiting
	// the execution of the "or".
	JumpAddress int
}

func (n *Or) Type() Type {
	return Type{TypeBool}
}

func (n *Or) Walk(v Visitor) {
	v.Enter(n)
	n.LHS.Walk(v)
	v.Event(n, EventAfterLogicalBinaryOp)
	n.RHS.Walk(v)
	v.Leave(n)
}
