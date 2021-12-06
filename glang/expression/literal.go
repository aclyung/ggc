package expression

import (
	"almeng.com/glang/syntax"
	"almeng.com/glang/token"
)

type LiteralExpressionSyntax struct {
	LiteralToken SyntaxToken
}

func NewliteralExpressionSyntax(literalToken SyntaxToken) *LiteralExpressionSyntax {
	syntax := &LiteralExpressionSyntax{literalToken}
	return syntax
}

func (LiteralSyntax *LiteralExpressionSyntax) IsKindValid() bool {
	kind := LiteralSyntax.Kind()
	return kind == token.BOOL || kind == token.INT || kind == token.FLOAT

}

func (LiteralSyntax *LiteralExpressionSyntax) Kind() token.Token {
	return LiteralSyntax.LiteralToken.Kind()
}

func (LiteralSyntax *LiteralExpressionSyntax) Type() syntax.Type {
	return syntax.ExpLiteral
}

func (LiteralSyntax *LiteralExpressionSyntax) Value() interface{} {
	return LiteralSyntax.LiteralToken.Value
}

func (LiteralSyntax *LiteralExpressionSyntax) GetChildren() []syntax.ExpressionSyntax {
	return []syntax.ExpressionSyntax{LiteralSyntax.LiteralToken}
}
