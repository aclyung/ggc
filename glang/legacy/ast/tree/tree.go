package tree

import (
	"almeng.com/glang/legacy/expression"
	"almeng.com/glang/legacy/general"
	"almeng.com/glang/legacy/general/Text"
	"almeng.com/glang/legacy/parser"
)

// AST structure

type Tree struct {
	Source      *Text.Source
	Diagnostics general.Diags
	Root        expression.CompilationUnit
}

// Returns new SyntaxTree
// Requires Diagnostics, root syntax, and EOF SyntaxToken

func NewSyntaxTree(source *Text.Source) Tree {
	parser := parser.Parser(source)
	root := parser.ParseCompilationUnit()
	diag := parser.Diag
	return Tree{source, diag, root}
}
