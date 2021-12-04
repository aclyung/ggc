package binding

import (
	"reflect"

	"almeng.com/glang/binding/boundNode"
)

type BoundUnaryExpression struct {
	Oper    BoundUnaryOperator
	Operand boundNode.BoundExpression
}

func NewBoundUnaryExpression(op BoundUnaryOperator, operand boundNode.BoundExpression) *BoundUnaryExpression {
	return &BoundUnaryExpression{op, operand}
}

func (unary *BoundUnaryExpression) Kind() boundNode.BoundNodeKind {
	return boundNode.Unary
}

func (unary *BoundUnaryExpression) Type() reflect.Kind {
	return unary.Oper.ResultType
}
