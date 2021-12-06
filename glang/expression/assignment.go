package expression

import (
	"almeng.com/glang/syntax"
)

type AssignmentExpressionSyntax struct {
	Syntax
	Ident       SyntaxToken
	AssignToken SyntaxToken
	Expression  syntax.ExpressionSyntax
}

func NewAssigmentExpressionSyntax(ident SyntaxToken, tok SyntaxToken, exp syntax.ExpressionSyntax) *AssignmentExpressionSyntax {
	e := NewSyntax(tok.Kind(), syntax.ExpAssign, ident, tok, exp)
	assign := &AssignmentExpressionSyntax{
		e,
		ident,
		tok,
		exp,
	}
	return assign
}
