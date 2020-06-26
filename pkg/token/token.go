/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package token

// Kind represents the type of a token. I would call this Type if "type" wasn't
// a reserved word in Go. So, there we have it, "Kind kind".
type Kind int

const (
	// Single-character tokens.
	KindLeftParen    Kind = iota // (
	KindRightParen               // )
	KindLeftBrace                // {
	KindRightBrace               // }
	KindLeftBracket              // [
	KindRightBracket             // ]
	KindComma                    // ,
	KindDot                      // .
	KindMinus                    // -
	KindPlus                     // +
	KindSlash                    // /
	KindStar                     // *
	KindColon                    // :
	KindTilde                    // ~
	KindAt                       // @
	KindHat                      // ^

	// One or two character tokens.
	KindEqual        // =
	KindEqualEqual   // ==
	KindBangEqual    // !=
	KindGreater      // >
	KindGreaterEqual // >=
	KindLess         // <
	KindLessEqual    // <=

	// Literals.
	KindIdentifier
	KindStringLiteral
	KindNumberLiteral

	// Keywords.
	KindAlias    // alias
	KindAnd      // and
	KindBnum     // bnum
	KindBool     // bool
	KindBreak    // break
	KindCase     // case
	KindClass    // class
	KindContinue // continue
	KindDo       // do
	KindElse     // else
	KindElseif   // elseif
	KindEnd      // end
	KindEnum     // enum
	KindFalse    // false
	KindFloat    // float
	KindFor      // for
	KindFunction // function
	KindGosub    // gosub
	KindGoto     // goto
	KindIf       // if
	KindIn       // in
	KindInt      // int
	KindListen   // listen
	KindMap      // map
	KindMeta     // meta
	KindNil      // nil
	KindNot      // not
	KindOr       // or
	KindPassage  // passage
	KindPrint    // print (temporary?!)
	KindReturn   // return
	KindSay      // say
	KindString   // string
	KindStruct   // struct
	KindSuper    // super
	KindSwitch   // switch
	KindThen     // then
	KindTrue     // true
	KindVars     // vars
	KindVoid     // void
	KindWhile    // while

	// Special tokens.
	KindError
	KindEOF // end-of-file

	// Not really a token.
	NumberOfKinds
)

// String converts a TokenKind to its string representation. Returns an empty
// string if an invalid kind value is passed.
func (kind Kind) String() string { // nolint:funlen,gocyclo
	switch kind {
	case KindLeftParen:
		return "KindLeftParen"
	case KindRightParen:
		return "KindRightParen"
	case KindLeftBrace:
		return "KindLeftBrace"
	case KindRightBrace:
		return "KindRightBrace"
	case KindLeftBracket:
		return "KindLeftBracket"
	case KindRightBracket:
		return "KindRightBracket"
	case KindComma:
		return "KindComma"
	case KindDot:
		return "KindDot"
	case KindMinus:
		return "KindMinus"
	case KindPlus:
		return "KindPlus"
	case KindSlash:
		return "KindSlash"
	case KindStar:
		return "KindStar"
	case KindColon:
		return "KindColon"
	case KindTilde:
		return "KindTilde"
	case KindAt:
		return "KindAt"
	case KindHat:
		return "KindHat"
	case KindEqual:
		return "KindEqual"
	case KindEqualEqual:
		return "KindEqualEqual"
	case KindBangEqual:
		return "KindBangEqual"
	case KindGreater:
		return "KindGreater"
	case KindGreaterEqual:
		return "KindGreaterEqual"
	case KindLess:
		return "KindLess"
	case KindLessEqual:
		return "KindLessEqual"
	case KindIdentifier:
		return "KindIdentifier"
	case KindStringLiteral:
		return "KindStringLiteral"
	case KindNumberLiteral:
		return "KindNumberLiteral"
	case KindAlias:
		return "KindAlias"
	case KindAnd:
		return "KindAnd"
	case KindBnum:
		return "KindBnum"
	case KindBool:
		return "KindBool"
	case KindBreak:
		return "KindBreak"
	case KindCase:
		return "KindCase"
	case KindClass:
		return "KindClass"
	case KindContinue:
		return "KindContinue"
	case KindDo:
		return "KindDo"
	case KindElse:
		return "KindElse"
	case KindElseif:
		return "KindElseif"
	case KindEnd:
		return "KindEnd"
	case KindEnum:
		return "KindEnum"
	case KindFalse:
		return "KindFalse"
	case KindFunction:
		return "KindFunction"
	case KindGosub:
		return "KindGosub"
	case KindGoto:
		return "KindGoto"
	case KindIf:
		return "KindIf"
	case KindIn:
		return "KindIn"
	case KindInt:
		return "KindInt"
	case KindFloat:
		return "KindFloat"
	case KindFor:
		return "KindFor"
	case KindListen:
		return "KindListen"
	case KindMap:
		return "KindMap"
	case KindMeta:
		return "KindMeta"
	case KindNil:
		return "KindNil"
	case KindNot:
		return "KindNot"
	case KindOr:
		return "KindOr"
	case KindString:
		return "KindString"
	case KindPassage:
		return "KindPassage"
	case KindPrint:
		return "KindPrint"
	case KindReturn:
		return "KindReturn"
	case KindSay:
		return "KindSay"
	case KindSuper:
		return "KindSuper"
	case KindStruct:
		return "KindStruct"
	case KindSwitch:
		return "KindSwitch"
	case KindThen:
		return "KindThen"
	case KindTrue:
		return "KindTrue"
	case KindVars:
		return "KindVars"
	case KindVoid:
		return "KindVoid"
	case KindWhile:
		return "KindWhile"
	case KindError:
		return "KindError"
	case KindEOF:
		return "KindEOF"
	}

	return ""
}

// A Token is a token. You know, one of these thingies the scanner generates and
// the compiler consumes.
type Token struct {
	// Kind is the kind of the token.
	Kind Kind

	// Lexeme is the text that makes up the token. It usually is just a slice of
	// the source code string. Error tokens, though, will use this to store the
	// error message as new string.
	Lexeme string

	// The line number where the token came from.
	Line int
}
