package binding

import (
	"almeng.com/glang/binding/boundNode"
	"almeng.com/glang/expression"
	"almeng.com/glang/syntax"
)

func (b *Binder) BindUnaryExpression(exp syntax.ExpressionSyntax) boundNode.BoundExpression {
	syn := exp.(*expression.UnaryExpressionSyntax)
	operand := b.Bind(syn.Operand)
	operKind := BindUnaryOperator(exp.Kind(), operand.Type()) //BindUnaryExpressionKind(syn.OperatorToken.Kind(), operand.Type())
	if operKind == IllegalUnaryOperator {
		oper := syn.OperatorToken
		operandType := operand.Type().String()
		b.Diag.UndefinedUnaryOperator(oper.Span, oper.Text, operandType)

		return operand
	}
	return NewBoundUnaryExpression(operKind, operand)
}
