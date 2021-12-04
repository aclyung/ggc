package binding

import (
	"almeng.com/glang/glang/binding/boundNode"
	"almeng.com/glang/glang/expression"
	"almeng.com/glang/glang/general"
	"almeng.com/glang/glang/syntax"
)

func (b *Binder) BindUnaryExpression(exp syntax.ExpressionSyntax) boundNode.BoundExpression {
	syn := exp.(*expression.UnaryExpressionSyntax)
	operand := b.Bind(syn.Operand)
	operKind := BindUnaryOperator(exp.Kind(), operand.Type().Kind()) //BindUnaryExpressionKind(syn.OperatorToken.Kind(), operand.Type())
	if operKind == IllegalUnaryOperator {
		b.Diag.Diagnose("Unary Operator '"+syn.OperatorToken.Text+"' is not defined for type "+operand.Type().String()+"", general.ERROR)
		return operand
	}
	return NewBoundUnaryExpression(operKind, operand)
}
