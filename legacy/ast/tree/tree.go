package tree

import (
	"main/legacy/expression"
	"main/legacy/general"
	"main/legacy/general/Text"
	"main/legacy/parser"
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
