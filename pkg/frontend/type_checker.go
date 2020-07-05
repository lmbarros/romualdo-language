/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package frontend

import (
	"fmt"

	"gitlab.com/stackedboxes/romulang/pkg/ast"
)

// typeChecker is a node visitor that implements type checking.
type typeChecker struct {
	// errors collects all type errors detected.
	errors []string

	// nodeStack is used to keep track of the nodes being processed. The current
	// on is on the top.
	nodeStack []ast.Node
}

func (tc *typeChecker) Enter(node ast.Node) {
	tc.nodeStack = append(tc.nodeStack, node)

	switch n := node.(type) {
	case *ast.Binary:
		tc.checkTypeBinary(n)
	}

}

func (tc *typeChecker) Leave(ast.Node) {
	tc.nodeStack = tc.nodeStack[:len(tc.nodeStack)-1]
}

func (tc *typeChecker) checkTypeBinary(node *ast.Binary) {
	switch node.Operator {
	case "<", "<=", ">", ">=":
		if node.LHS.Type().Tag != ast.TypeFloat {
			tc.error("Operator %v expects numeric operands; got '%v' (a %v)",
				node.Operator, node.LHS.Lexeme(), node.LHS.Type())
		}
	}
}

// error reports an error.
func (tc *typeChecker) error(format string, a ...interface{}) {
	tc.errors = append(tc.errors,
		fmt.Sprintf("[line %v]: %v", tc.currentLine(), fmt.Sprintf(format, a...)))
}

// currentLine returns the source code line corresponding to whatever we are
// currently compiling.
func (tc *typeChecker) currentLine() int {
	if len(tc.nodeStack) == 0 {
		return -1 // TODO: Hack for that forced RETURN we generate out of no real node.
	}
	return tc.nodeStack[len(tc.nodeStack)-1].Line()
}
