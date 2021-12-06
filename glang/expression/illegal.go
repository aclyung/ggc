package expression

import (
	"almeng.com/glang/syntax"
	"almeng.com/glang/token"
)

type IllegalExpressionSyntax struct{}

func NewIllegalExpressionSyntax() *IllegalExpressionSyntax {
	return &IllegalExpressionSyntax{}
}

func (eof *IllegalExpressionSyntax) Kind() token.Token {
	return token.ILLEGAL
}

func (eof *IllegalExpressionSyntax) Type() syntax.Type {
	return syntax.ILLEGAL
}

func (eof *IllegalExpressionSyntax) Value() interface{} {
	return nil
}

func (eof *IllegalExpressionSyntax) GetChildren() []syntax.ExpressionSyntax {
	return []syntax.ExpressionSyntax{}
}
