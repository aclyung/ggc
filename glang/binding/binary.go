package binding

import (
	"reflect"

	"almeng.com/glang/binding/boundNode"
)

type BoundBinaryExpression struct {
	Left  boundNode.BoundExpression
	Oper  BoundBinaryOperator
	Right boundNode.BoundExpression
}

func NewBoundBinaryExpression(left boundNode.BoundExpression, op BoundBinaryOperator, right boundNode.BoundExpression) *BoundBinaryExpression {
	return &BoundBinaryExpression{left, op, right}
}

func (unary *BoundBinaryExpression) Kind() boundNode.Kind {
	return boundNode.Binary
}

func (unary *BoundBinaryExpression) Type() reflect.Kind {
	return unary.Oper.ResultType
}
