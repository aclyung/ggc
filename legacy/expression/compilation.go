package expression

import (
	syntax2 "main/legacy/syntax"
	"main/legacy/token"
)

type CompilationUnit struct {
	Syntax
	Expression syntax2.ExpressionSyntax
	EOF        SyntaxToken
}

func NewCompilationUnit(exp syntax2.ExpressionSyntax, eof SyntaxToken) CompilationUnit {
	s := NewSyntax(token.ILLEGAL, syntax2.CompilationUnit, exp, eof)
	return CompilationUnit{s, exp, eof}
}
