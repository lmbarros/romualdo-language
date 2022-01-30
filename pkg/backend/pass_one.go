/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2022 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package backend

import (
	"gitlab.com/stackedboxes/romulang/pkg/ast"
	"gitlab.com/stackedboxes/romulang/pkg/bytecode"
)

// codeGeneratorPassOne populates the globals pool and creates the Chunks where
// the bytecode will be eventually written to.
//
// This implements the ast.Visitor interface.
type codeGeneratorPassOne struct {
	codeGenerator *codeGenerator
}

//
// The ast.Visitor interface
//

func (cg *codeGeneratorPassOne) Enter(node ast.Node) {
	if _, ok := node.(*ast.Block); ok {
		cg.codeGenerator.beginScope()
	}
	if cg.codeGenerator.scopeDepth > 0 {
		return
	}

	switch n := node.(type) {
	case *ast.VarDecl:
		// Global variable
		created := cg.codeGenerator.csw.SetGlobal(n.Name, cg.codeGenerator.valueFromNode(n.Initializer))
		if !created {
			cg.codeGenerator.ice(
				"duplicate definition of global name '%v' during pass one",
				n.Name)
		}

	case *ast.FunctionDecl:
		// Add a new Chunk for this function.
		bytecode.AddChunk(cg.codeGenerator.csw, cg.codeGenerator.debugInfo, n)
		created := cg.codeGenerator.csw.SetGlobal(n.Name, cg.codeGenerator.valueFromNode(n))
		if !created {
			cg.codeGenerator.ice(
				"duplicate definition of global name '%v' during pass one",
				n.Name)
		}
	}
}

func (cg *codeGeneratorPassOne) Leave(node ast.Node) {
	if _, ok := node.(*ast.Block); ok {
		cg.codeGenerator.endScope()
	}

	if cg.codeGenerator.scopeDepth > 0 {
		return
	}
}

func (cg *codeGeneratorPassOne) Event(node ast.Node, event int) {
}
