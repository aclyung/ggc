package node

type Type int

const (
	type_beg Type = iota
	ExpBi
	ExpNum
	ExpParen
	SyntaxToken
	type_end

	Oper_beg
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	// AND     // &
	// OR      // |
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
	Oper_end
)

var NodeType = [...]string{
	ExpBi:       "BinaryExpressionNode",
	ExpNum:      "NumberExpressionNode",
	ExpParen:    "ParenthesisExpressionNode",
	SyntaxToken: "SyntaxTokenNode",
	ADD:         "ADD",
	SUB:         "SUB",
	MUL:         "MUL",
	QUO:         "QUO",
	REM:         "REM",

	// AND:     "&",
	// OR:      "|",
	// XOR:     "^",
	// SHL:     "<<",
	// SHR:     ">>",
	// AND_NOT: "&^",

	ADD_ASSIGN: "ADD_ASSIGN",
	SUB_ASSIGN: "SUB_ASSIGN",
	MUL_ASSIGN: "MUL_ASSIGN",
	QUO_ASSIGN: "QUO_ASSIGN",
	REM_ASSIGN: "REM_ASSIGN",

	// AND_ASSIGN:     "&=",
	// OR_ASSIGN:      "|=",
	// XOR_ASSIGN:     "^=",
	// SHL_ASSIGN:     "<<=",
	// SHR_ASSIGN:     ">>=",
	// AND_NOT_ASSIGN: "&^=",

	LAND: "LAND",
	LOR:  "LOR",
	//ARROW: "<-",
	INC: "INC",
	DEC: "DEC",

	EQL:    "EQL",
	LSS:    "LSS",
	GTR:    "GTR",
	ASSIGN: "ASSIGN",
	NOT:    "NOT",

	NEQ:    "NEQ",
	LEQ:    "LEQ",
	GEQ:    "GEQ",
	DEFINE: "DEFINE",
	//ELLIPSIS: "...",

	LPAREN: "LPAREN",
	LBRACK: "LBRACK",
	LBRACE: "LBRACE",
	COMMA:  "COMMA",
	PERIOD: "PERIOD",

	RPAREN:    "RPAREN",
	RBRACK:    "RBRACK",
	RBRACE:    "RBRACE",
	SEMICOLON: "SEMICOLON",
	COLON:     "COLON",
}

var operators map[string]Type

func (t Type) String() string {
	return NodeType[t]
}
