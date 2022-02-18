package syntax

import (
	"unicode"
)

type any = interface{}

type Lexer struct {
	text  string
	pos   int
	token Token
	start int
	value any
}

const EOFCHAR = '\000'

// NewLexer returns a new lexer.
func NewLexer(source string) Lexer {
	return Lexer{text: source, pos: 0}
}

func (lex Lexer) current() rune {
	if lex.pos >= len(lex.text) {
		return EOFCHAR
	}
	return rune(lex.text[lex.pos])
}

// lookahead returns the next rune in the text.
func (lex *Lexer) lookahead() rune {
	if lex.pos+1 >= len(lex.text) {
		return EOFCHAR
	}
	return rune(lex.text[lex.pos+1])
}

// advance moves the lexer forward one rune.
func (lex *Lexer) next() {
	lex.pos++
}

// skipWhitespace advances the lexer to the next non-whitespace rune.
func (lex *Lexer) skipWhitespace() {
	for unicode.IsSpace(lex.current()) {
		lex.next()
	}
}

func (lex *Lexer) Lex() {
	lex.skipWhitespace()
	lex.start = lex.pos
	lex.token = ILLEGAL
	lex.value = nil

	switch lex.current() {
	case EOFCHAR:
		lex.token = EOF
	}
}

func lower(ch rune) rune     { return ('a' - 'A') | ch } // returns lower-case ch iff ch is ASCII letter
func isLetter(ch rune) bool  { return 'a' <= lower(ch) && lower(ch) <= 'z' || ch == '_' }
func isDecimal(ch rune) bool { return '0' <= ch && ch <= '9' }
