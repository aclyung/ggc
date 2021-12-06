package evaluator

import (
	"almeng.com/glang/binding"
	"almeng.com/glang/binding/boundNode"
	_ "almeng.com/glang/expression"
	"almeng.com/glang/general"
	_ "almeng.com/glang/lexer"
	_ "almeng.com/glang/syntax"
	_ "almeng.com/glang/token"
	_ "fmt"
	"reflect"
)

type Evaluator struct {
	root boundNode.BoundExpression
	vars *map[general.VariableSymbol]boundNode.BoundExpression
	diag general.Diags
}

func NewEvaluator(root boundNode.BoundExpression, vars *map[general.VariableSymbol]boundNode.BoundExpression) *Evaluator {
	return &Evaluator{root, vars, general.NewDiag()}
}

func (e *Evaluator) Evaluate() boundNode.BoundExpression {
	return e.ExpressionEvaluation(e.root)
}

func castNumber(l boundNode.BoundExpression, r boundNode.BoundExpression) (lfloat float64, lint int64, rfloat float64, rint int64, isInt bool) {
	lfloat, rfloat, lint, rint = 0, 0, 0, 0
	left, right := l.(*binding.BoundLiteralExpression), r.(*binding.BoundLiteralExpression)
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

func returnToken(val interface{}) boundNode.BoundExpression {
	return binding.NewBoundLiteralExpression(val)
}

func (e *Evaluator) ExpressionEvaluation(root boundNode.BoundExpression) boundNode.BoundExpression {
	expType := root.Kind()
	switch expType {
	case boundNode.Literal:
		return e.EvaluateLiteral(root)
	case boundNode.Variable:
		return e.EvaluateVariable(root)
	case boundNode.Assign:
		return e.EvaluateAssignment(root)
	case boundNode.Unary:
		return e.EvaluateUnary(root)
	case boundNode.Binary:
		return e.EvaluateBinary(root)
	}
	return nil
}

func (e *Evaluator) EvaluateLiteral(root boundNode.BoundExpression) boundNode.BoundExpression {
	val := root.(*binding.BoundLiteralExpression)
	return val
}

func (e *Evaluator) EvaluateVariable(root boundNode.BoundExpression) boundNode.BoundExpression {
	v := root.(*binding.BoundVariableExpression)
	vars := *(e.vars)
	name := v.Variable
	val := vars[name]
	return val
}

func (e *Evaluator) EvaluateAssignment(root boundNode.BoundExpression) boundNode.BoundExpression {
	a := root.(*binding.BoundAssignmentExpression)
	vars := *(e.vars)
	val := e.ExpressionEvaluation(a.Expression)
	vars[a.Variable] = val
	return binding.InvalidLiteralExpression
}

func (e *Evaluator) EvaluateUnary(root boundNode.BoundExpression) boundNode.BoundExpression {
	u := root.(*binding.BoundUnaryExpression)
	operand := e.ExpressionEvaluation(u.Operand)
	var rtn boundNode.BoundExpression
	switch u.Oper.OperKind {
	case binding.Identity:
		rtn = operand
	case binding.Negation:
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
	case binding.NOT:
		rtn = returnToken(!operand.(*binding.BoundLiteralExpression).Value.(bool))
	}
	return rtn
}

func (e *Evaluator) EvaluateBinary(root boundNode.BoundExpression) boundNode.BoundExpression {
	b := root.(*binding.BoundBinaryExpression)
	var oper binding.BoundBinaryOperKind = b.Oper.Oper
	left, right := b.Left, b.Right
	left = e.ExpressionEvaluation(left)
	right = e.ExpressionEvaluation(right)

	if binding.IsLogical(oper) {
		l, r := left.(*binding.BoundLiteralExpression).Value, right.(*binding.BoundLiteralExpression).Value
		var res bool
		switch oper {
		case binding.LAND:
			res = l.(bool) && r.(bool)
		case binding.LOR:
			res = l.(bool) || r.(bool)
		case binding.EQL:
			res = l == r
		case binding.NEQ:
			res = l != r
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
	case binding.LSS:
		if isInt {
			val := lint < rint
			return returnToken(val)
		}
		val := resL < resR
		return returnToken(val)
	case binding.LEQ:
		if isInt {
			val := lint <= rint
			return returnToken(val)
		}
		val := resL <= resR
		return returnToken(val)
	case binding.GTR:
		if isInt {
			val := lint > rint
			return returnToken(val)
		}
		val := resL > resR
		return returnToken(val)
	case binding.GEQ:
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
