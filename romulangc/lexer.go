package main

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/regex"
)

// romualdoLexer creates and returns a lexer for the Romualdo Language. It
// panics in case of errors.
func romualdoLexer() lexer.Definition {
	lexer, err := regex.New(`
		Blank = \s
		Comment = #[^\n\r]*
		String = "([^"\\]|\\")*"
		Int = (\+|-)?[1-9][0-9]*
		Ident = [a-zA-Z][a-zA-Z0-9_]*
		Op = =
		Colon = :
		Parens = \(|\)
		At = @
		Keyword = end|meta|passage|storyworld|vars
	`)

	if err != nil {
		panic(err)
	}

	return lexer
}
