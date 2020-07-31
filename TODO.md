# TODO

* BUG: `nodes.go`, `func (n *Binary) Type()`: Need to take the "new" types
  (float, bnum) into account.

## Big things

* Given some thought to the blend operator syntax. I am currently using `a~b~c`
  to blends from `a` to `b` weighted by `c`, but this doesn't look good for
  something like `-0.4b~0.3b~-0.15b`. It would require parentheses for blend
  within blend.
    * Must also think about its precedence. I just used anything that didn't
      felt too wrong.

## Not so big, but not small either

* Implement serialization and deserialization of `CompiledStoryworld`.
* Testing
    * Add some kind of end-to-end testing: if we run this Storyworld with this
      input, the output must be this one.
    * I am sure there are more things that can be unit tested.
* Implement smarter storage of line numbers in the Chunk. Something more
  efficient than storing one line number per instruction.
* On the VM, I currently use floats to represent bnums. Works nicely, except
  when converting a bnum to a string, in which case I'd like to have something
  like "0.1b" instead of just "0.1".

## Smallish improvements

* Should accept bnums for comparison operators
* Move the common visitor stuff (like keeping the current node) to some reusable
  `struct`.
* Implement a decent compiler UI
    * Print AST should be an option to it.
    * Disassemble the code (either from a binary or the just compiled code),
      too.

## Things to benchmark

* Type switches and `reflect.TypeOf` versus an explicit type tag on
  `bytecode.Value`.
