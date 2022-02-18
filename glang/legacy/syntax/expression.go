package syntax

import (
	"almeng.com/glang/legacy/token"
)

type ExpressionSyntax interface {
	Kind() token.Token
	Type() Type
	GetChildren() []ExpressionSyntax
}
