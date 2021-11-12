package lexer

import (
	_ "fmt"
	"strconv"
	"unicode"

	"almeng.com/glang/general"
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

func NewLexer(text string) *Lexer {
	return &Lexer{text: text, position: 0}
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

	cur := lexer.current()
	pos := lexer.position
	
	lexer.next()
	switch cur {
	case '+':
		return Token(token.ADD, pos, "+", nil)
	case '-':
		return Token(token.SUB, pos, "-", nil)
	case '*':
		return Token(token.MUL, pos, "*", nil)
	case '/':
		return Token(token.QUO, pos, "/", nil)
	case '(':
		return Token(token.LPAREN, pos, "(", nil)
	case ')':
		return Token(token.RPAREN, pos, ")", nil)
	default:
		return Token(token.ILLEGAL, pos, lexer.text[pos:pos+1], nil)
	}
}

func Token(token token.Token, position int, text string, value interface{}) *SyntaxToken {

	return &SyntaxToken{
		Token:    token,
		Position: position,
		Text:     text,
		Value:    value,
	}
}
