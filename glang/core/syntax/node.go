package syntax

type Node interface {
	Pos() Pos
	aNode()
}

type node struct {
	pos Pos
}

func (n *node) Pos() Pos { return n.pos }
func (*node) aNode()     {}

type File struct {
	//PkgName
	DeclList []Decl
	EOF      Pos
	node
}

type (
	Decl interface {
		Node
		aDecl()
	}

	OperDecl struct {
		Group                *Group
		NameList             []*Name
		Return, TypeL, TypeR Expr
		decl
	}

	VarDecl struct {
		Group    *Group // nil means not part of a group
		NameList *Name
		//Type     Expr // nil means no type
		Values Expr // nil means no values
		decl
	}
)

type decl struct{ node }

func (*decl) aDecl() {}

func NewName(pos Pos, value string) *Name {
	n := new(Name)
	n.pos = pos
	n.Value = value
	return n
}

type (
	Expr interface {
		Node
		aExpr()
	}

	// Placeholder for an expression that failed to parse
	// correctly and where we can't provide a better node.
	BadExpr struct {
		expr
	}

	// Value
	Name struct {
		Value string
		expr
	}

	// Value
	BasicLit struct {
		Value string
		Kind  LitKind
		Bad   bool // true means the gotLiteral Value has syntax errors
		expr
	}

	Operation struct {
		Op   Operator
		X, Y Expr // Y == nil means unary expression
		expr
	}

	ParenExpr struct {
		X Expr
		expr
	}
)

type expr struct{ node }

func (*expr) aExpr() {}

type Group struct {
	_ int // not empty so we are guaranteed different Group instances
}
