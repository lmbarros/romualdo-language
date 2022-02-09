/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2022 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package backend

import (
	"fmt"

	"gitlab.com/stackedboxes/romulang/pkg/ast"
	"gitlab.com/stackedboxes/romulang/pkg/bytecode"
)

// GenerateCode generates the bytecode for a given AST.
func GenerateCode(root ast.Node) (
	chunk *bytecode.CompiledStoryworld,
	debugInfo *bytecode.DebugInfo,
	err error) {

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

	passOne := &codeGeneratorPassOne{
		codeGenerator: &codeGenerator{
			csw:       bytecode.NewCompiledStoryworld(),
			debugInfo: &bytecode.DebugInfo{},
			nodeStack: make([]ast.Node, 0, 64),
		},
	}
	root.Walk(passOne)

	if len(passOne.codeGenerator.nodeStack) > 0 {
		return nil, nil, &codeGeneratorError{
			msg: "Internal Compiler Error: node stack not empty between passes",
		}
	}

	passTwo := &codeGeneratorPassTwo{
		codeGenerator: &codeGenerator{
			csw:       passOne.codeGenerator.csw,
			debugInfo: passOne.codeGenerator.debugInfo,
			nodeStack: passOne.codeGenerator.nodeStack,
		},
		currentChunkIndex: -1, // start with an invalid value, for easier debugging
	}
	root.Walk(passTwo)
	return passTwo.codeGenerator.csw, passTwo.codeGenerator.debugInfo, nil
}

// codeGeneratorError is a type used in panics to report an error in code
// generation.
type codeGeneratorError struct {
	msg string
}

func (e *codeGeneratorError) Error() string {
	return e.msg
}

// codeGenerator contains the code that is common among the actual code
// generation steps.
type codeGenerator struct {
	// csw is the CompiledStoryworld being generated.
	csw *bytecode.CompiledStoryworld

	// debugInfo is the DebugInfo corresponding to the CompiledStoryworld being
	// generated.
	debugInfo *bytecode.DebugInfo

	// nodeStack is used to keep track of the nodes being processed. The current
	// one is on the top.
	nodeStack []ast.Node

	// scopeDepth keeps track of the current scope depth we are in. Level 0 is
	// the global scope, and each nested block is one scope level deeper.
	scopeDepth int
}

//
// Other functions
//

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

// beginScope gets called when we enter into a new scope.
func (cg *codeGenerator) beginScope() {
	cg.scopeDepth++
}

// endScope gets called when we leave a scope.
func (cg *codeGenerator) endScope() {
	cg.scopeDepth--
}

// pushIntoNodeStack pushes a given node to the node stack.
func (cg *codeGenerator) pushIntoNodeStack(node ast.Node) {
	cg.nodeStack = append(cg.nodeStack, node)
}

// popFromNodeStack pops a node from the node stack.
func (cg *codeGenerator) popFromNodeStack() {
	cg.nodeStack = cg.nodeStack[:len(cg.nodeStack)-1]
}

// valueFrom node retruns a Value from a given Node (that holds a literal
// value).
func (cg *codeGenerator) valueFromNode(node ast.Node) bytecode.Value {
	switch n := node.(type) {
	case *ast.StringLiteral:
		return bytecode.NewValueString(n.Value)
	case *ast.BoolLiteral:
		return bytecode.NewValueBool(n.Value)
	case *ast.IntLiteral:
		return bytecode.NewValueInt(n.Value)
	case *ast.FloatLiteral:
		return bytecode.NewValueFloat(n.Value)
	case *ast.BNumLiteral:
		// BNums are internally represented as floats.
		return bytecode.NewValueFloat(n.Value)
	case *ast.FunctionDecl:
		return bytecode.NewValueFunction(n.ChunkIndex)
	default:
		cg.ice("Unexpected node of type %T", node)
	}
	return bytecode.Value{}
}

// currentLine returns the source code line corresponding to whatever we are
// currently compiling.
func (cg *codeGenerator) currentLine() int {
	return cg.nodeStack[len(cg.nodeStack)-1].Line()
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
