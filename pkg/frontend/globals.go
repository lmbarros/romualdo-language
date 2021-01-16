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

// extractGlobalTypes extract the types of all globally-declared variable in the
// sw Storyworld.
func extractGlobalTypes(sw *ast.Storyworld) map[string]ast.Type {
	types := map[string]ast.Type{}
	for _, decl := range sw.Declarations {
		if vars, ok := decl.(*ast.VarsBlock); ok {
			for _, v := range vars.Vars {
				types[v.Name] = v.Type()
			}
		}
	}
	return types
}

// globalTypeSetter is a node visitor that sets the type of all VarRef nodes
// that reference global variables.
type globalTypeSetter struct {
	// errors collects all type errors detected.
	errors []string

	// nodeStack is used to keep track of the nodes being processed. The current
	// one is on the top.
	nodeStack []ast.Node

	// GlobalTypes maps global variables by names to their types. Must be set
	// before
	GlobalTypes map[string]ast.Type
}

//
// The Visitor interface
//
func (tc *globalTypeSetter) Enter(node ast.Node) {
	tc.nodeStack = append(tc.nodeStack, node)

	if n, ok := node.(*ast.VarRef); ok {
		// All variables are global for now.
		t, ok := tc.GlobalTypes[n.Name]
		if !ok {
			tc.error("Undeclared global variable '%v'.", n.Name)
		}
		n.VarType = t
	}
}

func (tc *globalTypeSetter) Leave(ast.Node) {
	tc.nodeStack = tc.nodeStack[:len(tc.nodeStack)-1]
}

// error reports an error.
func (tc *globalTypeSetter) error(format string, a ...interface{}) {
	tc.errors = append(tc.errors,
		fmt.Sprintf("[line %v]: %v", tc.currentLine(), fmt.Sprintf(format, a...)))
}

// currentLine returns the source code line corresponding to whatever we are
// currently analyzing.
func (tc *globalTypeSetter) currentLine() int {
	if len(tc.nodeStack) == 0 {
		return -1 // TODO: Hack for that forced RETURN we generate out of no real node.
	}
	return tc.nodeStack[len(tc.nodeStack)-1].Line()
}
