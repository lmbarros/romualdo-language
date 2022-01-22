/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2022 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package bytecode

import (
	"errors"
	"io"
)

// DebugInfo contains debug information matching a CompiledStoryworld. All
// information that is not strictly necessary to run a Storyworld but is useful
// for debugging, producing better error reporting, etc, belongs here. This can
// be serialized and deserialized.
type DebugInfo struct {
	// ChunksNames contains the names of the functions on a CompiledStoryworld.
	// There is one entry for each entry in the corresponding
	// CompiledStoryworld.Chunks.
	ChunksNames []string

	// The source code line that generated each instruction of each Chunk. This
	// must be interpreted like this: ChunksLines[chunkIndex][codeIndex]
	// contains the source code line that generated the bytecode at
	// CompiledStoryworld.Chunks[chunkIndex].Code[codeIndex]. Notice that we
	// have one entry for each entry in Code. Very space-inefficient, but very
	// simple.
	ChunksLines [][]int
}

// ReadDebugInfo deserializes a DebugInformation, reading the binary data from
// r.
func ReadDebugInfo(r io.Reader) (*DebugInfo, error) {
	return nil, errors.New("not implemented yet")
}

// WriteTo serializes a DebugInfo, writing the binary data to w.
func (di *DebugInfo) WriteTo(w io.Writer) (n int64, err error) {
	// TODO: not implemented yet.
	return 0, errors.New("not implemented yet")
}
