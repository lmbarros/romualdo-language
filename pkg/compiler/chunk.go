/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package compiler

type OpCode uint8

const (
	OpReturn OpCode = iota
)

type Chunk struct {
	Code []uint8
}
