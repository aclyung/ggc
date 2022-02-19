package syntax

type Type int

const (
	ILLEGAL Type = iota
	EOF
	ExpLiteral
	ExpBinary
	ExpUnary
	ExpParen
	ExpAssign
	ExpName

	CompilationUnit

	Token
)

var NodeType = [...]string{
	ILLEGAL:         "Illegal",
	EOF:             "EOF",
	ExpLiteral:      "LiteralExpressionNode",
	ExpBinary:       "BinaryExpressionNode",
	ExpUnary:        "UnaryExpressionNode",
	ExpParen:        "ParenthesisExpressionNode",
	ExpAssign:       "VariableAssignExpressionNode",
	ExpName:         "IdentifierExpressionNode",
	CompilationUnit: "CompilationUnit",
	Token:           "SyntaxTokenNode",
}

func (t Type) String() string {
	return NodeType[t]
}
