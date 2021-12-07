package parser

import (
	"almeng.com/glang/expression"
	"almeng.com/glang/syntax"
	"almeng.com/glang/token"
)

func (p *parser) parseBinaryExpression(parentPrecedence int) syntax.ExpressionSyntax {
	var left syntax.ExpressionSyntax
	uPrec := p.current().Kind().Precedence()
	switch p.current().Kind() {
	case token.NOT, token.ADD, token.SUB:
		uPrec = token.UnaryPrec
	default:
		break

	}
	if uPrec != token.LowestPrec && uPrec >= parentPrecedence {
		oper := p.NextToken()
		operand := p.parseBinaryExpression(uPrec)
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
		right := p.parseBinaryExpression(precedence)
		left = expression.NewBinaryExpressionSyntax(left, oper, right)
	}
	return left
}

func (p *parser) ParseExpression(parentPrecedence int) syntax.ExpressionSyntax {
	return p.ParseAssignmentExpression(parentPrecedence)
}

func (p *parser) ParseAssignmentExpression(parentPrecedence int) syntax.ExpressionSyntax {
	pars := *p
	if pars.NextToken().Kind() == token.IDENT && pars.NextToken().Kind() == token.ASSIGN {
		ident := p.NextToken()
		operTok := p.NextToken()
		right := p.ParseAssignmentExpression(0)
		return expression.NewAssigmentExpressionSyntax(ident, operTok, right)
	}
	return p.parseBinaryExpression(0)
}

func (p *parser) ParsePrevExpression() syntax.ExpressionSyntax {
	switch p.current().Kind() {
	case token.LPAREN:
		return p.ParseParenthesizedExpression()
	case token.BOOL, token.INT, token.FLOAT:
		return p.ParseLiteralExpression()
	case token.IDENT:
		return p.ParseIdentifierExpression()
	case token.EOF:
		return expression.NewEOFExpressionSyntax()
	default:
		return expression.NewIllegalExpressionSyntax()
	}
}

func (p *parser) ParseParenthesizedExpression() syntax.ExpressionSyntax {
	left := p.NextToken()
	express := p.ParseExpression(0)
	right := p.MatchToken(token.RPAREN)
	return expression.NewParenExpressionSyntax(left, express, right)
}

func (p *parser) ParseLiteralExpression() syntax.ExpressionSyntax {
	val := p.MatchToken(p.current().Kind())
	return expression.NewLiteralExpressionSyntax(val)
}

func (p *parser) ParseIdentifierExpression() syntax.ExpressionSyntax {
	ident := p.NextToken()
	return expression.NewNameExpressionSyntax(ident)
}
