package expression

import (
	"almeng.com/glang/general/TextSpan"
	"almeng.com/glang/syntax"
	"almeng.com/glang/token"
)

type SyntaxToken struct {
	token    token.Token
	Position int
	Text     string
	Value    interface{}
	Span     TextSpan.TextSpan
}

func (tok SyntaxToken) GetChildren() []syntax.ExpressionSyntax {
	return []syntax.ExpressionSyntax{}
}

func (tok SyntaxToken) Type() syntax.Type {
	return syntax.Token
}

func (tok SyntaxToken) Kind() token.Token {
	return tok.token
}

func NewSyntaxToken(token token.Token, position int, text string, value interface{}) *SyntaxToken {
	return &SyntaxToken{token, position, text, value, TextSpan.Span(position, position+len(text))}
}
