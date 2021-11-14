package syntax

import (
	"almeng.com/glang/glang/lexer"
	node2 "almeng.com/glang/glang/parser/node"
	"almeng.com/glang/glang/token"
)

type BinaryExpressionSyntax struct {
	OperatorToken lexer.SyntaxToken
	Left          node2.ExpressionSyntax
	Right         node2.ExpressionSyntax
}

func (b BinaryExpressionSyntax) Kind() token.Token {
	return b.OperatorToken.Token
}
func (b BinaryExpressionSyntax) Type() node2.Type {
	return node2.ExpBi
}
func (b BinaryExpressionSyntax) GetChildren() []node2.ExpressionSyntax {
	return []node2.ExpressionSyntax{b.Left, b.OperatorToken, b.Right}
}

func NewBinaryExpressionSyntax(left node2.ExpressionSyntax, opToken lexer.SyntaxToken, right node2.ExpressionSyntax) *BinaryExpressionSyntax {
	bi := &BinaryExpressionSyntax{
		opToken,
		left,
		right}
	return bi
}
