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

//
// The Visitor interface
//
func (tc *typeChecker) Enter(node ast.Node) {
	tc.nodeStack = append(tc.nodeStack, node)

	switch n := node.(type) {
	case *ast.Binary:
		tc.checkBinary(n)
	case *ast.Unary:
		tc.checkTypeUnary(n)
	}

}

func (tc *typeChecker) Leave(ast.Node) {
	tc.nodeStack = tc.nodeStack[:len(tc.nodeStack)-1]
}

//
// Type checking
//

// checkBinary checks for typing errors in a binary operator.
func (tc *typeChecker) checkBinary(node *ast.Binary) {
	switch node.Operator {
	case "<", "<=", ">", ">=":
		if node.LHS.Type().Tag != ast.TypeFloat {
			tc.error("Operator %v expects numeric operands; got a %v on the left-hand side",
				node.Operator, node.LHS.Type())
		}
		if node.RHS.Type().Tag != ast.TypeFloat {
			tc.error("Operator %v expects numeric operands; got a %v on the right-hand side",
				node.Operator, node.RHS.Type())
		}
	case "==", "!=":
		if node.LHS.Type().Tag != node.RHS.Type().Tag {
			tc.error("Operator %v expects operands of same type; got a %v and a %v",
				node.Operator, node.LHS.Type(), node.RHS.Type())
		}
	case "+":
		if node.LHS.Type().Tag != ast.TypeFloat && node.LHS.Type().Tag != ast.TypeString {
			tc.error("Operator %v expects either strings or float operands; got a %v on the left-hand side",
				node.Operator, node.LHS.Type())
		}
		if node.RHS.Type().Tag != ast.TypeFloat && node.RHS.Type().Tag != ast.TypeString {
			tc.error("Operator %v expects either strings or float operands; got a %v on the right-hand side",
				node.Operator, node.RHS.Type())
		}
		if node.LHS.Type().Tag != node.RHS.Type().Tag {
			tc.error("Operator %v expects operands of same type; got a %v and a %v",
				node.Operator, node.LHS.Type(), node.RHS.Type())
		}
	default:
		if node.LHS.Type().Tag != ast.TypeFloat {
			tc.error("Operator %v expects float operands; got a %v on the left-hand side",
				node.Operator, node.LHS.Type())
		}
		if node.RHS.Type().Tag != ast.TypeFloat {
			tc.error("Operator %v expects float operands; got a %v on the right-hand side",
				node.Operator, node.RHS.Type())
		}
	}

}

// checkUnary checks for typing errors in a unary operator.
func (tc *typeChecker) checkTypeUnary(node *ast.Unary) {
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
