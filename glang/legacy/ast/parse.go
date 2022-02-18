package ast

import (
	"almeng.com/glang/legacy/ast/tree"
	"almeng.com/glang/legacy/expression"
	"almeng.com/glang/legacy/general/Text"
	"almeng.com/glang/legacy/lexer"
	"almeng.com/glang/legacy/token"
)

// ParseCompilationUnit text into AST

func ParseTree(text string, handler *Text.Source) tree.Tree {
	handler.From(text)
	return ParseSource(handler)
}

func ParseSource(source *Text.Source) tree.Tree {
	return tree.NewSyntaxTree(source)
}

func ParseTokens(text string) []expression.SyntaxToken {
	source := &Text.Source{}
	return ParseTokensSource(source)
}
func ParseTokensSource(source *Text.Source) []expression.SyntaxToken {
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
