package binding

import (
	"almeng.com/glang/glang/token"
	"reflect"
)

type BoundUnaryOperKind = int

const (
	Identity BoundUnaryOperKind = iota
	Negation
	NOT
)

type BoundUnaryOperator struct {
	Token       token.Token
	OperKind    BoundUnaryOperKind
	OperandKind reflect.Kind
}

var IllegalUnaryOperator = NewBoundUnaryOperator(token.ILLEGAL, ILLEGAL, reflect.Invalid)

var unaryOpers = []BoundUnaryOperator{
	NewBoundUnaryOperator(token.NOT, NOT, reflect.Bool),
	NewBoundUnaryOperator(token.ADD, ADD, reflect.Int64),
	NewBoundUnaryOperator(token.SUB, SUB, reflect.Int64),
	NewBoundUnaryOperator(token.ADD, ADD, reflect.Float64),
	NewBoundUnaryOperator(token.SUB, SUB, reflect.Float64),
}

func NewBoundUnaryOperator(tok token.Token, kind BoundUnaryOperKind, operandType reflect.Kind) BoundUnaryOperator {
	return BoundUnaryOperator{tok, kind, operandType}
}

func BindUnaryOperator(tok token.Token, operandType reflect.Kind) BoundUnaryOperator {
	for _, op := range unaryOpers {
		if op.Token == tok && op.OperandKind == operandType {
			return op
		}
	}
	return IllegalUnaryOperator
}
