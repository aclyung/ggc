package node

type Type int

const (
	type_beg Type = iota
	ExpBi
	ExpNum
	SyntaxToken
	type_end
)

var NodeType = [...]string {
	ExpBi: "BinaryExpressionNode",
	ExpNum: "NumberExpressionNode",
	SyntaxToken:"SyntaxTokenNode",
}

func (t Type) String() string {
	return NodeType[t]
}




