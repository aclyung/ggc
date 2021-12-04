package binding

import (
	"almeng.com/glang/glang/binding/boundNode"
	"almeng.com/glang/glang/expression"
	"almeng.com/glang/glang/general"
	"almeng.com/glang/glang/syntax"
	"reflect"
)

type Binder struct {
	Diag general.Diags
}

func NewBinder() *Binder {
	return &Binder{general.NewDiag()}
}

func (b *Binder) Bind(exp syntax.ExpressionSyntax) boundNode.BoundExpression {
	switch exp.Type() {
	case syntax.ExpLiteral:
		return b.BindLiteralExpression(exp.(*expression.Literal))
	case syntax.ExpBinary:
		return b.BindBinaryExpression(exp.(*expression.BinaryExpressionSyntax))
	case syntax.ExpUnary:
		return b.BindUnaryExpression(exp.(*expression.UnaryExpressionSyntax))
	}
	return nil
}

func isNumber(p reflect.Type) bool {
	switch p.Kind() {
	case reflect.Int64, reflect.Float64:
		return true
	}
	return false
}

func isBool(p reflect.Type) bool {
	return p.Kind() == reflect.Bool
}
