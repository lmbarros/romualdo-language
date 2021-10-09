/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2021 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package frontend

import (
	"fmt"
	"os"

	"gitlab.com/stackedboxes/romulang/pkg/ast"
)

// Parse parses and type checks a given Romualdo Language source code and
// returns its AST (Abstract Syntax Tree).
func Parse(source string) ast.Node {
	p := newParser(source)
	root := p.parse()
	if root == nil {
		return nil
	}

	// Assorted semantic checks (but no type checks)
	sc := &semanticChecker{}
	root.Walk(sc)
	if len(sc.errors) > 0 {
		for _, e := range sc.errors {
			fmt.Fprintf(os.Stderr, "%v\n", e)
		}
		return nil
	}

	// Look for undeclared variables, set types of global variables references
	globalTypes := extractGlobalTypes(root)
	vts := &variableTypeSetter{
		GlobalTypes: globalTypes,
	}
	root.Walk(vts)
	if len(vts.errors) > 0 {
		for _, e := range vts.errors {
			fmt.Fprintf(os.Stderr, "%v\n", e)
		}
		return nil
	}

	// Type checking
	tc := &typeChecker{}
	root.Walk(tc)
	if len(tc.errors) > 0 {
		for _, e := range tc.errors {
			fmt.Fprintf(os.Stderr, "%v\n", e)
		}
		return nil
	}

	return root
}
