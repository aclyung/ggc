package binding

import (
	"reflect"

	boundNode2 "almeng.com/glang/legacy/binding/boundNode"
	"almeng.com/glang/legacy/general"
)

type BoundAssignmentExpression struct {
	Variable   general.VariableSymbol
	Expression boundNode2.BoundExpression
}

func NewBoundAssignmentExpression(variable general.VariableSymbol, exp boundNode2.BoundExpression) *BoundAssignmentExpression {
	return &BoundAssignmentExpression{variable, exp}
}

func (assign *BoundAssignmentExpression) Kind() boundNode2.Kind {
	return boundNode2.Assign
}

func (assign *BoundAssignmentExpression) Type() reflect.Kind {
	return reflect.Invalid
}
