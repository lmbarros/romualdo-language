/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2022 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package ast

const (
	// EventAfterIfCondition is emitted right after the condition of an "if"
	// statement has been visited.
	EventAfterIfCondition = iota

	// EventAfterThenBlock is emitted wight after the "then" block (that is, the
	// block executed when the "if" condition is true) has been visited.
	EventAfterThenBlock

	// EventBeforeElse is emitted right before we visit the "else" part of an
	// "if" statement. This is not emitted for "if" statements that don't have
	// an "else".
	EventBeforeElse

	// EventAfterElse is emitted right after we visit the "else" part of an "if"
	// statement. This is not emitted for "if" statements that don't have an
	// "else".
	EventAfterElse

	// EventAfterWhileCondition is emitted right after the condition of a
	// "while" statement has been visited.
	EventAfterWhileCondition

	// EventAfterLogicalBinaryOp is emitted right after we visit the left-hand
	// side of a logical binary operator (namely, "and" or "or").
	EventAfterLogicalBinaryOp
)

// A Visitor has all the methods needed to traverse a Romualdo AST.
type Visitor interface {
	// Enter is called when entering a node during the traversal.
	Enter(node Node)

	// Event is called for special "events" during the tree traversal that we
	// might need to handle specially. The event argument is one of the Event*
	// constants.
	Event(node Node, event int)

	// Leave is called when leaving a node during the traversal.
	Leave(node Node)
}
