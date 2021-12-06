package binding

import (
	"almeng.com/glang/binding/boundNode"
	"almeng.com/glang/expression"
	"almeng.com/glang/general"
	"almeng.com/glang/general/TextSpan"
	"almeng.com/glang/syntax"
)

type Binder struct {
	Diag      general.Diags
	Variables *map[general.VariableSymbol]boundNode.BoundExpression
}

func NewBinder(vars *map[general.VariableSymbol]boundNode.BoundExpression) *Binder {
	return &Binder{general.NewDiag(), vars}
}

// bind ExpressionSyntax to BoundExpression

func (b *Binder) Bind(exp syntax.ExpressionSyntax) boundNode.BoundExpression {
	switch exp.Type() {
	case syntax.ExpLiteral:
		return b.BindLiteralExpression(exp.(*expression.LiteralExpressionSyntax))
	case syntax.ExpBinary:
		return b.BindBinaryExpression(exp.(*expression.BinaryExpressionSyntax))
	case syntax.ExpUnary:
		return b.BindUnaryExpression(exp.(*expression.UnaryExpressionSyntax))
	case syntax.ExpParen:
		return b.BindParenthesisExpression(exp.(*expression.ParenExpressionSyntax))
	case syntax.ExpAssign:
		return b.BindAssignExpression(exp.(*expression.AssignmentExpressionSyntax))
	case syntax.ExpName:
		return b.BindIdentExpression(exp.(*expression.NameExpressionSyntax))
	case syntax.EOF:
		b.Diag.Diagnose(TextSpan.Span(0, 0), "EOF", general.ERROR)
		return NewBoundEOFExpression()
	case syntax.ILLEGAL:
		b.Diag.Diagnose(TextSpan.Span(0, 0), "Illegal", general.ERROR)
		return NewBoundEOFExpression()
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
	var val general.VariableSymbol
	exist := false
	for k, _ := range vars {
		if k.Name == name {
			val = k
			exist = true
			break
		}
	}
	newVal := *general.NewVariableSymbol(name, boundExpression.Type())

	if (exist && val.Type == boundExpression.Type()) || !exist {
		return NewBoundAssignmentExpression(newVal, boundExpression)
	}
	b.Diag.VariableTypeMisMatch(exp.AssignToken.Span, name, val.Type.String(), boundExpression.Type().String())
	return NewBoundLiteralExpression(int64(0))
}

func (b *Binder) BindIdentExpression(exp *expression.NameExpressionSyntax) boundNode.BoundExpression {
	name := exp.Ident.Text
	vars := *(b.Variables)
	var val general.VariableSymbol
	ok := false
	for k, _ := range vars {
		if k.Name == name {
			val = k
			ok = true
			break
		}
	}
	if !ok {
		b.Diag.UndefinedIdentifier(exp.Ident.Span, name)
		return NewBoundLiteralExpression(int64(0))
	}
	return NewBoundVariableExpression(val)
}

// implemented in operator bindings

//func isNumber(p reflect.Type) bool {
//	switch p.Kind() {
//	case reflect.Int64, reflect.Float64:
//		return true
//	}
//	return false
//}
//
//func isBool(p reflect.Type) bool {
//	return p.Kind() == reflect.Bool
//}
