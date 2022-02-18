package binding

import (
	"almeng.com/glang/legacy/binding/boundNode"
	"almeng.com/glang/legacy/expression"
	"almeng.com/glang/legacy/syntax"
)

func (b *Binder) BindLiteralExpression(exp syntax.ExpressionSyntax) boundNode.BoundExpression {
	lit := exp.(*expression.LiteralExpressionSyntax)
	if lit.IsKindValid() {
		return NewBoundLiteralExpression(lit.Value())
	}
	return InvalidLiteralExpression
	//return NewBoundLiteralExpression(int64(0))
}
