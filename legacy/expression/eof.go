package expression

import (
	"main/legacy/syntax"
	"main/legacy/token"
)

type EOFExpressionSyntax struct {
	Syntax
}

func NewEOFExpressionSyntax() *EOFExpressionSyntax {
	e := NewSyntax(token.EOF, syntax.EOF)
	return &EOFExpressionSyntax{e}
}
