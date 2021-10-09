/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2021 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package frontend

import (
	"fmt"

	"gitlab.com/stackedboxes/romulang/pkg/ast"
)

// semanticChecker is a node visitor that implements assorted semantic checks.
type semanticChecker struct {
	// errors collects all type errors detected.
	errors []string

	// nodeStack is used to keep track of the nodes being processed. The current
	// on is on the top.
	nodeStack []ast.Node

	// firstGlobalsBlock contains the first globals block found. This is used to
	// detect multiple of these blocks (which is forbidden).
	firstGlobalsBlock *ast.GlobalsBlock

	// globalVariables maps the global variable names already declared to the
	// line where they were declared. Used to detect duplicates.
	globalVariables map[string]int
}

//
// The Visitor interface
//
func (sc *semanticChecker) Enter(node ast.Node) {
	sc.nodeStack = append(sc.nodeStack, node)

	switch n := node.(type) {
	case *ast.Storyworld:
		sc.globalVariables = map[string]int{}

	case *ast.GlobalsBlock:
		// At the base of the stack we have the Storyworld itself, so a globals
		// block would be the second node on the stack.
		if len(sc.nodeStack) == 2 {
			sc.checkDuplicateGlobalsBlock(n)
		}

	case *ast.VarDecl:
		sc.checkVarInitializer(n)

		if sc.isInsideGlobalsBlock() {
			sc.checkDuplicateGlobalVariable(n)
		}
	}
}

func (sc *semanticChecker) Leave(ast.Node) {
	sc.nodeStack = sc.nodeStack[:len(sc.nodeStack)-1]
}

//
// Semantic checking
//

// checkDuplicateGlobalsBlock checks if another globals block was already
// declared (which is forbidden).
func (sc *semanticChecker) checkDuplicateGlobalsBlock(node *ast.GlobalsBlock) {
	if sc.firstGlobalsBlock != nil {
		sc.error(
			"Only one 'globals' block is allowed. Found another one at line %v.",
			sc.firstGlobalsBlock.Line())
		return
	}

	sc.firstGlobalsBlock = node
}

// checkVarInitializer checks if the variable initializer is some literal value.
func (sc *semanticChecker) checkVarInitializer(node *ast.VarDecl) {
	switch node.Initializer.(type) {
	case *ast.StringLiteral, *ast.BoolLiteral, *ast.IntLiteral,
		*ast.FloatLiteral, *ast.BNumLiteral:
		break
	default:
		sc.error("Currently variables must be initialized with a literal value.")
	}
}

// checkDuplicateGlobalVariable checks if another global variable with the same
// name was already declared.
func (sc *semanticChecker) checkDuplicateGlobalVariable(node *ast.VarDecl) {
	line, found := sc.globalVariables[node.Name]
	if found {
		sc.error("There is already a global variable named '%v' declared at line %v.",
			node.Name, line)
		return
	}

	sc.globalVariables[node.Name] = node.LineNumber
}

// error reports an error.
func (sc *semanticChecker) error(format string, a ...interface{}) {
	sc.errors = append(sc.errors,
		fmt.Sprintf("[line %v]: %v", sc.currentLine(), fmt.Sprintf(format, a...)))
}

// currentLine returns the source code line corresponding to whatever we are
// currently analyzing.
func (sc *semanticChecker) currentLine() int {
	if len(sc.nodeStack) == 0 {
		return -1 // TODO: Hack for that forced RETURN we generate out of no real node.
	}
	return sc.nodeStack[len(sc.nodeStack)-1].Line()
}

// isInsideGlobalsBlock checks if we are currently inside a globals block.
func (sc *semanticChecker) isInsideGlobalsBlock() bool {
	for _, node := range sc.nodeStack {
		_, ok := node.(*ast.GlobalsBlock)
		if ok {
			return true
		}
	}
	return false
}
