package parser

import (
	"almeng.com/glang/glang/expression"
	"almeng.com/glang/glang/syntax"
	"almeng.com/glang/glang/token"
)

func (p *parser) parseExpression(parentPrecedence int) syntax.ExpressionSyntax {
	var left syntax.ExpressionSyntax
	uPrec := p.current().Kind().Precedence()
	switch p.current().Kind() {
	case token.ADD:
		fallthrough
	case token.SUB:
		uPrec = token.UnaryPrec
	default:
		break

	}
	if uPrec != token.LowestPrec && uPrec >= parentPrecedence {
		oper := p.NextToken()
		operand := p.parseExpression(uPrec)
		left = expression.NewUnaryExpressionSyntax(oper, operand)
	} else {
		left = p.ParsePrevExpression()
	}
	for {
		precedence := p.current().Kind().Precedence()
		if precedence == 0 || precedence <= parentPrecedence {
			break
		}
		oper := p.NextToken()
		right := p.parseExpression(precedence)
		left = expression.NewBinaryExpressionSyntax(left, oper, right)
	}
	return left
}

func (p *parser) ParsePrevExpression() syntax.ExpressionSyntax {
	tok := p.current().Kind()
	var numTok expression.SyntaxToken

	switch tok {
	case token.LPAREN:
		left := p.NextToken()
		express := p.parseExpression(0)
		right := p.MatchToken(token.RPAREN)
		return expression.NewParenExpressionSyntax(left, express, right)
	case token.BOOL:
		val := p.MatchToken(token.BOOL)
		return expression.NewliteralExpressionSyntax(val)
	case token.INT:
		numTok = p.MatchToken(token.INT)
	case token.FLOAT:
		numTok = p.MatchToken(token.FLOAT)
	default:
		numTok = p.MatchToken(token.EOF)
	}
	return expression.NewliteralExpressionSyntax(numTok)
}
