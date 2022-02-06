/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2022 Leandro Motta Barros                                     *
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

	// csw is the compiled storyworld we are executing.
	csw *bytecode.CompiledStoryworld

	// debugInfo contains the debug information corresponding to csw.
	// TODO: Make this optional. If nil, issue less friendly error messages,
	// etc.
	debugInfo *bytecode.DebugInfo

	// stack is the VM stack, used for storing values during interpretation.
	stack *Stack

	// frames is the stack of call frames. It has one entry for every function
	// that has started running bit hasn't returned yet.
	frames []*callFrame

	// The current call frame (the one on top of VM.frames).
	frame *callFrame
}

// New returns a new Virtual Machine.
func New() *VM {
	return &VM{
		stack: &Stack{},
	}
}

// currentChunk returns the chunk currently being executed.
func (vm *VM) currentChunk() *bytecode.Chunk {
	return vm.csw.Chunks[vm.frame.function.ChunkIndex]
}

// currentLines returns the map from instruction to source code lines for the
// chunk currently being executed. Returns nil if vm.debugInfo == nil.
func (vm *VM) currentLines() []int {
	if vm.debugInfo == nil {
		return nil
	}
	return vm.debugInfo.ChunksLines[vm.frame.function.ChunkIndex]
}

// readByte reads a byte from the current Chunk.
func (vm *VM) readByte() byte {
	index := vm.frame.ip
	vm.frame.ip++
	return vm.currentChunk().Code[index]
}

// Interpret interprets a given compiled Storyworld.
// TODO: DebugInfo should be optional.
func (vm *VM) Interpret(csw *bytecode.CompiledStoryworld, di *bytecode.DebugInfo) bool {
	vm.csw = csw
	vm.debugInfo = di

	// TODO: Eventually, we'll start from a Passage, not a function.
	f := bytecode.Function{ChunkIndex: csw.FirstChunk}
	vm.frames = append(vm.frames, &callFrame{
		function: f,
		stack:    vm.stack.createView(),
	})
	vm.frame = vm.frames[0]

	r := vm.run()
	if vm.stack.size() != 0 {
		vm.runtimeError("Stack size should be zero after execution, was %v.", vm.stack.size())
	}
	return r
}

// NewInternedValueString creates a new Value initialized to the interned string
// value v. Emphasis on "interned": if there is already some other string value
// equal to v on this VM, we'll reuse that same memory in the returned value.
//
// TODO: Someday, when I have more stuff working, do some benchmarking. Remove
// this call to intern() and see if the performance/memory difference is
// significant in typical usage.
func (vm *VM) NewInternedValueString(v string) bytecode.Value {
	s := vm.csw.Strings.Intern(v)
	return bytecode.NewValueString(s)
}

// run runs the code in vm.chunk.
func (vm *VM) run() bool { // nolint: funlen, gocyclo, gocognit
	for {
		if vm.DebugTraceExecution {
			fmt.Print("          ")

			for _, v := range vm.stack.data {
				fmt.Printf("[ %v ]", v)
			}

			fmt.Print("\n")

			vm.csw.DisassembleInstruction(vm.currentChunk(), os.Stdout, vm.frame.ip, vm.currentLines())
		}

		instruction := vm.currentChunk().Code[vm.frame.ip]
		vm.frame.ip++

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
			jumpOffset := int8(vm.readByte())
			vm.frame.ip += int(jumpOffset)

		case bytecode.OpJumpLong:
			jumpOffset := bytecode.DecodeSInt32(vm.currentChunk().Code[vm.frame.ip:])
			vm.frame.ip += jumpOffset + 4

		case bytecode.OpJumpIfFalse:
			jumpOffset := int8(vm.readByte())
			cond := vm.pop()
			if cond.IsBool() && !cond.AsBool() {
				vm.frame.ip += int(jumpOffset)
			}

		case bytecode.OpJumpIfFalseNoPop:
			jumpOffset := int8(vm.readByte())
			if vm.peek(0).IsBool() && !vm.peek(0).AsBool() {
				vm.frame.ip += int(jumpOffset)
			}

		case bytecode.OpJumpIfFalseLong:
			jumpOffset := bytecode.DecodeSInt32(vm.currentChunk().Code[vm.frame.ip:])
			vm.frame.ip += 4
			cond := vm.pop()
			if cond.IsBool() && !cond.AsBool() {
				vm.frame.ip += jumpOffset
			}

		case bytecode.OpJumpIfFalseNoPopLong:
			jumpOffset := bytecode.DecodeSInt32(vm.currentChunk().Code[vm.frame.ip:])
			vm.frame.ip += 4
			if vm.peek(0).IsBool() && !vm.peek(0).AsBool() {
				vm.frame.ip += jumpOffset
			}

		case bytecode.OpJumpIfTrueNoPop:
			jumpOffset := vm.readByte()
			if vm.peek(0).IsBool() && vm.peek(0).AsBool() {
				vm.frame.ip += int(jumpOffset)
			}

		case bytecode.OpJumpIfTrueNoPopLong:
			jumpOffset := bytecode.DecodeSInt32(vm.currentChunk().Code[vm.frame.ip:])
			vm.frame.ip += 4
			if vm.peek(0).IsBool() && vm.peek(0).AsBool() {
				vm.frame.ip += jumpOffset
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
				vm.runtimeError("Unexpected type on conversion to int: %T", v.Value)
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
				vm.runtimeError("Unexpected type on conversion to float: %T", v.Value)
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
				vm.runtimeError("Unexpected type on conversion to bnum: %T", v.Value)
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
				vm.runtimeError("Unexpected type on conversion to string: %T", v.Value)
			}

		case bytecode.OpPrint:
			v := vm.pop()
			fmt.Printf("%v\n", v)

		case bytecode.OpReadGlobal:
			value := vm.readGlobal()
			vm.push(value)

		case bytecode.OpReadLocal:
			index := vm.readByte()
			value := vm.frame.stack.at(int(index))
			vm.push(value)

		case bytecode.OpWriteGlobal:
			value := vm.top()
			vm.writeGlobal(value)

		case bytecode.OpWriteLocal:
			value := vm.top()
			index := vm.readByte()
			vm.frame.stack.setAt(int(index), value)

		default:
			vm.runtimeError("Unexpected instruction: %v", instruction)
		}
	}
}

// readConstant reads a single-byte constant index from the chunk bytecode and
// returns the corresponding constant value.
func (vm *VM) readConstant() bytecode.Value {
	index := vm.readByte()
	constant := vm.csw.Constants[index]
	return constant
}

// readConstant reads a three-byte constant index from the chunk bytecode and
// returns the corresponding constant value.
func (vm *VM) readLongConstant() bytecode.Value {
	chunk := vm.currentChunk()
	index := bytecode.DecodeUInt31(chunk.Code[vm.frame.ip:])
	constant := vm.csw.Constants[index]
	vm.frame.ip += 4
	return constant
}

// readGlobal reads a single-byte global index from the chunk bytecode and
// returns the corresponding global variable value.
func (vm *VM) readGlobal() bytecode.Value {
	value := vm.csw.Globals[vm.currentChunk().Code[vm.frame.ip]]
	vm.frame.ip++
	return value.Value
}

// writeGlobal sets the value of a global variable to value. For the variable,
// reads a single-byte from the chunk bytecode and uses it as the index into the
// globals table.
func (vm *VM) writeGlobal(value bytecode.Value) {
	vm.csw.Globals[vm.currentChunk().Code[vm.frame.ip]].Value = value
	vm.frame.ip++
}

// push pushes a value into the VM stack.
func (vm *VM) push(value bytecode.Value) {
	vm.stack.push(value)
}

// top returns the value on the top of the VM stack (without removing it).
// Panics on underflow.
func (vm *VM) top() bytecode.Value {
	return vm.stack.top()
}

// pop pops a value from the VM stack and returns it. Panics on underflow.
func (vm *VM) pop() bytecode.Value {
	return vm.stack.pop()
}

// peek returns a value on the stack that is a given distance from the top.
// Passing 0 means "give me the value on the top of the stack". The stack is not
// changed at all.
func (vm *VM) peek(distance int) bytecode.Value {
	return vm.stack.peek(distance)
}

// runtimeError stops the execution and reports a runtime error with a given
// message and fmt.Printf-like arguments.
//
// TODO: I need to think better about error handling in Romualdo. Especially
// those runtime errors that should be (an AFAIK are) caught in compile-time.
func (vm *VM) runtimeError(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	line := vm.currentLines()[vm.frame.ip-1]
	fmt.Fprintf(os.Stderr, "[line %d] in script\n", line)
	panic("runtimeError() called")
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

// callFrame contains the information needed at runtime about an ongoing
// function call.
type callFrame struct {
	// function is the function running.
	// TODO: Smells like this could be a Chunk. (Would be better for when
	// implementing Passages)
	function bytecode.Function

	// ip is the instruction pointer, which points to the next instruction to be
	// executed (it's an index into function's chunk).
	ip int

	// stack is the first index into the VM stack that this function can
	// use.
	stack *StackView
}
