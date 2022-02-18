package binding

import (
	"reflect"

	"almeng.com/glang/legacy/binding/boundNode"
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
