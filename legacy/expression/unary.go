package expression

import (
	syntax2 "main/legacy/syntax"
)

type UnaryExpressionSyntax struct {
	Syntax
	OperatorToken SyntaxToken
	Operand       syntax2.ExpressionSyntax
}

func NewUnaryExpressionSyntax(opToken SyntaxToken, operand syntax2.ExpressionSyntax) *UnaryExpressionSyntax {
	e := NewSyntax(opToken.Kind(), syntax2.ExpUnary, opToken, operand)
	unary := &UnaryExpressionSyntax{
		e,
		opToken,
		operand}
	return unary
}
