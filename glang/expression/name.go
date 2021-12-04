package expression

import (
	"almeng.com/glang/syntax"
	"almeng.com/glang/token"
)

type NameExpressionSyntax struct {
	Ident SyntaxToken
}

func (b NameExpressionSyntax) Kind() token.Token {
	return token.IDENT
}
func (b NameExpressionSyntax) Type() syntax.Type {
	return syntax.ExpName
}
func (b NameExpressionSyntax) GetChildren() []syntax.ExpressionSyntax {
	return []syntax.ExpressionSyntax{b.Ident}
}

func NewNameExpressionSyntax(ident SyntaxToken) *NameExpressionSyntax {
	assign := &NameExpressionSyntax{
		ident,
	}
	return assign
}
