package expression

import (
	"fmt"

	"main/legacy/general"
	syntax2 "main/legacy/syntax"
	"main/legacy/token"
)

type Syntax struct {
	kind     token.Token
	ExpType  syntax2.Type
	children []syntax2.ExpressionSyntax
}

func (e *Syntax) String() string {
	return sprintTree(e, "", true)
}
func sprintTree(nod syntax2.ExpressionSyntax, indent string, isLast bool) string {
	str := ""
	var mark string

	if isLast {
		mark = "└────"
	} else {
		mark = "├────"
	}
	str += indent + mark
	nodeType := nod.Type()

	if nodeType == syntax2.Token {
		str += general.ColorText(nodeType.String()+" ", general.Blue)
	} else {
		str += general.ColorText(nodeType.String()+" ", general.Yellow)
	}
	if nodeType == syntax2.Token {
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

func (e Syntax) GetChildren() []syntax2.ExpressionSyntax {
	return e.children
}

func (e Syntax) Type() syntax2.Type {
	return e.ExpType
}

func (e Syntax) Kind() token.Token {
	return e.kind
}

func NewSyntax(kind token.Token, t syntax2.Type, c ...syntax2.ExpressionSyntax) Syntax {
	return Syntax{kind, t, c}
}
