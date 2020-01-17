package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprint(os.Stderr, "Usage: romulangc <file>\n")
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

	printer := &ASTPrinter{}
	ast.Walk(printer)
	fmt.Printf("%v", printer.result)
}
