package parser

import (
	"bufio"
	"fmt"
	"os"

	"almeng.com/glang/ast/tree"
	"almeng.com/glang/expression"
	"almeng.com/glang/general"
	"almeng.com/glang/lexer"
	"almeng.com/glang/token"
)

type parser struct {
	tokens   []expression.SyntaxToken
	position int
	Diag     general.Diags
	//Diagnostics []general.Diag
}

func Parser(text string) *parser {
	pars := &parser{}
	lex := lexer.NewLexer(text)
	pars.tokens = *new([]expression.SyntaxToken)
	pars.position = 0

	for {
		tok := lex.Lex()
		if tok == nil {
			continue
		}
		if tok.Kind() != token.ILLEGAL {
			pars.tokens = append(pars.tokens, *tok)
		}
		if tok.Kind() == token.EOF {
			break
		}
	}
	pars.Diag = lex.Diag
	//pars.Diagnostics = lex.Diagnostics
	//fmt.Println(pars.tokens)
	//fmt.Println(pars)
	return pars
}

// parse to AST

func (p *parser) Parse() tree.Tree {
	exp := p.ParseExpression(0)
	eoftok := p.MatchToken(token.EOF)
	return tree.NewSyntaxTree(p.Diag, exp, eoftok)
}

// returns current token and moves forward

func (p *parser) NextToken() expression.SyntaxToken {
	current := p.current()
	p.position++
	return current
}

//

func (p *parser) MatchToken(tok token.Token) expression.SyntaxToken {
	cur := p.current()
	if cur.Kind() == tok {
		return p.NextToken()
	}
	wanted := fmt.Sprint(tok)
	got := fmt.Sprint(cur.Kind())
	p.Diag.UnexpectedToken(p.current().Span, wanted, got)
	return *(expression.NewSyntaxToken(p.current().Kind(), p.position, "", nil))
}

//
//func (p *parser) Diagnose(text string, l general.Level) {
//	diag := general.Diag{text, l}
//	p.Diagnostics = append(p.Diagnostics, diag)
//}

func isOper(tok token.Token) bool {
	if token.Operator_beg < tok && tok <= token.REM {
		return true
	}
	return false
}

func (pars *parser) current() expression.SyntaxToken { return pars.peek(0) }

func (pars *parser) peek(offset int) expression.SyntaxToken {
	index := pars.position + offset
	if index >= len(pars.tokens) {
		return pars.tokens[len(pars.tokens)-1]
	}

	return pars.tokens[index]
}

func newBuffer(file *os.File) *bufio.Scanner {
	return bufio.NewScanner(file)
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
