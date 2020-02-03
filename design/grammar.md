# The Romualdo Language Commented Grammar

The Romualdo language (Romulang for short) is not used to create programs, it is
used to create *storyworlds*. A storyworld is a sequence of declarations (which
can all be in a single file or in split over multiple files.)

```ebnf
storyworld = declaration* ;
```

## Declarations

```ebnf
declaration = metaBlock
            | varBlock
            | typeDecl
            | functionDecl
            | passageDecl
            ;
```

### Metadata block

Metadata blocks are used to attach data to the storyworld itself and to
passages. At global level they are immutable and not very interesting -- but the
versions of the global `metaBlock`s are used to version the whole storyworld.
Every update of a storyworld needs to have a new `metaBlock` with a new version.
The first version is `1`, all subsequent versions must be incremented by one.

At local (passage) level, metadata behave kinda like `static` variables in
traditional languages: they are mutable information that belongs to the passage
itself, not to a particular "instance" of the passage. They are specially useful
because one can select passages according to the values of its metadata.

```ebnf
metaBlock = "meta" "@" INTEGER varDecl* "end" ;

varDecl = IDENTIFIER ":" type "=" expression ;
```

### Variable declarations

For implementation simplicity (*sigh*), all variables must be declared in
`varBlock`s. They can appear at global scope (that is, as a top-level
declaration) and at local scope (inside functions or passages).

At global scope they must include a version. There can be only one `varBlock` of
any given version in the set of files passed to the compiler. A global
`varBlock` of a given version can add new globals, but cannot remove existing
ones. It is also OK to redeclare a variable with the same name and type but with
a possibly changed initializer (the new initializer is used when starting a new
story).

At local scope, the `varBlock` must not include a version.

```ebnf
varBlock = "vars" ( "@" INTEGER )? varDecl* "end" ;
```

### Type declarations

These are the user-defined types (more on them later):

* Alias: basically a new name for an existing type. They are a different type,
  though, and you cannot pass an alias for the aliased without an explicit
  conversion.
* Enum: a type whose value must come from a collection of possible named values.
* Struct: A type that aggregates various fields of certain types.

```ebnf
typeDecl = aliasDecl
         | enumDecl
         | structDecl ;

aliasDecl = "alias" IDENTIFIER type ;

enumDecl = "enum" IDENTIFIER IDENTIFIER* "end" ;

structDecl = "struct" IDENTIFIER varDecl* "end";
```

### Functions

Functions are like the functions in programming languages. Take arguments, do
stuff, return value.

Functions are not versioned.

```ebnf
functionDecl = "function" IDENTIFIER "(" parameters? ")" ":" type
               statement*
               "end" ;

parameters = IDENTIFIER ":" type ( "," IDENTIFIER ":" type )*
```

### Passages

Passages are the building block of stories. Recursive and fractal. The whole
story and its smallest part can both be passages. They work a lot like
functions, but:

* They can have meta blocks (the meta block itself can't contain a version).
* They can use the special statements for interaction and for calling other
  passages.
* They can be interrupted (to handle user input).
* They are versioned (because a player can upgrade the storyworld and then load
  a saved version coming from an older version of the storyworld).
* Maybe more stuff I forgot. 🙂

```ebnf
passageDecl = "passage" IDENTIFIER "@" INTEGER "(" parameters? ")" ":" type
              metaBlock?
              statement*
              "end"
```

## Types

Romualdo is strongly-typed (well, mostly).

```ebnf
type = "bool"
     | "int"
     | "float"
     | "bnum"
     | "string"
     | "map"
     | "void"
     | type "[" "]"
     | "function" "(" ( type ( "," type )* )? ")" ":" type
     | qualifiedIdentifier ;

qualifiedIdentifier = IDENTIFIER ( "." IDENTIFIER )*
```

The supported types are:

* `bool`: Booleans, true or false, no surprise here, right?
* `int`: Integer number, size not formally specified for now (but c'mon, this is
  the 21st century, no less than 32 bits is a reasonable assumption);
* `float`: Floating point number, most likely a IEEE 754 binary64 (double
  precision) number (but you should not count on that).
* `bnum`: Chris Crawford's bounded numbers, which I hope will be nice for doing
  story things like character models (that's what `bnum`s were designed for,
  anyway).
* `string`: A string of text. Should support Unicode, though I know enough
  Unicode to know that there must be several Unicode things not supported
  (yet?).
* `map`: A JSON-like thing, mapping string keys to other values. Also "the type
  unsafe corners of the language". Handy for communication with the external
  world, though.
* `void`: A non-type. Used when a type is formally required, but is not really
  needed (like the return value of a function that doesn't return anything).
* Arrays: a sequence of zero or more elements of the same type. `int[]` is an
  array of `int`s, `string[]` is an array of `string`s, and so on.
* Functions: functions taking a certain set of parameters and returning a
  certain type.
* User-defined types: that's why we have that `qualifiedIdentifier` in the list
  of types. It is "qualified" instead of a regular `IDENTIFIER` because it may
  contain a namespace.

## Statements

Statements are language constructs that do stuff. They don't have a value.

```ebnf
statement = expression
          | whileStmt
          | ifStmt
          | returnStmt
          | gotoStmt
          | sayStmt ;

whileStmt = "while" expression "do"
            statement*
            "done" ;

ifStmt = "if" expression "then" statement*
         ( "elseif" expression "then" statement* )*
         ( "else" statement* )?
         "done" ;

returnStmt = "return" expression? ;

gotoStmt = "goto" qualifiedIdentifier "(" arguments? ")" ;

sayStmt = "say" expression ;

arguments = expression ( "," expression )*
```

Some notes about the statements:

* Expressions can be used as statements. Depending on the expression this can be
  useful (an assignment often is used by itself) or useless (an expression like
  `1 + 1` by itself serves no purpose -- but is considered valid nevertheless).
* Nothing surprising about `while` loops: execute a sequence of statements as
  long as a given expression evaluates to `true`.
* Nothing surprising with `if`s either.
* Ditto for `return`s.
* Romualdo has a `goto` statement, but not exactly that one [considered
  harmful](https://homepages.cwi.nl/~storm/teaching/reader/Dijkstra68.pdf)! Our
  `goto`s are used to stop the execution of the current passage and move the
  control to another passage.
* The `say` statement is used to send information to the driver program that is
  running the storyworld. Typically, it is used to describe events that happened
  in the story and need to be somehow shown to the player (the *how* in the
  *somehow* is responsibility of the driver, not of Romualdo). The `expression`
  after the `say` keyword must evaluate to a `map`.
* `for` loops are the most notable absence here. I want to support things like
  `for i in range(0, 10) do ... end` and `for t in arrayOfThings do ... end`,
  but then I'd have to store some additional state (the `range()` result, the
  current pointer into `arrayOfThings`) somewhere and don't know where this
  somewhere would be. Maybe to implement `for` loops I'll need to relax my
  variable declaration rules (now I require them to be all declared in a single
  `vars` block). Anyway, for now, `for` loops are not available at all.

## Expressions

Expressions evaluate to a value. The different levels of precedence are encoded
in the grammar (which makes the grammar weirder to look at, but will hopefully
translate more directly to the implementation).

```ebnf
expression = assignment ;

assignment = qualifiedIdentifier "=" assignment
           | logicOr ;

logicOr = logicAnd ( "or" logicAnd )* ;

logicAnd = equality ( "and" equality )* ;

equality = comparison ( ( "!=" | "==" ) comparison )* ;

comparison = addition ( ( ">" | ">=" | "<" | "<=" ) addition )* ;

addition = multiplication ( ( "-" | "+" ) multiplication )* ;

multiplication = exponentiation ( ( "/" | "*" ) exponentiation )* ;

exponentiation = unary ( "^" exponentiation )* ;

unary = ( "not" | "-" | "+" ) unary
      | call ;

call = primary ( "(" arguments? ")" | "." qualifiedIdentifier | "[" expression "]" )* ;

primary = "true" | "false"
        | FLOAT
        | INTEGER
        | STRING
        | "[" ( expression ( "," expression )* ","? )? "]"
        | "{" ( mapEntry   ( "," mapEntry   )* ","? )? "}"
        | qualifiedIdentifier
        | "(" expression ")"
        | "gosub" qualifiedIdentifier "(" arguments? ")"
        | "listen" expression ;

mapEntry = IDENTIFIER "=" expression ;
```

Notes about expressions:

* In the `call` rule, I am using `qualifiedIdentifier` in a way that is not
  semantically correct: it will match both real qualified identifiers and
  chains of `struct`member accesses.
* The `gosub` expression must appear in a passage. It calls another passage, and
  returns with that passage's return value.
* The `listen` expression is used to get input from the player. Its `expression`
  argument must evaluate to a `map`, which represents the alternatives the
  player has. `listen` transfer the control to the driver program, which can
  access the data from this `map`, show the player alternatives, get a choice
  from the player and call the storyworld again passing the player choice (at
  this point the storyworld gets the control again). The player choice is the
  `listen` expression value, and is always a `map`;
* Logical operators `and` and `or` have short-circuited evaluation.
* Note the syntax for literal arrays and maps.

## Lexical Grammar

```ebnf
EOL = ( "\n" | "\r" )+ ;

COMMENT = "#" ⟨anything except EOL⟩ EOL ;

NONZERO_DIGIT = "1" ... "9" ;

DIGIT = "0" | NONZERO_DIGIT ;

INTEGER = NONZERO_DIGIT DIGIT* ;

FLOAT = DIGIT+ "." DIGIT+
      | DIGIT+ ( "." DIGIT+ )? ("e" | "E") ( "+" | "-" )? DIGIT+ ;

STRING = '"' ( ⟨anything except '"' or "\\" ⟩ | '\\"' ) '"';

LETTER_LIKE = ⟨Unicode Letter⟩ | ⟨Unicode Emoji⟩ | "_" ;

IDENTIFIER = LETTER_LIKE ( LETTER_LIKE | DIGIT )*
```

Notes:

* Identifiers that start with an upper case Unicode letter (LU category) are
  public (accessible everywhere in the storyworld). Other identifiers are
  visible only locally (in the file where they are declared).
* Strings will probably be extended. On thing I want to add are localized
  strings. Not sure about the syntax and semantics, but I want them!

## Final words

This grammar is, to a large extent, inspired by the [Lox
grammar](http://www.craftinginterpreters.com/appendix-i.html) in Robert
Nystrom's great *Crafting Interpreters*.

The Romualdo syntax is, in no small extent, inspired by
[Lua](https://www.lua.org).
