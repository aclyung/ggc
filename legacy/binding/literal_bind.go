package binding

import (
	"main/legacy/binding/boundNode"
	"main/legacy/expression"
	"main/legacy/syntax"
)

func (b *Binder) BindLiteralExpression(exp syntax.ExpressionSyntax) boundNode.BoundExpression {
	lit := exp.(*expression.LiteralExpressionSyntax)
	if lit.IsKindValid() {
		return NewBoundLiteralExpression(lit.Value())
	}
	return InvalidLiteralExpression
	//return NewBoundLiteralExpression(int64(0))
}
