/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests Kind to string conversion. Looks useless, but I actually got some
// missing cases with it!
func TestKindString(t *testing.T) { // nolint:funlen
	assert.Equal(t, "", Kind(-1).String())

	assert.Equal(t, "KindLeftParen", KindLeftParen.String())
	assert.Equal(t, "KindRightParen", KindRightParen.String())
	assert.Equal(t, "KindLeftBrace", KindLeftBrace.String())
	assert.Equal(t, "KindRightBrace", KindRightBrace.String())
	assert.Equal(t, "KindLeftBracket", KindLeftBracket.String())
	assert.Equal(t, "KindRightBracket", KindRightBracket.String())
	assert.Equal(t, "KindComma", KindComma.String())
	assert.Equal(t, "KindContinue", KindContinue.String())
	assert.Equal(t, "KindDot", KindDot.String())
	assert.Equal(t, "KindMinus", KindMinus.String())
	assert.Equal(t, "KindPlus", KindPlus.String())
	assert.Equal(t, "KindSlash", KindSlash.String())
	assert.Equal(t, "KindStar", KindStar.String())
	assert.Equal(t, "KindColon", KindColon.String())
	assert.Equal(t, "KindTilde", KindTilde.String())
	assert.Equal(t, "KindAt", KindAt.String())
	assert.Equal(t, "KindHat", KindHat.String())
	assert.Equal(t, "KindEqual", KindEqual.String())
	assert.Equal(t, "KindEqualEqual", KindEqualEqual.String())
	assert.Equal(t, "KindBangEqual", KindBangEqual.String())
	assert.Equal(t, "KindGreater", KindGreater.String())
	assert.Equal(t, "KindGreaterEqual", KindGreaterEqual.String())
	assert.Equal(t, "KindLess", KindLess.String())
	assert.Equal(t, "KindLessEqual", KindLessEqual.String())
	assert.Equal(t, "KindIdentifier", KindIdentifier.String())
	assert.Equal(t, "KindStringLiteral", KindStringLiteral.String())
	assert.Equal(t, "KindNumberLiteral", KindNumberLiteral.String())
	assert.Equal(t, "KindAlias", KindAlias.String())
	assert.Equal(t, "KindAnd", KindAnd.String())
	assert.Equal(t, "KindBnum", KindBnum.String())
	assert.Equal(t, "KindBool", KindBool.String())
	assert.Equal(t, "KindBreak", KindBreak.String())
	assert.Equal(t, "KindCase", KindCase.String())
	assert.Equal(t, "KindClass", KindClass.String())
	assert.Equal(t, "KindDo", KindDo.String())
	assert.Equal(t, "KindElse", KindElse.String())
	assert.Equal(t, "KindElseif", KindElseif.String())
	assert.Equal(t, "KindEnd", KindEnd.String())
	assert.Equal(t, "KindEnum", KindEnum.String())
	assert.Equal(t, "KindFalse", KindFalse.String())
	assert.Equal(t, "KindFloat", KindFloat.String())
	assert.Equal(t, "KindFor", KindFor.String())
	assert.Equal(t, "KindFunction", KindFunction.String())
	assert.Equal(t, "KindGosub", KindGosub.String())
	assert.Equal(t, "KindGoto", KindGoto.String())
	assert.Equal(t, "KindIf", KindIf.String())
	assert.Equal(t, "KindIn", KindIn.String())
	assert.Equal(t, "KindInt", KindInt.String())
	assert.Equal(t, "KindListen", KindListen.String())
	assert.Equal(t, "KindMap", KindMap.String())
	assert.Equal(t, "KindMeta", KindMeta.String())
	assert.Equal(t, "KindNil", KindNil.String())
	assert.Equal(t, "KindNot", KindNot.String())
	assert.Equal(t, "KindOr", KindOr.String())
	assert.Equal(t, "KindPassage", KindPassage.String())
	assert.Equal(t, "KindPrint", KindPrint.String())
	assert.Equal(t, "KindReturn", KindReturn.String())
	assert.Equal(t, "KindSay", KindSay.String())
	assert.Equal(t, "KindString", KindString.String())
	assert.Equal(t, "KindStruct", KindStruct.String())
	assert.Equal(t, "KindSuper", KindSuper.String())
	assert.Equal(t, "KindSwitch", KindSwitch.String())
	assert.Equal(t, "KindThen", KindThen.String())
	assert.Equal(t, "KindTrue", KindTrue.String())
	assert.Equal(t, "KindVars", KindVars.String())
	assert.Equal(t, "KindVoid", KindVoid.String())
	assert.Equal(t, "KindWhile", KindWhile.String())
	assert.Equal(t, "KindError", KindError.String())
	assert.Equal(t, "KindEOF", KindEOF.String())
}
