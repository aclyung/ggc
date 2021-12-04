package syntax

type Type int

const (
	type_beg Type = iota
	ExpNum
	ExpBinary
	ExpUnary
	ExpParen
	Token
	type_end
)

var NodeType = [...]string{
	ExpNum:    "NumberExpressionNode",
	ExpBinary: "BinaryExpressionNode",
	ExpUnary:  "UnaryExpressionNode",
	ExpParen:  "ParenthesisExpressionNode",
	Token:     "SyntaxTokenNode",
}

func (t Type) String() string {
	return NodeType[t]
}
