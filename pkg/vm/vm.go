/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2021 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package vm

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"gitlab.com/stackedboxes/romulang/pkg/bytecode"
)

// VM is a Romualdo Virtual Machine.
type VM struct {
	// Set DebugTraceExecution to true to make the VM disassemble the code as it
	// runs through it.
	DebugTraceExecution bool

	// chunk is the Chunk containing the code to execute.
	chunk *bytecode.Chunk

	// ip is the instruction pointer, which points to the next instruction to be
	// executed (it's an index into chunk.Code).
	ip int

	// stack is the VM stack, used for storing values during interpretation.
	stack []bytecode.Value

	// strings is the store of interned strings used by the VM. All strings are
	// added here.
	strings *bytecode.StringInterner
}

// New returns a new Virtual Machine.
func New() *VM {
	return &VM{
		strings: bytecode.NewStringInterner(),
	}
}

// Interpret interprets a given program, passed as the source code.
func (vm *VM) Interpret(chunk *bytecode.Chunk) bool {
	vm.chunk = chunk
	vm.strings = chunk.Strings
	r := vm.run()
	if len(vm.stack) != 0 {
		panic(fmt.Sprintf("Stack size should be zero after execution, was %v.", len(vm.stack)))
	}
	return r
}

// NewInternedValueString creates a new Value initialized to the interned string
// value v. Emphasis on "interned": if there is already some other string value
// equal to v on this VM, we'll reuse that same memory in the returned value.
//
// TODO: Someday, when I have more stuff working, do some benchmarking. Remove
// this call to intern() and see if the performance/memory different is
// significant in typical usage.
func (vm *VM) NewInternedValueString(v string) bytecode.Value {
	s := vm.strings.Intern(v)
	return bytecode.NewValueString(s)
}

// run runs the code in vm.chunk.
func (vm *VM) run() bool { // nolint: funlen, gocyclo, gocognit
	for {
		if vm.DebugTraceExecution {
			fmt.Print("          ")

			for _, v := range vm.stack {
				fmt.Printf("[ %v ]", v)
			}

			fmt.Print("\n")

			vm.chunk.DisassembleInstruction(os.Stdout, vm.ip)
		}

		instruction := vm.chunk.Code[vm.ip]
		vm.ip++

		switch instruction {
		case bytecode.OpNop:
			break

		case bytecode.OpConstant:
			constant := vm.readConstant()
			vm.push(constant)

		case bytecode.OpConstantLong:
			constant := vm.readLongConstant()
			vm.push(constant)

		case bytecode.OpTrue:
			vm.push(bytecode.NewValueBool(true))

		case bytecode.OpFalse:
			vm.push(bytecode.NewValueBool(false))

		case bytecode.OpPop:
			vm.pop()

		case bytecode.OpEqual:
			b := vm.pop()
			a := vm.pop()
			vm.push(bytecode.NewValueBool(bytecode.ValuesEqual(a, b)))

		case bytecode.OpNotEqual:
			b := vm.pop()
			a := vm.pop()
			vm.push(bytecode.NewValueBool(!bytecode.ValuesEqual(a, b)))

		case bytecode.OpGreater:
			a, b, ok := vm.popTwoUnboundedNumberOperands()
			if !ok {
				return false
			}
			vm.push(bytecode.NewValueBool(a > b))

		case bytecode.OpGreaterEqual:
			a, b, ok := vm.popTwoUnboundedNumberOperands()
			if !ok {
				return false
			}
			vm.push(bytecode.NewValueBool(a >= b))

		case bytecode.OpLess:
			a, b, ok := vm.popTwoUnboundedNumberOperands()
			if !ok {
				return false
			}
			vm.push(bytecode.NewValueBool(a < b))

		case bytecode.OpLessEqual:
			a, b, ok := vm.popTwoUnboundedNumberOperands()
			if !ok {
				return false
			}
			vm.push(bytecode.NewValueBool(a <= b))

		case bytecode.OpAdd:
			switch {
			case vm.peek(0).IsString() && vm.peek(1).IsString():
				a, b, ok := vm.popTwoStringOperands()
				if !ok {
					return false
				}
				vm.push(vm.NewInternedValueString(a + b))

			case vm.peek(0).IsInt() && vm.peek(1).IsInt():
				a, b, ok := vm.popTwoIntOperands()
				if !ok {
					return false
				}
				vm.push(bytecode.NewValueInt(a + b))

			default:
				a, b, ok := vm.popTwoUnboundedNumberOperands()
				if !ok {
					return false
				}
				vm.push(bytecode.NewValueFloat(a + b))
			}

		case bytecode.OpAddBNum:
			a, b, ok := vm.popTwoFloatOperands()
			if !ok {
				return false
			}
			a = boundedInverseTransform(a)
			b = boundedInverseTransform(b)
			vm.push(bytecode.NewValueFloat(boundedTransform(a + b)))

		case bytecode.OpSubtract:
			if vm.peek(0).IsInt() && vm.peek(1).IsInt() {
				a, b, ok := vm.popTwoIntOperands()
				if !ok {
					return false
				}
				vm.push(bytecode.NewValueInt(a - b))
			} else {
				a, b, ok := vm.popTwoUnboundedNumberOperands()
				if !ok {
					return false
				}
				vm.push(bytecode.NewValueFloat(a - b))
			}

		case bytecode.OpSubtractBNum:
			a, b, ok := vm.popTwoFloatOperands()
			if !ok {
				return false
			}
			a = boundedInverseTransform(a)
			b = boundedInverseTransform(b)
			vm.push(bytecode.NewValueFloat(boundedTransform(a - b)))

		case bytecode.OpMultiply:
			if vm.peek(0).IsInt() && vm.peek(1).IsInt() {
				a, b, ok := vm.popTwoIntOperands()
				if !ok {
					return false
				}
				vm.push(bytecode.NewValueInt(a * b))
			} else {
				a, b, ok := vm.popTwoUnboundedNumberOperands()
				if !ok {
					return false
				}
				vm.push(bytecode.NewValueFloat(a * b))
			}

		case bytecode.OpDivide:
			a, b, ok := vm.popTwoUnboundedNumberOperands()
			if !ok {
				return false
			}
			vm.push(bytecode.NewValueFloat(a / b))

		case bytecode.OpPower:
			a, b, ok := vm.popTwoUnboundedNumberOperands()
			if !ok {
				return false
			}
			vm.push(bytecode.NewValueFloat(math.Pow(a, b)))

		case bytecode.OpBlend:
			x, y, weight, ok := vm.popThreeFloatOperands()
			if !ok {
				return false
			}
			uWeight := 1 - ((1 - weight) / 2)
			result := y*uWeight + x*(1-uWeight)
			vm.push(bytecode.NewValueFloat(result))

		case bytecode.OpJump:
			jumpOffset := vm.chunk.Code[vm.ip]
			vm.ip += int(jumpOffset + 1)

		case bytecode.OpJumpIfFalse:
			jumpOffset := vm.chunk.Code[vm.ip]
			vm.ip++
			cond := vm.pop()
			if cond.IsBool() && !cond.AsBool() {
				vm.ip += int(jumpOffset)
			}

		case bytecode.OpNot:
			if !vm.peek(0).IsBool() {
				vm.runtimeError("Operand must be a Boolean value.")
				return false
			}
			vm.push(bytecode.NewValueBool(!vm.pop().AsBool()))

		case bytecode.OpNegate:
			switch {
			case vm.peek(0).IsInt():
				vm.push(bytecode.NewValueInt(-vm.pop().AsInt()))

			case vm.peek(0).IsFloat():
				vm.push(bytecode.NewValueFloat(-vm.pop().AsFloat()))

			default:
				vm.runtimeError("Operand must be a number.")
				return false
			}

		case bytecode.OpReturn:
			return true

		case bytecode.OpToInt:
			if !vm.peek(0).IsInt() {
				vm.runtimeError("Default value for conversion to int must be an integer number.")
				return false
			}
			d := vm.pop().AsInt()
			v := vm.pop()

			switch {
			case v.IsInt():
				vm.push(v)
			case v.IsFloat():
				vm.push(bytecode.NewValueInt(int64(v.AsFloat())))
			case v.IsBool():
				r := int64(0)
				if v.AsBool() {
					r = 1
				}
				vm.push(bytecode.NewValueInt(r))
			case v.IsString():
				r, err := strconv.ParseInt(v.AsString(), 10, 64)
				if err != nil {
					r = d
				}
				vm.push(bytecode.NewValueInt(r))
			default:
				panic(fmt.Sprintf("Unexpected type on conversion to int: %T", v))
			}

		case bytecode.OpToFloat:
			if !vm.peek(0).IsFloat() {
				vm.runtimeError("Default value for conversion to float must be a floating point number.")
				return false
			}
			d := vm.pop().AsFloat()
			v := vm.pop()

			switch {
			case v.IsFloat():
				vm.push(v)
			case v.IsInt():
				vm.push(bytecode.NewValueFloat(float64(v.AsInt())))
			case v.IsBool():
				r := float64(0.0)
				if v.AsBool() {
					r = 1.0
				}
				vm.push(bytecode.NewValueFloat(r))
			case v.IsString():
				r, err := strconv.ParseFloat(v.AsString(), 64)
				if err != nil {
					r = d
				}
				vm.push(bytecode.NewValueFloat(r))
			default:
				panic(fmt.Sprintf("Unexpected type on conversion to float: %T", v))
			}

		case bytecode.OpToBNum:
			if !vm.peek(0).IsFloat() {
				vm.runtimeError("Default value for conversion to bnum must be a floating point number.")
				return false
			}
			d := vm.pop().AsFloat()
			v := vm.pop()

			switch {
			case v.IsFloat():
				r := v.AsFloat()
				if r <= 0.0 || r >= 1.0 {
					r = d
				}
				vm.push(bytecode.NewValueFloat(r))
			case v.IsString():
				r, err := strconv.ParseFloat(v.AsString(), 64)
				if err != nil || r <= 0.0 || r >= 1.0 {
					r = d
				}
				vm.push(bytecode.NewValueFloat(r))
			default:
				panic(fmt.Sprintf("Unexpected type on conversion to bnum: %T", v))
			}

		case bytecode.OpToString:
			v := vm.pop()
			switch {
			case v.IsString():
				vm.push(v)
			case v.IsFloat():
				r := strconv.FormatFloat(v.AsFloat(), 'f', -1, 64)
				vm.push(vm.NewInternedValueString(r))
			case v.IsInt():
				r := strconv.FormatInt(v.AsInt(), 10)
				vm.push(vm.NewInternedValueString(r))
			case v.IsBool():
				r := strconv.FormatBool(v.AsBool())
				vm.push(vm.NewInternedValueString(r))
			default:
				panic(fmt.Sprintf("Unexpected type on conversion to float: %T", v))
			}

		case bytecode.OpPrint:
			v := vm.pop()
			fmt.Printf("%v\n", v)

		case bytecode.OpReadGlobal:
			value := vm.readGlobal()
			vm.push(value)

		case bytecode.OpReadLocal:
			value := vm.stack[vm.chunk.Code[vm.ip]]
			vm.ip++
			vm.push(value)

		case bytecode.OpWriteGlobal:
			value := vm.top()
			vm.writeGlobal(value)

		case bytecode.OpWriteLocal:
			value := vm.top()
			index := vm.chunk.Code[vm.ip]
			vm.ip++
			vm.stack[index] = value

		default:
			panic(fmt.Sprintf("Unexpected instruction: %v", instruction))
		}
	}
}

// readConstant reads a single-byte constant index from the chunk bytecode and
// returns the corresponding constant value.
func (vm *VM) readConstant() bytecode.Value {
	constant := vm.chunk.Constants[vm.chunk.Code[vm.ip]]
	vm.ip++

	return constant
}

// readConstant reads a three-byte constant index from the chunk bytecode and
// returns the corresponding constant value.
func (vm *VM) readLongConstant() bytecode.Value {
	index := bytecode.DecodeUInt31(vm.chunk.Code[vm.ip:])
	constant := vm.chunk.Constants[index]
	vm.ip += 4
	return constant
}

// readGlobal reads a single-byte global index from the chunk bytecode and
// returns the corresponding global variable value.
func (vm *VM) readGlobal() bytecode.Value {
	value := vm.chunk.Globals[vm.chunk.Code[vm.ip]]
	vm.ip++
	return value.Value
}

// writeGlobal sets the value of a global variable to value. For the variable,
// reads a single-byte from the chunk bytecode and uses it as the index into the
// globals table.
func (vm *VM) writeGlobal(value bytecode.Value) {
	vm.chunk.Globals[vm.chunk.Code[vm.ip]].Value = value
	vm.ip++
}

// push pushes a value into the VM stack.
func (vm *VM) push(value bytecode.Value) {
	vm.stack = append(vm.stack, value)
}

// top returns the value on the top of the VM stack (without removing it).
// Panics on underflow.
func (vm *VM) top() bytecode.Value {
	return vm.stack[len(vm.stack)-1]
}

// pop pops a value from the VM stack and returns it. Panics on underflow.
func (vm *VM) pop() bytecode.Value {
	top := vm.top()
	vm.stack = vm.stack[:len(vm.stack)-1]
	return top
}

// peek returns a value on the stack that is a given distance from the top.
// Passing 0 means "give me the value on the top of the stack". The stack is not
// changed at all.
func (vm *VM) peek(distance int) bytecode.Value {
	return vm.stack[len(vm.stack)-1-distance]
}

// runtimeError stops the execution and reports a runtime error with a given
// message and fmt.Printf-like arguments.
func (vm *VM) runtimeError(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	line := vm.chunk.Lines[vm.ip-1]
	fmt.Fprintf(os.Stderr, "[line %d] in script\n", line)
}

// popTwoIntOperands pops and returns two values from the stack, assumed to be
// integers, to be used as operands of a binary operator.
func (vm *VM) popTwoIntOperands() (a, b int64, ok bool) {
	if !vm.peek(0).IsInt() || !vm.peek(1).IsInt() {
		vm.runtimeError("Operands must be integer numbers.")
		return
	}
	b = vm.pop().AsInt()
	a = vm.pop().AsInt()
	ok = true
	return
}

// popTwoFloatOperands pops and returns two values from the stack, assumed to be
// floating point numbers, to be used as operands of a binary operator.
func (vm *VM) popTwoFloatOperands() (a, b float64, ok bool) {
	if !vm.peek(0).IsFloat() || !vm.peek(1).IsFloat() {
		vm.runtimeError("Operands must be floating point numbers.")
		return
	}
	b = vm.pop().AsFloat()
	a = vm.pop().AsFloat()
	ok = true
	return
}

// popThreeFloatOperands pops and returns three values from the stack, assumed
// to be floating point numbers, to be used as operands of an operator.
func (vm *VM) popThreeFloatOperands() (a, b, c float64, ok bool) {
	if !vm.peek(0).IsFloat() || !vm.peek(1).IsFloat() || !vm.peek(2).IsFloat() {
		vm.runtimeError("Operands must be floating point numbers.")
		return
	}
	c = vm.pop().AsFloat()
	b = vm.pop().AsFloat()
	a = vm.pop().AsFloat()
	ok = true
	return
}

// popTwoUnboundedNumberOperands pops and returns two values from the stack,
// assumed to be integers ot floats, to be used as operands of a binary
// operator.
func (vm *VM) popTwoUnboundedNumberOperands() (a, b float64, ok bool) {
	b, ok = vm.popUnboundedNumberOperand()
	if !ok {
		return
	}

	a, ok = vm.popUnboundedNumberOperand()
	if !ok {
		return
	}

	ok = true
	return
}

// popUnboundedNumberOperand pops and returns twoone values from the stack,
// assumed to be and integer ot float, to be used as and operand of some
// operator.
func (vm *VM) popUnboundedNumberOperand() (v float64, ok bool) {
	switch {
	case vm.peek(0).IsFloat():
		return vm.pop().AsFloat(), true
	case vm.peek(0).IsInt():
		return float64(vm.pop().AsInt()), true
	default:
		vm.runtimeError("Operands must be integer or floating-point numbers.")
		return 0.0, false
	}
}

// popTwoStringOperands pops and returns two values from the stack, assumed to
// be strings to be used as operands of a binary operator.
func (vm *VM) popTwoStringOperands() (a, b string, ok bool) {
	if !vm.peek(0).IsString() || !vm.peek(1).IsString() {
		vm.runtimeError("Operands must be strings.")
		return
	}
	b = vm.pop().AsString()
	a = vm.pop().AsString()
	ok = true
	return
}

// boundedTransform transforms an unbounded number to a bounded one. See "Chris
// Crawford on Interactive Storytelling, 2nd Ed." page 184.
func boundedTransform(unboundedNumber float64) float64 {
	if unboundedNumber > 0 {
		return 1 - (1 / (1 + unboundedNumber))
	}
	return (1 / (1 - unboundedNumber)) - 1
}

// boundedInverseTransform transforms a bounded number to an unbounded one. See
// "Chris Crawford on Interactive Storytelling, 2nd Ed." page 184.
func boundedInverseTransform(boundedNumber float64) float64 {
	if boundedNumber > 0 {
		return (1 / (1 - boundedNumber)) - 1
	}
	return 1 - (1 / (1 + boundedNumber))
}
