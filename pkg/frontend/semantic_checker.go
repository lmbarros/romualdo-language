/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2022 Leandro Motta Barros                                     *
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

	// mainFunctionLine contains the line number where the main function was
	// found. A value of 0 means that main wasn't found yet.
	mainFunctionLine int
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
			sc.checkDuplicateGlobalName(n.Name, n.BaseNode)
		}

	case *ast.FunctionDecl:
		sc.checkDuplicateGlobalName(n.Name, n.BaseNode)

		if n.Name != "main" {
			break
		}
		if sc.mainFunctionLine != 0 {
			sc.error("Duplicate function main at line %v. The first one was at line %v.",
				n.LineNumber, sc.mainFunctionLine)
			break
		}
		sc.mainFunctionLine = n.LineNumber
	}
}

func (sc *semanticChecker) Leave(n ast.Node) {
	sc.nodeStack = sc.nodeStack[:len(sc.nodeStack)-1]

	if _, ok := n.(*ast.Storyworld); ok {
		if sc.mainFunctionLine == 0 {
			sc.error("Function 'main' not found.")
		}
	}
}

func (sc *semanticChecker) Event(node ast.Node, event int) {
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

// checkDuplicateGlobalName checks if something with the same name was already
// declared at the global scope. If this is a new globa, it also adds the name
// to the list of known globals, taking the corresponding line number from node.
func (sc *semanticChecker) checkDuplicateGlobalName(name string, node ast.BaseNode) {
	line, found := sc.globalVariables[name]
	if found {
		sc.error("The name '%v' was already globally declared at at line %v.", name, line)
		return
	}

	sc.globalVariables[name] = node.LineNumber
}

// error reports an error.
func (sc *semanticChecker) error(format string, a ...interface{}) {
	sc.errors = append(sc.errors,
		fmt.Sprintf("[line %v]: %v", sc.currentLine(), fmt.Sprintf(format, a...)))
}

// currentLine returns the source code line corresponding to whatever we are
// currently analyzing.
func (sc *semanticChecker) currentLine() int {
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
