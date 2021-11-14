package syntax

import (
	"almeng.com/glang/glang/lexer"
	node2 "almeng.com/glang/glang/parser/node"
	"almeng.com/glang/glang/token"
)

type ParenExpressionSyntax struct {
	LParen     lexer.SyntaxToken
	Expression node2.ExpressionSyntax
	RParen     lexer.SyntaxToken
}

func NewParenExpressionSyntax(lparen lexer.SyntaxToken, exp node2.ExpressionSyntax, rparen lexer.SyntaxToken) *ParenExpressionSyntax {
	return &ParenExpressionSyntax{lparen, exp, rparen}
}

func (p ParenExpressionSyntax) Kind() token.Token {
	return token.LPAREN
}

func (p ParenExpressionSyntax) Type() node2.Type {
	return node2.ExpParen
}

func (p ParenExpressionSyntax) GetChildren() []node2.ExpressionSyntax {
	return []node2.ExpressionSyntax{p.LParen, p.Expression, p.RParen}
}
