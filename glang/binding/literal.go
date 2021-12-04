package binding

import (
	"reflect"

	"almeng.com/glang/binding/boundNode"
)

type BoundLiteralExpression struct {
	Value interface{}
}

var InvalidLiteraExpression = BoundLiteralExpression{Value: nil}

func (b *BoundLiteralExpression) Kind() boundNode.BoundNodeKind {
	return boundNode.Literal
}

func (b *BoundLiteralExpression) Type() reflect.Kind {
	if *b == InvalidLiteraExpression {
		return reflect.Invalid
	}
	return reflect.TypeOf(b.Value).Kind()
}

func NewBoundLiteralExpression(value interface{}) *BoundLiteralExpression {
	return &BoundLiteralExpression{value}
}
