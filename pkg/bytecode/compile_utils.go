/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2022 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package bytecode

import "gitlab.com/stackedboxes/romulang/pkg/ast"

// AddChunk adds a new Chunk to csw and the corresponding debug information to
// di. Also sets funcDecl.ChunkIndex. Returns the new Chunk.
func AddChunk(csw *CompiledStoryworld, di *DebugInfo, funcDecl *ast.FunctionDecl) *Chunk {
	funcDecl.ChunkIndex = len(csw.Chunks)
	newChunk := &Chunk{}
	csw.Chunks = append(csw.Chunks, newChunk)
	di.ChunksNames = append(di.ChunksNames, funcDecl.Name)
	di.ChunksLines = append(di.ChunksLines, []int{})
	return newChunk
}
