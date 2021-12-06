package expression

import (
	"almeng.com/glang/syntax"
)

type BinaryExpressionSyntax struct {
	Syntax
	OperatorToken SyntaxToken
	Left          syntax.ExpressionSyntax
	Right         syntax.ExpressionSyntax
}

func NewBinaryExpressionSyntax(left syntax.ExpressionSyntax, opToken SyntaxToken, right syntax.ExpressionSyntax) *BinaryExpressionSyntax {
	e := NewSyntax(opToken.Kind(), syntax.ExpBinary, opToken, left, right)
	bi := &BinaryExpressionSyntax{
		e,
		opToken,
		left,
		right}
	return bi
}
