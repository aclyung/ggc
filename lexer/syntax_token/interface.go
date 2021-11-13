package interfaces

import "almeng.com/glang/token"

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
