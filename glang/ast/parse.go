package ast

import (
	"almeng.com/glang/ast/tree"
	"almeng.com/glang/expression"
	"almeng.com/glang/lexer"
	"almeng.com/glang/parser"
	"almeng.com/glang/token"
)

func ParseTree(text string) tree.Tree {
	pars := parser.Parser(text)
	return pars.Parse()
}
func ParseTokens(text string) []expression.SyntaxToken {
	lex := lexer.NewLexer(text)
	toks := make([]expression.SyntaxToken, 0)
	for {
		tok := lex.Lex()
		if tok == nil {
			continue
		}
		if tok.Kind() == token.EOF {
			break
		}
		toks = append(toks, *tok)
	}
	return toks
}
