package binding

import (
	"almeng.com/glang/binding/boundNode"
	_ "almeng.com/glang/syntax"
	"reflect"
)

type BoundAssignmentExpression struct {
	Name       string
	Expression boundNode.BoundExpression
}

func NewBoundAssignmentExpression(name string, exp boundNode.BoundExpression) *BoundAssignmentExpression {
	return &BoundAssignmentExpression{name, exp}
}

func (assign *BoundAssignmentExpression) Kind() boundNode.BoundNodeKind {
	return boundNode.Assign
}

func (assign *BoundAssignmentExpression) Type() reflect.Kind {
	return reflect.Invalid
}
