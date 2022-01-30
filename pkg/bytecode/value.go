/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2022 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package bytecode

import (
	"fmt"
	"reflect"
)

// A ValueKind represents one of the types a value in the Romualdo Virtual
// Machine can have. This is the type from the perspective of the VM (in the
// sense that user-defined types are obviously not directly represented here).
// We use "kind" in the name because "type" is a keywork in Go.
type ValueKind int

const (
	// ValueFloat identifies a floating-point value. In this implementation,
	// they are repsented by a 64-bit IEEE 754 number, but I'd argue that if you
	// depend on the exact representation, Romualdo is not the right tool for
	// you.
	//
	// On the VM level, this is also used for Bounded Numbers. The compiler
	// knows the difference between a float and a bnum, but for the VM they are
	// all floats.
	ValueFloat ValueKind = iota

	// ValueInt identifies a signed integer value. In this implementation, they
	// are 64-bit. Other implementations may use different representations, but
	// let's all agree the number will be at least 32-bit.
	ValueInt

	// ValueBool identifies a Boolean value.
	ValueBool

	// ValueString identifies a string value.
	ValueString

	// ValueFunction identifies a function value.
	ValueFunction
)

// Function is the runtime representation of a function. We don't include any
// sort of information about return and parameter types because type-checking is
// all done statically at compile-time.
type Function struct {
	// ChunkIndex points to the Chunk that contains this function's bytecode.
	// It's an index into the global slice of Chunks.
	ChunkIndex int
}

// Value is a Romualdo language value.
type Value struct {
	Value interface{}
}

// NewValueFloat creates a new Value initialized to the floating-point number
// v.
func NewValueFloat(v float64) Value {
	return Value{
		Value: v,
	}
}

// NewValueInt creates a new Value initialized to the integer number v.
func NewValueInt(v int64) Value {
	return Value{
		Value: v,
	}
}

// NewValueBool creates a new Value initialized to the Boolean value v.
func NewValueBool(v bool) Value {
	return Value{
		Value: v,
	}
}

// NewValueString creates a new Value initialized to the string value v.
//
// This shall not be called directly by most client code. Instead, use the
// wrappers that go by the name of NewInternedValueString (as methods on objects
// around the code base). These wrappers use interning to hopefully make the
// string handling more efficient.
func NewValueString(v string) Value {
	return Value{
		Value: v,
	}
}

// NewValueFunction creates a new Value of type Function, that will run the code at the given Chunk index.
func NewValueFunction(index int) Value {
	return Value{
		Value: Function{
			ChunkIndex: index,
		},
	}
}

// AsFloat returns this Value's value, assuming it is a floating-point number.
func (v Value) AsFloat() float64 {
	return v.Value.(float64)
}

// AsInt returns this Value's value, assuming it is an integer number.
func (v Value) AsInt() int64 {
	return v.Value.(int64)
}

// AsBool returns this Value's value, assuming it is a Boolean value.
func (v Value) AsBool() bool {
	return v.Value.(bool)
}

// AsString returns this Value's value, assuming it is a string value.
func (v Value) AsString() string {
	return v.Value.(string)
}

// AsFunction returns this Value's value, assuming it is a function value.
func (v Value) AsFunction() Function {
	return v.Value.(Function)
}

// IsFloat checks if the value contains a floating-point number.
func (v Value) IsFloat() bool {
	_, ok := v.Value.(float64)
	return ok
}

// IsInt checks if the value contains an integer number.
func (v Value) IsInt() bool {
	_, ok := v.Value.(int64)
	return ok
}

// IsBool checks if the value contains a Boolean value.
func (v Value) IsBool() bool {
	_, ok := v.Value.(bool)
	return ok
}

// IsString checks if the value contains a string value.
func (v Value) IsString() bool {
	_, ok := v.Value.(string)
	return ok
}

// IsFunction checks if the value contains a function value.
func (v Value) IsFunction() bool {
	_, ok := v.Value.(Function)
	return ok
}

// String converts the value to a string.
func (v Value) String() string {
	switch vv := v.Value.(type) {
	case float64:
		return fmt.Sprintf("%g", vv)
	case int64:
		return fmt.Sprintf("%d", vv)
	case bool:
		return fmt.Sprintf("%v", vv)
	case string:
		return fmt.Sprintf("%v", vv)
	case Function:
		// TODO: Would be nice to include the function name if we had the debug
		// information around. Hard to access this info from here, though.
		return fmt.Sprintf("<function %p>", vv)
	default:
		return fmt.Sprintf("<Unexpected type %T>", vv)
	}
}

// ValuesEqual checks if a and b are considered equal.
func ValuesEqual(a, b Value) bool {
	if reflect.TypeOf(a.Value) != reflect.TypeOf(b.Value) {
		return false
	}

	switch va := a.Value.(type) {
	case bool:
		return va == b.Value.(bool)
	case float64:
		return va == b.Value.(float64)
	case int64:
		return va == b.Value.(int64)
	case string:
		return va == b.Value.(string)
	case Function:
		// TODO: Not sure if makes sense, but for now let's consider that two
		// functions are the same if they have the same bytecode.
		return va.ChunkIndex == b.Value.(Function).ChunkIndex

	default:
		panic(fmt.Sprintf("Unexpected Value type: %T", va))
	}
}
