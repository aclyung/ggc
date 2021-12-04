package binding

import (
	"almeng.com/glang/binding/boundNode"
	"almeng.com/glang/general"
	_ "almeng.com/glang/syntax"
	"reflect"
)

type BoundAssignmentExpression struct {
	Variable   general.VariableSymbol
	Expression boundNode.BoundExpression
}

func NewBoundAssignmentExpression(variable general.VariableSymbol, exp boundNode.BoundExpression) *BoundAssignmentExpression {
	return &BoundAssignmentExpression{variable, exp}
}

func (assign *BoundAssignmentExpression) Kind() boundNode.BoundNodeKind {
	return boundNode.Assign
}

func (assign *BoundAssignmentExpression) Type() reflect.Kind {
	return reflect.Invalid
}
