package lexer

import (
	"fmt"
	"strconv"
	"unicode"

	"almeng.com/glang/general"
	"almeng.com/glang/parser/node"
	"almeng.com/glang/token"
)

type SyntaxToken struct {
	Token    token.Token
	Position int
	Text     string
	Value    interface{}
}
type Lexer struct {
	text     string
	position int
}

const EOFCHAR = '\000'

func (tok SyntaxToken) GetChildren() []node.ExpressionSyntax {
	return []node.ExpressionSyntax{}
}

func (tok SyntaxToken) Type() node.Type {
	return node.SyntaxToken
}

func (tok SyntaxToken) Kind() token.Token{
	return tok.Token
}

func NewLexer(text string) Lexer {
	fmt.Print()
	return Lexer{text: text, position: 0}
}

func (lexer *Lexer) current() rune {
	if lexer.position >= len(lexer.text) {
		return EOFCHAR
	}
	return rune(lexer.text[lexer.position])
}

func (lexer *Lexer) next() {
	lexer.position++
}

func (lexer *Lexer) NextToken() *SyntaxToken {
	if lexer.position >= len(lexer.text) {
		return Token(token.EOF, lexer.position, string(EOFCHAR), nil)
	}
	if lexer.current() == '"' {
		beg := lexer.position
		for {
			lexer.next()
			if lexer.current() == '"' {
				lexer.next()
				break
			}
			if lexer.current() == EOFCHAR {
				lexer.next()
				return Token(token.ILLEGAL, lexer.position-1, lexer.text[beg:lexer.position-1], general.Err("Unexpectedly face EOF"))
			}
		}
		text := lexer.text[beg:lexer.position]
		value := text[1 : len(text)-1]
		return Token(token.STRING, beg, text, value)
	}

	if lexer.current() == '\'' {
		beg := lexer.position
		for {
			lexer.next()
			cur := lexer.current()
			if cur == '\'' {
				lexer.next()
				break
			}
			if cur == EOFCHAR {
				lexer.next()
				return Token(token.ILLEGAL, lexer.position-1, lexer.text[beg:lexer.position-1], general.Err("Unexpectedly face EOF"))
			}
		}

		text := lexer.text[beg:lexer.position]
		if len(text) == 2 || len(text) > 3 {
			return Token(token.ILLEGAL, beg, text, general.Err("Illegal rune literal"))
		}
		value := int(rune(lexer.text[beg+1]))
		return Token(token.INT, beg, text, value)

	}

	if unicode.IsDigit(lexer.current()) {
		beg := lexer.position
		for unicode.IsDigit(lexer.current()) {
			lexer.next()
		}
		text := lexer.text[beg:lexer.position]
		value, err := strconv.Atoi(text)
		general.ErrCheck(err)
		return Token(token.INT, beg, text, value)

	}
	if unicode.IsSpace(lexer.current()) {
		for unicode.IsSpace(lexer.current()) {
			lexer.next()
		}
		return nil
	}

	cur := string(lexer.current())
	pos := lexer.position

	lexer.next()
	tok := token.LookOperUp(cur)
	if tok.IsOperator() {
		ncur := string(lexer.current())
		ntok := token.LookOperUp(ncur)
		op := cur + ncur
		optok := token.LookOperUp(op)
		if ntok.IsOperator() && optok.IsOperator() {
			lexer.next()
			return Token(optok, pos, op, nil)
		}
		return Token(tok, pos, cur, nil)
	}
	return Token(token.ILLEGAL, pos, lexer.text[pos:pos+1], nil)
	// switch cur {
	// case '+':
	// 	return Token(token.ADD, pos, "+", nil)
	// case '-':
	// 	return Token(token.SUB, pos, "-", nil)
	// case '*':
	// 	return Token(token.MUL, pos, "*", nil)
	// case '/':
	// 	return Token(token.QUO, pos, "/", nil)
	// case '(':
	// 	return Token(token.LPAREN, pos, "(", nil)
	// case ')':
	// 	return Token(token.RPAREN, pos, ")", nil)
	// default:
	// 	return Token(token.ILLEGAL, pos, lexer.text[pos:pos+1], nil)
	// }
}

func Token(token token.Token, position int, text string, value interface{}) *SyntaxToken {

	return &SyntaxToken{
		Token:    token,
		Position: position,
		Text:     text,
		Value:    value,
	}
}
