package expression

type Type int

const (
	type_beg Type = iota
	ExpBi
	ExpNum
	ExpParen
	SyntaxToken
	type_end
)

var NodeType = [...]string{
	ExpBi:       "BinaryExpressionNode",
	ExpNum:      "NumberExpressionNode",
	ExpParen:    "ParenthesisExpressionNode",
	SyntaxToken: "SyntaxTokenNode",
}

var operators map[string]Type

func (t Type) String() string {
	return NodeType[t]
}
