package main

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
)

// SourceFile contains all the declarations found in a single Romualdo Language
// source file.
type SourceFile struct {
	// Namespace is the namespace in which all the declarations are. The
	// namespace is derived from the file path. Declarations in a file on the
	// compilation root would be in the global namespace. Declarations in a file
	// located at `compilationRoot/foo/Bar` would be in the `foo.bar` namespace.
	// Notice that the namespace is always in lower case.
	Namespace string

	// Declarations are the declarations found in the source file.
	Declarations []*Declaration `@@*`
}

// Declaration is any of the declarations making up a Romualdo program.
type Declaration struct {
	Pos        lexer.Position
	Storyworld *Storyworld `  @@`
	// TypeDecl   *TypeDecl   `| @@`
	// Function   *Function   `| @@`
	// Passage    *Passage    `| @@`
}

// Storyworld represents a `storyworld` declaration. In a valid Romualdo
// program, there can be only one of them (but this is checked by the semantic
// analysis, not by the parser.
type Storyworld struct {
	Pos             lexer.Position
	Storyworld      *string            `"storyworld"`
	StoryworldBlock []*StoryworldBlock `@@*`
	End             *string            `"end"`
}

// StoryworldBlock is any block that can be inside a `storyworld` block.
type StoryworldBlock struct {
	Pos  lexer.Position
	Meta *Meta `  @@`
	Vars *Vars `| @@`
}

// Meta represents a `meta` block.
type Meta struct {
	Pos       lexer.Position
	Meta      *string      `"meta"`
	MetaEntry []*MetaEntry `@@*`
	End       *string      `"end"`
}

// MetaEntry represents one entry in a `meta` block.
type MetaEntry struct {
	Pos   lexer.Position
	Name  *string `@Ident "="`
	Value *string `@String`
}

type Vars struct {
	Pos     lexer.Position
	Meta    *string    `"vars"`
	VarDecl []*VarDecl `@@*`
	End     *string    `"end"`
}

type VarDecl struct {
	Pos   lexer.Position
	Name  *string `@Ident ":"`
	Type  *string `@Ident "="`
	Value *string `@String`
}

// TypeDecl is a declaration of a user-defined type.
type TypeDecl struct {
	Pos   lexer.Position
	Alias *Alias `  @@`
	Enum  *Enum  `| @@`
	Class *Class `| @@`
}

type Alias struct {
}

type Enum struct {
}

type Class struct {
}

type Function struct {
}

type Passage struct {
}

// Parse parses a given string (assumed to be Romualdo source code) and returns
// the resulting AST and an error.
func Parse(input string) (*SourceFile, error) {
	parser, err := participle.Build(
		&SourceFile{},
		participle.Lexer(romualdoLexer()),
		participle.Elide("Comment", "Blank"))

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
