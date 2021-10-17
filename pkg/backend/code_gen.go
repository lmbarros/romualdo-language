/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2021 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package backend

import (
	"fmt"
	"math"

	"gitlab.com/stackedboxes/romulang/pkg/ast"
	"gitlab.com/stackedboxes/romulang/pkg/bytecode"
)

// GenerateCode generates the bytecode for a given AST.
func GenerateCode(root ast.Node) (chunk *bytecode.Chunk, err error) {
	cg := &codeGenerator{
		chunk:     bytecode.NewChunk(),
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

	// nodeStack is used to keep track of the nodes being processed. The current
	// on is on the top.
	nodeStack []ast.Node

	// locals holds the local variables currently in scope.
	locals []local

	// scopeDepth keeps track of the current scope depth we are in.
	//
	// TODO: How to interpret it? What is level zero? Global? Right at the start
	// of  a function declaration, is it at level one then?
	scopeDepth int
}

//
// The ast.Visitor interface
//

func (cg *codeGenerator) Enter(node ast.Node) {
	cg.nodeStack = append(cg.nodeStack, node)

	if _, ok := node.(*ast.Block); ok {
		cg.beginScope()
	}
}

func (cg *codeGenerator) Leave(node ast.Node) { // nolint: funlen, gocyclo
	switch n := node.(type) {
	case *ast.Storyworld:
		break

	case *ast.FloatLiteral:
		cg.emitConstant(bytecode.NewValueFloat(n.Value))

	case *ast.IntLiteral:
		cg.emitConstant(bytecode.NewValueInt(n.Value))

	case *ast.BNumLiteral:
		// For the VM, a BNum is just a float.
		cg.emitConstant(bytecode.NewValueFloat(n.Value))

	case *ast.BoolLiteral:
		if n.Value {
			cg.emitBytes(bytecode.OpTrue)
		} else {
			cg.emitBytes(bytecode.OpFalse)
		}

	case *ast.StringLiteral:
		cg.emitConstant(cg.newInternedValueString(n.Value))

	case *ast.Unary:
		switch n.Operator {
		case "+":
			break // no-op
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
			// If the type checker did its job, we can look only to the LHS here
			if n.LHS.Type().Tag == ast.TypeBNum {
				cg.emitBytes(bytecode.OpAddBNum)
			} else {
				cg.emitBytes(bytecode.OpAdd)
			}
		case "-":
			// If the type checker did its job, we can look only to the LHS here
			if n.LHS.Type().Tag == ast.TypeBNum {
				cg.emitBytes(bytecode.OpSubtractBNum)
			} else {
				cg.emitBytes(bytecode.OpSubtract)
			}
		case "*":
			cg.emitBytes(bytecode.OpMultiply)
		case "/":
			cg.emitBytes(bytecode.OpDivide)
		case "^":
			cg.emitBytes(bytecode.OpPower)
		default:
			cg.ice("unknown binary operator: %v", n.Operator)
		}

	case *ast.Blend:
		cg.emitBytes(bytecode.OpBlend)

	case *ast.TypeConversion:
		switch n.Operator {
		case "int":
			cg.emitBytes(bytecode.OpToInt)
		case "float":
			cg.emitBytes(bytecode.OpToFloat)
		case "bnum":
			cg.emitBytes(bytecode.OpToBNum)
		case "string":
			cg.emitBytes(bytecode.OpToString)
		default:
			cg.ice("unknown type conversion operator: %v", n.Operator)
		}

	case *ast.IfStmt:
		break

	case *ast.BuiltInFunction:
		if n.Function != "print" {
			cg.ice("only %q is supported, got %q", "print", n.Function)
		}
		cg.emitBytes(bytecode.OpPrint)

	case *ast.GlobalsBlock:
		break

	case *ast.VarDecl:
		if cg.isInsideGlobalsBlock() {
			// Global variable
			created := cg.chunk.SetGlobal(n.Name, cg.valueFromNode(n.Initializer))
			if !created {
				cg.ice(
					"duplicate definition of global variable '%v' on code generation",
					n.Name)
			}
			// Pop the global value that now is stored in cg.chunk.Globals:
			cg.emitBytes(bytecode.OpPop)
		} else {
			// Local variable
			if len(cg.locals) == 256 {
				cg.error("Currently only up to 255 global variables are supported.")
			}

			for _, local := range cg.locals {
				if local.name == n.Name {
					cg.error("Local variable %q already defined. Shadowing not allowed.", n.Name)
				}
			}

			cg.locals = append(cg.locals, local{name: n.Name, depth: cg.scopeDepth})
		}

	case *ast.VarRef:
		localIndex := cg.resolveLocal(n.Name)
		if localIndex < 0 {
			// It's a global
			i := cg.chunk.GetGlobalIndex(n.Name)
			if i < 0 {
				cg.ice("global variable '%v' not found in the globals pool", n.Name)
			}
			if i > 255 {
				// TODO: Can this even happen? I guess GetGlobalIndex will never return
				// anything over 255.
				cg.error("Currently only up to 255 global variables are supported.")
			}
			cg.emitBytes(bytecode.OpReadGlobal, byte(i))
		} else {
			// It's a local
			cg.emitBytes(bytecode.OpReadLocal, byte(localIndex))
		}

	case *ast.Assignment:
		localIndex := cg.resolveLocal(n.VarName)
		if localIndex < 0 {
			// It's a global
			i := cg.chunk.GetGlobalIndex(n.VarName)
			if i < 0 {
				cg.error("Global variable '%v' not declared.", n.VarName)
			}
			if i > 255 {
				// TODO: Can this even happen? I guess GetGlobalIndex will never return
				// anything over 255.
				cg.error("Currently only up to 255 global variables are supported.")
			}
			cg.emitBytes(bytecode.OpWriteGlobal, byte(i))
		} else {
			// It's a local
			cg.emitBytes(bytecode.OpWriteLocal, byte(localIndex))
		}

	case *ast.ExpressionStmt:
		cg.emitBytes(bytecode.OpPop)

	case *ast.Block:
		cg.endScope()

	default:
		cg.ice("unknown node type: %T", n)
	}

	cg.nodeStack = cg.nodeStack[:len(cg.nodeStack)-1]
}

func (cg *codeGenerator) Event(node ast.Node, event int) {
	if n, ok := node.(*ast.IfStmt); ok {
		switch event {

		// We initially emit a short jumps with placeholder jump offsets. We
		// update the jump offsets once we know the size of the code block that
		// will be jumped over.
		case ast.EventAfterIfCondition:
			n.IfJumpAddress = len(cg.chunk.Code)
			cg.emitBytes(bytecode.OpJumpIfFalse, 0x00)

		case ast.EventAfterThenBlock:
			addressToPatch := n.IfJumpAddress
			jumpOffset := len(cg.chunk.Code) - addressToPatch - 2
			if jumpOffset > math.MaxInt8 || jumpOffset < math.MinInt8 {
				cg.error("Jump offset of %v is longer than supported.", jumpOffset)
			}
			cg.chunk.Code[addressToPatch+1] = uint8(jumpOffset)

		case ast.EventBeforeElse:
			n.ElseJumpAddress = len(cg.chunk.Code)
			cg.emitBytes(bytecode.OpJump, 0x00)

			// Re-patch the "if" jump address, because the "else" block will
			// generate an additional jump (which takes two bytes).
			addressToPatch := n.IfJumpAddress
			jumpOffset := int(cg.chunk.Code[addressToPatch+1]) + 2
			if jumpOffset > math.MaxInt8 || jumpOffset < math.MinInt8 {
				cg.error("Jump offset of %v is longer than supported.", jumpOffset)
			}
			cg.chunk.Code[addressToPatch+1] = uint8(jumpOffset)

		case ast.EventAfterElse:
			addressToPatch := n.ElseJumpAddress
			jumpOffset := len(cg.chunk.Code) - addressToPatch - 2
			if jumpOffset > math.MaxInt8 || jumpOffset < math.MinInt8 {
				cg.error("Jump offset of %v is longer than supported.", jumpOffset)
			}
			cg.chunk.Code[addressToPatch+1] = uint8(jumpOffset)
		}
	}
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
		b0, b1, b2 := bytecode.UIntToThreeBytes(constantIndex)
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

// beginScope gets called when we enter into a new scope.
func (cg *codeGenerator) beginScope() {
	cg.scopeDepth++
}

// endScope gets called when we leave a scope.
func (cg *codeGenerator) endScope() {
	cg.scopeDepth--

	for len(cg.locals) > 0 && cg.locals[len(cg.locals)-1].depth > cg.scopeDepth {
		cg.emitBytes(bytecode.OpPop)
		cg.locals = cg.locals[:len(cg.locals)-1]
	}
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

// newInternedValueString creates a new Value initialized to the interned string
// value v. Emphasis on "interned": if there is already some other string value
// equal to v on this VM, we'll reuse that same memory in the returned value.
func (cg *codeGenerator) newInternedValueString(v string) bytecode.Value {
	s := cg.currentChunk().Strings.Intern(v)
	return bytecode.NewValueString(s)
}

func (cg *codeGenerator) valueFromNode(node ast.Node) bytecode.Value {
	switch n := node.(type) {
	case *ast.StringLiteral:
		return bytecode.Value{Value: n.Value}
	case *ast.BoolLiteral:
		return bytecode.Value{Value: n.Value}
	case *ast.IntLiteral:
		return bytecode.Value{Value: n.Value}
	case *ast.FloatLiteral:
		return bytecode.Value{Value: n.Value}
	case *ast.BNumLiteral:
		return bytecode.Value{Value: n.Value}
	default:
		cg.ice("Unexpected node of type %T", node)
	}
	return bytecode.Value{}
}

// isInsideGlobalsBlock checks if we are currently inside a globals block.
func (cg *codeGenerator) isInsideGlobalsBlock() bool {
	for _, node := range cg.nodeStack {
		_, ok := node.(*ast.GlobalsBlock)
		if ok {
			return true
		}
	}
	return false
}

// resolveLocal finds the index into the locals array of the local variable
// named name.
func (cg *codeGenerator) resolveLocal(name string) int {
	for i, local := range cg.locals {
		if local.name == name {
			return i
		}
	}

	return -1
}
