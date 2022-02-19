package evaluator

import (
	"reflect"

	binding2 "main/legacy/binding"
	boundNode2 "main/legacy/binding/boundNode"
	general2 "main/legacy/general"
)

type Evaluator struct {
	root boundNode2.BoundExpression
	vars *map[general2.VariableSymbol]boundNode2.BoundExpression
	diag general2.Diags
}

func NewEvaluator(root boundNode2.BoundExpression, vars *map[general2.VariableSymbol]boundNode2.BoundExpression) *Evaluator {
	return &Evaluator{root, vars, general2.NewDiag()}
}

func (e *Evaluator) Evaluate() boundNode2.BoundExpression {
	return e.ExpressionEvaluation(e.root)
}

func castNumber(l boundNode2.BoundExpression, r boundNode2.BoundExpression) (lfloat float64, lint int64, rfloat float64, rint int64, isInt bool) {
	lfloat, rfloat, lint, rint = 0, 0, 0, 0
	left, right := l.(*binding2.BoundLiteralExpression), r.(*binding2.BoundLiteralExpression)
	var lval, rval interface{} = left.Value, right.Value

	isSame := left.Type() == right.Type()
	isInt = isSame && left.Type() == reflect.Int64

	if isSame {
		goto same
	} else {
		goto different
	}

same:
	if isInt {
		goto same_int
	}
	goto same_float

different:
	switch left.Value.(type) {
	case float64:
		goto float
	default:
		goto int
	}

same_float:
	lfloat, rfloat = left.Value.(float64), right.Value.(float64)
	return

same_int:
	lint, rint = lval.(int64), rval.(int64)
	return

float:
	lfloat = lval.(float64)
	rint = rval.(int64)
	return
int:
	lint = lval.(int64)
	rfloat = rval.(float64)
	return
}

//isInt = (left.Type() == right.Type()) && left.Type() == reflect.Int64
//if isInt {
//	lint, rint = lval.(int64), rval.(int64)
//	return
//}
//
//isTypeSame := left.Type() == right.Type()
//if isTypeSame {
//	lfloat, rfloat = left.Value.(float64), right.Value.(float64)
//	return
//}

//switch left.Value.(type) {
//case float64:
//	lfloat = lval.(float64)
//	rint = rval.(int64)
//default:
//	lint = lval.(int64)
//	rfloat = rval.(float64)
//}
//return

func returnToken(val interface{}) boundNode2.BoundExpression {
	return binding2.NewBoundLiteralExpression(val)
}

func (e *Evaluator) ExpressionEvaluation(root boundNode2.BoundExpression) boundNode2.BoundExpression {
	expType := root.Kind()
	switch expType {
	case boundNode2.Literal:
		return e.EvaluateLiteral(root)
	case boundNode2.Variable:
		return e.EvaluateVariable(root)
	case boundNode2.Assign:
		return e.EvaluateAssignment(root)
	case boundNode2.Unary:
		return e.EvaluateUnary(root)
	case boundNode2.Binary:
		return e.EvaluateBinary(root)
	}
	return nil
}

func (e *Evaluator) EvaluateLiteral(root boundNode2.BoundExpression) boundNode2.BoundExpression {
	val := root.(*binding2.BoundLiteralExpression)
	return val
}

func (e *Evaluator) EvaluateVariable(root boundNode2.BoundExpression) boundNode2.BoundExpression {
	v := root.(*binding2.BoundVariableExpression)
	vars := *(e.vars)
	name := v.Variable
	val := vars[name]
	return val
}

func (e *Evaluator) EvaluateAssignment(root boundNode2.BoundExpression) boundNode2.BoundExpression {
	a := root.(*binding2.BoundAssignmentExpression)
	vars := *(e.vars)
	val := e.ExpressionEvaluation(a.Expression)
	vars[a.Variable] = val
	return binding2.InvalidLiteralExpression
}

func (e *Evaluator) EvaluateUnary(root boundNode2.BoundExpression) boundNode2.BoundExpression {
	u := root.(*binding2.BoundUnaryExpression)
	operand := e.ExpressionEvaluation(u.Operand)
	var rtn boundNode2.BoundExpression
	switch u.Oper.OperKind {
	case binding2.Identity:
		rtn = operand
	case binding2.Negation:
		minus := returnToken(int64(-1))
		lfloat, lint, rfloat, rint, isInt := castNumber(minus, operand)
		if isInt {
			val := lint * rint
			rtn = returnToken(val)
			break
		}
		resL := lfloat + float64(lint)
		resR := rfloat + float64(rint)
		val := resL * resR
		rtn = returnToken(val)
	case binding2.NOT:
		rtn = returnToken(!operand.(*binding2.BoundLiteralExpression).Value.(bool))
	}
	return rtn
}

func (e *Evaluator) EvaluateBinary(root boundNode2.BoundExpression) boundNode2.BoundExpression {
	b := root.(*binding2.BoundBinaryExpression)
	var oper binding2.BoundBinaryOperKind = b.Oper.Oper
	left, right := b.Left, b.Right
	left = e.ExpressionEvaluation(left)
	right = e.ExpressionEvaluation(right)

	if binding2.IsLogical(oper) {
		l, r := left.(*binding2.BoundLiteralExpression).Value, right.(*binding2.BoundLiteralExpression).Value
		var res bool
		switch oper {
		case binding2.LAND:
			res = l.(bool) && r.(bool)
		case binding2.LOR:
			res = l.(bool) || r.(bool)
		case binding2.EQL:
			res = l == r
		case binding2.NEQ:
			res = l != r
		}
		return returnToken(res)
	}

	lfloat, lint, rfloat, rint, isInt := castNumber(left, right)
	resL := lfloat + float64(lint)
	resR := rfloat + float64(rint)

	switch oper {
	case binding2.ADD:
		if isInt {
			val := lint + rint
			return returnToken(val)
		}
		val := resL + resR
		return returnToken(val)
	case binding2.SUB:
		if isInt {
			val := lint - rint
			return returnToken(val)
		}
		val := resL - resR
		return returnToken(val)
	case binding2.MUL:
		if isInt {
			val := lint * rint
			return returnToken(val)
		}
		val := resL * resR
		return returnToken(val)
	case binding2.QUO:
		if isInt {
			val := lint / rint
			return returnToken(val)
		}
		val := resL / resR
		return returnToken(val)
	case binding2.LSS:
		if isInt {
			val := lint < rint
			return returnToken(val)
		}
		val := resL < resR
		return returnToken(val)
	case binding2.LEQ:
		if isInt {
			val := lint <= rint
			return returnToken(val)
		}
		val := resL <= resR
		return returnToken(val)
	case binding2.GTR:
		if isInt {
			val := lint > rint
			return returnToken(val)
		}
		val := resL > resR
		return returnToken(val)
	case binding2.GEQ:
		if isInt {
			val := lint >= rint
			return returnToken(val)
		}
		val := resL >= resR
		return returnToken(val)
	default:
		return nil
	}
}
