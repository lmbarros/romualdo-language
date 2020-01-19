package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "Usage: romulangc [flags] <file>\n")
		return
	}

	fileContent, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading '%v': %v\n", os.Args[1], err.Error())
		return
	}

	ast, err := Parse(string(fileContent))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing: %v\n", err.Error())
		return
	}

	// Do work
	actionPrintAST := false

	flag.BoolVar(&actionPrintAST, "printAST", false,
		"Print the AST instead of compiling")

	if actionPrintAST {
		printAST(ast)
	} else {
		compile(ast)
	}
}

func compile(ast *SourceFile) {
	compiler := &GDScriptBackend{}
	ast.Walk(compiler)
	fmt.Printf("%v", compiler.result)
}

func printAST(ast *SourceFile) {
	printer := &ASTPrinter{}
	ast.Walk(printer)
	fmt.Printf("%v", printer.result)
}
