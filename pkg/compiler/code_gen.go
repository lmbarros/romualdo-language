/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package compiler

import (
	"fmt"
	"math"

	"gitlab.com/stackedboxes/romulang/pkg/ast"
	"gitlab.com/stackedboxes/romulang/pkg/bytecode"
)

func GenerateCode(root ast.Node) (chunk *bytecode.Chunk, err error) {
	cg := &codeGenerator{
		chunk:     &bytecode.Chunk{},
		nodeStack: make([]ast.Node, 0, 64),
	}

	defer func() {
		if r := recover(); r != nil {
			chunk = nil
			if e, ok := r.(*codeGeneratorError); ok {
				err = e
				return
			}
			panic(fmt.Sprintf("Unexpected error type: %T", r))
		}
	}()

	root.Walk(cg)
	cg.emitBytes(bytecode.OpReturn)
	return cg.chunk, nil
}

// codeGeneratorError is a type used in panics to report an error in code
// generation.
type codeGeneratorError struct {
	msg string
}

func (e *codeGeneratorError) Error() string {
	return e.msg
}

// codeGenerator is a visitor that generates a compiled Chunk from an AST.
type codeGenerator struct {
	// Chunk is the chunk of bytecode being generated.
	chunk *bytecode.Chunk

	// currentNode is the node currently being processed.
	nodeStack []ast.Node
}

//
// The ast.Visitor interface
//

func (cg *codeGenerator) Enter(node ast.Node) {
	cg.nodeStack = append(cg.nodeStack, node)
}

func (cg *codeGenerator) Leave(node ast.Node) {
	switch n := node.(type) {
	case *ast.FloatLiteral:
		cg.emitConstant(bytecode.NewValueFloat(n.Value))

	case *ast.IntLiteral:
		// TODO

	case *ast.BoolLiteral:
		// TODO

	case *ast.StringLiteral:
		cg.emitConstant(bytecode.NewValueString(n.Value))

	case *ast.Unary:
		switch n.Operator {
		case "-":
			cg.emitBytes(bytecode.OpNegate)
		case "not":
			cg.emitBytes(bytecode.OpNot)
		default:
			cg.ice("unknown unary operator: %v", n.Operator)
		}

	case *ast.Binary:
		switch n.Operator {
		case "!=":
			cg.emitBytes(bytecode.OpNotEqual)
		case "==":
			cg.emitBytes(bytecode.OpEqual)
		case ">":
			cg.emitBytes(bytecode.OpGreater)
		case ">=":
			cg.emitBytes(bytecode.OpGreaterEqual)
		case "<":
			cg.emitBytes(bytecode.OpLess)
		case "<=":
			cg.emitBytes(bytecode.OpLessEqual)
		case "+":
			cg.emitBytes(bytecode.OpAdd)
		case "-":
			cg.emitBytes(bytecode.OpSubtract)
		case "*":
			cg.emitBytes(bytecode.OpMultiply)
		case "/":
			cg.emitBytes(bytecode.OpDivide)
		case "^":
			cg.emitBytes(bytecode.OpPower)
		default:
			cg.ice("unknown binary operator: %v", n.Operator)
		}

	default:
		cg.ice("unknown node type: %T", n)
	}

	cg.nodeStack = cg.nodeStack[:len(cg.nodeStack)-1]

}

//
// Actual code generation
//

// currentLine returns the source code line corresponding to whatever we are
// currently compiling.
func (cg *codeGenerator) currentLine() int {
	if len(cg.nodeStack) == 0 {
		return -1 // TODO: Hack for that forced RETURN we generate out of no real node.
	}
	return cg.nodeStack[len(cg.nodeStack)-1].Line()
}

// currentChunk returns the current chunk we are compiling into.
func (cg *codeGenerator) currentChunk() *bytecode.Chunk {
	return cg.chunk
}

// emitBytes writes one or more bytes to the bytecode chunk being generated.
func (cg *codeGenerator) emitBytes(bytes ...byte) {
	for _, b := range bytes {
		cg.currentChunk().Write(b, cg.currentLine())
	}
}

// emitConstant emits the bytecode for a constant having a given value.
func (cg *codeGenerator) emitConstant(value bytecode.Value) {
	constantIndex := cg.makeConstant(value)
	if constantIndex <= math.MaxUint8 {
		cg.emitBytes(bytecode.OpConstant, byte(constantIndex))
	} else {
		b0, b1, b2 := bytecode.IntToThreeBytes(constantIndex)
		cg.emitBytes(bytecode.OpConstantLong, b0, b1, b2)
	}
}

// makeConstant adds value to the pool of constants and returns the index in
// which it was added. If there is already a constant with this value, its index
// is returned (hey, we don't need duplicate constants, right? They are
// constant, after all!)
func (cg *codeGenerator) makeConstant(value bytecode.Value) int {
	if i := cg.currentChunk().SearchConstant(value); i >= 0 {
		return i
	}

	constantIndex := cg.currentChunk().AddConstant(value)
	if constantIndex >= bytecode.MaxConstantsPerChunk {
		cg.error("Too many constants in one chunk.")
		return 0
	}

	return constantIndex
}

// error panics, reporting an error on the current node with a given error
// message.
func (cg *codeGenerator) error(format string, a ...interface{}) {
	e := &codeGeneratorError{
		msg: fmt.Sprintf("[line %v]: %v", cg.currentLine(),
			fmt.Sprintf(format, a...)),
	}
	panic(e)
}

// ice reports an internal compiler error.
func (cg *codeGenerator) ice(format string, a ...interface{}) {
	cg.error(fmt.Sprintf("Internal compiler error: %v", fmt.Sprintf(format, a...)))
}
