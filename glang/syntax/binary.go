package syntax

import (
	"almeng.com/glang/glang/lexer"
	"almeng.com/glang/glang/parser/node"
	"almeng.com/glang/glang/syntax/expression"
	"almeng.com/glang/glang/token"
)

type BinaryExpressionSyntax struct {
	OperatorToken lexer.SyntaxToken
	Left          node.ExpressionSyntax
	Right         node.ExpressionSyntax
}

func (b BinaryExpressionSyntax) Kind() token.Token {
	return b.OperatorToken.Token
}
func (b BinaryExpressionSyntax) Type() expression.Type {
	return expression.ExpBi
}
func (b BinaryExpressionSyntax) GetChildren() []node.ExpressionSyntax {
	return []node.ExpressionSyntax{b.Left, b.OperatorToken, b.Right}
}

func NewBinaryExpressionSyntax(left node.ExpressionSyntax, opToken lexer.SyntaxToken, right node.ExpressionSyntax) *BinaryExpressionSyntax {
	bi := &BinaryExpressionSyntax{
		opToken,
		left,
		right}
	return bi
}
