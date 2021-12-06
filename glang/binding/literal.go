package binding

import (
	"reflect"

	"almeng.com/glang/binding/boundNode"
)

type BoundLiteralExpression struct {
	Value interface{}
}

var InvalidLiteralExpression = &BoundLiteralExpression{Value: nil}

func (b *BoundLiteralExpression) Kind() boundNode.Kind {
	return boundNode.Literal
}

func (b *BoundLiteralExpression) Type() reflect.Kind {
	if b == InvalidLiteralExpression {
		return reflect.Invalid
	}
	return reflect.TypeOf(b.Value).Kind()
}

func NewBoundLiteralExpression(value interface{}) *BoundLiteralExpression {
	return &BoundLiteralExpression{value}
}
