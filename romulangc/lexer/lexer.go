package lexer

import (
	"strings"

	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/regex"
)

// New creates and returns a lexer for the Romualdo Language.
//
// It should never fail, but in case of internal bugs (like a badly defined
// lexical grammar), it will panic.
func New() lexer.Definition {

	const lexicalGrammarTemplate = `
		BLANK = \s
		COMMENT = #[^\n\r]*
		KEYWORD = alias|and|bnum|bool|else|elseif|end|enum|false|float|function|gosub|goto|if|int|listen|map|meta|not|or|passage|return|sat|string|struct|then|true|vars|void|while
		FLOAT = ([1-9]+.[1-9]+)|([1-9]+(.[1-9]+)?(e|E)(+|-)?[1-9]+)
		INTEGER = [1-9][0-9]*
		STRING = "([^"\\]|\\")*"
		IDENTIFIER = ⟨LETTER_LIKE⟩(⟨LETTER_LIKE⟩|[0-9])*
		SYMBOL2 = !=|==|>=|<=
		SYMBOL1 = @|:|=|\\(|\\)|,|\\[|\\]|{|}|<|>|\\+|-|/|\\*|^|.
	`

	lexicalGrammar := strings.ReplaceAll(lexicalGrammarTemplate,
		`⟨LETTER_LIKE⟩`,
		`(\pL|_|[\x{1F300}-\x{1F5FF}]|[\x{1F900}-\x{1F9FF}]|[\x{1FA70}-\x{1FAFF}]|[\x{1F600}-\x{1F64F}]|[\x{1F680}-\x{1F6FF}]|[\x{2600}-\x{26FF}]|[\x{2700}-\x{27BF}])`)

	lexer, err := regex.New(lexicalGrammar)

	if err != nil {
		panic(err)
	}

	return lexer
}
