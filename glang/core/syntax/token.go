package syntax

// Token is the set of lexical tokens of the Go programming language.
type Token int

const (
	_    Token = iota
	_EOF       // EOF

	// names and literals
	_Name    // name
	_Literal // Literal

	// operators and operations
	// _Operator is excluding '*' (_Star)
	_Operator // op
	_AssignOp // op=
	_IncOp    // opop
	_Assign   // =
	_Define   // :=
	_Star     // *

	// delimiters
	_Lparen  // (
	_Lbrack  // [
	_Lbrace  // {
	_Rparen  // )
	_Rbrack  // ]
	_Rbrace  // }
	_Comma   // ,
	_Semi    // ;
	_Colon   // :
	_Dot     // .
	_Comment // //

	// keywords
	keyword_beg
	_Import // import
	_If     // if
	_Else   // else
	_Space  // space
	_Var    // var
	_Const  // const
	_Type   // type
	_Oper   // oper
	_Func   // func
	_Return // return
	_For    // for
	_While  // while
	_Break  // break
	Operator_beg
	OperNot       // !
	OperAdd       // add
	OperSub       // sub
	OperMul       // mul
	OperDiv       // div
	OperRem       // rem
	OperEql       // eql
	OperGtr       // gtr
	Reversed_oper // rem
	OperRAdd      // radd
	OperRSub      // rsub
	OperRMul      // rmul
	OperRDiv      // rdiv
	OperRRem      // rrem
	OperREql      // reql
	OperRGtr      // rgtr
	Operator_end
	keyword_end
)

//	// keywords
//	_Case        // case
//	_Chan        // chan
//	_Const       // const
//	_Continue    // continue
//	_Default     // default
//	_Defer       // defer
//	_Fallthrough // fallthrough
//	_Go          // go
//	_Goto        // goto
//	_Interface   // interface
//	_Map         // map
//	_Range       // range
//	_Select      // select
//	_Struct      // struct
//	_Switch      // switch
//
//	// empty line comment to exclude it from .String
//	tokenCount //
//)

func (t Token) String() string {
	return tokenString[t]
}

var tokenString = map[Token]string{
	_EOF: "EOF",

	// names and literals
	_Name:    "name",
	_Literal: "Literal",

	// operators and operations
	// _Operator is excluding '*' (_Star)
	_Operator: "op",
	_AssignOp: "op=",
	_IncOp:    "opop",
	_Assign:   "=",
	_Define:   ":=",
	_Star:     "*",

	// delimiters
	_Lparen:  "(",
	_Lbrack:  "[",
	_Lbrace:  "{",
	_Rparen:  ")",
	_Rbrack:  "]",
	_Rbrace:  "}",
	_Comma:   ",",
	_Semi:    ";",
	_Colon:   ":",
	_Dot:     ".",
	_Comment: "//",

	_Var:     "var",
	_Const:   "const",
	_Type:    "type",
	_Import:  "import",
	_If:      "if",
	_Else:    "else",
	_Space:   "space",
	_Oper:    "oper",
	_Func:    "func",
	_Return:  "return",
	_For:     "for",
	_While:   "while",
	_Break:   "break",
	OperNot:  "not",
	OperAdd:  "add",
	OperSub:  "sub",
	OperMul:  "mul",
	OperDiv:  "div",
	OperEql:  "eql",
	OperGtr:  "gtr",
	OperRem:  "rem",
	OperRAdd: "radd",
	OperRSub: "rsub",
	OperRMul: "rmul",
	OperRDiv: "rdiv",
	OperRRem: "rrem",
	OperREql: "reql",
	OperRGtr: "rgtr",
}

var keywordToken map[Token]string

func (t Token) isKeyword() bool {
	return t > keyword_beg && t < keyword_end
}

func (t Token) isOperator() bool {
	return t > Operator_beg && t < Operator_end
}

func (t Token) IsReversedOper() bool {
	return t > Reversed_oper && t < Operator_end
}

func keyword(word string) Token {
	for tok, str := range tokenString {
		if str == word {
			return tok
		}
	}
	return _Name
}

type LitKind int

const (
	IntLit LitKind = iota
	FloatLit
	RuneLit
	StringLit
)

func (t Token) OperTokenToOperator() Operator {
	switch t {
	case OperEql:
		return Eql // ==
	case OperGtr:
		return Gtr // >
	case OperAdd:
		return Add // +
	case OperSub:
		return Sub // -
	case OperMul:
		return Mul // *
	case OperDiv:
		return Div // /
	case OperRem:
		return Rem // %
	}
	return 0
}

type Operator int

const (
	_ Operator = iota

	// Def is the : in :=
	Def // :
	Not // !

	// precOrOr
	OrOr // ||

	// precAndAnd
	AndAnd // &&

	// precCmp
	Eql // ==
	Neq // !=
	Lss // <
	Leq // <=
	Gtr // >
	Geq // >=

	// precAdd
	Add // +
	Sub // -
	//Or  // |
	//Xor // ^

	// precMul
	Mul // *
	Div // /
	Rem // %
	//And    // &
	//AndNot // &^
	//Shl    // <<
	//Shr    // >>
)

var _op = [...]string{
	Def:    ":",
	Not:    "!",
	OrOr:   "||",
	AndAnd: "&&",
	Eql:    "==",
	Neq:    "!=",
	Lss:    "<",
	Leq:    "<=",
	Gtr:    ">",
	Geq:    ">=",
	Add:    "+",
	Sub:    "-",
	//Or:     "|",
	//Xor:    "^",
	//Mul:    "*",
	//Div:    "/",
	Rem: "%",
	//And:    "&",
	//AndNot: "&^",
	//Shl:    "<<",
	//Shr:    ">>",
}

func (o Operator) String() string {
	return _op[o]
}

// Operator precedences
const (
	_ = iota
	precOrOr
	precAndAnd
	precCmp
	precAdd
	precMul
)
