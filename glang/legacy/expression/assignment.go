package expression

import (
	syntax2 "almeng.com/glang/legacy/syntax"
)

type AssignmentExpressionSyntax struct {
	Syntax
	Ident       SyntaxToken
	AssignToken SyntaxToken
	Expression  syntax2.ExpressionSyntax
}

func NewAssigmentExpressionSyntax(ident SyntaxToken, tok SyntaxToken, exp syntax2.ExpressionSyntax) *AssignmentExpressionSyntax {
	e := NewSyntax(tok.Kind(), syntax2.ExpAssign, ident, tok, exp)
	assign := AssignmentExpressionSyntax{
		e,
		ident,
		tok,
		exp,
	}
	return &assign
}
