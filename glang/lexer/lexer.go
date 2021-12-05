package lexer

import (
	"almeng.com/glang/expression"
	"almeng.com/glang/general/TextSpan"
	"strconv"
	"unicode"

	"almeng.com/glang/general"
	"almeng.com/glang/token"
)

type Lexer struct {
	text      string
	position  int
	Diag      general.Diags
	start     int
	tokenKind token.Token
	value     interface{}
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
	lex.start = lex.position
	lex.tokenKind = token.ILLEGAL
	lex.value = nil

	if unicode.IsSpace(lex.current()) {

		for unicode.IsSpace(lex.current()) {
			lex.next()
		}
		return nil

	} else if unicode.IsDigit(lex.current()) {

		lex.ReadNumberToken()

	} else if token.IsOper(string(lex.current())) {

		lex.ReadOperatorToken()

	} else if unicode.IsLetter(lex.current()) {

		lex.ReadLetter()

	} else {
		switch lex.current() {
		case EOFCHAR:
			return expression.NewSyntaxToken(token.EOF, lex.position, string(EOFCHAR), nil)
		case '\'':
			lex.ReadChar()
		case '"':
			lex.ReadString()
		default:
			lex.Diag.BadCharacter(TextSpan.Span(lex.start, lex.position), string(lex.current()))
			lex.position++
		}
	}
	t := lex.text[lex.start:lex.position]
	return expression.NewSyntaxToken(lex.tokenKind, lex.start, t, lex.value)

}

func (lex *Lexer) ReadNumberToken() {
	lex.tokenKind = token.INT
	for unicode.IsDigit(lex.current()) {
		lex.next()
	}
	if lex.current() == '.' {
		lex.next()
		lex.tokenKind = token.FLOAT
		for unicode.IsDigit(lex.current()) {
			lex.next()
		}
	}
	text := lex.text[lex.start:lex.position]
	if lex.tokenKind == token.INT {
		val, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			lex.Diag.InvalidNumber(TextSpan.Span(lex.start, lex.position), text, "int64")
			//lex.Diagnose("ERROR: the Number is not Valid int64.", general.ERROR)
		}
		lex.value = val
		return
	}
	val, err := strconv.ParseFloat(text, 8)
	if err != nil {
		lex.Diag.InvalidNumber(TextSpan.Span(lex.start, lex.position), text, "float64")
		//lex.Diagnose("ERROR: the Number is not Valid float64.", general.ERROR)
	}
	lex.value = val
}

func (lex *Lexer) ReadOperatorToken() {
	cur := string(lex.current())
	lex.next()
	nchar := string(lex.current())
	op := cur + nchar
	lex.value = nil
	if token.IsOper(nchar) && token.IsOper(op) {
		lex.next()
		lex.tokenKind = token.LookOperUp(op)
		return
	}
	lex.tokenKind = token.LookOperUp(cur)
	return

}

func (lex *Lexer) ReadLetter() {
	for unicode.IsLetter(lex.current()) {
		lex.next()
	}
	text := lex.text[lex.start:lex.position]
	isBool := text == "true" || text == "false"
	if isBool {
		lex.tokenKind = token.BOOL
		lex.value = text == "true"
		return
	}
	lex.tokenKind = token.IDENT
	lex.value = nil
	return
}

func (lex *Lexer) ReadString() {
	for {
		lex.next()
		if lex.current() == '"' {
			lex.next()
			break
		}
		if lex.current() == EOFCHAR {
			lex.next()
			lex.tokenKind = token.ILLEGAL
			lex.value = nil
			lex.Diag.BadCharacter(TextSpan.Span(lex.start, lex.position), "Unexpectedly faced EOF")
		}
	}
	text := lex.text[lex.start:lex.position]
	lex.value = text[1 : len(text)-1]
	lex.tokenKind = token.STRING
}

func (lex *Lexer) ReadChar() {
	for {
		lex.next()
		cur := lex.current()
		if cur == '\'' {
			lex.next()
			break
		}
		if cur == EOFCHAR {
			lex.next()
			lex.tokenKind = token.ILLEGAL
			lex.Diag.BadCharacter(TextSpan.Span(lex.start, lex.position), "Unexpectedly faced EOF")
		}
	}

	text := lex.text[lex.start:lex.position]
	if len(text) == 2 || len(text) > 3 {
		lex.tokenKind = token.ILLEGAL
		lex.Diag.BadCharacter(TextSpan.Span(lex.start, lex.position), "Illegal rune literal")
	}
	value := int(rune(lex.text[lex.start+1]))
	lex.value = value
	lex.tokenKind = token.CHAR

}
