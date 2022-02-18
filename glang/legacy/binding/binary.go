package binding

import (
	"reflect"

	boundNode2 "almeng.com/glang/legacy/binding/boundNode"
)

type BoundBinaryExpression struct {
	Left  boundNode2.BoundExpression
	Oper  BoundBinaryOperator
	Right boundNode2.BoundExpression
}

func NewBoundBinaryExpression(left boundNode2.BoundExpression, op BoundBinaryOperator, right boundNode2.BoundExpression) *BoundBinaryExpression {
	return &BoundBinaryExpression{left, op, right}
}

func (unary *BoundBinaryExpression) Kind() boundNode2.Kind {
	return boundNode2.Binary
}

func (unary *BoundBinaryExpression) Type() reflect.Kind {
	return unary.Oper.ResultType
}
