package evaluator

import (
	"almeng.com/glang/glang/lexer"
	"almeng.com/glang/glang/parser/node"
	"almeng.com/glang/glang/syntax"
	"almeng.com/glang/glang/syntax/expression"
	"almeng.com/glang/glang/token"
	"fmt"
)

func castNumber(l node.ExpressionSyntax, r node.ExpressionSyntax) (lfloat float64, lint int64, rfloat float64, rint int64, isInt bool) {
	lfloat, rfloat, lint, rint = 0, 0, 0, 0
	left, right := l.(lexer.SyntaxToken), r.(lexer.SyntaxToken)
	var lval, rval interface{} = left.Value, right.Value
	isInt = (left.Kind() == right.Kind()) && left.Kind() == token.INT
	if isInt {
		lint, rint = lval.(int64), rval.(int64)
		return
	}

	isTypeSame := left.Kind() == right.Kind()
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

func returnToken(tok token.Token, val interface{}) lexer.SyntaxToken {
	if tok == token.FLOAT {
		resval := val.(float64)
		return *lexer.Token(tok, 0, fmt.Sprint(resval), resval)
	}
	resval := val.(int64)
	return *lexer.Token(tok, 0, fmt.Sprint(resval), resval)
}

func ExpressionEvaluation(root node.ExpressionSyntax) node.ExpressionSyntax {
	expType := root.Type()
	switch expType {
	case expression.ExpNum:
		val := root.(*syntax.Literal)
		return val.NumberToken
	case expression.ExpBi:
		nod := root.(*syntax.BinaryExpressionSyntax)
		oper := nod.OperatorToken.Kind()
		left, right := nod.Left, nod.Right
		left = ExpressionEvaluation(left)
		right = ExpressionEvaluation(right)

		lfloat, lint, rfloat, rint, isInt := castNumber(left, right)

		switch oper {
		case token.ADD:
			resL := lfloat + float64(lint)
			resR := rfloat + float64(rint)
			if isInt {
				val := lint + rint
				return returnToken(token.INT, val)
			}
			val := resL + resR
			return returnToken(token.FLOAT, val)
		case token.SUB:
			resL := lfloat + float64(lint)
			resR := rfloat + float64(rint)
			if isInt {
				val := lint - rint
				return returnToken(token.INT, val)
			}
			val := resL - resR
			return returnToken(token.FLOAT, val)
		case token.MUL:
			resL := lfloat + float64(lint)
			resR := rfloat + float64(rint)
			if isInt {
				val := lint * rint
				return returnToken(token.INT, val)
			}
			val := resL * resR
			return returnToken(token.FLOAT, val)
		case token.QUO:
			resL := lfloat + float64(lint)
			resR := rfloat + float64(rint)
			if isInt {
				val := lint / rint
				return returnToken(token.INT, val)
			}
			val := resL / resR
			return returnToken(token.FLOAT, val)
		default:
			return nil
		}
	case expression.ExpParen:
		return ExpressionEvaluation(root.(*syntax.ParenExpressionSyntax).Expression)
	}
	return nil
}
