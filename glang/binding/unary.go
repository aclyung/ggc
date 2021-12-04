package binding

import (
	"almeng.com/glang/glang/binding/boundNode"
	"reflect"
)

type BoundUnaryExpression struct {
	Oper    BoundUnaryOperKind
	Operand boundNode.BoundExpression
}

func NewBoundUnaryExpression(operKind BoundUnaryOperKind, operand boundNode.BoundExpression) *BoundUnaryExpression {
	return &BoundUnaryExpression{operKind, operand}
}

func (unary *BoundUnaryExpression) Kind() boundNode.BoundNodeKind {
	return boundNode.Unary
}

func (unary *BoundUnaryExpression) Type() reflect.Type {
	return unary.Operand.Type()
}
