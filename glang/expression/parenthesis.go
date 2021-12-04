package expression

import (
	"almeng.com/glang/syntax"
	"almeng.com/glang/token"
)

type ParenExpressionSyntax struct {
	LParen     SyntaxToken
	Expression syntax.ExpressionSyntax
	RParen     SyntaxToken
}

func NewParenExpressionSyntax(lparen SyntaxToken, exp syntax.ExpressionSyntax, rparen SyntaxToken) *ParenExpressionSyntax {
	return &ParenExpressionSyntax{lparen, exp, rparen}
}

func (p ParenExpressionSyntax) Kind() token.Token {
	return token.LPAREN
}

func (p ParenExpressionSyntax) Type() syntax.Type {
	return syntax.ExpParen
}

func (p ParenExpressionSyntax) GetChildren() []syntax.ExpressionSyntax {
	return []syntax.ExpressionSyntax{p.LParen, p.Expression, p.RParen}
}
