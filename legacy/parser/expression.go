package parser

import (
	expression2 "main/legacy/expression"
	"main/legacy/general/Text"
	"main/legacy/syntax"
	"main/legacy/token"
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
		left = expression2.NewUnaryExpressionSyntax(oper, operand)
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
		left = expression2.NewBinaryExpressionSyntax(left, oper, right)
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
		return expression2.NewAssigmentExpressionSyntax(ident, operTok, right)
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
		if len(p.tokens)-3 >= 0 {
			p.Diag.UnexpectedToken(Text.Span(0, 0), p.peek(-2).Kind().String(), token.EOF.String())
		} else {
			p.Diag.UnexpectedToken(Text.Span(0, 0), "Expression", token.EOF.String())
		}
		return expression2.NewEOFExpressionSyntax()
	default:
		return expression2.NewIllegalExpressionSyntax()
	}
}

func (p *parser) ParseParenthesizedExpression() syntax.ExpressionSyntax {
	left := p.NextToken()
	express := p.ParseExpression(0)
	right := p.MatchToken(token.RPAREN)
	return expression2.NewParenExpressionSyntax(left, express, right)
}

func (p *parser) ParseLiteralExpression() syntax.ExpressionSyntax {
	val := p.MatchToken(p.current().Kind())
	return expression2.NewLiteralExpressionSyntax(val)
}

func (p *parser) ParseIdentifierExpression() syntax.ExpressionSyntax {
	ident := p.NextToken()
	return expression2.NewNameExpressionSyntax(ident)
}
