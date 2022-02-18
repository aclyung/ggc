package expression

import (
	"almeng.com/glang/legacy/syntax"
	"almeng.com/glang/legacy/token"
)

type IllegalExpressionSyntax struct {
	Syntax
}

func NewIllegalExpressionSyntax() *IllegalExpressionSyntax {
	e := NewSyntax(token.ILLEGAL, syntax.ILLEGAL)
	return &IllegalExpressionSyntax{e}
}
