# The Romualdo Language

A programming language for Interactive Storytelling.

Probably should look more like a domain-specific language than a general-purpose
language, but since I don't know the domain well it will probably lean too much
to the general-purpose side. Romualdo 2.0 will probably look considerably
different.

## Notes to self

In order to not have to relearn this for the next I spend 5 months without
looking at this code, here are some common steps to make a language change:

* Add a new AST `Node` subtype at `pkg/ast/nodes.go`.
* Change the parser at `pkg/frontend/parser.go`; some new function will return a
  node of this new type.
* Maybe add some new semantic checks at `pkg/frontend/semantic_checker.go`.
* Maybe add some new type checks at `pkg/frontend/type_checker.go`.
* If a new opcode is needed:
    * Document it at `doc/instruction-set.md`.
    * Add it at `pkg/bytecode/chunk.go`.
    * Generate code for this new opcode at `pkg/backend/code_gen.go`.
    * Add code to interpret it at `pkg/vm/vm.go`.

And to add a new opcode:

* Document the new opcode in `doc/instruction-set.md`.
* Add the opcode constant to `pkg/bytecode/chunk.go`
* Implement the opcode runtime behavior in `pkg/vm/vm.go`.
* Implement the disassembling of the opcode in `pkg/bytecode/chunk.go`.

## Credits

* The Romualdo Language syntax is in no small extent inspired by
  [Lua](http://www.lua.org).
* The implementation of the compiler and virtual machine are strongly based on
  Bob Nystrom's excellent [Crafting
  Interpreters](http://www.craftinginterpreters.com). And maybe this is an
  understatement.
