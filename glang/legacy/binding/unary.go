package binding

import (
	"reflect"

	boundNode2 "almeng.com/glang/legacy/binding/boundNode"
)

type BoundUnaryExpression struct {
	Oper    BoundUnaryOperator
	Operand boundNode2.BoundExpression
}

func NewBoundUnaryExpression(op BoundUnaryOperator, operand boundNode2.BoundExpression) *BoundUnaryExpression {
	return &BoundUnaryExpression{op, operand}
}

func (unary *BoundUnaryExpression) Kind() boundNode2.Kind {
	return boundNode2.Unary
}

func (unary *BoundUnaryExpression) Type() reflect.Kind {
	return unary.Oper.ResultType
}
