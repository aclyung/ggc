package expression

import (
	"almeng.com/glang/general/Text"
	"almeng.com/glang/syntax"
	"almeng.com/glang/token"
)

type SyntaxToken struct {
	Syntax
	token    token.Token
	Position int
	Text     string
	Value    interface{}
	Span     Text.TextSpan
}

func NewSyntaxToken(token token.Token, position int, text string, value interface{}) *SyntaxToken {
	e := NewSyntax(token, syntax.Token)
	return &SyntaxToken{e, token, position, text, value, Text.Span(position, position+len(text))}
}
