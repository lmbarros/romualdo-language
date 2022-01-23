/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2021 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package ast

import (
	"fmt"
	"strings"
)

// A TypeTag identifies a type as seen by Romulang.
type TypeTag int

const (
	// TypeInvalid is used to represend an invalid type.
	TypeInvalid TypeTag = -1

	// TypeVoid identifies a void type (or rather nontype).
	TypeVoid = iota

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

	// TypeFunction identifies a function type. (The actual complete type of a
	// function includes its parameter types and return type.)
	TypeFunction
)

// Type describes a type. It includes a type tag and all the additional
// information needed to discern between different types that happen to share
// the same type tag (for example, two functions may have the same type tag, but
// they still might be of different types depending on their parameters and
// return types).
type Type struct {
	// Tag is the type tag. Think of it as a "high-level" type.
	Tag TypeTag

	// ParameterTypes is a slice with the types the function parameters. Valid
	// only of Tag == TypeFunction.
	ParameterTypes []*Type

	// ReturnType is the type of the function return value. Valid only if
	// Tag == TypeFunction.
	ReturnType *Type
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
	case TypeFunction:
		paramTypes := []string{}
		for _, paramType := range t.ParameterTypes {
			paramTypes = append(paramTypes, paramType.ReturnType.String())
		}
		return "function(" + strings.Join(paramTypes, ",") + "):" + t.ReturnType.String()
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
	return t.Tag == TypeInt || t.Tag == TypeFloat
}
