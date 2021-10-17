# TODO

## Bug Fix

* Cannot initialize a variable with a negative value!

## Big things

* Given some thought to the blend operator syntax. I am currently using `a~b~c`
  to blends from `a` to `b` weighted by `c`, but this doesn't look good for
  something like `-0.4b~0.3b~-0.15b`. It would require parentheses for blend
  within blend.
    * Must also think about its precedence. I just used anything that didn't
      felt too wrong.
* Right now, `Chunk.SetGlobal()` and `Chunk.GetGlobalIndex()` look linearly into
  the array of globals. I might want to add a map from names to indices at some
  point. In this case, would not need to store the global name in
  `Chunk.Globals`.
    * Or, if possible: store globals on the stack (they would be created and
      initialized by some new opcode that would be called at the start of the
      generated code), in which case I guess globals and locals would be pretty
      much the same thing from the VM point-of-view.
* I'd like to add constants to the language at some point.

## Not so big, but not small either

* Implement serialization and deserialization of `CompiledStoryworld`.
* Testing
    * Add some kind of end-to-end testing: if we run this Storyworld with this
      input, the output must be this one.
        * A language specification in prose, with examples that can be
          automatically extracted and used as test cases would be cool.
    * I am sure there are more things that can be unit tested.
* Implement smarter storage of line numbers in the Chunk. Something more
  efficient than storing one line number per instruction.
* On the VM, I currently use floats to represent bnums. Works nicely, except
  when converting a bnum to a string, in which case I'd like to have something
  like "0.1b" instead of just "0.1".
* Handling of globals: I might be able to refer to them by a in integer index
  instead of by name. (Which would be faster.)
* This is about the spec: to avoid confusing users, I say that a local variable
  cannot shadow a previously declared local. But what about global variables?
  I'd like to be consistent, but it's weird if creating a new global breaks
  previously working functions because of name clashes. Would also make
  impractical using function libraries, because a function in a library may or
  may not work depending on the globals I have on my storyworld. But if locals
  can shadow globals, maybe I'll want some syntax to "force access" the global
  one. Must think about this. (Also relevant: [section
  22.4.2](http://www.craftinginterpreters.com/local-variables.html#another-scope-edge-case).)
* Add a `POPN` instruction to pop *n* values from the stack at once. Use it to
  optimize the cleanup of locals when exiting of scopes.
* Avoid that linear search when resolving local variables.

## Smallish improvements

* Remove duplication between `typeChecker` and `semanticChecker`.
* Allow more than 255 globals. 2^24 should be enough, as with constants.
* Allow more than 255 locals. 2^24 should be enough, as with constants.
* I don't like the discrepancy in the naming of opcodes `READ_GLOBAL` and
  `CONSTANT`. They are both doing kind of the same thing, but only one has the
  `READ` prefix. (On the other hand, I guess there will never be an opcode to
  write a constant, so this is not really wrong.)
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
