package expression

import (
	"almeng.com/glang/syntax"
	"almeng.com/glang/token"
)

type IllegalExpressionSyntax struct {
	Syntax
}

func NewIllegalExpressionSyntax() *IllegalExpressionSyntax {
	e := NewSyntax(token.ILLEGAL, syntax.ILLEGAL)
	return &IllegalExpressionSyntax{e}
}
