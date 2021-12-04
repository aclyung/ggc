package syntax

type Type int

const (
	type_beg Type = iota
	ExpLiteral
	ExpBinary
	ExpUnary
	ExpParen
	ExpAssign
	ExpName
	Token
	type_end
)

var NodeType = [...]string{
	ExpLiteral: "NumberExpressionNode",
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
