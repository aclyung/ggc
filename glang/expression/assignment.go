package expression

import (
	"almeng.com/glang/syntax"
	"almeng.com/glang/token"
)

type AssignmentExpressionSyntax struct {
	Ident       SyntaxToken
	AssignToken SyntaxToken
	Expression  syntax.ExpressionSyntax
}

func (b AssignmentExpressionSyntax) Kind() token.Token {
	return b.AssignToken.Kind()
}
func (b AssignmentExpressionSyntax) Type() syntax.Type {
	return syntax.ExpAssign
}
func (b AssignmentExpressionSyntax) GetChildren() []syntax.ExpressionSyntax {
	return []syntax.ExpressionSyntax{b.Ident, b.Expression, b.AssignToken}
}

func NewAssigmentExpressionSyntax(ident SyntaxToken, tok SyntaxToken, exp syntax.ExpressionSyntax) *AssignmentExpressionSyntax {
	assign := &AssignmentExpressionSyntax{
		ident,
		tok,
		exp,
	}
	return assign
}
