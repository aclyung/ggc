package expression

import (
	"main/legacy/syntax"
	"main/legacy/token"
)

type IllegalExpressionSyntax struct {
	Syntax
}

func NewIllegalExpressionSyntax() *IllegalExpressionSyntax {
	e := NewSyntax(token.ILLEGAL, syntax.ILLEGAL)
	return &IllegalExpressionSyntax{e}
}
