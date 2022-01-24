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

// extractGlobalTypes extract the types of all globally-declared variable in the
// sw Storyworld.
//
// We need to do this on a separate step because the globals block can appear
// after the code that uses it.
func extractGlobalTypes(sw *ast.Storyworld) map[string]*ast.Type {
	types := map[string]*ast.Type{}
	for _, decl := range sw.Declarations {
		if globals, ok := decl.(*ast.GlobalsBlock); ok {
			for _, v := range globals.Vars {
				types[v.Name] = v.Type()
			}
		}
	}
	return types
}

// local represents a local variable.
type local struct {
	// name is the local variable name.
	name string

	// depth is the nesting level (AKA scope depth) of the local variable.
	depth int

	// varType is the type of the local variable.
	varType *ast.Type
}

// variableTypeSetter is a node visitor that sets the type of all VarRef nodes.
type variableTypeSetter struct {
	// errors collects all type errors detected.
	errors []string

	// nodeStack is used to keep track of the nodes being processed. The current
	// one is on the top.
	nodeStack []ast.Node

	// GlobalTypes maps global variables by names to their types. Must be set
	// before using the visitor.
	GlobalTypes map[string]*ast.Type

	// localTypes contains all local variables currently in scope. The visitor
	// keepts this up-to-date as it traverses the parse tree.
	localTypes []local

	// inGlobals tells if we are we inside a globals block.
	inGlobals bool

	// scopeDepth tells our current scope depth (AKA levels of nesting).
	scopeDepth int
}

//
// The Visitor interface
//
func (ts *variableTypeSetter) Enter(node ast.Node) {
	ts.nodeStack = append(ts.nodeStack, node)

	switch n := node.(type) {
	case *ast.VarRef:
		localIndex := ts.resolveLocal(n.Name)
		if localIndex < 0 {
			t, ok := ts.GlobalTypes[n.Name]
			if !ok {
				ts.error("Undeclared global variable '%v'.", n.Name)
			}
			n.VarType = t
		} else {
			t := ts.localTypes[localIndex].varType
			n.VarType = t
		}

	case *ast.Block:
		ts.scopeDepth++

	case *ast.GlobalsBlock:
		ts.inGlobals = true

	case *ast.VarDecl:
		if ts.inGlobals {
			break
		}
		ts.localTypes = append(ts.localTypes, local{name: n.Name, depth: ts.scopeDepth, varType: n.Type()})
	}
}

func (ts *variableTypeSetter) Leave(node ast.Node) {
	switch node.(type) {
	case *ast.GlobalsBlock:
		ts.inGlobals = false

	case *ast.Block:
		for i, lv := range ts.localTypes {
			if lv.depth == ts.scopeDepth {
				ts.localTypes = ts.localTypes[:i]
				break
			}
		}
		ts.scopeDepth--
	}

	ts.nodeStack = ts.nodeStack[:len(ts.nodeStack)-1]
}

func (ts *variableTypeSetter) Event(node ast.Node, event int) {
}

// error reports an error.
func (ts *variableTypeSetter) error(format string, a ...interface{}) {
	ts.errors = append(ts.errors,
		fmt.Sprintf("[line %v]: %v", ts.currentLine(), fmt.Sprintf(format, a...)))
}

// currentLine returns the source code line corresponding to whatever we are
// currently analyzing.
func (ts *variableTypeSetter) currentLine() int {
	if len(ts.nodeStack) == 0 {
		return -1 // TODO: Hack for that forced RETURN we generate out of no real node.
	}
	return ts.nodeStack[len(ts.nodeStack)-1].Line()
}

// resolveLocal returns the index into ts.localTypes of the local variable
// passed as parameter. If not found, returns -1.
func (ts *variableTypeSetter) resolveLocal(name string) int {
	for i, lv := range ts.localTypes {
		if lv.name == name {
			return i
		}
	}
	return -1
}
