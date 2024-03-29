/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2022 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package ast

// A Node is a node in Romualdo's AST (Abstract Syntax Tree).
type Node interface {
	// Type returns the type of Node.
	Type() *Type

	// Line returns the line of code that produced this node.
	Line() int

	// Walk is used to traverse the AST using the visitor v. Must start by
	// calling v.Enter(), then visit all subnodes (by calling their Walk()
	// methods), and finish by calling v.Leave().
	Walk(v Visitor)
}
