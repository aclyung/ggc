package syntax

import (
	"fmt"
	"io"
	"strings"
)

type parser struct {
	errh ErrHandler
	lexer
	indent  string
	first   error
	verbose bool
}

// nil means error has occured
func (p *parser) EOF() *File {
	if p.verbose {
		defer p.trace("file")()
	}

	// SourceFile = Space ";" { TopLevelDecl ";" } .
	f := new(File)
	f.pos = p.pos()
	if !p.got(_Space) {
		p.error("space must be declared first")
		return nil
	}
	f.SpaceName = p.name()
	p.print("space: " + f.SpaceName.Value)
	p.want(_Semi)

	// TopLevelDecl = Declaration | FuncDecl | OperDecl .
	for p.token != _EOF {
		switch p.token {
		case _Type:
			p.next()
			f.DeclList = p.appendGroup(f.DeclList, p.typeDecl)
		case _Var:
			p.next()
			f.DeclList = p.appendGroup(f.DeclList, p.varDecl)

		case _Func:
			p.next()
			f.DeclList = p.appendGroup(f.DeclList, p.funcDeclOrNil)
		case _Oper:
			p.next()
			f.DeclList = p.appendGroup(f.DeclList, p.operDecl)
		//case _Type:
		//	p.errorAt(p.pos(), "WARNING: declaration statement not implemented yet: "+Red+p.Token.String()+Reset)
		//	p.next()
		// Throwing exception here
		case _Semi:
			p.next()
		default:
			str := p.token.String()
			if p.token == _Name {
				str += "(" + p.segment() + ")"
			}
			p.errorAt(p.pos(), "ERROR: non-declaration statement outside function body: "+Red+str+Reset)
			p.next()
		}
	}
	return f
}

func (p *parser) trace(msg string) func() {
	p.print(msg + ":")
	const tab = "    "
	p.indent += tab

	return func() {
		p.indent = p.indent[:len(p.indent)-len(tab)]
		if err := recover(); err != nil {
			panic(err)
		}
		//p.print("")
	}
}

var line = -1

func (p *parser) print(msg string) {
	if !p.verbose {
		return
	}
	if line != p.line {
		fmt.Printf("line %-4d%s%s\n", p.line, p.indent, msg)
	} else {
		fmt.Printf("         %s%s\n", p.indent, msg)
	}
	line = p.line
}

// Testing Literal
func (p *parser) gotLiteral() error {
	print("Literal: " + p.lit + "\n")
	p.next()
	p.want(_EOF)
	return nil
}

func (p *parser) want(tok Token) {
	if !p.got(tok) {
		p.syntaxError(fmt.Sprintf("expected %s, got %s", tok, p.token))
	}
}

func (p *parser) got(tok Token) bool {
	if p.token == tok {
		p.next()
		return true
	}
	return false
}

func (p *parser) init(r io.Reader, errh ErrHandler) {
	p.errh = errh
	p.lexer.init(r,
		func(line, col int, msg string) {
			p.errorAt(p.posAt(line, col), msg)

		},
	)
}

func tokstring(tok Token) string {
	switch tok {
	case _Comma:
		return "comma"
	case _Semi:
		return "semicolon or newline"
	}
	return tok.String()
}

// ----------------------------------------------------------------------------
// Error handling
func (p *parser) pos() Pos                { return p.posAt(p.line, p.col) }
func (p *parser) posAt(line, col int) Pos { return NewPos(line, col) }
func (p *parser) error(msg string)        { p.errorAt(p.pos(), msg) }
func (p *parser) errorAt(pos Pos, msg string) {
	err := Error{pos, msg}
	if p.errh == nil {
		println(Yellow + err.Msg + Reset)
		return
	}
	p.errh(err)
}
func (p *parser) syntaxError(msg string) { p.syntaxErrorAt(p.pos(), msg) }

func (p *parser) syntaxErrorAt(pos Pos, msg string) {
	if p.verbose {
		p.print(Yellow + "syntax error: " + msg + Reset)
	}

	//if p.Token == _EOF && p.first != nil {
	//	return // avoid meaningless follow-up errors
	//}

	// add punctuation etc. as needed to msg
	switch {
	case msg == "":
		// nothing to do
	case strings.HasPrefix(msg, "in "), strings.HasPrefix(msg, "at "), strings.HasPrefix(msg, "after "):
		msg = " " + msg
	case strings.HasPrefix(msg, "expecting "):
		msg = ", " + msg
	default:
		// plain error - we don't care about current Token
		p.errorAt(pos, "syntax error: "+msg)
		return
	}

	// determine Token string
	var tok string
	switch p.token {
	case _Name, _Semi:
		tok = p.lit
	case _Literal:
		tok = "gotLiteral " + p.lit
	case _Operator:
		tok = p.op.String()
	case _AssignOp:
		tok = p.op.String() + "="
	case _IncOp:
		tok = p.op.String()
		tok += tok
	default:
		tok = tokstring(p.token)
	}

	p.errorAt(pos, "syntax error: unexpected "+tok+msg)
}

const stopset uint64 = 1<<_If |
	1<<_Var

func (p *parser) gotAssign() bool {
	switch p.token {
	case _Define:
		p.error("expecting =")
		fallthrough
	case _Assign:
		p.next()
		return true
	}
	return false
}

// ----------------------------------------------------------------------------
// Declarations
func (p *parser) appendGroup(list []Decl, f func(group *Group) Decl) []Decl {
	if x := f(nil); x != nil {
		list = append(list, x)
	}
	return list
}

// TypeSpec = identifier [ TypeParams ] [ "=" ] Type .
func (p *parser) typeDecl(group *Group) Decl {
	if p.verbose {
		defer p.trace("typeDecl")()
	}

	d := new(TypeDecl)
	d.pos = p.pos()
	d.Group = group

	d.Name = p.name()
	d.Alias = p.gotAssign()
	d.Type = p.typeOrNil()

	if d.Type == nil {
		d.Type = p.badExpr()
		p.syntaxError("in type declaration")
	} else if p.verbose {
		p.print("id: " + d.Name.Value)
		p.print("type: " + d.Type.(*Name).Value)
	}
	return d
}

// VarDecl = "var" identifier ( Type [ "=" Expr ] | "=" Expr ) .
func (p *parser) varDecl(group *Group) Decl {
	if p.verbose {
		defer p.trace("varDecl")()
	}

	d := new(VarDecl)
	d.pos = p.pos()
	d.Group = group

	d.NameList = p.name()
	p.print("id: " + d.NameList.Value)
	if p.gotAssign() {
		d.Values = p.expr()
	} else {
		if p.token != _Name {
			p.syntaxError("expecting name")
			p.next()
			return nil
		}

		d.Type = p.name()
		p.print("type: " + d.Type.(*Name).Value)
	}

	return d
}

// TypeDecl =

// FuncDecl = "func" FuncName Signature FuncBody .
// FuncName = identifier .
func (p *parser) funcDeclOrNil(group *Group) Decl {
	if p.verbose {
		defer p.trace("funcDecl")()
	}

	// func name(name type) type {Body}
	d := new(FuncDecl)
	d.pos = p.pos()
	d.Group = group

	if p.token != _Name {
		p.errorAt(p.pos(), "expecting name")
		return nil
	}

	//function name
	d.Name = p.name()
	p.print("id: " + d.Name.Value)

	// Signature
	d.Param, d.Return = p.funcType()

	// FuncBody
	if p.token == _Lbrace {
		d.Body = p.funcBody()
	}
	return d
}

// OperDecl = "oper" Receiver OperName OperOperand ReturnType OperBody .
// Receiver = "(" Param ")" .
// OperName =
//	"add" | "sub" | "mul" | "div" | "mod" |
//	"radd" | "rsub" | "rmul" | "rdiv" | "rmod" .
// OperOperand = "(" Param ")" .
// ReturnType = Type .
// OperBody = FuncBody .
func (p *parser) operDecl(group *Group) Decl {
	if p.verbose {
		defer p.trace("operDecl")()
	}

	d := new(OperDecl)
	d.pos = p.pos()
	d.Group = group
	d.TypeL = p.singleParam()

	if !p.token.isOperator() {
		p.errorAt(p.pos(), "Unexpected Operator name")
		return nil
	}
	d.Oper = p.token
	p.next()
	p.print("oper type: " + d.Oper.String())
	d.TypeR = p.singleParam()
	p.print("operands: " + d.TypeL.Name.Value + " " + d.TypeR.Name.Value)
	if p.token != _Name {
		p.errorAt(p.pos(), "expecting type")
		return nil
	}
	d.Return = p.name()
	p.print("return type: " + d.Return.(*Name).Value)
	d.Body = p.funcBody()

	return d
}

// FuncBody = Block .
func (p *parser) funcBody() *BlockStmt {
	body := p.blockStmt("")
	return body
}

func (p *parser) funcType() ([]*Field, Expr) {
	params := make([]*Field, 0)
	p.want(_Lparen)
	params = p.paramlist()
	ftype := p.typeOrNil()
	if ftype != nil {
		p.print("return type: " + ftype.(*Name).Value)
	}
	return params, ftype
}

// ----------------------------------------------------------------------------
// Statements

// SimpleStmt = EmptyStmt | ExpressionStmt | IncDecStmt | Assignment | ShortVarDecl .
func (p *parser) simpleStmt(ls Expr, keyword Token) SimpleStmt {
	if p.verbose {
		defer p.trace("simpleStmt")()
	}

	if ls == nil {
		ls = p.expr()
	}

	pos := p.pos()
	switch p.token {
	case _AssignOp, _Assign:
		if p.verbose {
			defer p.trace("assignment")()
		}
		op := p.op
		p.next()
		return p.assignStmt(pos, op, ls, p.expr())
	//case _Define:
	//	if p.verbose {
	//		defer p.p.verbose("shortVarDecl")()
	//	}
	//	p.next()
	//	return p.
	default:
		if p.verbose {
			defer p.trace("exprStmt")()
		}
		s := new(ExprStmt)
		s.pos = ls.Pos()
		s.X = ls
		return s
	}

}

func (p *parser) declStmt(f func(*Group) Decl) *DeclStmt {
	if p.verbose {
		defer p.trace("declStmt")()
	}

	s := new(DeclStmt)
	s.pos = p.pos()

	p.next() // _Const, _Type, or _Var
	s.DeclList = p.appendGroup(nil, f)

	return s
}

// Assignment = Expr assign_op Expr .
// assign_op = [ ass_op | mul_op ] "=" .
func (p *parser) assignStmt(pos Pos, op Operator, lhs, rhs Expr) *AssignStmt {
	a := new(AssignStmt)
	a.pos = pos
	a.Op = op
	a.Lhs = lhs
	a.Rhs = rhs
	return a
}

// Block = "{" StatementList "}" .
func (p *parser) blockStmt(context string) *BlockStmt {
	if p.verbose {
		defer p.trace("blockStmt")()
	}
	s := new(BlockStmt)
	s.pos = p.pos()
	// people coming from C may forget that braces are mandatory in Go
	if !p.got(_Lbrace) {
		p.syntaxError("expecting '{'")
		return nil
	}
	s.StmtList = p.stmtList()

	s.Rbrace = p.pos()
	p.want(_Rbrace)

	return s
}

// StatementList = { Statement ";" } .
func (p *parser) stmtList() (l []Stmt) {
	if p.verbose {
		defer p.trace("stmtList")()
	}

	for p.token != _EOF && p.token != _Rbrace {
		s := p.stmtOrNil()
		if s == nil {
			break
		}
		l = append(l, s)
		// ";" is optional before "}"
		if !p.got(_Semi) && p.token != _Rbrace {
			p.syntaxError("at end of statement")
			p.got(_Semi) // avoid spurious empty statement
		}
	}
	return
}

// Statement =
// 		Declaration | SimpleStmt | ReturnStmt | BreakStmt | ContinueStmt |
//		Block | IfStmt | ForStmt .
func (p *parser) stmtOrNil() Stmt {
	if p.verbose {
		defer p.trace("stmt")()
	}

	if p.token == _Name {
		p.print("lhs:")
		lhs := p.expr()
		return p.simpleStmt(lhs, 0)
	}
	switch p.token {
	case _Var:
		return p.declStmt(p.varDecl)
	case _Lbrace:
		return p.blockStmt("")
	case _Literal, _Name:
		return p.simpleStmt(nil, 0)
	case _For:
		return p.forStmt()
	case _If:
		return p.ifStmt()
	case _Return:
		s := new(ReturnStmt)
		s.pos = p.pos()
		p.next()
		if p.token != _Semi && p.token != _Rbrace {
			s.Return = p.expr()
		}
		return s
	case _Semi:
		func() { defer p.trace("empty stmt")() }()
		s := new(EmptyStmt)
		s.pos = p.pos()
		return s
	}
	return nil
}

// ----------------------------------------------------------------------------
// Expressions

func (p *parser) expr() Expr {
	if p.verbose {
		defer p.trace("expr")()
	}

	return p.binaryExpr(0)
}

// Expr = UnaryExpr | Expr binary_op Expr .
func (p *parser) binaryExpr(prec int) Expr {
	// don't p.verbose binaryExpr - only leads to overly nested p.verbose output

	x := p.unaryExpr()
	for (p.token == _Operator) && p.prec > prec {
		t := new(Operation)
		t.pos = p.pos()
		t.Op = p.op
		tprec := p.prec
		p.print("operator(" + t.Op.String() + ")")
		p.next()
		t.X = x
		t.Y = p.binaryExpr(tprec)
		x = t
	}
	return x
}

// UnaryExpr = PrimaryExpr | unary_op UnaryExpr .
func (p *parser) unaryExpr() Expr {
	if p.verbose {
		defer p.trace("unaryExpr")()
	}
	switch p.token {
	case _Operator:
		switch p.op {
		case Mul, Add, Sub, Not, Xor:
			x := new(Operation)
			x.pos = p.pos()
			x.Op = p.op
			p.next()
			x.X = p.unaryExpr()
			return x

		case And:
			x := new(Operation)
			x.pos = p.pos()
			x.Op = And
			p.next()
			// unaryExpr may have returned a parenthesized composite gotLiteral
			// (see comment in operand) - remove parentheses if any
			x.X = unparen(p.unaryExpr())
			return x
		}
	}
	return p.pexpr()
}

func (p *parser) operand() (rtn Expr) {
	if p.verbose {
		defer p.trace("operand")()
	}

	rtn = &BadExpr{}
	tok := p.token.String()
	switch p.token {
	case _Name:
		rtn = p.name()
		p.print(tok + "(" + rtn.(*Name).Value + ")")
	case _Literal:
		lit := p.literal()
		rtn = lit
		p.print(tok + "(" + lit.Value + ")")
	}
	return
}

// PrimaryExpr =
// 	Operand |
// 	PrimaryExpr Selector |
// 	PrimaryExpr Call .
//
// Selector       = "." identifier .
// Call			  = "(" [ ExprList ] ")" .
func (p *parser) pexpr() Expr {
	if p.verbose {
		defer p.trace("pexpr")()
	}
	x := p.operand()

loop:
	for {
		pos := p.pos()
		switch p.token {
		case _Dot:
			p.next()
			switch p.token {
			case _Name:
				// pexpr '.' sym
				t := new(SelectorExpr)
				t.pos = pos
				t.X = x
				t.Sel = p.name()
				x = t

			default:
				p.syntaxError("expecting name or (")
			}

		case _Lparen:
			t := new(CallExpr)
			t.pos = pos
			t.Func = x
			t.ArgList = p.argList()
			x = t

		default:
			break loop
		}
	}

	return x
}

// ----------------------------------------------------------------------------
// Types
func (p *parser) typeOrNil() Expr {
	if p.token == _Name {
		return p.name()
	}
	return nil
}

func (p *parser) literal() *BasicLit {
	if p.token == _Literal {
		b := new(BasicLit)
		b.pos = p.pos()
		b.Value = p.lit
		b.Kind = p.kind
		b.Bad = p.bad
		p.next()
		return b
	}
	return nil
}

func (p *parser) singleParam() *Field {
	param := new(Field)
	if !p.got(_Lparen) {
		p.syntaxError("expecting '('")
		return nil
	}
	first := true
recv:
	if p.token != _Name {
		str := "type"
		if first {
			str = "receiver"
		}
		p.syntaxError("expecting " + str)
		return nil
	}
	name := p.name()
	if first {
		param.Name = name
		first = false
		goto recv
	}
	param.Type = name
	p.want(_Rparen)
	return param
}

func (p *parser) paramlist() []*Field {
	list := make([]*Field, 0)
	none := "none"
	str := " "
redo:
	param := new(Field)
	switch p.token {
	case _Name:
		none = ""
		param.Name = p.name()
		if p.token == _Name {
			ptype := p.typeOrNil()
			str += none + param.Name.Value + "(" + ptype.(*Name).Value + ") "
			param.Type = ptype
			list = append(list, param)
			switch p.token {
			case _Comma:
				p.next()
				goto redo
			case _Rparen:
				p.next()
				p.print("params:" + str)
				return list
			default:
				p.syntaxError("expecting comma or ')'")
				p.next()
				return nil
			}
		} else {
			p.syntaxError("expecting type")
			p.next()
			return nil
		}
	case _Rparen:
		p.next()
		return nil
	default:
		p.syntaxError("expecting parameter or ')'")
		p.next()
		return nil
	}
}

func (p *parser) argList() []Expr {
	if p.verbose {
		defer p.trace("argList")()
	}
	list := make([]Expr, 0)
	p.want(_Lparen)
	if !p.got(_Rparen) {
		list = append(list, p.expr())
		for !p.got(_Rparen) {
			p.want(_Comma)
			list = append(list, p.expr())
		}
	}

	return list
}

// ----------------------------------------------------------------------------
// Common
func (p *parser) name() *Name {
	// no tracing to avoid overly p.verbose output

	if p.token == _Name {
		n := NewName(p.pos(), p.lit)
		p.next()
		return n
	}

	n := NewName(p.pos(), "_")
	p.error("expecting name")
	return n
}

func (p *parser) nameList(first *Name) []*Name {
	if p.verbose {
		defer p.trace("nameList")()
	}

	l := []*Name{first}
	for p.got(_Comma) {
		l = append(l, p.name())
	}

	return l
}

func (p *parser) forStmt() Stmt {
	if p.verbose {
		defer p.trace("forStmt")()
	}

	s := new(ForStmt)
	s.pos = p.pos()

	s.Init, s.Cond, s.Post = p.header(_For)
	s.Body = p.blockStmt("for clause")

	return s
}

func (p *parser) header(keyword Token) (init SimpleStmt, cond Expr, post SimpleStmt) {
	p.want(keyword)
	if p.token == _Lbrace {
		if keyword == _If {
			p.syntaxError("missing condition in if statement")
			cond = p.badExpr()
		}
		return
	}

	if p.token != _Semi {
		// accept potential varDecl but complain
		if p.got(_Var) {
			p.syntaxError(fmt.Sprintf("var declaration not allowed in %s initializer", tokstring(keyword)))
		}
		init = p.simpleStmt(nil, keyword)
	}
	var condStmt SimpleStmt
	var semi struct {
		pos Pos
		lit string // valid if pos.IsKnown()
	}
	if p.token != _Lbrace {
		if p.token == _Semi {
			semi.pos = p.pos()
			semi.lit = p.lit
			p.next()
		} else {
			// asking for a '{' rather than a ';' here leads to a better error message
			p.want(_Lbrace)
		}
		if keyword == _For {
			if p.token != _Semi {
				if p.token == _Lbrace {
					p.syntaxError("expecting for loop condition")
					goto done
				}
				condStmt = p.simpleStmt(nil, 0 /* range not permitted */)
			}
			p.want(_Semi)
			if p.token != _Lbrace {
				post = p.simpleStmt(nil, 0 /* range not permitted */)
				if a, _ := post.(*AssignStmt); a != nil && a.Op == Def {
					p.syntaxErrorAt(a.Pos(), "cannot declare in post statement of for loop")
				}
			}
		} else if p.token != _Lbrace {
			condStmt = p.simpleStmt(nil, keyword)
		}
	} else {
		condStmt = init
		init = nil
	}
done:
	// unpack condStmt
	switch s := condStmt.(type) {
	case nil:
		if keyword == _If && semi.pos.IsKnown() {
			if semi.lit != "semicolon" {
				p.syntaxErrorAt(semi.pos, fmt.Sprintf("unexpected %s, expecting { after if clause", semi.lit))
			} else {
				p.syntaxErrorAt(semi.pos, "missing condition in if statement")
			}
			b := new(BadExpr)
			b.pos = semi.pos
			cond = b
		}
	case *ExprStmt:
		cond = s.X
	default:
		p.syntaxErrorAt(s.Pos(), fmt.Sprintf("cannot use %s as value", s))
	}
	return
}

func (p *parser) badExpr() *BadExpr {
	b := new(BadExpr)
	b.pos = p.pos()
	return b
}

func (p *parser) ifStmt() *IfStmt {
	if p.verbose {
		defer p.trace("ifStmt")()
	}
	s := new(IfStmt)
	s.pos = p.pos()
	_, s.Cond, _ = p.header(_If)
	s.Block = p.blockStmt("If clause")
	if p.got(_Else) {
		switch p.token {
		case _If:
			s.Else = p.ifStmt()
		case _Lbrace:
			s.Else = p.blockStmt("")
		default:
			p.syntaxError("else must be followed by if or statement block")
		}
	}
	return s
}

func unparen(x Expr) Expr {
	for {
		p, ok := x.(*ParenExpr)
		if !ok {
			break
		}
		x = p.X
	}
	return x
}
