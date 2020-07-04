# TODO

## Strategy for AST implementation

* Make the `compileFn` return an AST node (but otherwise keep the same
  behavior).
* Implement an AST printer visitor, to test the AST infrastructure.
* Implement an `Executable` (or something) type, serializable, that is my
  "redistributable" format.
* Reimplement code generation as an AST visitor that generates an `Executable`.
* Implement static type checking for the `"abc" > 123` case.
* Implement static type checking for everything else.

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

* Document the instruction set.

## Smallish improvements

* All easy improvements done!

## Things to benchmark

* Type switches and `reflect.TypeOf` versus an explicit type tag on
  `bytecode.Value`.
