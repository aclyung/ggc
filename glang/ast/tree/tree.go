package tree

import (
	"almeng.com/glang/expression"
	"almeng.com/glang/general"
	"almeng.com/glang/syntax"
)

type Tree struct {
	Diagnostics general.Diags
	Root        syntax.ExpressionSyntax
	EOF         expression.SyntaxToken
}

func NewSyntaxTree(diag general.Diags, root syntax.ExpressionSyntax, eof expression.SyntaxToken) Tree {
	return Tree{diag, root, eof}
}
