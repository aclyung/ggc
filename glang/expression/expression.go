package expression

import (
	"almeng.com/glang/syntax"
	"almeng.com/glang/token"
)

type Syntax struct {
	kind     token.Token
	ExpType  syntax.Type
	children []syntax.ExpressionSyntax
}

func (e *Syntax) GetChildren() []syntax.ExpressionSyntax {
	return e.children
}

func (e *Syntax) Type() syntax.Type {
	return e.ExpType
}

func (e *Syntax) Kind() token.Token {
	return e.kind
}

func NewSyntax(kind token.Token, t syntax.Type, c ...syntax.ExpressionSyntax) Syntax {
	return Syntax{kind, t, c}
}
