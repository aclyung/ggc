package expression

import (
	syntax2 "main/legacy/syntax"
)

type BinaryExpressionSyntax struct {
	Syntax
	OperatorToken SyntaxToken
	Left          syntax2.ExpressionSyntax
	Right         syntax2.ExpressionSyntax
}

func NewBinaryExpressionSyntax(left syntax2.ExpressionSyntax, opToken SyntaxToken, right syntax2.ExpressionSyntax) *BinaryExpressionSyntax {
	e := NewSyntax(opToken.Kind(), syntax2.ExpBinary, opToken, left, right)
	bi := &BinaryExpressionSyntax{
		e,
		opToken,
		left,
		right}
	return bi
}
