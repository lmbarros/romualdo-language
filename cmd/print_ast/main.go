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
	case *ast.Storyworld:
		ap.builder.WriteString("Storyworld\n")
	case *ast.GlobalsBlock:
		ap.builder.WriteString("GlobalsBlock\n")
	case *ast.Block:
		ap.builder.WriteString("Block\n")
	case *ast.IfStmt:
		ap.builder.WriteString("If\n")
	case *ast.ExpressionStmt:
		ap.builder.WriteString("ExpressionStmt\n")
	case *ast.VarDecl:
		ap.builder.WriteString(fmt.Sprintf("VarDecl [%v: %v]\n", n.Name, n.Type()))
	case *ast.BuiltInFunction:
		ap.builder.WriteString(fmt.Sprintf("BuildInFunction [%v]\n", n.Function))
	case *ast.VarRef:
		ap.builder.WriteString(fmt.Sprintf("VarRef [%v: %v]\n", n.Name, n.Type()))
	case *ast.FloatLiteral:
		ap.builder.WriteString(fmt.Sprintf("FloatLiteral [%v]\n", n.Value))
	case *ast.BNumLiteral:
		ap.builder.WriteString(fmt.Sprintf("BNumLiteral [%v]\n", n.Value))
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
	case *ast.And:
		ap.builder.WriteString("And\n")
	case *ast.Or:
		ap.builder.WriteString("Or\n")
	case *ast.TypeConversion:
		ap.builder.WriteString(fmt.Sprintf("TypeConversion [%v]\n", n.Type()))
	case *ast.Assignment:
		ap.builder.WriteString(fmt.Sprintf("Assignment [%v]\n", n.VarName))
	case *ast.WhileStmt:
		ap.builder.WriteString("WhileStmt")
	case *ast.FunctionDecl:
		ap.builder.WriteString(fmt.Sprintf("FunctionDecl [%v(%v):%v]\n", n.Name, n.Parameters, n.ReturnType))
	default:
		panic(fmt.Sprintf("Unexpected node type: %T", n))
	}

	ap.indentLevel++
}

func (ap *ASTPrinter) Leave(ast.Node) {
	ap.indentLevel--
}

func (ap *ASTPrinter) Event(node ast.Node, event int) {
}

// indent returns a string good for indenting code level levels deep.
func indent(level int) string {
	return strings.Repeat("\t", level)
}
