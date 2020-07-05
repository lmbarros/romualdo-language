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

	"gitlab.com/stackedboxes/romulang/pkg/backend"
	"gitlab.com/stackedboxes/romulang/pkg/frontend"
	"gitlab.com/stackedboxes/romulang/pkg/vm"
)

const (
	exitCodeSuccess = iota
	exitCodeCompilationError
	exitCodeInterpretationError
)

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: romulangc <file>\n")
		os.Exit(1)
	}

	runFile(os.Args[1])
}

func runFile(path string) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %v: %v\n", path, err)
		os.Exit(1)
	}

	root := frontend.Parse(string(source))
	if root == nil {
		os.Exit(exitCodeCompilationError)
	}

	chunk, err := backend.GenerateCode(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(exitCodeCompilationError)
	}

	theVM := vm.New()
	if !theVM.Interpret(chunk) {
		os.Exit(exitCodeInterpretationError)
	}

	os.Exit(exitCodeSuccess)
}
