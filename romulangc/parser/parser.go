package parser

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"gitlab.com/stackedboxes/romulang/romulangc/ast"
)

// SourceFile contains all the declarations found in a single Romualdo Language
// source file.
type SourceFile struct {
	Pos lexer.Position

	// Namespace is the namespace in which all the declarations are. The
	// namespace is derived from the file path. Declarations in a file on the
	// compilation root would be in the global namespace. Declarations in a file
	// located at `compilationRoot/foo/Bar` would be in the `foo.bar` namespace.
	// Notice that the namespace is always in lower case.
	Namespace string

	// Declarations are the declarations found in the source file.
	Declarations []*ast.Declaration `@@*`
}

// Parse parses a given string (assumed to be Romualdo source code) and returns
// the resulting AST and an error.
func Parse(input string) (*SourceFile, error) {
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
