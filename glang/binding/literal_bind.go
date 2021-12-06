package binding

import (
	"almeng.com/glang/binding/boundNode"
	"almeng.com/glang/expression"
	"almeng.com/glang/syntax"
)

func (b *Binder) BindLiteralExpression(exp syntax.ExpressionSyntax) boundNode.BoundExpression {
	lit := exp.(*expression.LiteralExpressionSyntax)
	if lit.IsKindValid() {
		return NewBoundLiteralExpression(lit.Value())
	}
	return InvalidLiteralExpression
	//return NewBoundLiteralExpression(int64(0))
}
