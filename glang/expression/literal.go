package expression

import (
	"almeng.com/glang/syntax"
	"almeng.com/glang/token"
)

type Literal struct {
	LiteralToken SyntaxToken
}

func NewliteralExpressionSyntax(literalToken SyntaxToken) *Literal {
	syntax := &Literal{literalToken}
	return syntax
}

func (LiteralSyntax *Literal) IsKindValid() bool {
	kind := LiteralSyntax.Kind()
	return kind == token.BOOL || kind == token.INT || kind == token.FLOAT

}

func (LiteralSyntax *Literal) Kind() token.Token {
	return LiteralSyntax.LiteralToken.Kind()
}

func (LiteralSyntax *Literal) Type() syntax.Type {
	return syntax.ExpLiteral
}

func (LiteralSyntax *Literal) Value() interface{} {
	return LiteralSyntax.LiteralToken.Value
}

func (LiteralSyntax *Literal) GetChildren() []syntax.ExpressionSyntax {
	return []syntax.ExpressionSyntax{LiteralSyntax.LiteralToken}
}
