/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2021 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package bytecode

import (
	"errors"
	"io"
)

// CSWMagic is the "magic number" identifying a Romualdo Compiled Storyworld. It
// is comprised of the "RmldCSW" string followed by a SUB character.
var CSWMagic = []byte{0x52, 0x6D, 0x6C, 0x64, 0x43, 0x53, 0x57, 0x1A}

// CSWVersion is the current version of a Romualdo Compiled Storyworld.
const CSWVersion byte = 0

// CompiledStoryworld is a compiled, binary version of a Romualdo Language
// "program". It can be serialized and deserialized.
//
// You probably don't create one manually. Instead, either run the compiler to
// generate one from source code, or read one from a file.
//
// Overall file format:
//
// - Magic
//
// - 8-bit version (currently 0)
//
// - 32-bit size (binary data size in bytes, stored in little endian)
//
// - 32-bit CRC32 of the binary data (using the IEEE polynomial, stored in
//   little endian)
//
// - Binary data
type CompiledStoryworld struct {
	// Chunk is the Chunk of bytecode containing the compiled data.
	Chunk *Chunk
}

// ReadCompiledStoryworld deserializes a CompiledStoryworld, reading the binary
// data from r.
func ReadCompiledStoryworld(r io.Reader) (*CompiledStoryworld, error) {
	// TODO: not implemented yet. This should probably use a string interner
	// passed as parameter to allow all strings to be interned (both constants
	// and the ones created in run-time).
	return nil, errors.New("not implemented yet")
}

// WriteTo serializes a CompiledStoryworld, writing the binary data to w.
func (csw *CompiledStoryworld) WriteTo(w io.Writer) (n int64, err error) {
	// TODO: not implemented yet.
	return 0, errors.New("not implemented yet")
}
