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
		switch n := decl.(type) {
		case *ast.GlobalsBlock:
			for _, v := range n.Vars {
				types[v.Name] = v.Type()
			}

		case *ast.FunctionDecl:
			paramTypes := []*ast.Type{}
			for _, t := range n.Parameters {
				paramTypes = append(paramTypes, t.Type)
			}
			types[n.Name] = &ast.Type{
				Tag:            ast.TypeFunction,
				ReturnType:     n.ReturnType,
				ParameterTypes: paramTypes,
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

// variableTypeSetter is a node visitor that sets the type of all nodes of types
// VarRef, FunctionCall and Assignment.
type variableTypeSetter struct {
	// errors collects all type errors detected.
	errors []string

	// nodeStack is used to keep track of the nodes being processed. The current
	// one is on the top.
	nodeStack []ast.Node

	// globalTypes maps global variables names to their types. Must be set
	// before using the visitor.
	globalTypes map[string]*ast.Type

	// localTypes contains all local variables currently in scope. The visitor
	// keeps this up-to-date as it traverses the parse tree.
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
		n.VarType = ts.resolveType(n.Name)

	case *ast.FunctionCall:
		// FIXME: Using n.Function.Name for the function name is wrong: what if
		// we assign the function to a variable with a different name?
		n.FunctionType = ts.globalTypes[n.Function.Name]

	case *ast.FunctionDecl:
		for _, param := range n.Parameters {
			ts.localTypes = append(ts.localTypes, local{name: param.Name, depth: ts.scopeDepth, varType: param.Type})
		}

	case *ast.Assignment:
		n.VarType = ts.resolveType(n.VarName)

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

// resolveType returns the type associated with name in the current scope. If not
// found, returns nil.
func (ts *variableTypeSetter) resolveType(name string) *ast.Type {
	localIndex := ts.resolveLocal(name)
	if localIndex < 0 {
		t, ok := ts.globalTypes[name]
		if !ok {
			ts.error("Undeclared name '%v'.", name)
			return nil
		}
		return t
	} else {
		t := ts.localTypes[localIndex].varType
		return t
	}
}
