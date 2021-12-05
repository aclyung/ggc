package test

import (
	"almeng.com/glang/ast"
	assert "almeng.com/glang/general/Assert"
	"almeng.com/glang/token"
	"fmt"
	"testing"
)

func Token_Text(t token.Token) {
	text := fmt.Sprint(t)
	if t.IsLiteral() {
		return
	}
	if len(text) == 0 {
		return
	}

	toks := ast.ParseTokens(text)
	assert.Single(toks)
	tok := toks[0]

	assert.Equal(t, tok.Kind())
	assert.Equal(text, tok.Text)
}

func TestToken_Name(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		for _, v := range tokenData {
			Token_Text(v.kind)
		}
	})
}
