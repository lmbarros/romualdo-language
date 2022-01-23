/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020-2022 Leandro Motta Barros                                     *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package frontend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests Kind to string conversion. Looks useless, but I actually got some
// missing cases with it!
func TestKindString(t *testing.T) { // nolint:funlen
	assert.Equal(t, "", tokenKind(-1).String())

	assert.Equal(t, "tokenKindLeftParen", tokenKindLeftParen.String())
	assert.Equal(t, "tokenKindRightParen", tokenKindRightParen.String())
	assert.Equal(t, "tokenKindLeftBrace", tokenKindLeftBrace.String())
	assert.Equal(t, "tokenKindRightBrace", tokenKindRightBrace.String())
	assert.Equal(t, "tokenKindLeftBracket", tokenKindLeftBracket.String())
	assert.Equal(t, "tokenKindRightBracket", tokenKindRightBracket.String())
	assert.Equal(t, "tokenKindComma", tokenKindComma.String())
	assert.Equal(t, "tokenKindContinue", tokenKindContinue.String())
	assert.Equal(t, "tokenKindDot", tokenKindDot.String())
	assert.Equal(t, "tokenKindMinus", tokenKindMinus.String())
	assert.Equal(t, "tokenKindPlus", tokenKindPlus.String())
	assert.Equal(t, "tokenKindSlash", tokenKindSlash.String())
	assert.Equal(t, "tokenKindStar", tokenKindStar.String())
	assert.Equal(t, "tokenKindColon", tokenKindColon.String())
	assert.Equal(t, "tokenKindTilde", tokenKindTilde.String())
	assert.Equal(t, "tokenKindAt", tokenKindAt.String())
	assert.Equal(t, "tokenKindHat", tokenKindHat.String())
	assert.Equal(t, "tokenKindEqual", tokenKindEqual.String())
	assert.Equal(t, "tokenKindEqualEqual", tokenKindEqualEqual.String())
	assert.Equal(t, "tokenKindBangEqual", tokenKindBangEqual.String())
	assert.Equal(t, "tokenKindGreater", tokenKindGreater.String())
	assert.Equal(t, "tokenKindGreaterEqual", tokenKindGreaterEqual.String())
	assert.Equal(t, "tokenKindLess", tokenKindLess.String())
	assert.Equal(t, "tokenKindLessEqual", tokenKindLessEqual.String())
	assert.Equal(t, "tokenKindIdentifier", tokenKindIdentifier.String())
	assert.Equal(t, "tokenKindStringLiteral", tokenKindStringLiteral.String())
	assert.Equal(t, "tokenKindIntLiteral", tokenKindIntLiteral.String())
	assert.Equal(t, "tokenKindFloatLiteral", tokenKindFloatLiteral.String())
	assert.Equal(t, "tokenKindBNumLiteral", tokenKindBNumLiteral.String())
	assert.Equal(t, "tokenKindAlias", tokenKindAlias.String())
	assert.Equal(t, "tokenKindAnd", tokenKindAnd.String())
	assert.Equal(t, "tokenKindBNum", tokenKindBNum.String())
	assert.Equal(t, "tokenKindBool", tokenKindBool.String())
	assert.Equal(t, "tokenKindBreak", tokenKindBreak.String())
	assert.Equal(t, "tokenKindCase", tokenKindCase.String())
	assert.Equal(t, "tokenKindClass", tokenKindClass.String())
	assert.Equal(t, "tokenKindDo", tokenKindDo.String())
	assert.Equal(t, "tokenKindElse", tokenKindElse.String())
	assert.Equal(t, "tokenKindElseif", tokenKindElseif.String())
	assert.Equal(t, "tokenKindEnd", tokenKindEnd.String())
	assert.Equal(t, "tokenKindEnum", tokenKindEnum.String())
	assert.Equal(t, "tokenKindFalse", tokenKindFalse.String())
	assert.Equal(t, "tokenKindFloat", tokenKindFloat.String())
	assert.Equal(t, "tokenKindFor", tokenKindFor.String())
	assert.Equal(t, "tokenKindFunction", tokenKindFunction.String())
	assert.Equal(t, "tokenKindGlobals", tokenKindGlobals.String())
	assert.Equal(t, "tokenKindGosub", tokenKindGosub.String())
	assert.Equal(t, "tokenKindGoto", tokenKindGoto.String())
	assert.Equal(t, "tokenKindIf", tokenKindIf.String())
	assert.Equal(t, "tokenKindIn", tokenKindIn.String())
	assert.Equal(t, "tokenKindInt", tokenKindInt.String())
	assert.Equal(t, "tokenKindListen", tokenKindListen.String())
	assert.Equal(t, "tokenKindMap", tokenKindMap.String())
	assert.Equal(t, "tokenKindMeta", tokenKindMeta.String())
	assert.Equal(t, "tokenKindNil", tokenKindNil.String())
	assert.Equal(t, "tokenKindNot", tokenKindNot.String())
	assert.Equal(t, "tokenKindOr", tokenKindOr.String())
	assert.Equal(t, "tokenKindPassage", tokenKindPassage.String())
	assert.Equal(t, "tokenKindReturn", tokenKindReturn.String())
	assert.Equal(t, "tokenKindSay", tokenKindSay.String())
	assert.Equal(t, "tokenKindString", tokenKindString.String())
	assert.Equal(t, "tokenKindStruct", tokenKindStruct.String())
	assert.Equal(t, "tokenKindSuper", tokenKindSuper.String())
	assert.Equal(t, "tokenKindSwitch", tokenKindSwitch.String())
	assert.Equal(t, "tokenKindThen", tokenKindThen.String())
	assert.Equal(t, "tokenKindTrue", tokenKindTrue.String())
	assert.Equal(t, "tokenKindVar", tokenKindVar.String())
	assert.Equal(t, "tokenKindVoid", tokenKindVoid.String())
	assert.Equal(t, "tokenKindWhile", tokenKindWhile.String())
	assert.Equal(t, "tokenKindError", tokenKindError.String())
	assert.Equal(t, "tokenKindEOF", tokenKindEOF.String())
}
