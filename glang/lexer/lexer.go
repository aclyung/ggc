package lexer

import (
	"almeng.com/glang/general/TextSpan"
	"strconv"
	"unicode"

	"almeng.com/glang/expression"
	"almeng.com/glang/general"
	"almeng.com/glang/token"
)

type Lexer struct {
	text     string
	position int
	Diag     general.Diags
	//Diagnostics []general.Diag
}

const EOFCHAR = '\000'

func NewLexer(text string) Lexer {
	return Lexer{text: text, position: 0, Diag: general.NewDiag()} // ,Diagnostics: make([]general.Diag, 0)}
}

func (lex *Lexer) current() rune {
	if lex.position >= len(lex.text) {
		return EOFCHAR
	}
	return rune(lex.text[lex.position])
}

func (lex *Lexer) next() {
	lex.position++
}

func (lex *Lexer) Lex() *expression.SyntaxToken {
	if lex.position >= len(lex.text) {
		return Token(token.EOF, lex.position, string(EOFCHAR), nil)
	}
	if lex.current() == '"' {
		beg := lex.position
		for {
			lex.next()
			if lex.current() == '"' {
				lex.next()
				break
			}
			if lex.current() == EOFCHAR {
				lex.next()
				return Token(token.ILLEGAL, lex.position-1, lex.text[beg:lex.position-1], general.Err("Unexpectedly face EOF"))
			}
		}
		text := lex.text[beg:lex.position]
		value := text[1 : len(text)-1]
		return Token(token.STRING, beg, text, value)
	}

	if lex.current() == '\'' {
		beg := lex.position
		for {
			lex.next()
			cur := lex.current()
			if cur == '\'' {
				lex.next()
				break
			}
			if cur == EOFCHAR {
				lex.next()
				return Token(token.ILLEGAL, lex.position-1, lex.text[beg:lex.position-1], general.Err("Unexpectedly face EOF"))
			}
		}

		text := lex.text[beg:lex.position]
		if len(text) == 2 || len(text) > 3 {
			return Token(token.ILLEGAL, beg, text, general.Err("Illegal rune literal"))
		}
		value := int(rune(lex.text[beg+1]))
		return Token(token.INT, beg, text, value)

	}

	if unicode.IsDigit(lex.current()) {
		beg := lex.position
		tok := token.INT
		for unicode.IsDigit(lex.current()) {
			lex.next()
		}
		if lex.current() == '.' {
			lex.next()
			tok = token.FLOAT
			for unicode.IsDigit(lex.current()) {
				lex.next()
			}
		}
		text := lex.text[beg:lex.position]
		if tok == token.INT {
			value, err := strconv.ParseInt(text, 10, 64)
			if err != nil {
				lex.Diag.InvalidNumber(TextSpan.Span(beg, lex.position), text, "int64")
				//lex.Diagnose("ERROR: the Number is not Valid int64.", general.ERROR)
			}
			return Token(tok, beg, text, value)
		}
		value, err := strconv.ParseFloat(text, 8)
		if err != nil {
			lex.Diag.InvalidNumber(TextSpan.Span(beg, lex.position), text, "float64")
			//lex.Diagnose("ERROR: the Number is not Valid float64.", general.ERROR)
		}
		return Token(tok, beg, text, value)
	}
	if unicode.IsSpace(lex.current()) {
		for unicode.IsSpace(lex.current()) {
			lex.next()
		}
		return nil
	}
	if unicode.IsLetter(lex.current()) {
		beg := lex.position
		for unicode.IsLetter(lex.current()) {
			lex.next()
		}
		text := lex.text[beg:lex.position]
		isBool := text == "true" || text == "false"
		if isBool {
			return Token(token.BOOL, beg, text, text == "true")
		}
		return Token(token.IDENT, beg, text, nil)

	}
	cur := string(lex.current())
	pos := lex.position

	lex.next()
	tok := token.LookOperUp(cur)
	if tok.IsOperator() {
		ncur := string(lex.current())
		ntok := token.LookOperUp(ncur)
		op := cur + ncur
		optok := token.LookOperUp(op)
		if ntok.IsOperator() && optok.IsOperator() {
			lex.next()
			return Token(optok, pos, op, nil)
		}
		return Token(tok, pos, cur, nil)
	}

	lex.Diag.BadCharacter(TextSpan.Span(pos, lex.position), cur)
	return Token(token.ILLEGAL, pos, lex.text[pos:pos+1], nil)

}

//func (lex *Lexer) Diagnose(text string, l general.Level) {
//	diag := general.Diag{text, l}
//	lex.Diagnostics = append(lex.Diagnostics, diag)
//}
func Token(token token.Token, position int, text string, value interface{}) *expression.SyntaxToken {

	return expression.NewSyntaxToken(token, position, text, value)
}
