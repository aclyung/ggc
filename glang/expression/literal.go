package expression

import (
	"almeng.com/glang/glang/syntax"
	"almeng.com/glang/glang/token"
)

type Literal struct {
	NumberToken SyntaxToken
}

func NewliteralExpressionSyntax(literalToken SyntaxToken) *Literal {
	syntax := &Literal{literalToken}
	return syntax
}

func (LiteralSyntax *Literal) Kind() token.Token {
	return LiteralSyntax.NumberToken.Kind()
}

func (LiteralSyntax *Literal) Type() syntax.Type {
	return syntax.ExpNum
}

func (LiteralSyntax *Literal) Value() interface{} {
	return LiteralSyntax.NumberToken.Value
}

func (LiteralSyntax *Literal) GetChildren() []syntax.ExpressionSyntax {
	return []syntax.ExpressionSyntax{LiteralSyntax.NumberToken}
}
