# TODO

## Strategy for AST implementation

* ~~Make the `compileFn` return an AST node (but otherwise keep the same
  behavior).~~
* ~~Implement an AST printer visitor, to test the AST infrastructure.~~
* ~~Reimplement code generation as an AST visitor that generates an
  `Executable`.~~
* Implement static type checking for the `"abc" > 123` case.
* Implement static type checking for everything else.
* Add support for `int`s and `bnum`s.
* Add support for type conversions. Something like this:
    * `int("124")`: `string` to `int`. Returns zero if invalid.
    * `int("124", -999)`: `string` to `int`. Returns -999 if invalid.
* Implement serialization and deserialization of `CompiledStoryworld`.

## Big things

* Implement static typing.
    * Integers x floats could be a good initial case.
    * But then, I have strings in already.
    * I guess the way to go here is:
        * Produce an AST
        * Each AST node has a type
        * When generating code, look at the types on the AST to check for type
          errors and to know what type-specific opcodes to use.

## Not so big, but not small either

* Testing
    * Add some kind of end-to-end testing: if we run this Storyworld with this
      input, the output must be this one.
    * I am sure there are more things that can be unit tested.
* Maybe I'll want to keep string literals stored in a more organized way, with
  the string contents themselves in an easy to serialize format. (Because since
  I'll want to dump my bytecode to a ready-to-consume format). One key point for
  this is `compiler.go`, function `stringLiteral()`.
* Implement smarter storage of line numbers in the Chunk. Something more
  efficient than storing one line number per instruction.
* Document the instruction set.

## Smallish improvements

* Implement a decent compiler UI
    * Print AST should be an option to it.

## Things to benchmark

* Type switches and `reflect.TypeOf` versus an explicit type tag on
  `bytecode.Value`.
