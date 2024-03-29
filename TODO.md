# TODO

## Plan

* Implement:
    * Passages.
    * `say`, `listen`, `goto`, and `gosub`.
* Leave the book aside for a while and focus on tooling as I envision it:
    * Separate compilation and execution.
    * Split debugging information to a separate file.
    * Use command-line arguments or "commands" to:
        * Compile to bytecode.
        * Run bytecode.
        * Enable or disable trace execution.
        * Disassemble a binary.
    * Some side-by-side display of source code and assembly would be nice.
    * Allow to split the storyworld into multiple files.
    * Add a proper test suite.
    * Review wording of error messages. Maybe include an error code always (good
      to add test cases that are expected to fail in some particular way).
    * Romualdo syntax highlighting for VS Code would be cool.
* Then we can go back to the book and finish the language, adding more tests as
  we go.

## TODOs

* Bug: Cannot initialize a variable with a negative constant!
* One thing to consider: what's the performance hit of having only `LONG`
  instructions? It would make many things so much simpler...
* Give some thought to the blend operator syntax. I am currently using `a~b~c`
  to blends from `a` to `b` weighted by `c`, but:
    * It doesn't look good with negative numbers: `-0.4b~0.3b~-0.15b`.
    * It would require parentheses for blend within blend.
    * Must also think about its precedence. I just used anything that didn't
      felt too wrong.
    * Idea for syntax: `{a,b,c}`. Nests nicely, curly braces are available and
      also look a bit "blendy" with all this curliness.
* For completeness, we should have the `JUMP_IF_TRUE` and `JUMP_IF_TRUE_LONG`
  instructions. (We currently have only the `NO_POP` versions of them).
* Right now, `Chunk.SetGlobal()` and `Chunk.GetGlobalIndex()` look linearly into
  the array of globals. I might want to add a map from names to indices at some
  point. In this case, would not need to store the global name in
  `Chunk.Globals`.
    * Or, if possible: store globals on the stack (they would be created and
      initialized by some new opcode that would be called at the start of the
      generated code), in which case I guess globals and locals would be pretty
      much the same thing from the VM point-of-view.
* I'd like to add constants to the language at some point.
* Implement serialization and deserialization of `CompiledStoryworld`.
* Add `break` and `continue` to loops. (As I write this, only `while` loops
  exist, but who knows the future?)
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
* Remove duplication between `typeChecker` and `semanticChecker`.
* Allow more than 255 globals. 2^32 should be the way to go.
* Allow more than 255 locals. 2^32 should be the way to go.
* I don't like the discrepancy in the naming of opcodes `READ_GLOBAL` and
  `CONSTANT`. They are both doing kind of the same thing, but only one has the
  `READ` prefix. (On the other hand, I guess there will never be an opcode to
  write a constant, so this is not really wrong.)
* We currently do not detect infinite recursion. Given that Romualdo is meant to
  be embedded into other programs, a Storyworld should not be able to crash the
  host by pushing infinite call frames into call stack and exhausting all
  memory. But how to deal with this? Add a configurable limit to the stack size?
  Fine, but then, what can the host do?
    * For reference: `vm.VM.callValue()` or `vm.VM.callFunction()` would be
      natural places to add any check of this kind.
    * In the longer term, I am looking for better error handling in general.
      Ideally, somehow, the story should always be able to go on (if not very
      consistently), regardless of any bug in the Storyworld. Under this
      perspective, an overflowing call stack may be just another error to be
      handled.
* Should accept bnums for comparison operators
* Move the common visitor stuff (like keeping the current node) to some reusable
  `struct`.
* Implement a decent compiler UI
    * Print AST should be an option to it.
    * Disassemble the code (either from a binary or the just compiled code),
      too.
* Idea: a compiler flag that allows to ignore (or include) Passages with a
  certain metadata value. I guess this would be useful for implementing demos
  that contain only a subset of the complete Storyworld, plus maybe some
  additional or modified Passages (alternative implementations for the demo).

## Ideas

### Cleaner syntax for `say`

Imagine something like this:

```text
passage p(): void
    `Once upon a time
    `a knight with shining armor
    `left his castle in search of adventure.
end
```

The block prefixed with backtick characters is converted to a string. Then it is
filtered trough a certain function that converts from a string to a map. Then
this map is `say`d. In other words, that's a shortcut for a `say` statement.

Said filtering function could be user-provided (or a default is used). The
standard library could provide a number of useful ones.

## Things to benchmark

* Type switches and `reflect.TypeOf` versus an explicit type tag on
  `bytecode.Value`.
