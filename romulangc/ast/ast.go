package ast

import (
	"github.com/alecthomas/participle/lexer"
)

// SourceFile contains all the declarations found in a single Romualdo Language
// source file.
//
// Not technically part of the AST, but this is the best package to include it.
type SourceFile struct {
	Pos lexer.Position

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
	Pos          lexer.Position
	MetaBlock    *MetaBlock    `  @@`
	VarsBlock    *VarsBlock    `| @@`
	AliasDecl    *AliasDecl    `| @@`
	EnumDecl     *EnumDecl     `| @@`
	StructDecl   *StructDecl   `| @@`
	FunctionDecl *FunctionDecl `| @@`
	PassageDecl  *PassageDecl  `| @@`
}

// MetaBlock represents a `meta` block, containing metadata.
type MetaBlock struct {
	Pos     lexer.Position
	Version *int       `"meta" ( "@" @INTEGER )?`
	Entries []*VarDecl `@@* "end"`
}

// VarDecl represents a variable declaration.
type VarDecl struct {
	Pos          lexer.Position
	Name         *string     `@IDENTIFIER ":"`
	Type         *Type       `@@ ( "="`
	InitialValue *Expression `@@ )?`
}

// VarsBlock represents a block where variables are declared.
type VarsBlock struct {
	Pos     lexer.Position
	Version *int       `"vars" ( "@" @INTEGER )?`
	Vars    []*VarDecl `@@* "end"`
}

// AliasDecl represents an `alias` type declaration.
type AliasDecl struct {
	Pos  lexer.Position
	Name *string `"alias" @IDENTIFIER`
	Type *Type   `@@`
}

// EnumDecl represents an `enum` type declaration.
type EnumDecl struct {
	Pos      lexer.Position
	Name     *string   `"enum" @IDENTIFIER`
	Elements []*string `@IDENTIFIER* "end"`
}

// StructDecl represents a `struct` type declaration.
type StructDecl struct {
	Pos     lexer.Position
	Name    *string    `"struct" @IDENTIFIER`
	Members []*VarDecl `@@* "end"`
}

// FunctionDecl represents a function declaration.
type FunctionDecl struct {
	Pos        lexer.Position
	Name       *string      `"function" @IDENTIFIER`
	Parameters *Parameters  `"(" @@? ")"`
	ReturnType *Type        `":" @@`
	Body       []*Statement `@@* "end"`
}

// Parameters represents a sequence of one or more (formal) parameters.
type Parameters struct {
	Pos   lexer.Position
	First *Parameter   `@@`
	Rest  []*Parameter `( "," @@ )*` // xxxxxxxxxxxx does this work as intended?
}

// Parameter represents one (formal) parameter.
type Parameter struct {
	Pos  lexer.Position
	Name *string `@IDENTIFIER`
	Type *Type   `@@`
}

// PassageDecl represents the declaration of a passage.
type PassageDecl struct {
	Pos        lexer.Position
	Name       *string      `passage @IDENTIFIER`
	Version    *int         `"@" @INTEGER`
	Parameters *Parameters  `"(" @@ ")"`
	ReturnType *Type        `":" @@`
	Body       []*Statement `@@* "end"`
}

// Type represents a type.
type Type struct {
	Pos                 lexer.Position
	BuiltIn             *string              `  @("bool"|"int"|"float"|"bnum"|"string"|"map"|"void")`
	Array               *Type                `| "[" "]" @@`
	Function            *FunctionType        `| @@`
	QualifiedIdentifier *QualifiedIdentifier `| @@`
}

// FunctionType represents the type of a function.
type FunctionType struct {
	Pos            lexer.Position
	ParameterTypes *TypeList `"function" "(" @@? ")"`
	ReturnType     *Type     `":" @@`
}

// TypeList represents a list of one or more types.
type TypeList struct {
	Pos   lexer.Position
	First *Type   `@@`
	Rest  []*Type `( "," @@ )*` // xxxxxxxxxxxx does this work as intended?
}

// QualifiedIdentifier represents an identifier possibly within a namespace.
type QualifiedIdentifier struct {
	Pos   lexer.Position
	First *string   `@IDENTIFIER`
	Rest  []*string `( "," @@ )*` // xxxxxxxxxxxx does this work as intended? Should be @IDENTIFIER here?
}

// Statement represents, well, a statement.
type Statement struct {
	Pos        lexer.Position
	Expression *Expression `  @@`
	WhileStmt  *WhileStmt  `| @@`
	IfStmt     *IfStmt     `| @@`
	ReturnStmt *ReturnStmt `| @@`
	GotoStmt   *GotoStmt   `| @@`
	SayStmt    *SayStmt    `| @@`
}

// WhileStmt represents a `while` statement.
type WhileStmt struct {
	Pos       lexer.Position
	Condition *Expression  `"while" @@ "do"`
	Body      []*Statement `@@* "end"`
}

// IfStmt represents an `if` statement.
type IfStmt struct {
	Pos       lexer.Position
	Condition *Expression  `"if" @@ "then"`
	Then      []*Statement `@@*`
	Elseifs   []*Elseif    `elseif*`
	Else      []*Statement `( "else" @@* "end" )?`
}

// Elseif represents one `elseif` block of an `if` statement.
type Elseif struct {
	Pos       lexer.Position
	Condition *Expression  `"elseif" @@ "then"`
	Body      []*Statement `@@*`
}

// ReturnStmt represents a `return` statement.
type ReturnStmt struct {
	Pos  lexer.Position
	What *Expression `"return" @@`
}

// GotoStmt represents a `goto` statement.
type GotoStmt struct {
	Pos       lexer.Position
	Passage   *QualifiedIdentifier `"goto" @@`
	Arguments *Arguments           `"(" @@ ")"`
}

// SayStmt represents a `say` statement.
type SayStmt struct {
	Pos  lexer.Position
	What *Expression `"say" @@`
}

// Arguments represents a sequence of one or more arguments (AKA actual
// parameters).
type Arguments struct {
	Pos   lexer.Position
	First *Expression   `@@`
	Rest  []*Expression `( "," @@ )*` // xxxxxxxxxxxx does this work as intended?
}

// Expression represents an expression (things that evaluate to a value).
type Expression struct {
	Pos        lexer.Position
	Expression *Assignment `@@`
}

// Assignment represents an assignment expression.
type Assignment struct {
	Pos                 lexer.Position
	QualifiedIdentifier *QualifiedIdentifier `  ( @@`
	Value               *Assignment          `  "=" @@ )`
	Next                *LogicOr             `| @@`
}

// LogicOr represents a logic or expression.
type LogicOr struct {
	Pos   lexer.Position
	First *LogicAnd   `@@`
	Rest  []*LogicAnd `( "or" @@ )*` // xxxxxxxxxxxx does this work as intended?
}

// LogicAnd represents a logic and expression.
type LogicAnd struct {
	Pos   lexer.Position
	First *Equality   `@@`
	Rest  []*Equality `( "and" @@ )*` // xxxxxxxxxxxx does this work as intended?
}

// Equality represents an equality (or un-equality) expression.
// FIXME: I need to store the operator used!
type Equality struct {
	Pos   lexer.Position
	First *Comparison   `@@`
	Rest  []*Comparison `( ( "!=" | "==" ) @@ )*` // xxxxxxxxxxxx does this work as intended?
}

// Comparison represents a comparison expression.
// FIXME: I need to store the operator used!
type Comparison struct {
	Pos   lexer.Position
	First *Addition   `@@`
	Rest  []*Addition `( ( ">" | ">=" | "<" | "<=" ) @@ )*` // xxxxxxxxxxxx does this work as intended?
}

// Addition represents an addition expression.
// FIXME: I need to store the operator used!
type Addition struct {
	Pos   lexer.Position
	First *Multiplication   `@@`
	Rest  []*Multiplication `( ( "-" | "+" ) @@ )*` // xxxxxxxxxxxx does this work as intended?
}

// Multiplication represents a multiplication expression.
// FIXME: I need to store the operator used!
type Multiplication struct {
	Pos   lexer.Position
	First *Exponentiation   `@@`
	Rest  []*Exponentiation `( ( "/" | "*" ) @@ )*` // xxxxxxxxxxxx does this work as intended?
}

// Exponentiation represents an exponentiation expression.
type Exponentiation struct {
	Pos      lexer.Position
	Base     *Unary          `@@`
	Exponent *Exponentiation `( "^" @@ )*`
}

// Unary represents an expression with a unary operator.
type Unary struct {
	Pos      lexer.Position
	Operator *string `  ( "not" | "-" | "+" )`
	Operand  *Unary  `  @@`
	Call     *Call   `| @@`
}

// Call represents a call expression.
type Call struct {
	Pos         lexer.Position
	Primary     *Primary          `@@`
	Complements []*CallComplement `@@*`
}

// CallComplement represents a complement to a call expression (like the
// arguments for a function call or the indexing of an array).
type CallComplement struct {
	Pos                 lexer.Position
	Arguments           *Arguments           `  "(" @@? ")"`
	QualifiedIdentifier *QualifiedIdentifier `| "." @@`
	Index               *Expression          `| "[" @@ "]"`
}

// Primary represents a primary expression.
type Primary struct {
	Pos                     lexer.Position
	Bool                    *bool                `  "true" | "false"`
	Float                   *float64             `| @FLOAT`
	Integer                 *int                 `| @INTEGER`
	String                  *string              `| @STRING`
	Array                   *ArrayLiteral        `| @@`
	Map                     *MapLiteral          `| @@`
	QualifiedIdentifier     *QualifiedIdentifier `| @@`
	ParenthesizedExpression *Expression          `| "(" @@ ")"`
	Gosub                   *Gosub               `| @@`
	Listen                  *Expression          `"listen" @@`
}

// ArrayLiteral represents an array literal.
type ArrayLiteral struct {
	Pos   lexer.Position
	First *Expression   `"[" ( @@`
	Rest  []*Expression `( "," @@ )* ","? )? "]"` // xxxxxxxxxxxx does this work as intended?
}

// MapLiteral represents a map literal.
type MapLiteral struct {
	Pos   lexer.Position
	First *MapEntry   `"{" ( @@`
	Rest  []*MapEntry `( "," @@ )* ","? )? "}"` // xxxxxxxxxxxx does this work as intended?
}

// MapEntry represents one key-value pair in a map.
type MapEntry struct {
	Pos   lexer.Position
	Key   *string     `@@`
	Value *Expression `"=" @@`
}

// Gosub represents a gosub expression.
type Gosub struct {
	Pos                 lexer.Position
	QualifiedIdentifier *QualifiedIdentifier `@@`
	Arguments           *Arguments           `"(" @@ ")"`
}
