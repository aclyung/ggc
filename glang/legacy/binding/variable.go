package binding

import (
	"reflect"

	"almeng.com/glang/legacy/binding/boundNode"
	"almeng.com/glang/legacy/general"
)

type BoundVariableExpression struct {
	Variable general.VariableSymbol
}

func NewBoundVariableExpression(variable general.VariableSymbol) *BoundVariableExpression {
	return &BoundVariableExpression{variable}
}

func (variable *BoundVariableExpression) Kind() boundNode.Kind {
	return boundNode.Variable
}

func (variable *BoundVariableExpression) Type() reflect.Kind {
	return variable.Variable.Type
}
