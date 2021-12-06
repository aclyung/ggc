package expression

import (
	"almeng.com/glang/syntax"
	"almeng.com/glang/token"
)

type EOFExpressionSyntax struct{}

func NewEOFExpressionSyntax() *EOFExpressionSyntax {
	return &EOFExpressionSyntax{}
}

func (eof *EOFExpressionSyntax) IsKindValid() bool {
	kind := eof.Kind()
	return kind == token.BOOL || kind == token.INT || kind == token.FLOAT

}

func (eof *EOFExpressionSyntax) Kind() token.Token {
	return token.EOF
}

func (eof *EOFExpressionSyntax) Type() syntax.Type {
	return syntax.EOF
}

func (eof *EOFExpressionSyntax) Value() interface{} {
	return nil
}

func (eof *EOFExpressionSyntax) GetChildren() []syntax.ExpressionSyntax {
	return []syntax.ExpressionSyntax{}
}
