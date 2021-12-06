package expression

import (
	"almeng.com/glang/syntax"
	"almeng.com/glang/token"
)

type LiteralExpressionSyntax struct {
	Syntax
	LiteralToken SyntaxToken
}

func NewliteralExpressionSyntax(literalToken SyntaxToken) *LiteralExpressionSyntax {
	exp := NewSyntax(literalToken.Kind(), syntax.ExpLiteral, literalToken)
	return &LiteralExpressionSyntax{exp, literalToken}
}

func (LiteralSyntax *LiteralExpressionSyntax) IsKindValid() bool {
	kind := LiteralSyntax.Kind()
	return kind == token.BOOL || kind == token.INT || kind == token.FLOAT

}

func (LiteralSyntax *LiteralExpressionSyntax) Value() interface{} {
	return LiteralSyntax.LiteralToken.Value
}
