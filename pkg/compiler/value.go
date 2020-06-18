/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package compiler

import "fmt"

// Value is a Romualdo language value.
type Value float64

// String converts the value to a string.
func (v *Value) String() string {
	return fmt.Sprintf("%g", *v)
}
