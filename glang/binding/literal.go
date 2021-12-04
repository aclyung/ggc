package binding

import (
	"almeng.com/glang/glang/binding/boundNode"
	"reflect"
)

type BoundLiteralExpression struct {
	Value interface{}
}

func (b *BoundLiteralExpression) Kind() boundNode.BoundNodeKind {
	return boundNode.Literal
}

func (b *BoundLiteralExpression) Type() reflect.Type {
	return reflect.TypeOf(b.Value)
}

func NewBoundLiteralExpression(value interface{}) *BoundLiteralExpression {
	return &BoundLiteralExpression{value}
}
