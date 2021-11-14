package syntax

import (
	"almeng.com/glang/glang/lexer"
	node2 "almeng.com/glang/glang/parser/node"
	"almeng.com/glang/glang/token"
)

type NumberExpressionSyntax struct {
	NumberToken lexer.SyntaxToken
}

func NewNumberExpressionSyntax(numberToken lexer.SyntaxToken) *NumberExpressionSyntax {
	syntax := &NumberExpressionSyntax{numberToken}
	return syntax
}

func (numSyntax *NumberExpressionSyntax) Kind() token.Token {
	return numSyntax.NumberToken.Token
}

func (numSyntax *NumberExpressionSyntax) Type() node2.Type {
	return node2.ExpNum
}

func (numSyntax *NumberExpressionSyntax) Value() interface{} {
	return numSyntax.NumberToken.Value
}

func (numSyntax *NumberExpressionSyntax) GetChildren() []node2.ExpressionSyntax {
	return []node2.ExpressionSyntax{numSyntax.NumberToken}
}
