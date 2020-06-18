/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package compiler

import (
	"fmt"
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

	// The Chunk containing the code to execute.
	chunk *Chunk

	// The instruction pointer, which points to the next instruction to be
	// executed (it's an index into chunk.Code).
	ip int
}

// NewVM returns a new Virtual Machine.
func NewVM() *VM {
	return &VM{}
}

// Interpret interprets a given chunk of code.
func (vm *VM) Interpret(chunk *Chunk) InterpretResult {
	vm.chunk = chunk
	vm.ip = 0

	return vm.run()
}

// run runs the code in vm.chunk.
func (vm *VM) run() InterpretResult {
	for {
		if vm.DebugTraceExecution {
			vm.chunk.disassembleInstruction(os.Stdout, vm.ip)
		}

		instruction := vm.chunk.Code[vm.ip]
		vm.ip++

		switch instruction {
		case OpConstant:
			constant := vm.readConstant()
			fmt.Printf("%v\n", constant)

		case OpReturn:
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
