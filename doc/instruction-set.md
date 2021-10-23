# The Romualdo Virtual Machine Instruction Set

Not yet assigning a definitive value (or, er, "byte code") to each instruction,
but let's at least document what we can do.

## Assorted Topics

### Unbounded and Bounded Numbers

Romualdo has three types of numeric values: ints, floats and bnums. Whenever we
mention "unbounded numbers" in this document, we are talking about ints and
floats. This is contrast with bnums which are "bounded numbers".

### Operations Between Different Types

Essentially, the behavior of the VM matches the behavior of the language. In
general, operations between different types are not supported and values of
different types are considered different.

The only exception is when one operand is an int and the other is a float: in
this case, the int one is converted to a float and then the operation is
performed as if both operands were floats.

There is no automatic type conversion like this for any other type, not even for
bounded numbers.

All arithmetic operations between floats result in a float.

All arithmetic operations between bnums result in a bnum.

Most arithmetic operations between ints result in ints. The exceptions are
`DIVIDE` and `POWER`, which always yield float results.

### Immediate operands

Each instruction that has immediate operands interpret them in one of the few
possible ways described below. The description of each instruction tells which
of these interpretations it uses.

* **Unsigned byte.** The operand is a single byte, interpreted as an unsigned
  integer.
* **Signed byte.** The operand is a single byte, interpreted as a signed integer
  encoded in two's complement.
* **Unsigned 32-bit integer.** The operand is a 32-bit unsigned integer, stored
  in little-endian format.
* **Signed 32-bit integer.** The operand is a 32-bit signed integer, stored in
  little-endian byte order, encoded in two's complement.

## The Instructions

Instructions are listed in alphabetical order.

The fields "Pops" and "Pushes" describe the effects as perceived by the user,
not necessarily how the implementation works. For example, if you see some
instruction that pops a value and then pushes the same value back to the stack,
the implementation is free to leave the stack untouched.

### `ADD`

**Purpose:** Adds two unbounded numeric values.  
**Immediate Operands:** None.  
**Pops:** Two unbounded numeric values, *B* and *A*.  
**Pushes:** One value, the result of computing *A* + *B*.

### `ADD_BNUM`

**Purpose:** Adds two bounded numbers.  
**Immediate Operands:** None.  
**Pops:** Two bounded numbers, *B* and *A*.  
**Pushes:** One value, the result of computing the bounded sum *A* + *B*.

### `BLEND`

**Purpose:** Performs the blending operation on three bounded numbers.  
**Immediate Operands:** None.  
**Pops:** Three bounded numbers, *C*, *B* and *A*.  
**Pushes:** One value, the result of blending *A* and *B* with the weight *C*.

### `CONSTANT`

**Purpose:** Loads a constant with index in the [0, 255] interval.  
**Immediate Operands:** One byte *A*, interpreted as an index into the constant
pool.  
**Pops:** Nothing.  
**Pushes:** One value, the value of constant taken at the index *A* of the
constant pool.

### `CONSTANT_LONG`

**Purpose:** Loads a constant with index in the [0, 2^32] interval.  
**Immediate Operands:** A 32-bit unsigned integer.  
**Pops:** Nothing.  
**Pushes:** One value, the value of the constant at the index obtained from the
immediate operand.

If the constant you need is in the [0, 255] interval, it's better to use the
more efficient `CONSTANT` instruction.

### `DIVIDE`

**Purpose:** Divides two unbounded numeric values.  
**Immediate Operands:** None.  
**Pops:** Two unbounded numeric values, *B* and *A*.  
**Pushes:** One float value, the result of computing *A* / *B*.

The result is always a float, even if the result is a whole number.

### `EQUAL`

**Purpose:** Checks if two values are equal.  
**Immediate Operands:** None.  
**Pops:** Two values, *B* and *A*.  
**Pushes:** One Boolean value telling if *A* = *B*.

### `FALSE`

**Purpose:** Loads a `false` value.  
**Immediate Operands:** None.  
**Pops:** Nothing.  
**Pushes:** One Boolean value: `false`.

### `GREATER`

**Purpose:** Checks if a values is greater than another value.  
**Immediate Operands:** None.  
**Pops:** Two values, *B* and *A*.  
**Pushes:** One Boolean value telling if *A* > *B*.

### `GREATER_EQUAL`

**Purpose:** Checks if a values is greater than or equal to another value.  
**Immediate Operands:** None.  
**Pops:** Two values, *B* and *A*.  
**Pushes:** One Boolean value telling if *A* ≥ *B*.

### `JUMP`

**Purpose:** Jumps to a different location unconditionally.  
**Immediate Operands:** One signed byte, interpreted as the offset to jump.  
**Pops:** Nothing.  
**Pushes:** Nothing.  
**Other Effects:** Increments the instruction pointer by the amount taken as an
immediate operand. (The increment happens after this instruction and its operand
were fully read.)

### `JUMP_LONG`

**Purpose:** Jumps to a different location unconditionally, even if it is far
away.  
**Immediate Operands:** One signed 32-bit integer, interpreted as the offset to
jump.  
**Pops:** Nothing.  
**Pushes:** Nothing.  
**Other Effects:** Increments the instruction pointer by the amount taken as an
immediate operand. (The increment happens after this instruction and its operand
were fully read.)

If the jump offset fits into a signed 8-bit value, it is more efficient to use
`JUMP` instead.

### `JUMP_IF_FALSE`

**Purpose:** Jumps to a different location maybe.  
**Immediate Operands:** One signed byte, interpreted as the offset to jump.  
**Pops:** One Boolean value *A*.  
**Pushes:** Nothing.  
**Other Effects:** If *A* is a Boolean value and is false, increments the
instruction pointer by the amount taken as an immediate operand. (The increment
happens after this instruction and its operand were fully read.)

### `JUMP_IF_FALSE_LONG`

**Purpose:** Jumps to a different location maybe, even if it is far away.  
**Immediate Operands:** One signed 32-bit integer, interpreted as the offset to
jump.  
**Pops:** One Boolean value *A*.  
**Pushes:** Nothing.  
**Other Effects:** If *A* is a Boolean value and is false, increments the
instruction pointer by the amount taken as an immediate operand. (The increment
happens after this instruction and its operand were fully read.)

If the jump offset fits into a signed 8-bit value, it is more efficient to use
`JUMP_IF_FALSE` instead.

### `JUMP_IF_FALSE_NO_POP`

**Purpose:** Jumps to a different location maybe, leaving the stack intact.  
**Immediate Operands:** One signed byte, interpreted as the offset to jump.  
**Pops:** One Boolean value *A*.  
**Pushes:** The same value popped, *A*.  
**Other Effects:** If *A* is a Boolean value and is false, increments the
instruction pointer by the amount taken as an immediate operand. (The increment
happens after this instruction and its operand were fully read.)

This is equivalent to `JUMP_IF_FALSE`, but doesn't pop the tested value from the
stack.

### `JUMP_IF_FALSE_NO_POP_LONG`

**Purpose:** Jumps to a different location maybe, even if it is far away,
leaving the stack intact.  
**Immediate Operands:** One signed 32-bit integer, interpreted as the offset to
jump.  
**Pops:** One Boolean value *A*.  
**Pushes:** The same value popped, *A*.  
**Other Effects:** If *A* is a Boolean value and is false, increments the
instruction pointer by the amount taken as an immediate operand. (The increment
happens after this instruction and its operand were fully read.)

If the jump offset fits into a signed 8-bit value, it is more efficient to use
`JUMP_IF_FALSE_NO_POP` instead.

This is equivalent to `JUMP_IF_FALSE_LONG`, but doesn't pop the tested value
from the stack.

### `JUMP_IF_TRUE_NO_POP`

**Purpose:** Jumps to a different location maybe, leaving the stack intact.  
**Immediate Operands:** One signed byte, interpreted as the offset to jump.  
**Pops:** One Boolean value *A*.  
**Pushes:** The same value popped, *A*.  
**Other Effects:** If *A* is a Boolean value and is true, increments the
instruction pointer by the amount taken as an immediate operand. (The increment
happens after this instruction and its operand were fully read.)

### `JUMP_IF_TRUE_NO_POP_LONG`

**Purpose:** Jumps to a different location maybe, even if it is far away,
leaving the stack intact.  
**Immediate Operands:** One signed 32-bit integer, interpreted as the offset to
jump.  
**Pops:** One Boolean value *A*.  
**Pushes:** The same value popped, *A*.  
**Other Effects:** If *A* is a Boolean value and is true, increments the
instruction pointer by the amount taken as an immediate operand. (The increment
happens after this instruction and its operand were fully read.)

If the jump offset fits into a signed 8-bit value, it is more efficient to use
`JUMP_IF_TRUE_NO_POP` instead.

### `LESS`

**Purpose:** Checks if a values is less than another value.  
**Immediate Operands:** None.  
**Pops:** Two values, *B* and *A*.  
**Pushes:** One Boolean value, telling if *A* <> *B*.

### `LESS_EQUAL`

**Purpose:** Checks if a values is less than or equal to another value.  
**Immediate Operands:** None.  
**Pops:** Two values, *B* and *A*.  
**Pushes:** One Boolean value, telling if *A* ≤ *B*.

### `MULTIPLY`

**Purpose:** Multiplies two unbounded numeric values.  
**Immediate Operands:** None.  
**Pops:** Two unbounded numeric values, *B* and *A*.  
**Pushes:** One value, the result of computing *A* × *B*.

### `NEGATE`

**Purpose:** Performs arithmetic negation on numeric values.  
**Immediate Operands:** None.  
**Pops:** One numeric value, *A*.  
**Pushes:** One numeric value, -*A*.

Note that, unlike other arithmetic instructions, this one is shared between
bounded and unbounded numbers.

### `NOP`

**Purpose:** Does nothing.  
**Immediate Operands:** None.  
**Pops:** Nothing.  
**Pushes:** Nothing.

I can't really see any purpose for a no-op instruction in the Romualdo VM, but I
*really* wanted to have it. That's probably because of the tender memories I
have of `NOP` in the x86 architecture. Whatever.

### `NOT`

**Purpose:** Performs logical negation.  
**Immediate Operands:** None.  
**Pops:** One Boolean value, *A*.  
**Pushes:** One Boolean value, ¬*A*.

### `NOT_EQUAL`

**Purpose:** Checks if two values are different.  
**Immediate Operands:** None.  
**Pops:** Two values, *B* and *A*.  
**Pushes:** One Boolean value, telling if *A* ≠ *B*.

### `POP`

**Purpose:** Pops the value on the top of the stack.  
**Immediate Operands:** None.  
**Pops:** One values.  
**Pushes:** Nothing.

### `POWER`

**Purpose:** Raises an unbounded numeric value to the power of another unbounded
numeric value.  
**Immediate Operands:** None.  
**Pops:** Two unbounded numeric values, *B* and *A*.  
**Pushes:** One float value, the result of computing *A* to the *B*-th power.

AKA exponentiation.

### `PRINT`

**Purpose:** Prints a value.  
**Immediate Operands:** None.  
**Pops:** One value, the one to be printed.  
**Pushes:** Nothing.

Printing exists primarily for debugging or demo purposes. VM implementations
should try to provide a meaningful implementation, but it would not be a sin for
an implementation to make this a no-op.

### `READ_GLOBAL`

**Purpose:** Reads the value of a global variable.  
**Immediate Operands:** One byte *A*, interpreted as an index into the globals
pool.  
**Pops:** Nothing.  
**Pushes:** One value, the value of the global variable taken at the index *A*
of the globals pool.

### `READ_LOCAL`

**Purpose:** Reads the value of a local variable.  
**Immediate Operands:** One byte *A*, interpreted as the index into the stack
where the desired variable is stored.  
**Pops:** Nothing.  
**Pushes:** One value, the value of the local variable taken at the index *A*
of the stack.

### `RETURN`

TODO

### `SUBTRACT`

**Purpose:** Subtracts two unbounded numeric values.  
**Immediate Operands:** None.  
**Pops:** Two unbounded numeric values, *B* and *A*.  
**Pushes:** One value, the result of computing *A* - *B*.

### `SUBTRACT_BNUM`

**Purpose:** Subtracts two bounded numbers.  
**Immediate Operands:** None.  
**Pops:** Two bounded numbers, *B* and *A*.  
**Pushes:** One value, the result of computing the bounded subtraction *A* -
*B*.

### `TO_INT`

**Purpose:** Converts a value to an integer.  
**Immediate Operands:** None.  
**Pops:** Two values, *B* (the int to return if the conversion fails) and *A*
(the value to convert to an int).  
**Pushes:** One integer value, which is either *A* converted to an int or *B*
(if *A* cannot be converted to an int).

TODO: Define semantics. For example, what is `int(0.1b)`? Or `int(10.9)`?

### `TO_FLOAT`

**Purpose:** Converts a value to a floating-point number.  
**Immediate Operands:** None.  
**Pops:** Two values, *B* (the float to return if the conversion fails) and *A*
(the value to convert to a float).  
**Pushes:** One float value, which is either *A* converted to a float or *B*
(if *A* cannot be converted to a float).

TODO: Define semantics. For example, what is `int(0.1b)`?

### `TO_BNUM`

**Purpose:** Converts a value to a bounded number.  
**Immediate Operands:** None.  
**Pops:** Two values, *B* (the float to return if the conversion fails) and *A*
(the value to convert to a bnum).  
**Pushes:** One float value, which is either *A* converted to a bnum or *B*
(if *A* cannot be converted to a float within the bnum range).

Note that, in the VM, bnums are represented as floats. This instruction differs
from `TO_FLOAT` in that the conversion will fail (and thus *B* will be returned)
if the numeric conversion to float works but the result is outside of the bnum
valid range (-1, 1).

TODO: Define semantics.

### `TO_STRING`

**Purpose:** Converts a value to a string.  
**Immediate Operands:** None.  
**Pops:** One value, *A*, the value to convert to a string.  
**Pushes:** One value: the string representation of *A*.

### `TRUE`

**Purpose:** Loads a `true` value.  
**Immediate Operands:** None.  
**Pops:** Nothing.  
**Pushes:** One Boolean value: `true`.

### `WRITE_GLOBAL`

**Purpose:** Writes the value of a global variable.  
**Immediate Operands:** One byte *A*, interpreted as an index into the globals
pool.  
**Pops:** One value, the new value the global variable value at the
index *A* will be set to.  
**Pushes:** One value, the same that was popped.

### `WRITE_LOCAL`

**Purpose:** Writes the value of a local variable.  
**Immediate Operands:** One byte *A*, interpreted as an index into the stack
where the variable to be written is located at.  
**Pops:** One value, the new value the local variable.  
**Pushes:** One value, the same that was popped.
