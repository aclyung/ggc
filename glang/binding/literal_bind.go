package binding

import (
	"almeng.com/glang/glang/binding/boundNode"
	"almeng.com/glang/glang/expression"
	"almeng.com/glang/glang/syntax"
)

func (b *Binder) BindLiteralExpression(exp syntax.ExpressionSyntax) boundNode.BoundExpression {
	lit := exp.(*expression.Literal)
	if lit.IsKindValid() {
		return NewBoundLiteralExpression(lit.Value())
	}
	return NewBoundLiteralExpression(int64(0))
}
