package ast

import (
	"almeng.com/glang/glang/ast/tree"
	"almeng.com/glang/glang/parser"
)

func ParseTree(text string) tree.SyntaxTree {
	pars := parser.Parser(text)
	return pars.Parse()
}
