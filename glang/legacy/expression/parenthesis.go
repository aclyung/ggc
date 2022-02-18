package expression

import (
	syntax2 "almeng.com/glang/legacy/syntax"
	"almeng.com/glang/legacy/token"
)

type ParenExpressionSyntax struct {
	Syntax
	LParen     SyntaxToken
	Expression syntax2.ExpressionSyntax
	RParen     SyntaxToken
}

func NewParenExpressionSyntax(lparen SyntaxToken, exp syntax2.ExpressionSyntax, rparen SyntaxToken) *ParenExpressionSyntax {
	e := NewSyntax(token.LPAREN, syntax2.ExpParen, lparen, exp, rparen)
	return &ParenExpressionSyntax{e, lparen, exp, rparen}
}
