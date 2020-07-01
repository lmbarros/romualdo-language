/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package bytecode

import "fmt"

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
	ValueFloat ValueKind = iota

	// ValueInt identifies a signed integer value. In this implementation, they
	// are 64-bit. Other implementations may use different representations, but
	// let's all agree the number will be at least 32-bit.
	ValueInt

	// ValueBool identifies a Boolean value.
	ValueBool
)

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

// String converts the value to a string.
func (v Value) String() string {
	switch vv := v.Value.(type) {
	case float64:
		return fmt.Sprintf("%g", vv)
	case int64:
		return fmt.Sprintf("%d", vv)
	default:
		return fmt.Sprintf("<Unexpected type %T>", vv)
	}
}
