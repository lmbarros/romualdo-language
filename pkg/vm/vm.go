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
	"gitlab.com/stackedboxes/romulang/pkg/compiler"
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
func (vm *VM) Interpret(source string) InterpretResult {
	// TODO: Move the compilation out of here. The vm should read the bytecode
	// directly.
	c := compiler.New()
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

			vm.chunk.DisassembleInstruction(os.Stdout, vm.ip)
		}

		instruction := vm.chunk.Code[vm.ip]
		vm.ip++

		switch instruction {
		case bytecode.OpConstant:
			constant := vm.readConstant()
			vm.push(constant)

		case bytecode.OpAdd:
			b := vm.pop()
			a := vm.pop()
			vm.push(a + b)

		case bytecode.OpSubtract:
			b := vm.pop()
			a := vm.pop()
			vm.push(a - b)

		case bytecode.OpMultiply:
			b := vm.pop()
			a := vm.pop()
			vm.push(a * b)

		case bytecode.OpDivide:
			b := vm.pop()
			a := vm.pop()
			vm.push(a / b)

		case bytecode.OpPower:
			b := float64(vm.pop())
			a := float64(vm.pop())
			vm.push(bytecode.Value(math.Pow(a, b)))

		case bytecode.OpNegate:
			vm.push(-vm.pop())

		case bytecode.OpReturn:
			fmt.Println(vm.pop())
			return InterpretOK
		}
	}
}

// run runs the code in vm.chunk.
func (vm *VM) readConstant() bytecode.Value {
	constant := vm.chunk.Constants[vm.chunk.Code[vm.ip]]
	vm.ip++

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
