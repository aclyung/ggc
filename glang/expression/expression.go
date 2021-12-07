package expression

import (
	"almeng.com/glang/syntax"
	"almeng.com/glang/token"
	"fmt"
)

type Syntax struct {
	kind     token.Token
	ExpType  syntax.Type
	children []syntax.ExpressionSyntax
}

func (e *Syntax) String() string {
	return sprintTree(e, "", true)
}
func sprintTree(nod syntax.ExpressionSyntax, indent string, isLast bool) string {
	str := ""
	var mark string

	if isLast {
		mark = "└────"
	} else {
		mark = "├────"
	}
	str += indent + mark
	nodeType := nod.Type()
	str += nodeType.String() + " "
	if nodeType == syntax.Token {
		val := nod.(SyntaxToken)
		if val.Value != nil {
			str += fmt.Sprint(" ", val.Kind(), " | ", val.Value)
		} else {
			str += fmt.Sprint(" ", val.Kind())
		}
	}
	str += "\n"
	if isLast {
		indent += "     "
	} else {
		indent += "│    "
	}

	for i, v := range nod.GetChildren() {
		str += sprintTree(v, indent, i == len(nod.GetChildren())-1)
	}
	return str
}

func (e Syntax) GetChildren() []syntax.ExpressionSyntax {
	return e.children
}

func (e Syntax) Type() syntax.Type {
	return e.ExpType
}

func (e Syntax) Kind() token.Token {
	return e.kind
}

func NewSyntax(kind token.Token, t syntax.Type, c ...syntax.ExpressionSyntax) Syntax {
	return Syntax{kind, t, c}
}
