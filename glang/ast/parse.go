package ast

import (
	"almeng.com/glang/ast/tree"
	"almeng.com/glang/expression"
	"almeng.com/glang/general/Text"
	"almeng.com/glang/lexer"
	"almeng.com/glang/parser"
	"almeng.com/glang/token"
)

// Parse text into AST

func ParseTree(text string) tree.Tree {
	source := Text.From(text)
	return ParseSource(source)
}

func ParseSource(source Text.Source) tree.Tree {
	pars := parser.Parser(source)
	return pars.Parse()
}

func ParseTokens(text string) []expression.SyntaxToken {
	source := Text.From(text)
	return ParseTokensSource(source)
}
func ParseTokensSource(source Text.Source) []expression.SyntaxToken {
	lex := lexer.NewLexer(source)
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
