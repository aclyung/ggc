package binding

import (
	"almeng.com/glang/binding/boundNode"
	"reflect"
)

type BoundEOFExpression struct {
}

func (eof *BoundEOFExpression) Kind() boundNode.Kind {
	return boundNode.EOF
}

func (eof *BoundEOFExpression) Type() reflect.Kind {
	return reflect.Invalid
}

func NewBoundEOFExpression() *BoundEOFExpression {
	return &BoundEOFExpression{}
}
