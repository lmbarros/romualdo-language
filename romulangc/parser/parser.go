package parser

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
)

// SourceFile is a Romualdo Language source file.
type SourceFile struct {
	Declaration []*Declaration `@@*`
}

// Declaration is a declaration. Oh well.
type Declaration struct {
	Pos        lexer.Position
	TheContent *string `@Ident`
}

func ParseTest() *SourceFile {
	parser, err := participle.Build(&SourceFile{})
	if err != nil {
		panic(err)
	}

	sf := &SourceFile{}
	err = parser.ParseString(`
foo
bar baz
`, sf)

	if err != nil {
		panic(err)
	}

	return sf
}
