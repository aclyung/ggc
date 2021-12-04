package binding

import (
	"almeng.com/glang/glang/binding/boundNode"
	"reflect"
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

func (unary *BoundUnaryExpression) Type() reflect.Type {
	return unary.Operand.Type()
}
