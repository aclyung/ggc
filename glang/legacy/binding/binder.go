package binding

import (
	"almeng.com/glang/legacy/binding/boundNode"
	expression2 "almeng.com/glang/legacy/expression"
	general2 "almeng.com/glang/legacy/general"
	syntax2 "almeng.com/glang/legacy/syntax"
)

type Binder struct {
	Diag  general2.Diags
	Scope BoundScope
}

func NewBinder(parent Scope) *Binder {
	s := NewBoundScope(parent.(BoundScope))
	return &Binder{Diag: general2.NewDiag(), Scope: s}
}

// bind ExpressionSyntax to BoundExpression

func (b *Binder) Bind(exp syntax2.ExpressionSyntax) boundNode.BoundExpression {
	switch exp.Type() {
	case syntax2.ExpLiteral:
		return b.BindLiteralExpression(exp)
	case syntax2.ExpBinary:
		return b.BindBinaryExpression(exp)
	case syntax2.ExpUnary:
		return b.BindUnaryExpression(exp)
	case syntax2.ExpParen:
		return b.BindParenthesisExpression(exp)
	case syntax2.ExpAssign:
		return b.BindAssignExpression(exp)
	case syntax2.ExpName:
		return b.BindIdentExpression(exp)
	case syntax2.EOF:
		return NewBoundEOFExpression()
	case syntax2.ILLEGAL:
		//b.Diag.Diagnose(Text.Span(0, 0), "Illegal", general.ERROR)
		return NewBoundEOFExpression()
	}
	panic("Unexpected syntax")
}

func (b *Binder) BindParenthesisExpression(exp syntax2.ExpressionSyntax) boundNode.BoundExpression {
	return b.Bind(exp.(*expression2.ParenExpressionSyntax).Expression)
}

func (b *Binder) BindAssignExpression(exp syntax2.ExpressionSyntax) boundNode.BoundExpression {
	e := exp.(*expression2.AssignmentExpressionSyntax)
	name := e.Ident.Text
	boundExpression := b.Bind(e.Expression)
	vari := general2.NewVariableSymbol(name, boundExpression.Type())

	if !b.Scope.TryDeclare(vari) {
		b.Diag.IdentifierAlreadyDeclared(e.Ident.Span, name)
		return NewBoundLiteralExpression(int64(0))
	}
	vars := b.Scope.Variables()

	newVal := general2.NewVariableSymbol(name, boundExpression.Type())

	if vars[name].Type == boundExpression.Type() {
		return NewBoundAssignmentExpression(newVal, boundExpression)
	}
	b.Diag.VariableTypeMisMatch(e.AssignToken.Span, name, vars[name].Type.String(), boundExpression.Type().String())
	return NewBoundLiteralExpression(int64(0))
}

func (b *Binder) BindIdentExpression(exp syntax2.ExpressionSyntax) boundNode.BoundExpression {
	e := exp.(*expression2.NameExpressionSyntax)
	name := e.Ident.Text

	if !b.Scope.TryLookUp(name) {
		b.Diag.UndefinedIdentifier(e.Ident.Span, name)
		return NewBoundLiteralExpression(int64(0))
	}

	val := b.Scope.Variables()[name]
	return NewBoundVariableExpression(val)
}
