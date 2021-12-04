package token

import (
	"strconv"
	"unicode"
	"unicode/utf8"
)

// Token is the set of lexical tokens of the Go programming language.
type Token int

// The list of tokens.
const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	COMMENT

	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	literal_beg
	IDENT // main
	INT   // 12345
	FLOAT // 123.45
	//IMAG   // 123.45i
	CHAR   // 'a'
	STRING // "abc"
	BOOL   // true
	literal_end

	Operator_beg
	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND // &
	OR  // |
	// XOR     // ^
	// SHL     // <<
	// SHR     // >>
	// AND_NOT // &^

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	// AND_ASSIGN     // &=
	// OR_ASSIGN      // |=
	// XOR_ASSIGN     // ^=
	// SHL_ASSIGN     // <<=
	// SHR_ASSIGN     // >>=
	// AND_NOT_ASSIGN // &^=

	LAND // &&
	LOR  // ||
	//ARROW // <-
	INC // ++
	DEC // --

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !

	NEQ    // !=
	LEQ    // <=
	GEQ    // >=
	DEFINE // :=
	//ELLIPSIS // ...

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :
	Operator_end

	keyword_beg
	// Keywords
	BREAK //not supported
	//CASE //not supported
	//CHAN //remove

	//operation functions
	// PLS // plus
	// MIN // minus
	// TIM // times
	// DIV // div
	// MOD // mod

	CONST
	//CONTINUE //not supported

	//DEFAULT //not supported
	//DEFER
	ELSE
	//FALLTHROUGH //not supported
	FOR

	FUN
	// GO
	// GOTO
	IF
	// IMPORT

	// INTERFACE
	// MAP
	// PACKAGE
	// RANGE
	RETURN

	//SELECT
	// STRUCT
	//SWITCH
	// TYPE
	VAR
	keyword_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT: "IDENT",
	INT:   "INT",
	FLOAT: "FLOAT",
	//IMAG:   "IMAG",
	CHAR:   "CHAR",
	STRING: "STRING",
	BOOL:   "BOOL",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	AND: "&",
	OR:  "|",
	// XOR:     "^",
	// SHL:     "<<",
	// SHR:     ">>",
	// AND_NOT: "&^",

	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	QUO_ASSIGN: "/=",
	REM_ASSIGN: "%=",

	// AND_ASSIGN:     "&=",
	// OR_ASSIGN:      "|=",
	// XOR_ASSIGN:     "^=",
	// SHL_ASSIGN:     "<<=",
	// SHR_ASSIGN:     ">>=",
	// AND_NOT_ASSIGN: "&^=",

	LAND: "&&",
	LOR:  "||",
	//ARROW: "<-",
	INC: "++",
	DEC: "--",

	EQL:    "==",
	LSS:    "<",
	GTR:    ">",
	ASSIGN: "=",
	NOT:    "!",

	NEQ:    "!=",
	LEQ:    "<=",
	GEQ:    ">=",
	DEFINE: ":=",
	//ELLIPSIS: "...",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	COMMA:  ",",
	PERIOD: ".",

	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	COLON:     ":",

	BREAK: "break",
	//CASE:     "case",
	//CHAN:     "chan",
	CONST: "const",
	//CONTINUE: "continue",

	// PLS: "plus",
	// MIN: "minus",
	// TIM: "times",
	// DIV: "div",
	// MOD: "mod",

	//DEFAULT:     "default",
	//DEFER: "defer",
	ELSE: "else",
	//FALLTHROUGH: "fallthrough",
	FOR: "for",

	FUN: "fun",
	// GO:     "go",
	// GOTO:   "goto",
	IF: "if",
	// IMPORT: "import",

	// INTERFACE: "interface",
	// MAP:       "map",
	// PACKAGE:   "package",
	// RANGE:     "range",
	RETURN: "return",

	//SELECT: "select",
	// STRUCT: "struct",
	//SWITCH: "switch",
	// TYPE: "type",
	VAR: "var",
}

// String returns the string corresponding to the token tok.
// For operators, delimiters, and keywords the string is the actual
// token character sequence (e.g., for the token ADD, the string is
// "+"). For all other tokens the string corresponds to the token
// constant name (e.g. for the token IDENT, the string is "IDENT").
//

func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {

		s = "token(" + strconv.Itoa(int(tok)) + ")"

	}
	return s
}

// A set of constants for precedence-based expression parsing.
// Non-operators have lowest precedence, followed by operators
// starting with precedence 1 up to unary operators. The highest
// precedence serves as "catch-all" precedence for selector,
// indexing, and other operator and delimiter tokens.
//
const (
	LowestPrec  = 0 // non-operators
	UnaryPrec   = 6
	HighestPrec = 7
)

// Precedence returns the operator precedence of the binary
// operator op. If op is not a binary operator, the result
// is LowestPrecedence.
//
func (op Token) Precedence() int {
	switch op {
	case LOR:
		return 1
	case LAND:
		return 2
	case EQL, NEQ, LSS, LEQ, GTR, GEQ:
		return 3
	case ADD, SUB: //, OR, XOR:
		return 4
	case MUL, QUO, REM: //, SHL, SHR, AND, AND_NOT:
		return 5
	}
	return LowestPrec
}

var keywords, operators map[string]Token

func init() {
	keywords = make(map[string]Token)
	operators = make(map[string]Token)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
	for i := Operator_beg + 1; i < Operator_end; i++ {
		operators[tokens[i]] = i
	}
}

func LookOperUp(oper string) Token {
	if tok, is_oper := operators[oper]; is_oper {
		return tok
	}
	return ILLEGAL
}

// Lookup maps an identifier to its keyword token or IDENT (if not a keyword).
//
func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}

// Predicates

// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; it returns false otherwise.
//
func (tok Token) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
//
func (tok Token) IsOperator() bool { return Operator_beg < tok && tok < Operator_end }

// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
//
func (tok Token) IsKeyword() bool { return keyword_beg < tok && tok < keyword_end }

// IsExported reports whether name starts with an upper-case letter.
//

func IsExported(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(ch)
}

// IsKeyword reports whether name is a Go keyword, such as "func" or "return".
//
func IsKeyword(name string) bool {
	// TODO: opt: use a perfect hash function instead of a global map.
	_, ok := keywords[name]
	return ok
}

// func IsValid(text string) bool {
// 	_, ok := tokens[text]
// 	return ok
// }
// IsIdentifier reports whether name is a Go identifier, that is, a non-empty
// string made up of letters, digits, and underscores, where the first character
// is not a digit. Keywords are not identifiers.
//
func IsIdentifier(name string) bool {
	for i, c := range name {
		if !unicode.IsLetter(c) && c != '_' && (i == 0 || !unicode.IsDigit(c)) {
			return false
		}
	}
	return name != "" && !IsKeyword(name)
}
