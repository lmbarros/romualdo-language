/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package vm

import (
	"fmt"
	"math"
	"os"

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
}

// New returns a new Virtual Machine.
func New() *VM {
	return &VM{}
}

// Interpret interprets a given program, passed as the source code.
func (vm *VM) Interpret(chunk *bytecode.Chunk) bool {
	vm.chunk = chunk
	return vm.run()
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
				vm.push(bytecode.NewValueString(a + b))

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
			fmt.Println(vm.pop())
			return true

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
	index := bytecode.ThreeBytesToInt(
		vm.chunk.Code[vm.ip], vm.chunk.Code[vm.ip+1], vm.chunk.Code[vm.ip+2])

	constant := vm.chunk.Constants[index]
	vm.ip += 3
	return constant
}

// push pushes a value into the VM stack.
func (vm *VM) push(value bytecode.Value) {
	vm.stack = append(vm.stack, value)
}

// pop pops a value from the VM stack and returns it. Panics on underflow.
func (vm *VM) pop() bytecode.Value {
	top := vm.stack[len(vm.stack)-1]
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
func (vm *VM) popTwoIntOperands() (a int64, b int64, ok bool) {
	if !vm.peek(0).IsInt() || !vm.peek(1).IsInt() {
		vm.runtimeError("Operands must be integer numbers.")
		return
	}
	b = vm.pop().AsInt()
	a = vm.pop().AsInt()
	ok = true
	return
}

// popTwoUnboundedNumberOperands pops and returns two values from the stack,
// assumed to be integers ot floats, to be used as operands of a binary
// operator.
func (vm *VM) popTwoUnboundedNumberOperands() (a float64, b float64, ok bool) {
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
func (vm *VM) popTwoStringOperands() (a string, b string, ok bool) {
	if !vm.peek(0).IsString() || !vm.peek(1).IsString() {
		vm.runtimeError("Operands must be strings.")
		return
	}
	b = vm.pop().AsString()
	a = vm.pop().AsString()
	ok = true
	return
}
