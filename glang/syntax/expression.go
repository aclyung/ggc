package syntax

import "almeng.com/glang/token"

type ExpressionSyntax interface {
	Kind() token.Token
	Type() Type
	GetChildren() []ExpressionSyntax
}
