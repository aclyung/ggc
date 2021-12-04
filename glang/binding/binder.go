package binding

import (
	"reflect"

	"almeng.com/glang/binding/boundNode"
	"almeng.com/glang/expression"
	"almeng.com/glang/general"
	"almeng.com/glang/syntax"
)

type Binder struct {
	Diag      general.Diags
	Variables *map[string]boundNode.BoundExpression
}

func NewBinder(vars *map[string]boundNode.BoundExpression) *Binder {
	return &Binder{general.NewDiag(), vars}
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
	case syntax.ExpAssign:
		return b.BindAssignExpression(exp.(*expression.AssignmentExpressionSyntax))
	case syntax.ExpName:
		return b.BindNameExpression(exp.(*expression.NameExpressionSyntax))
	}
	panic("Unexpected syntax")
}

func (b *Binder) BindParenthesisExpression(exp *expression.ParenExpressionSyntax) boundNode.BoundExpression {
	return b.Bind(exp.Expression)
}

func (b *Binder) BindAssignExpression(exp *expression.AssignmentExpressionSyntax) boundNode.BoundExpression {
	var name = exp.Ident.Text
	boundExpression := b.Bind(exp.Expression)
	vars := *(b.Variables)
	variable, exist := vars[name]
	if (exist && variable.Type() == boundExpression.Type()) || !exist {
		return NewBoundAssignmentExpression(name, boundExpression)
	}
	b.Diag.VariableTypeMisMatch(exp.AssignToken.Span, name, variable.Type().String(), boundExpression.Type().String())
	return NewBoundLiteralExpression(int64(0))
}

func (b *Binder) BindNameExpression(exp *expression.NameExpressionSyntax) boundNode.BoundExpression {
	name := exp.Ident.Text
	vars := *(b.Variables)
	val, ok := vars[name]
	if !ok {
		b.Diag.UndefinedIdentifier(exp.Ident.Span, name)
		return NewBoundLiteralExpression(int64(0))
	}
	res := *(val.(*BoundLiteralExpression))
	return NewBoundVariableExpression(name, reflect.TypeOf(res.Value).Kind())
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
