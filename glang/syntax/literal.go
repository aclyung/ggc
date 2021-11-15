package syntax

import (
	"almeng.com/glang/glang/lexer"
	"almeng.com/glang/glang/parser/node"
	"almeng.com/glang/glang/syntax/expression"
	"almeng.com/glang/glang/token"
)

type Literal struct {
	NumberToken lexer.SyntaxToken
}

func NewliteralExpressionSyntax(literalToken lexer.SyntaxToken) *Literal {
	syntax := &Literal{literalToken}
	return syntax
}

func (LiteralSyntax *Literal) Kind() token.Token {
	return LiteralSyntax.NumberToken.Token
}

func (LiteralSyntax *Literal) Type() expression.Type {
	return expression.ExpNum
}

func (LiteralSyntax *Literal) Value() interface{} {
	return LiteralSyntax.NumberToken.Value
}

func (LiteralSyntax *Literal) GetChildren() []node.ExpressionSyntax {
	return []node.ExpressionSyntax{LiteralSyntax.NumberToken}
}
