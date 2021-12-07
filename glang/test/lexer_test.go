package test

import (
	"testing"

	"almeng.com/glang/ast"
	assert "almeng.com/glang/general/Assert"
	"almeng.com/glang/token"
)

//IDENT, "IDENT",
//INT,   "INT",
//FLOAT, "FLOAT",
//CHAR,   "CHAR",
//STRING, "STRING",
//BOOL,   "BOOL",
//ADD, "+",
//SUB, "-",
//MUL, "*",
//QUO, "/",
//REM, "%",
//AND, "&",
//OR,  "|",
//ADD_ASSIGN, "+=",
//SUB_ASSIGN, "-=",
//MUL_ASSIGN, "*=",
//QUO_ASSIGN, "/=",
//REM_ASSIGN, "%=",
//LAND, "&&",
//LOR,  "||",
//INC, "++",
//DEC, "--",
//EQL,    "==",
//LSS,    "<",
//GTR,    ">",
//ASSIGN, "=",
//NOT,    "!",
//NEQ,    "!=",
//LEQ,    "<=",
//GEQ,    ">=",
//DEFINE, ",=",
//LPAREN, "(",
//LBRACK, "[",
//LBRACE, "{",
//COMMA,  ",",
//PERIOD, ".",
//RPAREN,    ")",
//RBRACK,    "]",
//RBRACE,    "}",
//SEMICOLON, ";",
//COLON,     ",",

var tokenData = []struct {
	kind token.Token
	text string
}{
	{token.IDENT, "a"},
	{token.INT, "5342"},
	{token.FLOAT, "1323.12"},
	{token.CHAR, "'a'"},
	{token.STRING, "\"Heelooo\""},
	{token.BOOL, "true"},
	{token.ADD, "+"},
	{token.SUB, "-"},
	{token.MUL, "*"},
	{token.QUO, "/"},
	{token.REM, "%"},
	{token.AND, "&"},
	{token.OR, "|"},
	{token.ADD_ASSIGN, "+="},
	{token.SUB_ASSIGN, "-="},
	{token.MUL_ASSIGN, "*="},
	{token.QUO_ASSIGN, "/="},
	{token.REM_ASSIGN, "%="},
	{token.LAND, "&&"},
	{token.LOR, "||"},
	{token.INC, "++"},
	{token.DEC, "--"},
	{token.EQL, "=="},
	{token.LSS, "<"},
	{token.GTR, ">"},
	{token.ASSIGN, "="},
	{token.NOT, "!"},
	{token.NEQ, "!="},
	{token.LEQ, "<="},
	{token.GEQ, ">="},
	{token.DEFINE, ":="},
	{token.LPAREN, "("},
	{token.LBRACK, "["},
	{token.LBRACE, "{"},
	{token.COMMA, ","},
	{token.PERIOD, "."},
	{token.RPAREN, ")"},
	{token.RBRACK, "]"},
	{token.RBRACE, "}"},
	{token.SEMICOLON, ";"},
	{token.COLON, ":"},
}

func Lexer_Lex(kind token.Token, text string) {
	toks := ast.ParseTokens(text)
	assert.Single(toks)

	tok := toks[0]
	assert.Equal(kind, tok.Kind())
	assert.Equal(text, tok.Text)
}

func reqSeparator(kind1 token.Token, kind2 token.Token) bool {
	if kind1 == token.IDENT && kind2 == token.IDENT {
		return true
	}
	return false
}

func TestLexer(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		for _, v := range tokenData {
			Lexer_Lex(v.kind, v.text)
		}
	})
}
