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
	case syntax.ExpLiteral:
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
	if lit.IsKindValid() {
		return NewBoundLiteralExpression(lit.Value())
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
	if isNumber(Tleft) && isNumber(Tright) {
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
	}

	if isBool(Tleft) && isBool(Tright) {
		switch kind {
		case token.LAND:
			return LAND
		case token.LOR:
			return LOR
		}
	}

	return ILLEGAL
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

func BindUnaryExpressionKind(kind token.Token, operandType reflect.Type) BoundUnaryOperKind {
	if isNumber(operandType) {
		switch kind {
		case token.ADD:
			return Identity
		case token.SUB:
			return Negation
		}
	}
	if isBool(operandType) {
		switch kind {
		case token.NOT:
			return NOT
		}
	}
	return ILLEGAL
}
