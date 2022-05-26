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
	SpaceName *Name
	DeclList  []Decl
	EOF       Pos
	node
}

// Top Level Declarations
type (
	Decl interface {
		Node
		aDecl()
	}

	OperDecl struct {
		Group        *Group
		TypeL, TypeR *Field
		Oper         token
		Return       Expr
		Body         *BlockStmt
		decl
	}

	VarDecl struct {
		Group    *Group // nil means not part of a group
		NameList *Name
		Type     Expr // nil means no type
		Values   Expr // nil means no values
		decl
	}

	FuncDecl struct {
		Group  *Group // nil means not part of a group
		Param  []*Field
		Name   *Name // identifier
		Return Expr  // nil means no return type
		Body   *BlockStmt
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
	Stmt interface {
		Node
		aStmt()
	}

	ExprStmt struct {
		X Expr
		simpleStmt
	}

	EmptyStmt struct {
		simpleStmt
	}

	IncDecStmt struct {
		X   Expr
		Tok token
		simpleStmt
	}

	ContinueStmt struct {
		simpleStmt
	}

	BreakStmt struct {
		simpleStmt
	}

	ReturnStmt struct {
		Return Expr
		stmt
	}

	DeclStmt struct {
		DeclList []Decl
		stmt
	}

	AssignStmt struct {
		Lhs Expr
		Op  Operator
		Rhs Expr
		stmt
	}

	IfStmt struct {
		Cond  Expr
		Block *BlockStmt
		Else  Stmt
		stmt
	}

	ForStmt struct {
		Init simpleStmt
		Cond Expr
		Post simpleStmt
		Body *BlockStmt
		stmt
	}

	simpleStmt struct {
		stmt
	}

	BlockStmt struct {
		StmtList []Stmt
		Rbrace   Pos
		stmt
	}
)

type stmt struct{ node }

func (*stmt) aStmt() {}

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

	// X.Sel
	SelectorExpr struct {
		X   Expr
		Sel *Name
		expr
	}
	// Fun(ArgList[0], ArgList[1], ...)
	CallExpr struct {
		Fun     Expr
		ArgList []Expr // nil means no arguments
		expr
	}

	Field struct {
		Name *Name // nil means anonymous field/parameter (structs/parameters), or embedded element (interfaces)
		Type Expr  // field names declared in a list share the same Type (identical pointers)
		node
	}
)

type expr struct{ node }

func (*expr) aExpr() {}

type Group struct {
	_ int // not empty so we are guaranteed different Group instances
}
