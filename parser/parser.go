package parser

import (
	"almeng.com/glang/parser/node"
	"bufio"
	"fmt"
	"os"

	"almeng.com/glang/lexer"
	"almeng.com/glang/token"
)

type parser struct {
	tokens   []lexer.SyntaxToken
	position int
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
		if tok.Token == token.EOF {
			break
		}
		if tok.Token != token.ILLEGAL {
			pars.tokens = append(pars.tokens, *tok)
		}

	}
	fmt.Println(pars.tokens)
	fmt.Println(pars)
	return pars
}

func (p *parser ) NextToken() lexer.SyntaxToken {
	current := p.current()
	p.position++
	return current
}

func (p *parser) Match(tok token.Token) lexer.SyntaxToken {
	if p.current().Kind() == tok {
		return p.NextToken()
	}
	return *(lexer.Token(tok, p.position, "", nil))
}

func (p *parser) Parse() node.ExpressionSyntax {
	left:= p.ParsePrevExpression()
	for cur:= p.current().Kind(); cur ==token.ADD|| cur== token.SUB; {
		opTok := p.NextToken()
		if cur == opTok.Kind() {

		}
		right := p.ParsePrevExpression()
		left = BinaryExpressionSyntax(left,opTok,right)
		cur= p.current().Kind()
	}
	return left
}

func (p *parser) ParsePrevExpression() node.ExpressionSyntax {
	numTok := p.Match(token.INT)
	return NumberExpressionSyntax(numTok)
}

func (pars *parser) current() lexer.SyntaxToken { return pars.peek(0) }

func (pars *parser) peek(offset int) lexer.SyntaxToken {
	index := pars.position + offset
	if index >= len(pars.tokens) {
		return pars.tokens[len(pars.tokens)-1]
	}

	return pars.tokens[index]
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
