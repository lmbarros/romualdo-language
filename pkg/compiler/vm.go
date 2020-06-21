/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package compiler

import (
	"fmt"
	"math"
	"os"
)

// InterpretResult is the result of interpreting some Romualdo code.
type InterpretResult int

const (
	// InterpretOK is used to indicate that the interpretation worked without
	// errors.
	InterpretOK InterpretResult = iota

	// InterpretCompileError is used to indicate a compilation error.
	InterpretCompileError

	// InterpretRuntimeError is used to indicate a runtime error.
	InterpretRuntimeError
)

// VM is a Romualdo Virtual Machine.
type VM struct {
	// Set DebugTraceExecution to true to make the VM disassemble the code as it
	// runs through it.
	DebugTraceExecution bool

	// chunk is the Chunk containing the code to execute.
	chunk *Chunk

	// ip is the instruction pointer, which points to the next instruction to be
	// executed (it's an index into chunk.Code).
	ip int

	// stack is the VM stack, used for storing values during interpretation.
	stack []Value
}

// NewVM returns a new Virtual Machine.
func NewVM() *VM {
	return &VM{}
}

// Interpret interprets a given program, passed as the source code.
func (vm *VM) Interpret(source string) InterpretResult {
	c := &Compiler{}
	c.Compile(source)

	return InterpretOK
}

// run runs the code in vm.chunk.
func (vm *VM) run() InterpretResult { //nolint: gocyclo
	for {
		if vm.DebugTraceExecution {
			fmt.Print("          ")

			for _, v := range vm.stack {
				fmt.Printf("[ %v ]", v)
			}

			fmt.Print("\n")

			vm.chunk.disassembleInstruction(os.Stdout, vm.ip)
		}

		instruction := vm.chunk.Code[vm.ip]
		vm.ip++

		switch instruction {
		case OpConstant:
			constant := vm.readConstant()
			vm.push(constant)

		case OpAdd:
			b := vm.pop()
			a := vm.pop()
			vm.push(a + b)

		case OpSubtract:
			b := vm.pop()
			a := vm.pop()
			vm.push(a - b)

		case OpMultiply:
			b := vm.pop()
			a := vm.pop()
			vm.push(a * b)

		case OpDivide:
			b := vm.pop()
			a := vm.pop()
			vm.push(a / b)

		case OpPower:
			b := float64(vm.pop())
			a := float64(vm.pop())
			vm.push(Value(math.Pow(a, b)))

		case OpNegate:
			vm.push(-vm.pop())

		case OpReturn:
			fmt.Println(vm.pop())
			return InterpretOK
		}
	}
}

// run runs the code in vm.chunk.
func (vm *VM) readConstant() Value {
	constant := vm.chunk.Constants[vm.chunk.Code[vm.ip]]
	vm.ip++

	return constant
}

// push pushes a value into the VM stack.
func (vm *VM) push(value Value) {
	vm.stack = append(vm.stack, value)
}

// pop pops a value from the VM stack and returns it. Panics on underflow.
func (vm *VM) pop() Value {
	top := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]

	return top
}
