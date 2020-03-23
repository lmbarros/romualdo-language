package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"gitlab.com/stackedboxes/romulang/romulangc/ast"
)

func main() {
	// Command-line args
	actionPrintAST := false

	flag.BoolVar(&actionPrintAST, "printAST", false,
		"Print the AST instead of compiling")

	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Fprint(os.Stderr, "Usage: romulangc [flags] <file>\n")
		return
	}

	// Read input file and parse it
	fileContent, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading '%v': %v\n", flag.Arg(0), err.Error())
		return
	}

	ast, err := Parse(string(fileContent))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing: %v\n", err.Error())
		return
	}

	// Do work
	if actionPrintAST {
		printAST(ast)
	} else {
		compile(ast)
	}
}

func compile(ast *ast.SourceFile) {
	compiler := &GDScriptBackend{}
	ast.Walk(compiler)
	fmt.Printf("%v", compiler.result)
}

func printAST(ast *ast.SourceFile) {
	printer := &ASTPrinter{}
	ast.Walk(printer)
	fmt.Printf("%v", printer.result)
}
