package tree

import (
	"almeng.com/glang/glang/general"
	"almeng.com/glang/glang/lexer"
	"almeng.com/glang/glang/parser/node"
)

type SyntaxTree struct {
	Diagnostics []general.Diag
	Root        node.ExpressionSyntax
	EOF         lexer.SyntaxToken
}

func NewSyntaxTree(diag []general.Diag, root node.ExpressionSyntax, eof lexer.SyntaxToken) SyntaxTree {
	return SyntaxTree{diag, root, eof}
}
