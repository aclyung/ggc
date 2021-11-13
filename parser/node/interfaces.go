package node

import "almeng.com/glang/token"

type Node interface {

}

type ExpressionSyntax interface{
	Kind() token.Token
	Type() Type
	GetChildren() []ExpressionSyntax
}