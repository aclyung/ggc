package expression

import (
	"almeng.com/glang/syntax"
	"almeng.com/glang/token"
)

type ParenExpressionSyntax struct {
	Syntax
	LParen     SyntaxToken
	Expression syntax.ExpressionSyntax
	RParen     SyntaxToken
}

func NewParenExpressionSyntax(lparen SyntaxToken, exp syntax.ExpressionSyntax, rparen SyntaxToken) *ParenExpressionSyntax {
	e := NewSyntax(token.LPAREN, syntax.ExpParen, lparen, exp, rparen)
	return &ParenExpressionSyntax{e, lparen, exp, rparen}
}
