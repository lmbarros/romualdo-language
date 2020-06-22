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

	"gitlab.com/stackedboxes/romulang/pkg/vm"
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

	theVM := vm.New()
	result := theVM.Interpret(string(source))

	if result != vm.InterpretOK {
		os.Exit(1)
	}
}
