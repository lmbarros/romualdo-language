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

## The Instructions

Instructions are listed in alphabetical order.

### `ADD`

**Purpose:** Adds two unbounded numeric values.  
**Immediate Operands:** None.  
**Pops:** Two unbounded numeric values, B and A.  
**Pushes:** One value, the result of computing *A* + *B*.

### `ADD_BNUM`

**Purpose:** Adds two bounded numbers.  
**Immediate Operands:** None.  
**Pops:** Two bounded numbers, B and A.  
**Pushes:** One value, the result of computing the bounded sum of *A* and *B*.

### `CONSTANT`

**Purpose:** Loads a constant with index in the [0, 255] interval.  
**Immediate Operands:** One byte *A*, interpreted as an index into the constant
pool.  
**Pops:** Nothing.  
**Pushes:** One value, the value of constant taken at the index *A* of the
constant pool.

### `CONSTANT_LONG`

**Purpose:** Loads a constant with index in the [0, 16777215] interval.  
**Immediate Operands:** Three bytes, *A*, *B*, *C*, interpreted as a 24-bit
index into the constant pool. This value is stored in a little endian format (in
other words, *A* is the least significant byte, *C* is the most significant
one).  
**Pops:** Nothing.  
**Pushes:** One value, the value of the constant at the index obtained from the
immediate operands.

If the constant you need is in the [0, 255] interval, it's better to use the
more efficient `CONSTANT` instruction.

### `DIVIDE`

**Purpose:** Divides two unbounded numeric values.  
**Immediate Operands:** None.  
**Pops:** Two unbounded numeric values, B and A.  
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
**Pops:** Two unbounded numeric values, B and A.  
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

### `POWER`

**Purpose:** Raises an unbounded numeric value to the power of another unbounded
numeric value.  
**Immediate Operands:** None.  
**Pops:** Two unbounded numeric values, B and A.  
**Pushes:** One float value, the result of computing *A* to the *B*-th power.

AKA exponentiation.

### `RETURN`

TODO

### `SUBTRACT`

**Purpose:** Subtracts two unbounded numeric values.  
**Immediate Operands:** None.  
**Pops:** Two unbounded numeric values, B and A.  
**Pushes:** One value, the result of computing *A* - *B*.

### `SUBTRACT_BNUM`

**Purpose:** Subtracts two bounded numbers.  
**Immediate Operands:** None.  
**Pops:** Two bounded numbers, B and A.  
**Pushes:** One value, the result of computing the bounded subtraction *A* -
*B*.

### `TRUE`

**Purpose:** Loads a `true` value.  
**Immediate Operands:** None.  
**Pops:** Nothing.  
**Pushes:** One Boolean value: `true`.
