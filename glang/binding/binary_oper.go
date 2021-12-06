package binding

import (
	"reflect"

	"almeng.com/glang/token"
)

// Binary Operator Kind

type BoundBinaryOperKind = int

const (
	ILLEGAL BoundBinaryOperKind = iota
	ADD
	SUB
	MUL
	QUO

	LSS
	LEQ
	GTR
	GEQ

	log_beg
	LAND
	LOR
	EQL
	NEQ
	log_end
)

type BoundBinaryOperator struct {
	Token      token.Token
	Oper       BoundUnaryOperKind
	Left       reflect.Kind
	Right      reflect.Kind
	ResultType reflect.Kind
}

var IllegalBinaryOperator = NewBoundBinaryOperator(token.ILLEGAL, ILLEGAL, reflect.Invalid, reflect.Invalid)

// operator type binding

var BinaryOpers = []BoundBinaryOperator{
	NewBoundBinaryOperator(token.LAND, LAND, reflect.Bool, reflect.Bool),
	NewBoundBinaryOperator(token.LOR, LOR, reflect.Bool, reflect.Bool),

	NewBoundBinaryOperator(token.EQL, EQL, reflect.Interface, reflect.Bool),
	NewBoundBinaryOperator(token.NEQ, NEQ, reflect.Interface, reflect.Bool),

	NewBoundBinaryOperator(token.LSS, LSS, reflect.Interface, reflect.Bool),
	NewBoundBinaryOperator(token.LEQ, LEQ, reflect.Interface, reflect.Bool),
	NewBoundBinaryOperator(token.GTR, GTR, reflect.Interface, reflect.Bool),
	NewBoundBinaryOperator(token.GEQ, GEQ, reflect.Interface, reflect.Bool),

	NewBoundBinaryOperator(token.ADD, ADD, reflect.Int64, reflect.Int64),
	NewBoundBinaryOperator(token.SUB, SUB, reflect.Int64, reflect.Int64),
	NewBoundBinaryOperator(token.MUL, MUL, reflect.Int64, reflect.Int64),
	NewBoundBinaryOperator(token.QUO, QUO, reflect.Int64, reflect.Int64),

	NewBoundBinaryOperator(token.ADD, ADD, reflect.Float64, reflect.Float64),
	NewBoundBinaryOperator(token.SUB, SUB, reflect.Float64, reflect.Float64),
	NewBoundBinaryOperator(token.MUL, MUL, reflect.Float64, reflect.Float64),
	NewBoundBinaryOperator(token.QUO, QUO, reflect.Float64, reflect.Float64),
}

// define binary operator

func BinaryOperator(tok token.Token, kind BoundUnaryOperKind, left reflect.Kind, right reflect.Kind, resType reflect.Kind) BoundBinaryOperator {
	return BoundBinaryOperator{tok, kind, left, right, resType}
}

// binary operator defining interface

func NewBoundBinaryOperator(tok token.Token, kind BoundUnaryOperKind, operand reflect.Kind, resType reflect.Kind) BoundBinaryOperator {
	return BinaryOperator(tok, kind, operand, operand, resType)
}

// check Operator is defined for operands
// if defined, returns proper BoundBinaryOperator
// if not, returns IllegalBinaryOperator

func BindBinaryOperator(tok token.Token, left reflect.Kind, right reflect.Kind) BoundBinaryOperator {
	if left != right {
		return IllegalBinaryOperator
	}
	for _, op := range BinaryOpers {
		if op.Token == tok && op.Left == left && op.Right == right {
			return op
		}
		if op.Token == tok && op.Left == reflect.Interface {
			return op
		}
	}

	return IllegalBinaryOperator
}

// returns if given operator is logical operator or not

func IsLogical(kind BoundBinaryOperKind) bool {
	return log_beg < kind && kind < log_end
}
