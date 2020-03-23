package parser

import (
	"github.com/alecthomas/participle"
	"gitlab.com/stackedboxes/romulang/romulangc/ast"
)

// Parse parses a given string (assumed to be Romualdo source code) and returns
// the resulting AST and an error.
func Parse(input string) (*ast.SourceFile, error) {
	parser, err := participle.Build(
		&SourceFile{},
		participle.Lexer(romulangLexer.New()),
		participle.Elide("COMMENT", "BLANK"))

	if err != nil {
		return nil, err
	}

	result := &SourceFile{}
	err = parser.ParseString(input, result)

	if err != nil {
		return nil, err
	}

	return result, nil
}
