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
	Token
)

var NodeType = [...]string{
	ExpLiteral: "LiteralExpressionNode",
	ExpBinary:  "BinaryExpressionNode",
	ExpUnary:   "UnaryExpressionNode",
	ExpParen:   "ParenthesisExpressionNode",
	ExpAssign:  "VariableAssignExpressionNode",
	ExpName:    "IdentifierExpressionNode",
	Token:      "SyntaxTokenNode",
}

func (t Type) String() string {
	return NodeType[t]
}
