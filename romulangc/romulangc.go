package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

const expectedArgs = 2

func main() {
	if len(os.Args) != expectedArgs {
		fmt.Fprint(os.Stderr, "Usage: romulangc <file>\n")
		return
	}

	fileContent, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading '%v': %v\n", os.Args[1], err.Error())
	}

	_, err = Parse(string(fileContent))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing: %v\n", err.Error())
	}
}
