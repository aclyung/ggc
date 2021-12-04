package binding

import (
	"almeng.com/glang/glang/binding/boundNode"
	"reflect"
)

type BoundBinaryExpression struct {
	Left  boundNode.BoundExpression
	Oper  BoundBinaryOperator
	Right boundNode.BoundExpression
}

func NewBoundBinaryExpression(left boundNode.BoundExpression, op BoundBinaryOperator, right boundNode.BoundExpression) *BoundBinaryExpression {
	return &BoundBinaryExpression{left, op, right}
}

func (unary *BoundBinaryExpression) Kind() boundNode.BoundNodeKind {
	return boundNode.Binary
}

func (unary *BoundBinaryExpression) Type() reflect.Type {
	return unary.Left.Type()
}
