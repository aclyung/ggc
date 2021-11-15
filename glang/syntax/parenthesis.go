package syntax

import (
	"almeng.com/glang/glang/lexer"
	"almeng.com/glang/glang/parser/node"
	"almeng.com/glang/glang/syntax/expression"
	"almeng.com/glang/glang/token"
)

type ParenExpressionSyntax struct {
	LParen     lexer.SyntaxToken
	Expression node.ExpressionSyntax
	RParen     lexer.SyntaxToken
}

func NewParenExpressionSyntax(lparen lexer.SyntaxToken, exp node.ExpressionSyntax, rparen lexer.SyntaxToken) *ParenExpressionSyntax {
	return &ParenExpressionSyntax{lparen, exp, rparen}
}

func (p ParenExpressionSyntax) Kind() token.Token {
	return token.LPAREN
}

func (p ParenExpressionSyntax) Type() expression.Type {
	return expression.ExpParen
}

func (p ParenExpressionSyntax) GetChildren() []node.ExpressionSyntax {
	return []node.ExpressionSyntax{p.LParen, p.Expression, p.RParen}
}
