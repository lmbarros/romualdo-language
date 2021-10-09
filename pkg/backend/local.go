/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2021 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package backend

// local represents a local variable.
type local struct {
	// name is the local variable name.
	name string

	// depth is the nesting level (AKA scope depth) of the local variable.
	depth int
}
