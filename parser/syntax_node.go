package parser

import (
	"almeng.com/glang/lexer"
	"almeng.com/glang/parser/node"
	"almeng.com/glang/token"
)

type Node = node.Node


type SyntaxNode struct {
	token token.Token
}

// type ExpressionSyntax struct {
// 	SyntaxNode
// }

// Node interface implementations

func GetChildren(syntax node.ExpressionSyntax) []node.ExpressionSyntax{

	return []node.ExpressionSyntax{syntax}
}

func hi(n Node) interface{}{
	a:= make([]Node,0)
	return a
}

//func (nod SyntaxNode) Kind() token.Token {
//	return nod.token
//}

type NnumberExpressionSyntax struct {
	child_left int
	NumberToken lexer.SyntaxToken
}

func NumberExpressionSyntax(numberToken lexer.SyntaxToken) *NnumberExpressionSyntax {
	syntax := &NnumberExpressionSyntax{1, numberToken}
	return syntax
}

func (numSyntax *NnumberExpressionSyntax) Kind() token.Token{
	return token.INT
}

func (numSyntax NnumberExpressionSyntax) Type() node.Type {
	return node.ExpNum
}

func (numSyntax *NnumberExpressionSyntax) GetChildren() []node.ExpressionSyntax{
	return []node.ExpressionSyntax{numSyntax.NumberToken}
}



type binaryExpressionSyntax struct {
	OperatorToken lexer.SyntaxToken
	Left node.ExpressionSyntax
	Right node.ExpressionSyntax
}

func (b binaryExpressionSyntax) Kind() token.Token {
	return b.OperatorToken.Token
}
func (b binaryExpressionSyntax) Type() node.Type {
	return node.ExpBi
}
func (b binaryExpressionSyntax) GetChildren() []node.ExpressionSyntax {
	return []node.ExpressionSyntax{b.OperatorToken, b.Left, b.Right}
}

func BinaryExpressionSyntax(left node.ExpressionSyntax, opToken lexer.SyntaxToken, right node.ExpressionSyntax ) *binaryExpressionSyntax {
	bi := &binaryExpressionSyntax{
		opToken,
		left,
		right}
	return bi
}