package evaluator

import (
	"almeng.com/glang/glang/binding"
	"almeng.com/glang/glang/binding/boundNode"
	_ "almeng.com/glang/glang/expression"
	_ "almeng.com/glang/glang/lexer"
	_ "almeng.com/glang/glang/syntax"
	_ "almeng.com/glang/glang/token"
	_ "fmt"
	"reflect"
)

type Evaluator struct {
	root boundNode.BoundExpression
}

func NewEvaluator(root boundNode.BoundExpression) *Evaluator {
	return &Evaluator{root}
}

func (e *Evaluator) Evaluate() boundNode.BoundExpression {
	return ExpressionEvaluation(e.root)
}

func castNumber(l boundNode.BoundExpression, r boundNode.BoundExpression) (lfloat float64, lint int64, rfloat float64, rint int64, isInt bool) {
	lfloat, rfloat, lint, rint = 0, 0, 0, 0
	left, right := l.(*binding.BoundLiteralExpression), r.(*binding.BoundLiteralExpression)
	var lval, rval interface{} = left.Value, right.Value
	isInt = (left.Type() == right.Type()) && left.Type().Kind() == reflect.Int64
	if isInt {
		lint, rint = lval.(int64), rval.(int64)
		return
	}

	isTypeSame := left.Type() == right.Type()
	if isTypeSame {
		lfloat, rfloat = left.Value.(float64), right.Value.(float64)
		return
	}

	switch left.Value.(type) {
	case float64:
		lfloat = lval.(float64)
		rint = rval.(int64)
	default:
		lint = lval.(int64)
		rfloat = rval.(float64)
	}
	return
}

func returnToken(val interface{}) boundNode.BoundExpression {
	return binding.NewBoundLiteralExpression(val)
}

func ExpressionEvaluation(root boundNode.BoundExpression) boundNode.BoundExpression {
	expType := root.Kind()
	switch expType {
	case boundNode.Literal:
		val := root.(*binding.BoundLiteralExpression)
		return val
	case boundNode.Unary:
		u := root.(*binding.BoundUnaryExpression)
		operand := ExpressionEvaluation(u.Operand)
		switch u.Oper {
		case binding.Identity:
			return operand
		case binding.Negation:
			minus := returnToken(int64(-1))
			lfloat, lint, rfloat, rint, isInt := castNumber(minus, operand)
			if isInt {
				val := lint * rint
				return returnToken(val)
			}
			resL := lfloat + float64(lint)
			resR := rfloat + float64(rint)
			val := resL * resR
			return returnToken(val)
		case binding.NOT:
			return returnToken(!operand.(*binding.BoundLiteralExpression).Value.(bool))

		}
	case boundNode.Binary:
		nod := root.(*binding.BoundBinaryExpression)
		var oper binding.BoundBinaryOperKind = nod.Oper
		left, right := nod.Left, nod.Right
		left = ExpressionEvaluation(left)
		right = ExpressionEvaluation(right)

		if binding.IsLogical(oper) {
			Lleft, Lright := left.(*binding.BoundLiteralExpression).Value.(bool), right.(*binding.BoundLiteralExpression).Value.(bool)
			var res bool
			switch oper {
			case binding.LAND:
				res = Lleft && Lright
			case binding.LOR:
				res = Lleft || Lright
			}
			return returnToken(res)
		}

		lfloat, lint, rfloat, rint, isInt := castNumber(left, right)
		resL := lfloat + float64(lint)
		resR := rfloat + float64(rint)

		switch oper {
		case binding.ADD:
			if isInt {
				val := lint + rint
				return returnToken(val)
			}
			val := resL + resR
			return returnToken(val)
		case binding.SUB:
			if isInt {
				val := lint - rint
				return returnToken(val)
			}
			val := resL - resR
			return returnToken(val)
		case binding.MUL:
			if isInt {
				val := lint * rint
				return returnToken(val)
			}
			val := resL * resR
			return returnToken(val)
		case binding.QUO:
			if isInt {
				val := lint / rint
				return returnToken(val)
			}
			val := resL / resR
			return returnToken(val)

		default:
			return nil
		}
		//case syntax.ExpParen:
		//	return ExpressionEvaluation(root.(*expression.ParenExpressionSyntax).Expression)
	}
	return nil
}
