# TODO

## Big things

* Nothing!

## Not so big, but not small either

* Add support for the blend operator for `bnum`s.
* Add support for type conversions. Something like this:
    * `int("124")`: `string` to `int`. Returns zero if invalid.
    * `int("124", -999)`: `string` to `int`. Returns -999 if invalid.
* Implement serialization and deserialization of `CompiledStoryworld`.
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

## Smallish improvements

* Move the common visitor stuff (like keeping the current node) to some reusable
  `struct`.
* Implement a decent compiler UI
    * Print AST should be an option to it.
    * Disassemble the code (either from a binary or the just compiled code),
      too.

## Things to benchmark

* Type switches and `reflect.TypeOf` versus an explicit type tag on
  `bytecode.Value`.
