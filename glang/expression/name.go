package expression

import (
	"almeng.com/glang/syntax"
	"almeng.com/glang/token"
)

type NameExpressionSyntax struct {
	Syntax
	Ident SyntaxToken
}

func NewNameExpressionSyntax(ident SyntaxToken) *NameExpressionSyntax {
	e := NewSyntax(token.IDENT, syntax.ExpName, ident)
	assign := &NameExpressionSyntax{
		e,
		ident,
	}
	return assign
}
