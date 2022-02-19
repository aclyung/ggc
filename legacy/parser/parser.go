package parser

import (
	"bufio"
	"fmt"
	"os"

	expression2 "main/legacy/expression"
	"main/legacy/general"
	"main/legacy/general/Text"
	"main/legacy/lexer"
	"main/legacy/token"
)

type parser struct {
	source   *Text.Source
	tokens   []expression2.SyntaxToken
	position int
	Diag     general.Diags
	//Diagnostics []general.Diag
}

func Parser(source *Text.Source) *parser {
	pars := &parser{source: source}
	lex := lexer.NewLexer(source)
	pars.tokens = *new([]expression2.SyntaxToken)
	pars.position = 0

	for {
		tok := lex.Lex()
		if tok == nil {
			continue
		}
		//if tok.Kind() == token.ILLEGAL {
		pars.tokens = append(pars.tokens, *tok)

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

func (p *parser) ParseCompilationUnit() expression2.CompilationUnit {
	exp := p.ParseExpression(0)
	eoftok := p.MatchToken(token.EOF)
	return expression2.NewCompilationUnit(exp, eoftok)
}

// returns current token and moves forward

func (p *parser) NextToken() expression2.SyntaxToken {
	current := p.current()
	p.position++
	return current
}

//

func (p *parser) MatchToken(tok token.Token) expression2.SyntaxToken {
	cur := p.current()
	if cur.Kind() == tok {
		return p.NextToken()
	}
	wanted := fmt.Sprint(tok)
	got := fmt.Sprint(cur.Kind())
	span := p.current().Span
	if span.End >= p.position {
		span.End -= 1
	}
	p.Diag.UnexpectedToken(span, wanted, got)
	return *(expression2.NewSyntaxToken(p.current().Kind(), p.position, "", nil))
}

func (pars *parser) current() expression2.SyntaxToken { return pars.peek(0) }

func (pars *parser) peek(offset int) expression2.SyntaxToken {
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
