/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package ast

import "fmt"

// A TypeTag identifies a type as seen by Romulang.
type TypeTag int

const (
	// TypeVoid identifies a void type (or rather nontype).
	TypeVoid TypeTag = iota

	// TypeInt identifies an integer number type, AKA int.
	TypeInt

	// TypeFloat identifies a floating-point number type, AKA float.
	TypeFloat

	// TypeBNum identifies a bounded number number type, AKA bnum.
	TypeBNum

	// TypeBool identifies a Boolean type, AKA bool.
	TypeBool

	// TypeString identifies a string type.
	TypeString
)

// Type describes a type. It includes a type tag and all the additional
// information needed to discern between different types that happen to share
// the same type tag (for example, two functions may have the same type tag, but
// they still might be of different types depending on their parameters and
// return types).
type Type struct {
	// Tag is the type tag. Think of it as a "high-level" type.
	Tag TypeTag
}

// String converts a Type to a string that looks like what a user would see in
// his storyworld code.
func (t Type) String() string {
	switch t.Tag {
	case TypeVoid:
		return "void"
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
	case TypeBNum:
		return "bnum"
	case TypeBool:
		return "bool"
	case TypeString:
		return "string"
	default:
		panic(fmt.Sprintf("unexpected type tag: %v", t.Tag))
	}
}

// IsNumeric checks if the type is numeric, that is, an int, float ot bnum.
func (t Type) IsNumeric() bool {
	return t.Tag == TypeInt || t.Tag == TypeFloat || t.Tag == TypeBNum
}

// IsUnboundedNumeric checks if the type is an unbounded numeric type, that is,
// either an int or a float.
func (t Type) IsUnboundedNumeric() bool {
	return t.Tag == TypeInt || t.Tag == TypeFloat || t.Tag == TypeBNum
}
