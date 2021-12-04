package binding

import (
	"almeng.com/glang/binding/boundNode"
	"reflect"
)

type BoundVariableExpression struct {
	Name    string
	VarKind reflect.Kind
}

func NewBoundVariableExpression(name string, kind reflect.Kind) *BoundVariableExpression {
	return &BoundVariableExpression{name, kind}
}

func (variable *BoundVariableExpression) Kind() boundNode.BoundNodeKind {
	return boundNode.Variable
}

func (variable *BoundVariableExpression) Type() reflect.Kind {
	return variable.VarKind
}
