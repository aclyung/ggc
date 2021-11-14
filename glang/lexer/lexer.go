package lexer

import (
	"almeng.com/glang/glang/general"
	node2 "almeng.com/glang/glang/parser/node"
	"almeng.com/glang/glang/token"
	"strconv"
	"unicode"
)

type SyntaxToken struct {
	Token    token.Token
	Position int
	Text     string
	Value    interface{}
}
type Lexer struct {
	text        string
	position    int
	Diagnostics []general.Diag
}

const EOFCHAR = '\000'

func (tok SyntaxToken) GetChildren() []node2.ExpressionSyntax {
	return []node2.ExpressionSyntax{}
}

func (tok SyntaxToken) Type() node2.Type {
	if tok.Token.IsOperator() {
		return node2.Type(tok.Token-token.Operator_beg) + node2.Oper_beg
	}
	return node2.SyntaxToken
}

func (tok SyntaxToken) Kind() token.Token {
	return tok.Token
}

func NewLexer(text string) Lexer {
	return Lexer{text: text, position: 0, Diagnostics: make([]general.Diag, 0)}
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
		tok := token.INT
		for unicode.IsDigit(lexer.current()) {
			lexer.next()
		}
		if lexer.current() == '.' {
			lexer.next()
			tok = token.FLOAT
			for unicode.IsDigit(lexer.current()) {
				lexer.next()
			}
		}
		text := lexer.text[beg:lexer.position]
		if tok == token.INT {
			value, err := strconv.ParseInt(text, 10, 64)
			if err != nil {
				lexer.Diagnose("ERROR: the Number is not Valid int64.", general.ERROR)
			}
			return Token(tok, beg, text, value)
		}
		value, err := strconv.ParseFloat(text, 8)
		if err != nil {
			lexer.Diagnose("ERROR: the Number is not Valid float64.", general.ERROR)
		}
		return Token(tok, beg, text, value)
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

	lexer.Diagnose("ERROR: Illegal Charater '"+string(cur)+"'", general.ERROR)
	return Token(token.ILLEGAL, pos, lexer.text[pos:pos+1], nil)

}
func (lex *Lexer) Diagnose(text string, l general.Level) {
	diag := general.Diag{text, l}
	lex.Diagnostics = append(lex.Diagnostics, diag)
}
func Token(token token.Token, position int, text string, value interface{}) *SyntaxToken {

	return &SyntaxToken{
		Token:    token,
		Position: position,
		Text:     text,
		Value:    value,
	}
}
