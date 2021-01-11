/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package frontend

// tokenKind represents the type of a token. I would call this tokenType if
// "type" wasn't a reserved word in Go. So, there we have it, "tokenKind".
type tokenKind int

const (
	// Single-character tokens.
	tokenKindLeftParen    tokenKind = iota // (
	tokenKindRightParen                    // )
	tokenKindLeftBrace                     // {
	tokenKindRightBrace                    // }
	tokenKindLeftBracket                   // [
	tokenKindRightBracket                  // ]
	tokenKindComma                         // ,
	tokenKindDot                           // .
	tokenKindMinus                         // -
	tokenKindPlus                          // +
	tokenKindSlash                         // /
	tokenKindStar                          // *
	tokenKindColon                         // :
	tokenKindTilde                         // ~
	tokenKindAt                            // @
	tokenKindHat                           // ^

	// One or two character tokens.
	tokenKindEqual        // =
	tokenKindEqualEqual   // ==
	tokenKindBangEqual    // !=
	tokenKindGreater      // >
	tokenKindGreaterEqual // >=
	tokenKindLess         // <
	tokenKindLessEqual    // <=

	// Literals.
	tokenKindIdentifier
	tokenKindStringLiteral
	tokenKindIntLiteral
	tokenKindFloatLiteral
	tokenKindBNumLiteral

	// Keywords.
	tokenKindAlias    // alias
	tokenKindAnd      // and
	tokenKindBNum     // bnum
	tokenKindBool     // bool
	tokenKindBreak    // break
	tokenKindCase     // case
	tokenKindClass    // class
	tokenKindContinue // continue
	tokenKindDo       // do
	tokenKindElse     // else
	tokenKindElseif   // elseif
	tokenKindEnd      // end
	tokenKindEnum     // enum
	tokenKindFalse    // false
	tokenKindFloat    // float
	tokenKindFor      // for
	tokenKindFunction // function
	tokenKindGosub    // gosub
	tokenKindGoto     // goto
	tokenKindIf       // if
	tokenKindIn       // in
	tokenKindInt      // int
	tokenKindListen   // listen
	tokenKindMap      // map
	tokenKindMeta     // meta
	tokenKindNil      // nil
	tokenKindNot      // not
	tokenKindOr       // or
	tokenKindPassage  // passage
	tokenKindReturn   // return
	tokenKindSay      // say
	tokenKindString   // string
	tokenKindStruct   // struct
	tokenKindSuper    // super
	tokenKindSwitch   // switch
	tokenKindThen     // then
	tokenKindTrue     // true
	tokenKindVars     // vars
	tokenKindVoid     // void
	tokenKindWhile    // while

	// Special tokens.
	tokenKindError
	tokenKindEOF // end-of-file

	// Not really a token.
	numberOfTokenKinds
)

// String converts a tokenKind to its string representation. Returns an empty
// string if an invalid kind value is passed.
func (kind tokenKind) String() string { // nolint:funlen,gocyclo
	switch kind {
	case tokenKindLeftParen:
		return "tokenKindLeftParen"
	case tokenKindRightParen:
		return "tokenKindRightParen"
	case tokenKindLeftBrace:
		return "tokenKindLeftBrace"
	case tokenKindRightBrace:
		return "tokenKindRightBrace"
	case tokenKindLeftBracket:
		return "tokenKindLeftBracket"
	case tokenKindRightBracket:
		return "tokenKindRightBracket"
	case tokenKindComma:
		return "tokenKindComma"
	case tokenKindDot:
		return "tokenKindDot"
	case tokenKindMinus:
		return "tokenKindMinus"
	case tokenKindPlus:
		return "tokenKindPlus"
	case tokenKindSlash:
		return "tokenKindSlash"
	case tokenKindStar:
		return "tokenKindStar"
	case tokenKindColon:
		return "tokenKindColon"
	case tokenKindTilde:
		return "tokenKindTilde"
	case tokenKindAt:
		return "tokenKindAt"
	case tokenKindHat:
		return "tokenKindHat"
	case tokenKindEqual:
		return "tokenKindEqual"
	case tokenKindEqualEqual:
		return "tokenKindEqualEqual"
	case tokenKindBangEqual:
		return "tokenKindBangEqual"
	case tokenKindGreater:
		return "tokenKindGreater"
	case tokenKindGreaterEqual:
		return "tokenKindGreaterEqual"
	case tokenKindLess:
		return "tokenKindLess"
	case tokenKindLessEqual:
		return "tokenKindLessEqual"
	case tokenKindIdentifier:
		return "tokenKindIdentifier"
	case tokenKindStringLiteral:
		return "tokenKindStringLiteral"
	case tokenKindIntLiteral:
		return "tokenKindIntLiteral"
	case tokenKindFloatLiteral:
		return "tokenKindFloatLiteral"
	case tokenKindBNumLiteral:
		return "tokenKindBNumLiteral"
	case tokenKindAlias:
		return "tokenKindAlias"
	case tokenKindAnd:
		return "tokenKindAnd"
	case tokenKindBNum:
		return "tokenKindBNum"
	case tokenKindBool:
		return "tokenKindBool"
	case tokenKindBreak:
		return "tokenKindBreak"
	case tokenKindCase:
		return "tokenKindCase"
	case tokenKindClass:
		return "tokenKindClass"
	case tokenKindContinue:
		return "tokenKindContinue"
	case tokenKindDo:
		return "tokenKindDo"
	case tokenKindElse:
		return "tokenKindElse"
	case tokenKindElseif:
		return "tokenKindElseif"
	case tokenKindEnd:
		return "tokenKindEnd"
	case tokenKindEnum:
		return "tokenKindEnum"
	case tokenKindFalse:
		return "tokenKindFalse"
	case tokenKindFunction:
		return "tokenKindFunction"
	case tokenKindGosub:
		return "tokenKindGosub"
	case tokenKindGoto:
		return "tokenKindGoto"
	case tokenKindIf:
		return "tokenKindIf"
	case tokenKindIn:
		return "tokenKindIn"
	case tokenKindInt:
		return "tokenKindInt"
	case tokenKindFloat:
		return "tokenKindFloat"
	case tokenKindFor:
		return "tokenKindFor"
	case tokenKindListen:
		return "tokenKindListen"
	case tokenKindMap:
		return "tokenKindMap"
	case tokenKindMeta:
		return "tokenKindMeta"
	case tokenKindNil:
		return "tokenKindNil"
	case tokenKindNot:
		return "tokenKindNot"
	case tokenKindOr:
		return "tokenKindOr"
	case tokenKindString:
		return "tokenKindString"
	case tokenKindPassage:
		return "tokenKindPassage"
	case tokenKindReturn:
		return "tokenKindReturn"
	case tokenKindSay:
		return "tokenKindSay"
	case tokenKindSuper:
		return "tokenKindSuper"
	case tokenKindStruct:
		return "tokenKindStruct"
	case tokenKindSwitch:
		return "tokenKindSwitch"
	case tokenKindThen:
		return "tokenKindThen"
	case tokenKindTrue:
		return "tokenKindTrue"
	case tokenKindVars:
		return "tokenKindVars"
	case tokenKindVoid:
		return "tokenKindVoid"
	case tokenKindWhile:
		return "tokenKindWhile"
	case tokenKindError:
		return "tokenKindError"
	case tokenKindEOF:
		return "tokenKindEOF"
	}

	return ""
}

// A token is a token -- you know, one of these thingies the scanner generates
// and the compiler consumes.
type token struct {
	// kind is the kind of the token.
	kind tokenKind

	// Lexeme is the text that makes up the token. It usually is just a slice of
	// the source code string. Error tokens, though, will use this to store the
	// error message as new string.
	lexeme string

	// The line number where the token came from.
	line int
}
