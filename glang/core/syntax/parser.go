package syntax

import (
	"fmt"
	"io"
	"strings"
)

const trace = true

type parser struct {
	errh ErrHandler
	lexer
	indent string
	first  error
}

func (p *parser) EOF() error {
	if trace {
		defer p.trace("file")()
	}

	f := new(File)
	f.pos = p.pos()
	if !p.got(_Space) {
		p.error("space must be declared first")
		return Error{}
	}
	f.SpaceName = p.name()
	p.print("space: " + f.SpaceName.Value)
	p.want(_Semi)
	for p.token != _EOF {
		switch p.token {
		case _Var:
			p.next()
			f.DeclList = p.appendGroup(f.DeclList, p.varDecl)
		case _Literal:
			_ = p.gotLiteral()
		case _Semi:
			p.next()
		case _Oper, _Func:
			p.errorAt(p.pos(), "ERROR: cannot resolve token "+p.token.String()+"("+p.segment()+")")
			p.next()
		default:
			p.errorAt(p.pos(), "ERROR: non-declaration statement outside function body: "+p.token.String()+"("+p.segment()+")")
			p.next()
		}
	}
	return nil
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

func (p *parser) want(tok token) {
	if !p.got(tok) {
		print("error occurred: token unexpected.")
		p.errorf("expected %s, got %s", tok, p.token)
	}
}

func (p *parser) got(tok token) bool {
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

func (p *parser) posAt(line, col int) Pos {
	return NewPos(line, col)
}

func tokstring(tok token) string {
	switch tok {
	case _Comma:
		return "comma"
	case _Semi:
		return "semicolon or newline"
	}
	return tok.String()
}

func (p *parser) pos() Pos         { return p.posAt(p.line, p.col) }
func (p *parser) error(msg string) { p.errorAt(p.pos(), msg) }
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
	if trace {
		p.print("syntax error: " + msg)
	}

	if p.token == _EOF && p.first != nil {
		return // avoid meaningless follow-up errors
	}

	// add punctuation etc. as needed to msg
	switch {
	case msg == "":
		// nothing to do
	case strings.HasPrefix(msg, "in "), strings.HasPrefix(msg, "at "), strings.HasPrefix(msg, "after "):
		msg = " " + msg
	case strings.HasPrefix(msg, "expecting "):
		msg = ", " + msg
	default:
		// plain error - we don't care about current token
		p.errorAt(pos, "syntax error: "+msg)
		return
	}

	// determine token string
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

func (p *parser) appendGroup(list []Decl, f func(group *Group) Decl) []Decl {
	if x := f(nil); x != nil {
		list = append(list, x)
	}
	return list
}

func (p *parser) varDecl(group *Group) Decl {
	if trace {
		defer p.trace("varDecl")()
	}

	d := new(VarDecl)
	d.pos = p.pos()
	d.Group = group

	d.NameList = p.name()
	p.print("id: " + d.NameList.Value)
	if p.gotAssign() {
		d.Values = p.expr()
	}
	return d
}

func (p *parser) expr() Expr {
	if trace {
		defer p.trace("expr")()
	}

	return p.binaryExpr(nil, 0)
}

func (p *parser) binaryExpr(x Expr, prec int) Expr {
	// don't trace binaryExpr - only leads to overly nested trace output

	if x == nil {
		x = p.unaryExpr()
	}
	for (p.token == _Operator) && p.prec > prec {
		t := new(Operation)
		t.pos = p.pos()
		t.Op = p.op
		tprec := p.prec
		p.print("operator(" + t.Op.String() + ")")
		p.next()
		t.X = x
		t.Y = p.binaryExpr(nil, tprec)
		x = t
	}
	return x
}

func (p *parser) unaryExpr() Expr {
	if trace {
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
	return p.operand()
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

func (p *parser) name() *Name {
	// no tracing to avoid overly verbose output

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
	if trace {
		defer p.trace("nameList")()
	}

	l := []*Name{first}
	for p.got(_Comma) {
		l = append(l, p.name())
	}

	return l
}

func (p *parser) operand() (rtn Expr) {
	if trace {
		defer p.trace("operand")()
	}

	rtn = &BadExpr{}
	tok := p.token.String()
	switch p.token {
	case _Literal:
		lit := p.literal()
		rtn = lit
		p.print(tok + "(" + lit.Value + ")")
	}
	return
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
