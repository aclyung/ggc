package expression

import (
	"almeng.com/glang/glang/syntax"
	"almeng.com/glang/glang/token"
)

type BinaryExpressionSyntax struct {
	OperatorToken SyntaxToken
	Left          syntax.ExpressionSyntax
	Right         syntax.ExpressionSyntax
}

func (b BinaryExpressionSyntax) Kind() token.Token {
	return b.OperatorToken.Kind()
}
func (b BinaryExpressionSyntax) Type() syntax.Type {
	return syntax.ExpBinary
}
func (b BinaryExpressionSyntax) GetChildren() []syntax.ExpressionSyntax {
	return []syntax.ExpressionSyntax{b.Left, b.OperatorToken, b.Right}
}

func NewBinaryExpressionSyntax(left syntax.ExpressionSyntax, opToken SyntaxToken, right syntax.ExpressionSyntax) *BinaryExpressionSyntax {
	bi := &BinaryExpressionSyntax{
		opToken,
		left,
		right}
	return bi
}
