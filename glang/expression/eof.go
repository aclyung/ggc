package expression

import (
	"almeng.com/glang/syntax"
	"almeng.com/glang/token"
)

type EOFExpressionSyntax struct {
	Syntax
}

func NewEOFExpressionSyntax() *EOFExpressionSyntax {
	e := NewSyntax(token.EOF, syntax.EOF)
	return &EOFExpressionSyntax{e}
}
