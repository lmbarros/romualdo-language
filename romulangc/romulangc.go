package main

import (
	"fmt"

	"gitlab.com/stackedboxes/romulang/romulangc/parser"
)

func main() {
	fmt.Printf("Hello!\n")

	fmt.Println(*parser.ParseTest())
}
