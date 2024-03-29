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

// typeChecker is a node visitor that implements type checking.
type typeChecker struct {
	// errors collects all type errors detected.
	errors []string

	// nodeStack is used to keep track of the nodes being processed. The current
	// one is on the top.
	nodeStack []ast.Node
}

//
// The Visitor interface
//
func (tc *typeChecker) Enter(node ast.Node) {
	tc.nodeStack = append(tc.nodeStack, node)

	switch n := node.(type) {
	case *ast.Assignment:
		tc.checkAssignment(n)
	case *ast.Binary:
		tc.checkBinary(n)
	case *ast.Unary:
		tc.checkUnary(n)
	case *ast.Blend:
		tc.checkBlend(n)
	case *ast.FunctionCall:
		tc.checkFunctionCall(n)
	case *ast.ReturnStmt:
		tc.checkReturnStmt(n)
	case *ast.TypeConversion:
		tc.checkTypeConversion(n)
	case *ast.VarDecl:
		tc.checkVarType(n)
	case *ast.And:
		tc.checkAnd(n)
	case *ast.IfStmt:
		tc.checkIf(n)
	case *ast.WhileStmt:
		tc.checkWhile(n)
	}

}

func (tc *typeChecker) Leave(ast.Node) {
	tc.nodeStack = tc.nodeStack[:len(tc.nodeStack)-1]
}

func (tc *typeChecker) Event(node ast.Node, event int) {
}

//
// Type checking
//

// checkAssignment type checks an assignment operator.
func (tc *typeChecker) checkAssignment(node *ast.Assignment) {
	if node.VarType != node.Value.Type() {
		tc.error("Variable '%v' is of type %v, cannot assign a %v value to it.", node.VarName, node.VarType, node.Value.Type())
	}
}

// checkBinary type checks a binary operator.
func (tc *typeChecker) checkBinary(node *ast.Binary) {
	switch node.Operator {
	case "<", "<=", ">", ">=":
		// TODO: Why only unbounded? We should be able to compare BNums, right?
		if !node.LHS.Type().IsUnboundedNumeric() {
			tc.error("Operator %v expects numeric operands; got a %v on the left-hand side",
				node.Operator, node.LHS.Type())
		}
		if !node.RHS.Type().IsUnboundedNumeric() {
			tc.error("Operator %v expects numeric operands; got a %v on the right-hand side",
				node.Operator, node.RHS.Type())
		}

	case "==", "!=":
		// Values of the same type can be compared
		if node.LHS.Type().Tag == node.RHS.Type().Tag {
			return
		}

		// Unbounded numeric types can be compared
		if node.LHS.Type().IsUnboundedNumeric() && node.RHS.Type().IsUnboundedNumeric() {
			return
		}

		// Nothing else can be compared
		tc.error("Operator %v expects operands of same type or two unbounded numeric values; got a %v and a %v",
			node.Operator, node.LHS.Type(), node.RHS.Type())

	case "+":
		// It is OK to add two bounded numbers
		if node.LHS.Type().Tag == ast.TypeBNum && node.RHS.Type().Tag == ast.TypeBNum {
			return
		}

		// It is OK to add two unbounded numbers
		if node.LHS.Type().IsUnboundedNumeric() && node.RHS.Type().IsUnboundedNumeric() {
			return
		}

		// It is OK to add (ahem, concatenate) two strings
		if node.LHS.Type().Tag == ast.TypeString && node.RHS.Type().Tag == ast.TypeString {
			return
		}

		// Nothing else can be added
		tc.error("Operator %v cannot work with values of type %v and %v",
			node.Operator, node.LHS.Type(), node.RHS.Type())

	case "-":
		// It is OK to subtract two bounded numbers
		if node.LHS.Type().Tag == ast.TypeBNum && node.RHS.Type().Tag == ast.TypeBNum {
			return
		}

		// It is OK to subtract two unbounded numbers
		if node.LHS.Type().IsUnboundedNumeric() && node.RHS.Type().IsUnboundedNumeric() {
			return
		}

		// Nothing else can be subtracted
		tc.error("Operator %v cannot work with values of type %v and %v",
			node.Operator, node.LHS.Type(), node.LHS.Type())

	default:
		if !node.LHS.Type().IsUnboundedNumeric() {
			tc.error("Operator %v expects unbounded numeric operands; got a %v on the left-hand side",
				node.Operator, node.LHS.Type())
		}
		if !node.RHS.Type().IsUnboundedNumeric() {
			tc.error("Operator %v expects unbounded numeric operands; got a %v on the left-hand side",
				node.Operator, node.RHS.Type())
		}
	}
}

// checkAnd type checks an "and" operator.
func (tc *typeChecker) checkAnd(node *ast.And) {
	if node.LHS.Type().Tag != ast.TypeBool {
		tc.error("Operator 'and' expects Boolean operands; got a %v on the left-hand side",
			node.LHS.Type())
	}
	if node.RHS.Type().Tag != ast.TypeBool {
		tc.error("Operator 'and' expects Boolean operands; got a %v on the right-hand side",
			node.RHS.Type())
	}
}

// checkIf checks an if statement.
func (tc *typeChecker) checkIf(node *ast.IfStmt) {
	if node.Condition.Type().Tag != ast.TypeBool {
		tc.error("The condition of an 'if' must be Boolean.")
	}
}

// checkWhile checks a while statement.
func (tc *typeChecker) checkWhile(node *ast.WhileStmt) {
	if node.Condition.Type().Tag != ast.TypeBool {
		tc.error("The condition of a 'while' must be Boolean.")
	}
}

// checkUnary type checks a unary operator.
func (tc *typeChecker) checkUnary(node *ast.Unary) {
	switch node.Operator {
	case "not":
		if node.Operand.Type().Tag != ast.TypeBool {
			tc.error("Operator %v expects a bool operand; got a %v",
				node.Operator, node.Operand.Type())
		}

	case "-", "+":
		if !node.Operand.Type().IsNumeric() {
			tc.error("Operator %v expects a float operand; got a %v",
				node.Operator, node.Operand.Type())
		}
	}
}

// checkBlend type checks a blend operator.
func (tc *typeChecker) checkBlend(node *ast.Blend) {
	if node.X.Type().Tag != ast.TypeBNum {
		tc.error("The blend Operator expects bnum operands; got a %v as the first one",
			node.X.Type())
	}

	if node.Y.Type().Tag != ast.TypeBNum {
		tc.error("The blend Operator expects bnum operands; got a %v as the second one",
			node.Y.Type())
	}
	if node.Weight.Type().Tag != ast.TypeBNum {
		tc.error("The blend Operator expects bnum operands; got a %v as the third one",
			node.Weight.Type())
	}
}

// checkFunctionCall type checks a function call. As a bonus, it also checks the
// arity.
func (tc *typeChecker) checkFunctionCall(node *ast.FunctionCall) {
	numArgs := len(node.Arguments)
	numParams := len(node.FunctionType.ParameterTypes)
	if numArgs != numParams {
		tc.error("Function '%v' expects %v arguments, but got %v.", node.Function.Name, numParams, numArgs)
		return
	}

	for i, paramType := range node.FunctionType.ParameterTypes {
		// TODO: Will fail with function types. Need to implement a
		// `TypesEqual()` function.
		argType := node.Arguments[i].Type()
		if paramType != argType {
			tc.error("Function '%v' expects a %v as argument %v, but got a %v.",
				node.Function.Name, paramType, i+1, argType)
		}
	}
}

// checkReturnStmt type checks a return statement.
func (tc *typeChecker) checkReturnStmt(node *ast.ReturnStmt) {
	thisFunc := tc.innermostFunctionDecl()
	funcName := thisFunc.Name
	funcType := thisFunc.ReturnType

	if node.ReturnValue == nil {
		if thisFunc.ReturnType.Tag != ast.TypeVoid {
			tc.error("Function '%v' expects a return value of type  %v.", funcName, funcType)
		}
		return
	}

	returnType := node.ReturnValue.Type()

	if returnType != funcType {
		tc.error("Function '%v' expects a return value of type %v, got a %v.",
			funcName, funcType, returnType)
	}
}

// checkTypeConversion type checks type conversion operator.
func (tc *typeChecker) checkTypeConversion(node *ast.TypeConversion) {
	switch node.Operator {
	case "int":
		if node.Value.Type().Tag == ast.TypeBNum {
			tc.error("Cannot convert a bnum to an int")
		}

		if node.Default.Type().Tag != ast.TypeInt {
			tc.error("The default value for a conversion to int must be an int; got a %v",
				node.Default.Type())
		}
	case "float":
		if node.Default.Type().Tag != ast.TypeFloat {
			tc.error("The default value for a conversion to float must be a float; got a %v",
				node.Default.Type())
		}
	case "bnum":
		if node.Value.Type().Tag == ast.TypeBool {
			tc.error("Cannot convert a bool to a bnum")
		}

		if node.Value.Type().Tag == ast.TypeInt {
			tc.error("Cannot convert an int to a bnum")
		}

		if node.Default.Type().Tag != ast.TypeBNum {
			tc.error("The default value for a conversion to bnum must be a bnum; got a %v",
				node.Default.Type())
		}

	case "string":
		if node.Default.Type().Tag != ast.TypeString {
			tc.error("The default value for a conversion to string must be a string; got a %v",
				node.Default.Type())
		}
	}
}

// checkVarType type checks a variable declaration.
func (tc *typeChecker) checkVarType(node *ast.VarDecl) {
	if node.Type().Tag == ast.TypeVoid {
		tc.error("Cannot create a variable of type 'void'.")
		return
	}
	if node.Type().Tag != node.Initializer.Type().Tag {
		tc.error("Cannot initialize variable of type '%v' with a value of type '%v'.",
			node.Type(),
			node.Initializer.Type())
	}
}

// innermostFunctionDecl returns the innermost function declaration we are
// currently in.
func (tc *typeChecker) innermostFunctionDecl() *ast.FunctionDecl {
	for i := len(tc.nodeStack) - 1; i >= 0; i-- {
		if functionDecl, ok := tc.nodeStack[i].(*ast.FunctionDecl); ok {
			return functionDecl
		}
	}
	return nil // Can't happen
}

// error reports an error.
func (tc *typeChecker) error(format string, a ...interface{}) {
	tc.errors = append(tc.errors,
		fmt.Sprintf("[line %v]: %v", tc.currentLine(), fmt.Sprintf(format, a...)))
}

// currentLine returns the source code line corresponding to whatever we are
// currently analyzing.
func (tc *typeChecker) currentLine() int {
	return tc.nodeStack[len(tc.nodeStack)-1].Line()
}
