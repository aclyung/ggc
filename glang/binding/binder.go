package binding

import (
	"almeng.com/glang/binding/boundNode"
	"almeng.com/glang/expression"
	"almeng.com/glang/general"
	"almeng.com/glang/general/Text"
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
		return b.BindLiteralExpression(exp)
	case syntax.ExpBinary:
		return b.BindBinaryExpression(exp)
	case syntax.ExpUnary:
		return b.BindUnaryExpression(exp)
	case syntax.ExpParen:
		return b.BindParenthesisExpression(exp)
	case syntax.ExpAssign:
		return b.BindAssignExpression(exp)
	case syntax.ExpName:
		return b.BindIdentExpression(exp)
	case syntax.EOF:
		b.Diag.Diagnose(Text.Span(0, 0), "EOF", general.ERROR)
		return NewBoundEOFExpression()
	case syntax.ILLEGAL:
		b.Diag.Diagnose(Text.Span(0, 0), "Illegal", general.ERROR)
		return NewBoundEOFExpression()
	}
	panic("Unexpected syntax")
}

func (b *Binder) BindParenthesisExpression(exp syntax.ExpressionSyntax) boundNode.BoundExpression {
	return b.Bind(exp.(*expression.ParenExpressionSyntax).Expression)
}

func (b *Binder) BindAssignExpression(exp syntax.ExpressionSyntax) boundNode.BoundExpression {
	e := exp.(*expression.AssignmentExpressionSyntax)
	var name = e.Ident.Text
	boundExpression := b.Bind(e.Expression)
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
	b.Diag.VariableTypeMisMatch(e.AssignToken.Span, name, val.Type.String(), boundExpression.Type().String())
	return NewBoundLiteralExpression(int64(0))
}

func (b *Binder) BindIdentExpression(exp syntax.ExpressionSyntax) boundNode.BoundExpression {
	e := exp.(*expression.NameExpressionSyntax)
	name := e.Ident.Text
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
		b.Diag.UndefinedIdentifier(e.Ident.Span, name)
		return NewBoundLiteralExpression(int64(0))
	}
	return NewBoundVariableExpression(val)
}
