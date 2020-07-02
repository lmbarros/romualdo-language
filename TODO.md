# TODO

## Big things

* Implement static typing.
    * Integers x floats could be a good initial case.

## Not so big, but not small either

* Testing
    * Add some kind of end-to-end testing: if we run this Storyworld with this
      input, the output must be this one.
    * I am sure there are more things that can be unit tested.
* Document the instruction set.

## Smallish improvements

* Add the unary plus operator.
* Implement opcodes for `!=`, `<=`, `>=`.

## Things to benchmark

* Type switches and `reflect.TypeOf` versus an explicit type tag on
  `bytecode.Value`.
