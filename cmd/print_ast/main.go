/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gitlab.com/stackedboxes/romulang/pkg/ast"
	"gitlab.com/stackedboxes/romulang/pkg/frontend"
)

const (
	exitCodeSuccess = iota
	exitCodeCompilationError
)

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: print_ast <file>\n")
		os.Exit(1)
	}

	parseAndPrintAST(os.Args[1])
}

func parseAndPrintAST(path string) {

	source, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %v: %v\n", path, err)
		os.Exit(1)
	}

	root := frontend.Parse(string(source))
	if root == nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(exitCodeCompilationError)
	}

	ap := &ASTPrinter{}
	root.Walk(ap)
	fmt.Println(ap)

	os.Exit(exitCodeSuccess)
}

type ASTPrinter struct {
	indentLevel int
	builder     strings.Builder
}

func (ap *ASTPrinter) String() string {
	return ap.builder.String()
}

func (ap *ASTPrinter) Enter(node ast.Node) {
	ap.builder.WriteString(indent(ap.indentLevel))

	switch n := node.(type) {
	case *ast.FloatLiteral:
		ap.builder.WriteString(fmt.Sprintf("FloatLiteral [%v]\n", n.Value))
	case *ast.IntLiteral:
		ap.builder.WriteString(fmt.Sprintf("IntLiteral [%v]\n", n.Value))
	case *ast.BoolLiteral:
		ap.builder.WriteString(fmt.Sprintf("BoolLiteral [%v]\n", n.Value))
	case *ast.StringLiteral:
		ap.builder.WriteString(fmt.Sprintf("StringLiteral [%v]\n", n.Value))
	case *ast.Unary:
		ap.builder.WriteString(fmt.Sprintf("Unary [%v]\n", n.Operator))
	case *ast.Binary:
		ap.builder.WriteString(fmt.Sprintf("Binary [%v]\n", n.Operator))

	default:
		panic(fmt.Sprintf("Unexpected node type: %T", n))
	}

	ap.indentLevel++
}

func (ap *ASTPrinter) Leave(ast.Node) {
	ap.indentLevel--
}

// indent returns a string good for indenting code level levels deep.
func indent(level int) string {
	return strings.Repeat("\t", level)
}
