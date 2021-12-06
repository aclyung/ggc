package expression

import (
	"almeng.com/glang/syntax"
)

type UnaryExpressionSyntax struct {
	Syntax
	OperatorToken SyntaxToken
	Operand       syntax.ExpressionSyntax
}

func NewUnaryExpressionSyntax(opToken SyntaxToken, operand syntax.ExpressionSyntax) *UnaryExpressionSyntax {
	e := NewSyntax(opToken.Kind(), syntax.ExpUnary, opToken, operand)
	unary := &UnaryExpressionSyntax{
		e,
		opToken,
		operand}
	return unary
}
