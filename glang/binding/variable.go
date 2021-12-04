package binding

import (
	"almeng.com/glang/binding/boundNode"
	"almeng.com/glang/general"
	"reflect"
)

type BoundVariableExpression struct {
	Variable general.VariableSymbol
}

func NewBoundVariableExpression(variable general.VariableSymbol) *BoundVariableExpression {
	return &BoundVariableExpression{variable}
}

func (variable *BoundVariableExpression) Kind() boundNode.BoundNodeKind {
	return boundNode.Variable
}

func (variable *BoundVariableExpression) Type() reflect.Kind {
	return variable.Variable.Type
}
