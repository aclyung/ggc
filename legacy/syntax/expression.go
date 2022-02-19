package syntax

import (
	"main/legacy/token"
)

type ExpressionSyntax interface {
	Kind() token.Token
	Type() Type
	GetChildren() []ExpressionSyntax
}
