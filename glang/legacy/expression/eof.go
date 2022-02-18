package expression

import (
	"almeng.com/glang/legacy/syntax"
	"almeng.com/glang/legacy/token"
)

type EOFExpressionSyntax struct {
	Syntax
}

func NewEOFExpressionSyntax() *EOFExpressionSyntax {
	e := NewSyntax(token.EOF, syntax.EOF)
	return &EOFExpressionSyntax{e}
}
