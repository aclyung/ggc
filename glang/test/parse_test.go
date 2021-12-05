package test

import (
	"almeng.com/glang/ast"
	"almeng.com/glang/expression"
	assert "almeng.com/glang/general/Assert"
	"almeng.com/glang/general/Stack"
	"almeng.com/glang/syntax"
	"almeng.com/glang/token"
	"fmt"
	"log"
	"reflect"
	"testing"
)

var binary = []token.Token{
	token.ADD,
	token.SUB,
	token.MUL,
	token.QUO,
	token.REM,
	token.AND,
	token.OR,
	token.LAND,
	token.LOR,
	token.EQL,
	token.LSS,
	token.GTR,
	token.NEQ,
	token.LEQ,
	token.GEQ,
}

var unary = []token.Token{
	token.ADD,
	token.SUB,
}

type AssertTree struct {
	Tree      []syntax.ExpressionSyntax
	position  int
	Current   syntax.ExpressionSyntax
	hasErrors bool
}

func NewAssertTree(node syntax.ExpressionSyntax) *AssertTree {
	s := Flat(node)
	return &AssertTree{s, 0, s[0], false}
}

func Flat(node syntax.ExpressionSyntax) []syntax.ExpressionSyntax {
	stack := Stack.Stack(reflect.TypeOf(reflect.ValueOf(node)))
	stack.Push(node)
	rtn := make([]syntax.ExpressionSyntax, 0)
	for stack.Size() > 0 {
		n := stack.Pop().(syntax.ExpressionSyntax)
		fmt.Println(n.Type())
		rtn = append(rtn, n)
		s := n.GetChildren()
		for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
			s[i], s[j] = s[j], s[i]
		}

		for _, v := range s {
			stack.Push(v)
		}
	}
	return rtn
	//a := make([]syntax.ExpressionSyntax, 0)
	//a = append(a, node)
	//
	//index := len(a)-1
	//for  {
	//	if index < 0 {
	//		break
	//	}
	//	v:= a[index]
	//	index--
	//	s := v.GetChildren()
	//	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
	//		s[i], s[j] = s[j], s[i]
	//	}
	//
	//	for _, c := range s {
	//		a = append(a, c)
	//	}
	//}
	//fmt.Println(a)
	//return a
}

func (t *AssertTree) MarkFailed() bool {
	t.hasErrors = true
	return false
}

func (t *AssertTree) Status(superior_err interface{}) {
	defer func() {
		defer func() {
			recover()
		}()
		err := recover()
		if err != nil {
			log.Panic(superior_err)
		}
	}()
	if t.hasErrors {
		assert.False(t.Next())
	}
}

func (t *AssertTree) AssertToken(kind token.Token, text string) {
	defer func() {
		err := recover()
		if err != nil {
			t.hasErrors = true
		}
		t.Status(err)
	}()
	assert.True(t.Next())
	assert.Equal(t.Current.Type(), syntax.Token)
	tok := t.Current.(expression.SyntaxToken)
	assert.Equal(kind, tok.Kind())
	assert.Equal(text, tok.Text)
}

func (t *AssertTree) AssertNode(kind syntax.Type) {
	defer func() {
		err := recover()
		if err != nil {
			t.hasErrors = true
		}
		t.Status(err)
	}()
	assert.True(t.Next())
	assert.Equal(kind, t.Current.Type())
	assert.NotEqual(t.Current.Type(), syntax.Token)
}

func (t *AssertTree) Next() bool {
	if t.position >= len(t.Tree)-1 {
		return false
	}
	t.Current = t.Tree[t.position]
	t.position++

	return true
}

func Opers() [][]token.Token {
	opers := make([][]token.Token, 0)
	for _, v := range binary {
		for _, k := range unary {
			opers = append(opers, []token.Token{v, k})
		}
	}
	return opers
}

func Parser_binary_operator_precedence(op1 token.Token, op2 token.Token) {
	p1, p2 := op1.Precedence(), op2.Precedence()
	text1, text2 := op1.String(), op2.String()
	prompt := fmt.Sprintf("a %v b %v c", text1, text2)
	exp := ast.ParseTree(prompt).Root

	fmt.Println(prompt)
	if p1 >= p2 {
		e := NewAssertTree(exp)
		e.AssertNode(syntax.ExpBinary)
		e.AssertNode(syntax.ExpBinary)
		e.AssertNode(syntax.ExpName)
		e.AssertToken(token.IDENT, "a")
		e.AssertToken(op1, text1)
		e.AssertNode(syntax.ExpName)
		e.AssertToken(token.IDENT, "b")
		e.AssertToken(op2, text2)
		e.AssertNode(syntax.ExpName)
		e.AssertToken(token.IDENT, "c")

		// 	      op1
		//      /     \
		// 	  op2      c
		//  /     \
		// a	  b
	} else {
		// 	  op1
		//  /     \
		//a 	 op2
		//	   /     \
		//	  b		  c
		//assert.False(true)
		e := NewAssertTree(exp)
		e.AssertNode(syntax.ExpBinary)
		e.AssertNode(syntax.ExpName)
		e.AssertToken(token.IDENT, "a")
		e.AssertToken(op1, text1)
		e.AssertNode(syntax.ExpBinary)
		e.AssertNode(syntax.ExpName)
		e.AssertToken(token.IDENT, "b")
		e.AssertToken(op2, text2)
		e.AssertNode(syntax.ExpName)
		e.AssertToken(token.INT, "c")
	}
}

func Parser_Unary_operator_precedence(una token.Token, bi token.Token) {
	p1, p2 := una.Precedence(), bi.Precedence()
	text1, text2 := una.String(), bi.String()
	prompt := fmt.Sprintf("%v a %v b", text1, text2)
	exp := ast.ParseTree(prompt).Root

	fmt.Println(prompt)
	if p1 >= p2 {
		e := NewAssertTree(exp)
		e.AssertNode(syntax.ExpBinary)
		e.AssertNode(syntax.ExpUnary)
		e.AssertToken(una, text1)
		e.AssertNode(syntax.ExpName)
		e.AssertToken(token.IDENT, "a")
		e.AssertToken(bi, text2)
		e.AssertNode(syntax.ExpName)
		e.AssertToken(token.IDENT, "b")

	} else {
		// 	  op1
		//  /     \
		//a 	 op2
		//	   /     \
		//	  b		  c
		//assert.False(true)
		e := NewAssertTree(exp)
		e.AssertNode(syntax.ExpUnary)
		e.AssertToken(una, text1)
		e.AssertNode(syntax.ExpBinary)
		e.AssertNode(syntax.ExpName)
		e.AssertToken(token.IDENT, "a")
		e.AssertToken(bi, text2)
		e.AssertNode(syntax.ExpName)
		e.AssertToken(token.IDENT, "b")
	}
}

func TestParseBiPrecedence(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		opers := Opers()
		for _, v := range opers {
			Parser_binary_operator_precedence(v[0], v[1])
		}
	})
}

func TestParseUPrecedence(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		opers := Opers()
		for _, v := range opers {
			Parser_Unary_operator_precedence(v[1], v[0])
		}
	})
}
