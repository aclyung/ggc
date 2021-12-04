package ast

import (
	"almeng.com/glang/ast/tree"
	"almeng.com/glang/parser"
)

func ParseTree(text string) tree.Tree {
	pars := parser.Parser(text)
	return pars.Parse()
}
