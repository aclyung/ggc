package binding

import (
	"reflect"

	"main/legacy/binding/boundNode"
	"main/legacy/general"
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
