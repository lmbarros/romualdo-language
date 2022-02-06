/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2022 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package backend

import (
	"math"

	"gitlab.com/stackedboxes/romulang/pkg/ast"
	"gitlab.com/stackedboxes/romulang/pkg/bytecode"
)

// codeGeneratorPassTwo does the actual bytecode generation. It fills in the
// Chunks with bytecode.
//
// This implements the ast.Visitor interface.
type codeGeneratorPassTwo struct {
	codeGenerator *codeGenerator

	// locals holds the local variables currently in scope.
	locals []local
}

//
// The ast.Visitor interface
//

func (cg *codeGeneratorPassTwo) Enter(node ast.Node) {
	cg.codeGenerator.pushIntoNodeStack(node)

	switch n := node.(type) {
	case *ast.Block:
		cg.codeGenerator.beginScope()

	case *ast.WhileStmt:
		n.ConditionAddress = len(cg.currentChunk().Code)

	case *ast.FunctionDecl:
		// Even though the function body is already a Block that does the
		// scoping little dance, we do it also for function declarations -- here
		// and on Leave(), where we pop the descoped arguments. Not sure this is
		// the most elegant way to do it.
		cg.codeGenerator.beginScope()

		for _, param := range n.Parameters {
			if !cg.defineLocalVariable(param.Name) {
				break
			}
		}

	default:
		// nothing
	}
}

func (cg *codeGeneratorPassTwo) Leave(node ast.Node) { // nolint: funlen, gocyclo
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
			cg.codeGenerator.ice("unknown unary operator: %v", n.Operator)
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
			cg.codeGenerator.ice("unknown binary operator: %v", n.Operator)
		}

	case *ast.And:
		addressToPatch := n.JumpAddress
		jumpOffset := len(cg.currentChunk().Code) - addressToPatch - 2
		cg.patchJump(addressToPatch, jumpOffset)

	case *ast.Or:
		addressToPatch := n.JumpAddress
		jumpOffset := len(cg.currentChunk().Code) - addressToPatch - 2
		cg.patchJump(addressToPatch, jumpOffset)

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
			cg.codeGenerator.ice("unknown type conversion operator: %v", n.Operator)
		}

	case *ast.IfStmt:
		break

	case *ast.WhileStmt:
		// Emit the jump back to the start of the loop

		// FIXME: I think this -2 must be -5 if the patch below upgrades the
		// jump to a long jump.
		jumpOffset := n.ConditionAddress - len(cg.currentChunk().Code) - 2
		if jumpOffset >= math.MinInt8 {
			cg.emitBytes(bytecode.OpJump, byte(int8(jumpOffset)))
		} else {
			cg.emitBytes(bytecode.OpJumpLong)
			bytecode.EncodeSInt32(cg.currentChunk().Code, jumpOffset)
		}

		// Patch the jump that skips the body when the condition is false
		addressToPatch := n.SkipJumpAddress
		jumpOffset = len(cg.currentChunk().Code) - addressToPatch - 2
		cg.patchJump(addressToPatch, jumpOffset)

	case *ast.BuiltInFunction:
		if n.Function != "print" {
			cg.codeGenerator.ice("only %q is supported, got %q", "print", n.Function)
		}
		cg.emitBytes(bytecode.OpPrint)

	case *ast.GlobalsBlock:
		break

	case *ast.VarDecl:
		if cg.codeGenerator.isInsideGlobalsBlock() {
			// Globals were already handled by the globalsExtractor.
			break
		}
		cg.defineLocalVariable(n.Name)

	case *ast.VarRef:
		localIndex := cg.resolveLocal(n.Name)
		if localIndex < 0 {
			// It's a global
			i := cg.codeGenerator.csw.GetGlobalIndex(n.Name)
			if i < 0 {
				cg.codeGenerator.ice("global variable '%v' not found in the globals pool", n.Name)
			}
			if i > 255 {
				// TODO: Can this even happen? I guess GetGlobalIndex will never return
				// anything over 255.
				cg.codeGenerator.error("Currently only up to 255 global variables are supported.")
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
			i := cg.codeGenerator.csw.GetGlobalIndex(n.VarName)
			if i < 0 {
				cg.codeGenerator.error("Global variable '%v' not declared.", n.VarName)
			}
			if i > 255 {
				// TODO: Can this even happen? I guess GetGlobalIndex will never return
				// anything over 255.
				cg.codeGenerator.error("Currently only up to 255 global variables are supported.")
			}
			cg.emitBytes(bytecode.OpWriteGlobal, byte(i))
		} else {
			// It's a local
			cg.emitBytes(bytecode.OpWriteLocal, byte(localIndex))
		}

	case *ast.ExpressionStmt:
		cg.emitBytes(bytecode.OpPop)

	case *ast.Block:
		cg.codeGenerator.endScope()
		cg.popDescopedLocals()

	case *ast.FunctionDecl:
		// Nothing to do here. Functions were already handled by
		// globalsExtractor.

		// Here just create a function object referring to the Chunk of compiled
		// bytecode we just generated and store it in a global variable with the
		// function name.
		//
		// TODO: Eventually we'll support nested functions -- then this will
		// change.
		currentChunkIndex := len(cg.codeGenerator.csw.Chunks) - 1
		f := bytecode.Value{
			Value: bytecode.Function{
				ChunkIndex: currentChunkIndex,
			},
		}
		cg.codeGenerator.csw.SetGlobal(n.Name, f)

		// TODO: For now, we add an implicit return at the end of the function.
		// Later on we'll want to do that only if the function doesn't already
		// have a return statement at the end.
		cg.emitBytes(bytecode.OpReturn)

		// No need to worry about duplicate `main`s: the semantic checker
		// already verified this.
		if n.Name == "main" {
			cg.codeGenerator.csw.FirstChunk = currentChunkIndex
		}

		cg.codeGenerator.endScope()
		cg.popDescopedLocals()

	default:
		cg.codeGenerator.ice("unknown node type: %T", n)
	}

	cg.codeGenerator.popFromNodeStack()
}

func (cg *codeGeneratorPassTwo) Event(node ast.Node, event int) {
	switch n := node.(type) {
	case *ast.IfStmt:
		// We initially emit a short jump with placeholder jump offsets. We
		// update the jump offsets once we know the size of the code block that
		// will be jumped over.
		switch event {

		case ast.EventAfterIfCondition:
			n.IfJumpAddress = len(cg.currentChunk().Code)
			cg.emitBytes(bytecode.OpJumpIfFalse, 0x00)

		case ast.EventAfterThenBlock:
			addressToPatch := n.IfJumpAddress
			jumpOffset := len(cg.currentChunk().Code) - addressToPatch - 2
			cg.patchJump(addressToPatch, jumpOffset)

		case ast.EventBeforeElse:
			n.ElseJumpAddress = len(cg.currentChunk().Code)
			cg.emitBytes(bytecode.OpJump, 0x00)

			// Re-patch the "if" jump address, because the "else" block will
			// generate an additional jump (which takes two bytes).
			//
			// FIXME: Likely to have a bug here. What if this additional jump is
			// later patched to a long one, which takes 5 bytes?
			addressToPatch := n.IfJumpAddress
			jumpOffset := int(cg.currentChunk().Code[addressToPatch+1]) + 2
			cg.patchJump(addressToPatch, jumpOffset)

		case ast.EventAfterElse:
			addressToPatch := n.ElseJumpAddress
			jumpOffset := len(cg.currentChunk().Code) - addressToPatch - 2
			cg.patchJump(addressToPatch, jumpOffset)

		default:
			cg.codeGenerator.ice("Unexpected event while generating code for 'if' statement: %v", event)
		}

	case *ast.WhileStmt:
		if event != ast.EventAfterWhileCondition {
			cg.codeGenerator.ice("Unexpected event while generating code for 'while' statement: %v", event)
		}
		n.SkipJumpAddress = len(cg.currentChunk().Code)
		cg.emitBytes(bytecode.OpJumpIfFalse, 0x00)

	case *ast.And:
		if event != ast.EventAfterLogicalBinaryOp {
			cg.codeGenerator.ice("Unexpected event while generating code for 'and' expression: %v", event)
		}
		n.JumpAddress = len(cg.currentChunk().Code)
		cg.emitBytes(bytecode.OpJumpIfFalseNoPop, 0x00)
		cg.emitBytes(bytecode.OpPop)

	case *ast.Or:
		if event != ast.EventAfterLogicalBinaryOp {
			cg.codeGenerator.ice("Unexpected event while generating code for 'or' expression: %v", event)
		}
		n.JumpAddress = len(cg.currentChunk().Code)
		cg.emitBytes(bytecode.OpJumpIfTrueNoPop, 0x00)
		cg.emitBytes(bytecode.OpPop)
	}
}

//
// Actual code generation
//

// currentChunk returns the current chunk we are compiling into.
func (cg *codeGeneratorPassTwo) currentChunk() *bytecode.Chunk {
	// For now we don't support nested functions, so the current function is
	// always the last one in the list of chunks, because we deal with them
	// one-by-one: add the new chunk, compile to it, and go to the next
	// function.
	return cg.codeGenerator.csw.Chunks[len(cg.codeGenerator.csw.Chunks)-1]
}

// currentLines returns the current array mapping instructions to source code
// lines.
//
// TODO: Returning a pointer to a slice is ugly as hell, and leads to even
// uglier client code.
func (cg *codeGeneratorPassTwo) currentLines() *[]int {
	return &cg.codeGenerator.debugInfo.ChunksLines[len(cg.codeGenerator.debugInfo.ChunksLines)-1]
}

// emitBytes writes one or more bytes to the bytecode chunk being generated.
func (cg *codeGeneratorPassTwo) emitBytes(bytes ...byte) {
	for _, b := range bytes {
		cg.currentChunk().Write(b)
		lines := cg.currentLines()
		*lines = append(*lines, cg.codeGenerator.currentLine())
	}
}

// emitConstant emits the bytecode for a constant having a given value.
func (cg *codeGeneratorPassTwo) emitConstant(value bytecode.Value) {
	if cg.codeGenerator.isInsideGlobalsBlock() {
		// Globals are initialized directly from the initializer value from the
		// AST. No need to push the initializer value to the stack.
		return
	}

	constantIndex := cg.makeConstant(value)
	if constantIndex <= math.MaxUint8 {
		cg.emitBytes(bytecode.OpConstant, byte(constantIndex))
	} else {
		operandStart := len(cg.currentChunk().Code) + 1
		cg.emitBytes(bytecode.OpConstantLong, 0, 0, 0, 0)
		bytecode.EncodeUInt31(cg.currentChunk().Code[operandStart:], constantIndex)
	}
}

// makeConstant adds value to the pool of constants and returns the index in
// which it was added. If there is already a constant with this value, its index
// is returned (hey, we don't need duplicate constants, right? They are
// constant, after all!)
func (cg *codeGeneratorPassTwo) makeConstant(value bytecode.Value) int {
	if i := cg.codeGenerator.csw.SearchConstant(value); i >= 0 {
		return i
	}

	constantIndex := cg.codeGenerator.csw.AddConstant(value)
	if constantIndex >= bytecode.MaxConstantsPerChunk {
		cg.codeGenerator.error("Too many constants in one chunk.")
		return 0
	}

	return constantIndex
}

// defineLocalVariable creates a new local variable called name (in other words,
// this appends a proper entry to cg.locals). Assumes the corresponding value is
// on the stack already. Returns true on success. On error, emits a compilation
// error and returns false.
func (cg *codeGeneratorPassTwo) defineLocalVariable(name string) bool {
	if len(cg.locals) == 256 {
		cg.codeGenerator.error("Currently only up to 255 global variables are supported.")
		return false
	}

	for _, local := range cg.locals {
		if local.name == name {
			cg.codeGenerator.error("Local variable %q already defined. Shadowing not allowed.", name)
		}
	}

	cg.locals = append(cg.locals, local{name: name, depth: cg.codeGenerator.scopeDepth})
	return true
}

// popDescopedLocals pops all local variables declared on scopes deeper than the
// current scope depth.
func (cg *codeGeneratorPassTwo) popDescopedLocals() {
	for len(cg.locals) > 0 && cg.locals[len(cg.locals)-1].depth > cg.codeGenerator.scopeDepth {
		cg.emitBytes(bytecode.OpPop)
		cg.locals = cg.locals[:len(cg.locals)-1]
	}
}

// newInternedValueString creates a new Value initialized to the interned string
// value v. Emphasis on "interned": if there is already some other string value
// equal to v on this VM, we'll reuse that same memory in the returned value.
func (cg *codeGeneratorPassTwo) newInternedValueString(v string) bytecode.Value {
	s := cg.codeGenerator.csw.Strings.Intern(v)
	return bytecode.NewValueString(s)
}

// resolveLocal finds the index into the locals array of the local variable
// named name.
func (cg *codeGeneratorPassTwo) resolveLocal(name string) int {
	for i, local := range cg.locals {
		if local.name == name {
			return i
		}
	}

	return -1
}

// patchJump patches a jump instruction. This means two things. First, setting
// the operand of the jump instruction at addressToPatch to jumpOffset. Second,
// if a short jump instruction is currently used and the requested jump offset
// is larger than what a short jump supports, we "upgrade" the intruction to a
// long jump.
//
// The "upgrade to a long jump" does some memory copying to open up space for
// the longer operand used by long jumps, which is a bit unfortunate, but at
// least this is a compile-time, not a run-time cost. Also, this works because
// all jump offsets are relative, and the language doesn't support arbitrary
// jumps that could be broken when parts of the bytecode shift to give space for
// longer jump offsets.
func (cg *codeGeneratorPassTwo) patchJump(addressToPatch, jumpOffset int) {
	if jumpOffset > math.MaxInt32 || jumpOffset < math.MinInt32 {
		cg.codeGenerator.error("Jump offset of %v is larger than supported.", jumpOffset)
	}

	if cg.isShortJumpOpcode(cg.currentChunk().Code[addressToPatch]) {
		// Short jump instruction with short offset: just patch the offset
		if jumpOffset >= math.MinInt8 && jumpOffset <= math.MaxInt8 {
			cg.currentChunk().Code[addressToPatch+1] = uint8(jumpOffset)
			return
		}

		// Short jump instruction with a long offset: upgrade to a long jump.
		// The opcode of the long version is always one larger than the opcode
		// of the short version.
		cg.currentChunk().Code[addressToPatch]++

		// Move all bytecode starting from just after the jump instruction three
		// bytes "downslice", to open space for the longer jump offset.
		end := len(cg.currentChunk().Code)
		cg.currentChunk().Code = append(cg.currentChunk().Code, 0x00, 0x00, 0x00)
		copy(cg.currentChunk().Code[addressToPatch+4:], cg.currentChunk().Code[addressToPatch+1:end])
		lines := cg.currentLines()
		*lines = append(*lines, 0x00, 0x00, 0x00)
		copy((*lines)[addressToPatch+4:], (*lines)[addressToPatch+1:end])

		// Don't return yet, we'll patch the jump offset right after this if
		// block.
	}

	// Already using a long jump instruction, simply patch the jump offset.
	bytecode.EncodeSInt32(cg.currentChunk().Code[addressToPatch+1:], jumpOffset)
}

// Checks if opcode is one the jump instruction variations that use a single
// signed byte to represent the jump offset.
func (cg *codeGeneratorPassTwo) isShortJumpOpcode(opcode uint8) bool {
	return opcode == bytecode.OpJump ||
		opcode == bytecode.OpJumpIfFalse ||
		opcode == bytecode.OpJumpIfFalseNoPop ||
		opcode == bytecode.OpJumpIfTrueNoPop
}
