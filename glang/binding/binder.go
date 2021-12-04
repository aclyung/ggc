package binding

import (
	"almeng.com/glang/glang/binding/boundNode"
	"almeng.com/glang/glang/expression"
	"almeng.com/glang/glang/general"
	"almeng.com/glang/glang/syntax"
	"almeng.com/glang/glang/token"
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
	case syntax.ExpNum:
		return b.BindLiteralExpression(exp.(*expression.Literal))
	case syntax.ExpBinary:
		return b.BindBinaryExpression(exp.(*expression.BinaryExpressionSyntax))
	case syntax.ExpUnary:
		return b.BindUnaryExpression(exp.(*expression.UnaryExpressionSyntax))
	}
	return nil
}

func (b *Binder) BindLiteralExpression(exp syntax.ExpressionSyntax) boundNode.BoundExpression {
	lit := exp.(*expression.Literal)
	val := lit.Value()
	switch lit.Kind() {
	case token.INT:
		return NewBoundLiteralExpression(val)
	case token.FLOAT:
		return NewBoundLiteralExpression(val)
	}
	return NewBoundLiteralExpression(int64(0))
}
func (b *Binder) BindUnaryExpression(exp syntax.ExpressionSyntax) boundNode.BoundExpression {
	syn := exp.(*expression.UnaryExpressionSyntax)
	operand := b.Bind(syn.Operand)
	operKind := BindUnaryExpressionKind(syn.OperatorToken.Kind(), operand.Type())
	if operKind == ILLEGAL {
		b.Diag.Diagnose("Unary Operator '"+syn.OperatorToken.Text+"' is not defined for type "+operand.Type().String()+"", general.ERROR)
		return operand
	}
	return NewBoundUnaryExpression(operKind, operand)
}
func (b *Binder) BindBinaryExpression(exp syntax.ExpressionSyntax) boundNode.BoundExpression {
	biExp := exp.(*expression.BinaryExpressionSyntax)
	left := b.Bind(biExp.Left)
	right := b.Bind(biExp.Right)
	operKind := BindBinaryExpressionKind(biExp.Kind(), left.Type(), right.Type())
	if operKind == ILLEGAL {
		b.Diag.Diagnose("Binary Operator '"+biExp.OperatorToken.Text+"' is not defined for types "+left.Type().String()+" and "+right.Type().String(), general.ERROR)
		return left
	}
	return NewBoundBinaryExpression(left, operKind, right)
}

func BindBinaryExpressionKind(kind token.Token, Tleft reflect.Type, Tright reflect.Type) BoundBinaryOperKind {
	if Tleft.Kind() == Tright.Kind() {
		if !isTypeValid(Tleft) {
			return ILLEGAL
		}
	}
	if !(isTypeValid(Tleft) && isTypeValid(Tright)) {
		return ILLEGAL
	}

	switch kind {
	case token.ADD:
		return ADD
	case token.SUB:
		return SUB
	case token.MUL:
		return MUL
	case token.QUO:
		return QUO
	}
	return ILLEGAL
}

func isTypeValid(p reflect.Type) bool {
	return p.Kind() == reflect.Int64 || p.Kind() == reflect.Float64
}

func BindUnaryExpressionKind(kind token.Token, operandType reflect.Type) BoundUnaryOperKind {
	if !isTypeValid(operandType) {
		return ILLEGAL
	}

	switch kind {
	case token.ADD:
		return Identity
	case token.SUB:
		return Negation
	}
	return ILLEGAL
}
