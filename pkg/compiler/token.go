/******************************************************************************\
* The Romualdo Language                                                        *
*                                                                              *
* Copyright 2020 Leandro Motta Barros                                          *
* Licensed under the MIT license (see LICENSE.txt for details)                 *
\******************************************************************************/

package compiler

// TokenKind represents the type of a token. I would call this TokenType if
// "type" wasn't a reserved word in Go. So, there we have it, "token kind".
type TokenKind int

const (
	// Single-character tokens.
	TokenLeftParen    TokenKind = iota // (
	TokenRightParen                    // )
	TokenLeftBrace                     // {
	TokenRightBrace                    // }
	TokenLeftBracket                   // [
	TokenRightBracket                  // ]
	TokenComma                         // ,
	TokenDot                           // .
	TokenMinus                         // -
	TokenPlus                          // +
	TokenSlash                         // /
	TokenStar                          // *
	TokenColon                         // :
	TokenTilde                         // ~
	TokenAt                            // @
	TokenHat                           // ^

	// One or two character tokens.
	TokenEqual        // =
	TokenEqualEqual   // ==
	TokenBangEqual    // !=
	TokenGreater      // >
	TokenGreaterEqual // >=
	TokenLess         // <
	TokenLessEqual    // <=

	// Literals.
	TokenIdentifier
	TokenStringLiteral
	TokenNumberLiteral

	// Keywords.
	TokenAlias    // alias
	TokenBnum     // bnum
	TokenBool     // bool
	TokenAnd      // and
	TokenDo       // do
	TokenElse     // else
	TokenElseif   // elseif
	TokenEnd      // end
	TokenEnum     // enum
	TokenFalse    // false
	TokenFunction // function
	TokenGosub    // gosub
	TokenGoto     // goto
	TokenIf       // if
	TokenIn       // in
	TokenInt      // int
	TokenFloat    // float
	TokenFor      // for
	TokenListen   // listen
	TokenMap      // map
	TokenMeta     // meta
	TokenNot      // not
	TokenOr       // or
	TokenString   // string
	TokenPassage  // passage
	TokenPrint    // print (temporary?!)
	TokenReturn   // return
	TokenSay      // say
	TokenStruct   // struct
	TokenThen     // then
	TokenTrue     // true
	TokenVars     // vars
	TokenVoid     // void
	TokenWhile    // while

	// Special tokens.
	TokenError
	TokenEOF // end-of-file
)

// String converts a TokenKind to its string representation. Returns an empty
// string if an invalid kind value is passed.
func (kind TokenKind) String() string {
	switch kind {
	case TokenLeftParen:
		return "TokenLeftParen"
	case TokenRightParen:
		return "TokenRightParen"
	case TokenLeftBrace:
		return "TokenLeftBrace"
	case TokenRightBrace:
		return "TokenRightBrace"
	case TokenLeftBracket:
		return "TokenLeftBracket"
	case TokenRightBracket:
		return "TokenRightBracket"
	case TokenComma:
		return "TokenComma"
	case TokenDot:
		return "TokenDot"
	case TokenMinus:
		return "TokenMinus"
	case TokenPlus:
		return "TokenPlus"
	case TokenSlash:
		return "TokenSlash"
	case TokenStar:
		return "TokenStar"
	case TokenColon:
		return "TokenColon"
	case TokenTilde:
		return "TokenTilde"
	case TokenAt:
		return "TokenAt"
	case TokenHat:
		return "TokenHat"
	case TokenEqual:
		return "TokenEqual"
	case TokenEqualEqual:
		return "TokenEqualEqual"
	case TokenBangEqual:
		return "TokenBangEqual"
	case TokenGreater:
		return "TokenGreater"
	case TokenGreaterEqual:
		return "TokenGreaterEqual"
	case TokenLess:
		return "TokenLess"
	case TokenLessEqual:
		return "TokenLessEqual"
	case TokenIdentifier:
		return "TokenIdentifier"
	case TokenStringLiteral:
		return "TokenStringLiteral"
	case TokenNumberLiteral:
		return "TokenNumberLiteral"
	case TokenAlias:
		return "TokenAlias"
	case TokenBnum:
		return "TokenBnum"
	case TokenBool:
		return "TokenBool"
	case TokenAnd:
		return "TokenAnd"
	case TokenDo:
		return "TokenDo"
	case TokenElse:
		return "TokenElse"
	case TokenElseif:
		return "TokenElseif"
	case TokenEnd:
		return "TokenEnd"
	case TokenEnum:
		return "TokenEnum"
	case TokenFalse:
		return "TokenFalse"
	case TokenFunction:
		return "TokenFunction"
	case TokenGosub:
		return "TokenGosub"
	case TokenGoto:
		return "TokenGoto"
	case TokenIf:
		return "TokenIf"
	case TokenIn:
		return "TokenIn"
	case TokenInt:
		return "TokenInt"
	case TokenFloat:
		return "TokenFloat"
	case TokenFor:
		return "TokenFor"
	case TokenListen:
		return "TokenListen"
	case TokenMap:
		return "TokenMap"
	case TokenMeta:
		return "TokenMeta"
	case TokenNot:
		return "TokenNot"
	case TokenOr:
		return "TokenOr"
	case TokenString:
		return "TokenString"
	case TokenPassage:
		return "TokenPassage"
	case TokenPrint:
		return "TokenPrint"
	case TokenReturn:
		return "TokenReturn"
	case TokenSay:
		return "TokenSay"
	case TokenStruct:
		return "TokenStruct"
	case TokenThen:
		return "TokenThen"
	case TokenTrue:
		return "TokenTrue"
	case TokenVars:
		return "TokenVars"
	case TokenVoid:
		return "TokenVoid"
	case TokenWhile:
		return "TokenWhile"
	case TokenError:
		return "TokenError"
	case TokenEOF:
		return "TokenEOF"
	}

	return ""
}

// A Token is a token. You know, one of these thingies the scanner generates and
// the compiler consumes.
type Token struct {
	// Kind is the kind of the token.
	Kind TokenKind

	// Lexeme is the text that makes up the token. It usually is just a slice of
	// the source code string. Error tokens, though, will use this to store the
	// error message as new string.
	Lexeme string

	// The line number where the token came from.
	Line int
}
