package expression

import (
	"almeng.com/glang/syntax"
	"almeng.com/glang/token"
)

type UnaryExpressionSyntax struct {
	OperatorToken SyntaxToken
	Operand       syntax.ExpressionSyntax
}

func (unary UnaryExpressionSyntax) Kind() token.Token {
	return unary.OperatorToken.Kind()
}
func (unary UnaryExpressionSyntax) Type() syntax.Type {
	return syntax.ExpUnary
}
func (unary UnaryExpressionSyntax) GetChildren() []syntax.ExpressionSyntax {
	return []syntax.ExpressionSyntax{unary.OperatorToken, unary.Operand}
}

func NewUnaryExpressionSyntax(opToken SyntaxToken, operand syntax.ExpressionSyntax) *UnaryExpressionSyntax {
	unary := &UnaryExpressionSyntax{
		opToken,
		operand}
	return unary
}
