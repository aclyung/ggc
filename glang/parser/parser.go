package parser

import (
	"almeng.com/glang/glang/ast/tree"
	"almeng.com/glang/glang/general"
	"almeng.com/glang/glang/lexer"
	"almeng.com/glang/glang/parser/node"
	syntax2 "almeng.com/glang/glang/syntax"
	"almeng.com/glang/glang/token"
	"bufio"
	"fmt"
	"os"
)

type parser struct {
	tokens      []lexer.SyntaxToken
	position    int
	Diagnostics []general.Diag
}

func Parser(text string) *parser {
	pars := &parser{}
	lex := lexer.NewLexer(text)
	pars.tokens = *new([]lexer.SyntaxToken)
	pars.position = 0

	for {
		tok := lex.NextToken()
		if tok == nil {
			continue
		}
		if tok.Token != token.ILLEGAL {
			pars.tokens = append(pars.tokens, *tok)
		}
		if tok.Token == token.EOF {
			break
		}

	}
	pars.Diagnostics = lex.Diagnostics
	//fmt.Println(pars.tokens)
	//fmt.Println(pars)
	return pars
}

// Returns current Token and Moves forward
func (p *parser) NextToken() lexer.SyntaxToken {
	current := p.current()
	p.position++
	return current
}

func (p *parser) Diagnose(text string, l general.Level) {
	diag := general.Diag{text, l}
	p.Diagnostics = append(p.Diagnostics, diag)
}

func (p *parser) Match(tok token.Token) lexer.SyntaxToken {
	cur := p.current()
	if cur.Kind() == tok {
		return p.NextToken()
	}
	wanted := fmt.Sprint(tok)
	got := fmt.Sprint(cur.Kind())
	p.Diagnose("WARNING: Expected <"+wanted+">, got <"+got+">", general.WARN)
	return *(lexer.Token(p.current().Kind(), p.position, "", nil))
}

func (p *parser) ParseTerm() node.ExpressionSyntax {
	if len(p.tokens) == 0 {
		p.Diagnose("ERROR: Syntax analyzing failed. No Token has been provided", general.ERROR)
		return nil
	}
	left := p.ParseFactor()
	for cur := p.current().Kind(); cur == token.ADD || cur == token.SUB; {
		opTok := p.NextToken()
		right := p.ParseFactor()
		left = syntax2.NewBinaryExpressionSyntax(left, opTok, right)
		cur = p.current().Kind()
	}
	return left
}

func (p *parser) ParseFactor() node.ExpressionSyntax {
	if len(p.tokens) == 0 {
		p.Diagnose("ERROR: Syntax analyzing failed. No Token has been provided", general.ERROR)
		return nil
	}
	left := p.ParsePrevExpression()
	for cur := p.current().Kind(); cur == token.MUL || cur == token.QUO; {
		opTok := p.NextToken()
		right := p.ParsePrevExpression()
		left = syntax2.NewBinaryExpressionSyntax(left, opTok, right)
		cur = p.current().Kind()
	}
	return left
}

func isOper(tok token.Token) bool {
	if token.Operator_beg < tok && tok <= token.REM {
		return true
	}
	return false
}

func (p *parser) parseExpression() node.ExpressionSyntax {
	return p.ParseTerm()
}

func (p *parser) ParsePrevExpression() node.ExpressionSyntax {
	tok := p.current().Kind()
	var numTok lexer.SyntaxToken
	switch tok {
	case token.LPAREN:
		left := p.NextToken()
		express := p.parseExpression()
		right := p.Match(token.RPAREN)
		return syntax2.NewParenExpressionSyntax(left, express, right)
	case token.INT:
		numTok = p.Match(token.INT)
	case token.FLOAT:
		numTok = p.Match(token.FLOAT)
	default:
		numTok = p.Match(token.EOF)
	}
	return syntax2.NewNumberExpressionSyntax(numTok)
}

func (pars *parser) current() lexer.SyntaxToken { return pars.peek(0) }

func (pars *parser) peek(offset int) lexer.SyntaxToken {
	index := pars.position + offset
	if index >= len(pars.tokens) {
		return pars.tokens[len(pars.tokens)-1]
	}

	return pars.tokens[index]
}

func (p *parser) Parse() tree.SyntaxTree {
	exp := p.ParseTerm()
	eoftok := p.Match(token.EOF)
	return tree.NewSyntaxTree(p.Diagnostics, exp, eoftok)
}

// func Parser(file os.File) string {
// 	str := ""
// 	buf := newBuffer(&file)
// 	for buf.Scan() {
// 		text := buf.Text()
// 		str += text + ";\n"
// 		general.ErrCheck(buf.Err())
// 	}

// 	return str
// }

func newBuffer(file *os.File) *bufio.Scanner {
	return bufio.NewScanner(file)
}
