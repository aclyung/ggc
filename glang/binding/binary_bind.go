package binding

import (
	"almeng.com/glang/glang/binding/boundNode"
	"almeng.com/glang/glang/expression"
	"almeng.com/glang/glang/general"
	"almeng.com/glang/glang/syntax"
)

func (b *Binder) BindBinaryExpression(exp syntax.ExpressionSyntax) boundNode.BoundExpression {
	biExp := exp.(*expression.BinaryExpressionSyntax)
	left := b.Bind(biExp.Left)
	right := b.Bind(biExp.Right)
	operKind := BindBinaryOperator(biExp.Kind(), left.Type().Kind(), right.Type().Kind())
	if operKind == IllegalBinaryOperator {
		b.Diag.Diagnose("Binary Operator '"+biExp.OperatorToken.Text+"' is not defined for types "+left.Type().String()+" and "+right.Type().String(), general.ERROR)
		return left
	}
	return NewBoundBinaryExpression(left, operKind, right)
}
