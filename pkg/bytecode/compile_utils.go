/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2022 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package bytecode

// AddChunk adds a new Chunk to csw and the corresponding debug information to
// di. Returns the new Chunk.
func AddChunk(csw *CompiledStoryworld, di *DebugInfo, name string) *Chunk {
	newChunk := NewChunk()
	csw.Chunks = append(csw.Chunks, newChunk)
	di.ChunksNames = append(di.ChunksNames, name)
	di.ChunksLines = append(di.ChunksLines, []int{})
	return newChunk
}
