package binding

import (
	"main/legacy/binding/boundNode"
	"main/legacy/expression"
	general2 "main/legacy/general"
)

type GlobalScope interface {
	Previous() GlobalScope
	Diags() general2.Diags
	Variables() map[string]general2.VariableSymbol
	Expression() boundNode.BoundExpression
}

type BoundGlobalScope struct {
	Prev GlobalScope
	Diag general2.Diags
	Vars map[string]general2.VariableSymbol
	Exp  boundNode.BoundExpression
}

func BindBoundGlobalScope(unit expression.CompilationUnit) BoundGlobalScope {
	binder := NewBinder(nil)
	exp := binder.Bind(unit.Expression)
	vars := binder.Scope.Variables()
	diags := binder.Diag
	return NewBoundGlobalScope(nil, diags, vars, exp)
}

func NewBoundGlobalScope(prev GlobalScope, diags general2.Diags, vars map[string]general2.VariableSymbol, exp boundNode.BoundExpression) BoundGlobalScope {
	return BoundGlobalScope{prev, diags, vars, exp}
}

func (b BoundGlobalScope) Previous() GlobalScope {
	return b.Prev
}

func (b BoundGlobalScope) Diags() general2.Diags {
	return b.Diag
}

func (b BoundGlobalScope) Variables() map[string]general2.VariableSymbol {
	return b.Vars
}

func (b BoundGlobalScope) Expression() boundNode.BoundExpression {
	return b.Exp
}
