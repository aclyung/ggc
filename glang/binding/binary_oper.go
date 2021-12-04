package binding

import (
	"almeng.com/glang/glang/token"
	"reflect"
)

type BoundBinaryOperKind = int

const (
	ILLEGAL BoundBinaryOperKind = iota
	ADD
	SUB
	MUL
	QUO

	log_beg
	LAND
	LOR
	log_end
)

type BoundBinaryOperator struct {
	Token token.Token
	Oper  BoundUnaryOperKind
	Left  reflect.Kind
	Right reflect.Kind
}

var IllegalBinaryOperator = NewBoundBinaryOperator(token.ILLEGAL, ILLEGAL, reflect.Invalid)

var BinaryOpers = []BoundBinaryOperator{
	NewBoundBinaryOperator(token.LAND, LAND, reflect.Bool),
	NewBoundBinaryOperator(token.LOR, LOR, reflect.Bool),

	NewBoundBinaryOperator(token.ADD, ADD, reflect.Int64),
	NewBoundBinaryOperator(token.SUB, SUB, reflect.Int64),
	NewBoundBinaryOperator(token.MUL, MUL, reflect.Int64),
	NewBoundBinaryOperator(token.QUO, QUO, reflect.Int64),

	NewBoundBinaryOperator(token.ADD, ADD, reflect.Float64),
	NewBoundBinaryOperator(token.SUB, SUB, reflect.Float64),
	NewBoundBinaryOperator(token.MUL, MUL, reflect.Float64),
	NewBoundBinaryOperator(token.QUO, QUO, reflect.Float64),
}

func BinaryOperator(tok token.Token, kind BoundUnaryOperKind, left reflect.Kind, right reflect.Kind) BoundBinaryOperator {
	return BoundBinaryOperator{tok, kind, left, right}
}
func NewBoundBinaryOperator(tok token.Token, kind BoundUnaryOperKind, operand reflect.Kind) BoundBinaryOperator {
	return BinaryOperator(tok, kind, operand, operand)
}

func BindBinaryOperator(tok token.Token, left reflect.Kind, right reflect.Kind) BoundBinaryOperator {
	for _, op := range BinaryOpers {
		if op.Token == tok && op.Left == left && op.Right == right {
			return op
		}
	}
	return IllegalBinaryOperator
}

func IsLogical(kind BoundBinaryOperKind) bool {
	return log_beg < kind && kind < log_end
}
