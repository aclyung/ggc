package compiler

import (
	buitin "almeng.com/glang/core/builtin"
	"almeng.com/glang/core/syntax"
	"almeng.com/glang/global"
	"fmt"
	"github.com/almenglee/general"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"os"
	"strconv"
)

type DeclTypes interface {
	*syntax.FuncDecl | *syntax.VarDecl | *syntax.TypeDecl | *syntax.OperDecl
}

func Filter[T DeclTypes]() func(int, syntax.Decl) bool {
	return func(i int, v syntax.Decl) bool {
		_, ok := v.(T)
		return ok
	}
}

func Append[T DeclTypes](list *general.List[T]) func(syntax.Decl) {
	return func(decl syntax.Decl) {
		list.Append(decl.(T))
	}
}

func (c *Compiler) CompileFile(f *syntax.File) {
	decls := general.AsList(f.DeclList)

	var (
		Funcs = general.EmptyList[*syntax.FuncDecl]()
		Opers = general.EmptyList[*syntax.OperDecl]()
		Vars  = general.EmptyList[*syntax.VarDecl]()
		Types = general.EmptyList[*syntax.TypeDecl]()
	)

	decls.Filter(Filter[*syntax.FuncDecl]()).Each(Append(Funcs))
	decls.Filter(Filter[*syntax.OperDecl]()).Each(Append(Opers))
	decls.Filter(Filter[*syntax.VarDecl]()).Each(Append(Vars))
	decls.Filter(Filter[*syntax.TypeDecl]()).Each(Append(Types))

	space := f.SpaceName.Value
	if c.CurrentSpace == nil || c.CurrentSpace.Name() != space {
		ctx := c.Global.NewContext(ir.NewBlock(space))
		c.CurrentSpace = ctx
		c.Spaces.Append(ctx)
	}

	if space == "main" {
		if Funcs.Filter(func(i int, v *syntax.FuncDecl) bool { return v.Name.Value == "main" }).First() == nil {
			println("program entry point was not found")
			os.Exit(1)
		}
	}
	Types.Each(func(v *syntax.TypeDecl) { c.CompileType(space, v) })
	var FuncDefs = make([]*ir.Func, 0)
	Funcs.Each(func(v *syntax.FuncDecl) { FuncDefs = append(FuncDefs, c.DefineFunc(space, v)) })
	//Vars.Each(func(v *syntax.VarDecl) { c.CurrentSpace.vars[v.NameList.Value] = c.Module.new })
	Opers.Each(func(v *syntax.OperDecl) { c.CompileOper(space, v) })
	Funcs.Iter(func(i int, v *syntax.FuncDecl) { c.CompileFunc(space, FuncDefs[i], v) })
	//Types.Each(func(v *syntax.TypeDecl){ c.CompileType(space, v)})
}

type Context struct {
	*ir.Block
	parent *Context
	vars   map[string]value.Value
}

func NewContext(b *ir.Block) *Context {
	return &Context{
		Block:  b,
		parent: nil,
		vars:   make(map[string]value.Value),
	}
}

func (c *Context) NewContext(b *ir.Block) *Context {
	ctx := NewContext(b)
	ctx.parent = c
	return ctx
}

func (c *Context) lookupVariable(name string) value.Value {
	if v, ok := c.vars[name]; ok {
		return v
	} else if c.parent != nil {
		return c.parent.lookupVariable(name)
	} else {
		fmt.Printf("variable: `%s`\n", name)
		panic("no such variable")
	}
}

func RetType(m *ir.Module, name *syntax.Name) types.Type {
	t := general.AsList(m.TypeDefs).Filter(func(i int, e types.Type) bool {
		return e.Name() == name.Value
	})
	if t == nil {

	}
	return nil
}

func OperName(operator syntax.Token, L types.Type, R types.Type) (name string) {
	name = "oper."
	nameL, nameR := "_"+L.Name(), "_"+R.Name()
	if operator.IsReversedOper() {
		name += (operator - (syntax.Reversed_oper - syntax.Operator_beg)).String() + "."
		name += nameR + nameL
	} else {
		name += operator.String() + "."
		name += nameL + nameR
	}
	return
}

func (c *Compiler) CompileOper(space string, f *syntax.OperDecl) {
	r, l := c.QueryType(space, f.TypeR), c.QueryType(space, f.TypeL)
	c.Opers.Append(NewOperator(f.Oper, l, r, f.Body))
}

func (c *Compiler) DefineFunc(space string, f *syntax.FuncDecl) *ir.Func {
	m := c.Module
	name := f.Name.Value
	name = space + ".func." + name
	var ret types.Type
	ret = c.QueryType(space, f.Return)
	Func := m.NewFunc(name, ret)
	if Func.Name() == "main.func.main" {
		if !ret.Equal(types.Void) || f.Param != nil {
			fmt.Println("main function must have no arguments and no return values")
			os.Exit(1)
		}
		c.Main = Func
	}
	return Func
}

func (c *Compiler) CompileFunc(space string, f *ir.Func, def *syntax.FuncDecl) *ir.Func {
	var ret types.Type
	ret = c.QueryType(space, def.Return)
	Func := f
	last, ok := def.Body.StmtList[len(def.Body.StmtList)-1].(*syntax.ReturnStmt)
	if !ok {
		last = &syntax.ReturnStmt{Return: nil}
		def.Body.StmtList = append(def.Body.StmtList, last)
	}

	rtn := c.CompileExpr(ir.NewBlock(""), last.Return)
	var type_return types.Type = buitin.Void
	if rtn != nil {
		type_return = rtn.Type()
	}
	if !ret.Equal(type_return) {
		tName := type_return.String()
		if type_return == buitin.Void {
			tName = "nil"
		}
		println("expected return type " + ret.Name() + " but got " + tName)
		os.Exit(1)
	}
	c.CompileBody(Func, def.Body)
	return Func
}

func DefType[T interface {
	types.IntType |
		types.FloatType |
		types.ArrayType |
		types.StructType
}](t *T) *T {
	_type := new(T)
	*_type = *t
	return _type
}

func (c *Compiler) CompileType(space string, t *syntax.TypeDecl) {
	name := space + ".type." + t.Name.Value
	var tt types.Type
	if t.Type == nil {
		tt = types.Void
	} else {
		tt = c.QueryType(space, t.Type)
	}
	switch tt.(type) {
	case *types.IntType:
		c.Module.NewTypeDef(name, DefType(tt.(*types.IntType)))
	case *types.FloatType:
		c.Module.NewTypeDef(name, DefType(tt.(*types.FloatType)))
	case *types.ArrayType:
		c.Module.NewTypeDef(name, DefType(tt.(*types.ArrayType)))
	case *types.StructType:
		c.Module.NewTypeDef(name, DefType(tt.(*types.StructType)))
	}
}

type Query struct {
	Name       string
	IsInternal bool
}

func _qName(space string, t syntax.Expr) Query {
	switch t.(type) {
	case *syntax.Field:
		return _qName(space, t.(*syntax.Field).Type)
	case *syntax.Name:
		return Query{t.(*syntax.Name).Value, true}
	case *syntax.SelectorExpr:
		s := t.(*syntax.SelectorExpr)
		return Query{s.Sel.Value + _qName(space, s.X).Name, false}
	}
	panic("not expected Expr type")
}

func (c *Compiler) QueryType(space string, t syntax.Expr) types.Type {
	if t == nil {
		return buitin.Void
	}
	name := ""
	q := _qName(space, t)
	if q.IsInternal {
		name_internal := space + ".type." + q.Name
		typ := general.AsList(c.Module.TypeDefs).Filter(func(i int, _typ types.Type) bool { return _typ.Name() == name_internal }).First()
		if typ != nil {
			return *typ
		}
	}
	name += q.Name
	typ := general.AsList(c.Module.TypeDefs).Filter(func(i int, _typ types.Type) bool { return _typ.Name() == name }).First()
	if typ == nil {
		panic("no type name " + name + " was defined")
	}
	return *typ
}

var str *ir.InstGetElementPtr

func (c *Compiler) CompileBody(p *ir.Func, body *syntax.BlockStmt) *ir.Block {
	b := p.NewBlock("")
	for _, v := range body.StmtList {
		c.CompileStmt(b, v)
	}
	return b
}

func (c *Compiler) EvalOperation(p *ir.Block, oper *syntax.Operation) value.Value {
	if oper.Y == nil {
		return c.EvalUnary(oper.Op, oper.X)

	}
	x := c.CompileExpr(p, oper.X)
	y := c.CompileExpr(p, oper.Y)

	// TODO
	o := c.QueryOperator(oper.Op, x.Type(), y.Type())
	return o.Operation(c, p, x, y)

}

func (c *Compiler) EvalUnary(op syntax.Operator, x syntax.Expr) value.Value {
	return nil
}
func (c *Compiler) QueryOperator(oper syntax.Operator, L, R types.Type) *Operator {
	var op syntax.Token
	switch oper {
	case syntax.Add:
		op = syntax.OperAdd
	case syntax.Sub:
		op = syntax.OperSub
	case syntax.Mul:
		op = syntax.OperMul
	case syntax.Div:
		op = syntax.OperDiv
	case syntax.Rem:
		op = syntax.OperRem
	default:
		op = syntax.OperAdd
	}
	return *c.Opers.Filter(func(i int, v *Operator) bool { return v.Name() == OperName(op, L, R) }).First()
}

func (c *Compiler) CompileExpr(p *ir.Block, s syntax.Expr) value.Value {
	// TODO
	switch s.(type) {
	case *syntax.BadExpr:
		expr := s.(*syntax.BadExpr)
		_ = expr

	case *syntax.Name:
		expr := s.(*syntax.Name)
		_ = expr
		return constant.NewInt(types.I64, 1)

	case *syntax.BasicLit:
		expr := s.(*syntax.BasicLit)
		switch expr.Kind {
		case syntax.StringLit:
			return global.NewLocalString(p, expr.Value)
		case syntax.IntLit:
			v, err := strconv.ParseInt(expr.Value, 10, 64)
			if err != nil {
				panic("wrong value")
			}
			return constant.NewInt(types.I64, v)
		case syntax.FloatLit:
			v, err := strconv.ParseFloat(expr.Value, 64)
			if err != nil {
				panic("wrong value")
			}
			return constant.NewFloat(buitin.Float, v)
		}

	case *syntax.Operation:
		expr := s.(*syntax.Operation)
		return c.EvalOperation(p, expr)

	case *syntax.ParenExpr:
		expr := s.(*syntax.ParenExpr)
		_ = expr

	case *syntax.SelectorExpr:
		expr := s.(*syntax.SelectorExpr)
		_ = expr

	case *syntax.CallExpr:
		expr := s.(*syntax.CallExpr)
		name := expr.Func.(*syntax.Name).Value
		callee := *(general.AsList(c.Module.Funcs).Filter(func(i int, v *ir.Func) bool { return v.Name() == name }).First())
		var args []value.Value
		general.AsList(expr.ArgList).Each(func(v syntax.Expr) { args = append(args, c.CompileExpr(p, v)) })
		p.NewCall(callee, args...)

	case *syntax.Field:
		expr := s.(*syntax.Field)
		_ = expr

	}
	return nil
}

func (c *Compiler) CompileStmt(p *ir.Block, s syntax.Stmt) value.Value {
	switch s.(type) {
	case *syntax.ExprStmt:
		stmt := s.(*syntax.ExprStmt)
		return c.CompileExpr(p, stmt.X)
	case *syntax.EmptyStmt:
		return nil
	case *syntax.IncDecStmt:
		stmt := s.(*syntax.IncDecStmt)
		_ = stmt

	case *syntax.ContinueStmt:
		stmt := s.(*syntax.ContinueStmt)
		_ = stmt

	case *syntax.BreakStmt:
		stmt := s.(*syntax.BreakStmt)
		_ = stmt

	case *syntax.ReturnStmt:
		stmt := s.(*syntax.ReturnStmt)
		p.NewRet(c.CompileExpr(p, stmt.Return))
		//ir.NewRet()
	case *syntax.DeclStmt:
		stmt := s.(*syntax.DeclStmt)
		_ = stmt

	case *syntax.AssignStmt:
		stmt := s.(*syntax.AssignStmt)
		_ = stmt

	case *syntax.IfStmt:
		//stmt := s.(*syntax.IfStmt)
		//c.CompileIf()

	case *syntax.ForStmt:
		stmt := s.(*syntax.ForStmt)
		_init, _cond, _post := stmt.Init, stmt.Cond, stmt.Post
		_, _, _ = _init, _cond, _post

	case *syntax.BlockStmt:
		stmt := s.(*syntax.BlockStmt)
		c.CompileBody(p.Parent, stmt)
	}
	return nil
}
