package node

import (
	"almeng.com/glang/glang/syntax/expression"
	"almeng.com/glang/glang/token"
)

type Node interface {
}

type ExpressionSyntax interface {
	Kind() token.Token
	Type() expression.Type
	GetChildren() []ExpressionSyntax
}
