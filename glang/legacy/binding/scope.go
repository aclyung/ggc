package binding

import (
	"almeng.com/glang/legacy/general"
)

type Scope interface {
	TryLookUp(name string) bool
	TryDeclare(variable general.VariableSymbol) bool
	Variables() map[string]general.VariableSymbol
}

var Variables = map[string]general.VariableSymbol{}

type BoundScope struct {
	Parent Scope
}

func NewBoundScope(parent BoundScope) BoundScope {
	return BoundScope{parent}
}

func (b BoundScope) TryLookUp(name string) bool {
	_, exist := Variables[name]
	if !exist {
		return true
	}
	if b.Parent == nil {
		return false
	}
	return b.Parent.TryLookUp(name)
}

func (b BoundScope) Variables() map[string]general.VariableSymbol {
	return Variables
}

func (b BoundScope) TryDeclare(variable general.VariableSymbol) bool {
	name := variable.Name
	_, exist := Variables[name]
	if exist {
		return false
	}
	Variables[name] = variable
	return true
}
