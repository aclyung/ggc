package binding

import (
	"reflect"

	"almeng.com/glang/binding/boundNode"
	"almeng.com/glang/expression"
	"almeng.com/glang/general"
	"almeng.com/glang/syntax"
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
	case syntax.ExpParen:
		return b.BindParenthesisExpression(exp.(*expression.ParenExpressionSyntax))

	}
	return nil
}

func (b *Binder) BindParenthesisExpression(exp *expression.ParenExpressionSyntax) boundNode.BoundExpression {
	return b.Bind(exp.Expression)
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
