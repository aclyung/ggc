package syntax

import "almeng.com/glang/glang/token"

type ExpressionSyntax interface {
	Kind() token.Token
	Type() Type
	GetChildren() []ExpressionSyntax
}
