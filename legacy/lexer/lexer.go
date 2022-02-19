package lexer

import (
	"strconv"
	"unicode"

	"main/legacy/expression"
	"main/legacy/general"
	Text2 "main/legacy/general/Text"
	"main/legacy/token"
)

type Lexer struct {
	text      *Text2.Source
	position  int
	Diag      general.Diags
	start     int
	tokenKind token.Token
	value     interface{}
	//Diagnostics []general.Diag
}

const EOFCHAR = '\000'

func NewLexer(source *Text2.Source) Lexer {
	//text := source.Text
	return Lexer{text: source, position: 0, Diag: general.NewDiag()} // ,Diagnostics: make([]general.Diag, 0)}
}

func (lex *Lexer) current() rune {
	if lex.position >= lex.text.Length() {
		return EOFCHAR
	}
	return rune(lex.text.At(lex.position))
}

func (lex *Lexer) next() {
	lex.position++
}

func IllegalSyntax(position int) *expression.SyntaxToken {
	return expression.NewSyntaxToken(token.ILLEGAL, position, "", nil)
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

	} else if unicode.IsLetter(lex.current()) || lex.current() == '_' {

		lex.ReadLetter()

	} else {
		switch lex.current() {
		case EOFCHAR:
			return expression.NewSyntaxToken(token.EOF, lex.position, string(EOFCHAR), nil)
		case '\'':
			return lex.ReadChar()
		case '"':
			return lex.ReadString()
		default:
			lex.Diag.BadCharacter(Text2.Span(lex.start, lex.position), string(lex.current()))
			lex.position++
		}
	}
	t := lex.text.ToString(lex.start, lex.position)
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
	text := lex.text.ToString(lex.start, lex.position)
	if lex.tokenKind == token.INT {
		val, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			lex.Diag.InvalidNumber(Text2.Span(lex.start, lex.position), text, "int64")
			//lex.Diagnose("ERROR: the Number is not Valid int64.", general.ERROR)
		}
		lex.value = val
		return
	}
	val, err := strconv.ParseFloat(text, 8)
	if err != nil {
		lex.Diag.InvalidNumber(Text2.Span(lex.start, lex.position), text, "float64")
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
	for unicode.IsLetter(lex.current()) || lex.current() == '_' {
		lex.next()
	}
	text := lex.text.ToString(lex.start, lex.position)
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

func (lex *Lexer) ReadString() *expression.SyntaxToken {
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
			lex.Diag.BadCharacter(Text2.Span(lex.start, lex.position-1), "Unexpectedly faced EOF")
			return IllegalSyntax(lex.start)
		}
	}
	text := lex.text.ToString(lex.start, lex.position)
	lex.value = text[1 : len(text)-1]
	lex.tokenKind = token.STRING
	return expression.NewSyntaxToken(token.STRING, lex.start, text, lex.value)
}

func (lex *Lexer) ReadChar() *expression.SyntaxToken {
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
			lex.Diag.BadCharacter(Text2.Span(lex.start, lex.position-1), "Unexpectedly faced EOF")
			return IllegalSyntax(lex.start)
		}
	}

	text := lex.text.ToString(lex.start, lex.position)
	if len(text) == 2 || len(text) > 3 {
		lex.tokenKind = token.ILLEGAL
		lex.Diag.BadCharacter(Text2.Span(lex.start, lex.position), "Illegal rune literal")
	}
	value := int(lex.text.At(lex.start + 1))
	lex.value = value
	lex.tokenKind = token.CHAR
	return expression.NewSyntaxToken(lex.tokenKind, lex.start, text, lex.value)
}
