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
            | globalsBlock
            | aliasDecl
            | enumDecl
            | structDecl
            | functionDecl
            | passageDecl ;
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

varDecl = IDENTIFIER ":" type ( "=" expression )? ;
```

### Globals

A `globalsBlock` is used to declare and initialize global variables. Each
`globalsBlock` must include a version and there can be only one `globalsBlock`
of any given version in the set of files passed to the compiler. When compared
to a `globalsBlock` of the previous version, a given version can add new
globals, but cannot remove existing ones. It is also OK to redeclare a variable
with the same name and type but with a possibly changed initializer (the new
initializer is used when starting a new story).

```ebnf
globalsBlock = "globals" "@" INTEGER
            varDecl*
            "end" ;
```

### Type declarations

These are the user-defined types (more on them later):

* Alias: basically a new name for an existing type. They are a different type,
  though, and you cannot pass an alias for the aliased without an explicit
  conversion.
* Enum: a type whose value must come from a collection of possible named values.
* Struct: A type that aggregates various fields of certain types.

```ebnf
aliasDecl = "alias" IDENTIFIER type ;

enumDecl = "enum" IDENTIFIER IDENTIFIER* "end" ;

structDecl = "struct" IDENTIFIER varDecl* "end" ;
```

### Functions

Functions are like the functions in programming languages. Take arguments, do
stuff, return value.

Functions are not versioned.

```ebnf
functionDecl = "function" IDENTIFIER "(" parameters? ")" ":" type
               statement*
               "end" ;

parameters = parameter ( "," parameter )* ;

parameter = IDENTIFIER ":" type ;
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
              "end" ;
```

### Digression: Functions versus Passages

Why having both functions and Passages? Is this really needed? How do they
differ? I don't know, this area is still a design TODO, but here's what I have
in mind for the first implementation.

First, I think there's a clear conceptual difference between them; they have
different purposes. Passages are the structuring element of storyworlds, they
provide support for storytelling. Functions are for computation. Of course, at
least on the current design, they behave a lot like each other. But what if I
want to drift the design towards something significantly different? I think
there's potential on the idea that Passages and functions should be more
different, to make each one better in doing its job: either telling stories or
computing. The limits are not sharp, but to put in Crawfordian terms, a Passage
is more about talking and listening, while a function is about thinking
(especially heavy thinking).

Having functions and Passages as separate entities in the language will allow me
to evolve them in different directions, and this sounds like a good enough
reason to make them separate things from the start. It's mostly about the
future.

Back to the present. For now, functions and Passages will be a lot like each
other. Here are the differences I intend to implement as of now:

* Functions are "more dynamic": we can create anonymous functions on the spot,
  closures, etc. Passages are pretty much static: we must know at compile-time
  what are all the Passages. Each passage has possibly mutable metadata
  associated with it (kinda like `static` variables on a C function), and this
  metadata must be usable by algorithms that decide what Passage to run next. I
  think that everything gets much harder to reason about if we allow Passages to
  be "as dynamic as functions". So:
    * "Function lambdas" are OK, "Passage lambdas" are not OK.
    * Function closures are OK, Passage closures are not OK.
    * Functions don't have a metadata block, Passages do.
* Passages can use `say` and `listen`; functions can't.
* Passages are versioned, functions are not.
    * That's because functions can't use `listen` and therefore they are
      uninterruptible. We'll never have a function "in progress" in the
      serializable VM state.

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
     | "[" "]" type
     | functionType
     | qualifiedIdentifier ;

functionType = "function" "(" typeList? ")" ":" type ;

typeList = type ( "," type )* ;

qualifiedIdentifier = IDENTIFIER ( "." IDENTIFIER )* ;
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

TODO: About the type-unsafety of `map`s: We need to redesign `map` accesses in
such a way that the syntax always require a default value to be provided. So, if
the desired key is not there, we get a default value, not a runtime error. (This
might be tricky to implement, especially to have a nice syntax for chains of
nested `map`s.)

## Statements

Statements are language constructs that do stuff. They don't have a value.

```ebnf
statement = expression
          | varDeclStmt
          | blockStmt
          | whileStmt
          | ifStmt
          | returnStmt
          | gotoStmt
          | sayStmt ;

varDeclStmt = "var" varDecl ;

blockStmt = "do"
            statement*
            "end" ;

whileStmt = "while" expression "do"
            statement*
            "end" ;

ifStmt = "if" expression "then" statement*
         elseif*
         ( "else" statement* )?
         "end" ;

elseif = "elseif" expression "then" statement*

returnStmt = "return" expression? ;

gotoStmt = "goto" qualifiedIdentifier "(" arguments? ")" ;

sayStmt = "say" expression ;

arguments = expression ( "," expression )* ;
```

Some notes about the statements:

* Expressions can be used as statements. Depending on the expression this can be
  useful (a function call is often used for its side-effects only) or useless
  (an expression like `1 + 1` by itself serves no purpose -- but is considered
  valid nevertheless).
* Local variable declarations can appear anywhere (anywhere a statement can
  appear, *bien entendu*). A local variable exists from the point it is declared
  until the end of its scope. A local variable cannot shadow an existing local
  variable.
* The only purpose of `do`...`end` statements is to create blocks, which allow
  to control the lifetime of the enclosed local variables. I honestly didn't
  intend to have this on the language, but I added them to allow me having local
  variables before I have other block-defining statements. Maybe I'll remove it
  in the future.
* Nothing surprising about `while` loops: execute a sequence of statements as
  long as a given expression evaluates to `true`.
* Nothing surprising with `if`s either.
* Ditto for `return`s.
* Romualdo has a `goto` statement, but not exactly that one [considered
  harmful](https://homepages.cwi.nl/~storm/teaching/reader/Dijkstra68.pdf)! Our
  `goto`s are used to stop the execution of the current passage and move the
  control to another passage.
    * TODO: OK, but how does this work, really? If there is a `gosub` underneath
      the caller, it is expecting a return value of a certain type. So a `goto`
      can never go to a passage with a different return type. I ask myself: is
      `goto` any good? Why not using always `gosub` followed by a `return`? The
      call stack would remain more informative -- we'd know for sure how we got
      to that point.
        * TODO: But then, do I even need `gosub`? I mean, just call the passage
          like a function and be happy. I can introduce special syntax later if
          I want.
* The `say` statement is used to send information to the driver program that is
  running the storyworld. Typically, it is used to describe events that happened
  in the story and need to be somehow shown to the player (the *how* in the
  *somehow* is responsibility of the driver, not of Romualdo). The `expression`
  after the `say` keyword must evaluate to a `map`.
* `for` loops are the most notable absence here. I want to support things like
  `for i in range(0, 10) do ... end` and `for t in arrayOfThings do ... end`,
  but then I'd have to store some additional state (the `range()` result, the
  current pointer into `arrayOfThings`) somewhere and don't know where this
  somewhere would be. Probably in a local variable. Anyway, for now, `for` loops
  are not available at all.

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

call = primary callComplement* ;

callComplement = "(" arguments? ")" | "." qualifiedIdentifier | "[" expression "]" ;

primary = "true" | "false"
        | FLOAT
        | INTEGER
        | STRING
        | arrayLiteral
        | mapLiteral
        | qualifiedIdentifier
        | "(" expression ")"
        | gosub
        | "listen" expression ;

arrayLiteral = "[" ( expression ( "," expression )* ","? )? "]" ;

mapLiteral = "{" ( mapEntry   ( "," mapEntry   )* ","? )? "}" ;

gosub = "gosub" qualifiedIdentifier "(" arguments? ")" ;

mapEntry = IDENTIFIER "=" expression ;
```

Notes about expressions:

* TODO: Do I really want assignments as expressions? Maybe they should be
  statements, for an arguably cleaner language.
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
* Strings will probably be extended. One thing I want to add are localized
  strings. Not sure about the syntax and semantics, but I want them!
    * Simple unorthodox idea: use only backticks to enclose strings. This way,
      single and double quotes (often useful in storytelling!) are always
      available inside strings.
    * Or: backtick-delimited strings are the special ones, that are
      automatically exported to a translatable file and whatnot.

## Final words

This grammar is, to a large extent, inspired by the [Lox
grammar](http://www.craftinginterpreters.com/appendix-i.html) in Robert
Nystrom's great *Crafting Interpreters*.

The Romualdo syntax is, in no small extent, inspired by
[Lua](https://www.lua.org).
