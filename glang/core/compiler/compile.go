package compiler

import (
	"almeng.com/glang/core/builtin"
	"almeng.com/glang/core/syntax"
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
	var FuncDefs = make([]*Context, 0)
	Funcs.Each(func(v *syntax.FuncDecl) { FuncDefs = append(FuncDefs, c.DefineFunc(space, v)) })
	//Vars.Each(func(v *syntax.VarDecl) { c.CurrentSpace.vars[v.NameList.Value] = c.compileVar() })
	Opers.Each(func(v *syntax.OperDecl) { c.CompileOper(space, v) })
	Funcs.Iter(func(i int, v *syntax.FuncDecl) { c.CompileFunc(space, FuncDefs[i], v) })
	//Types.Each(func(v *syntax.TypeDecl){ c.CompileType(space, v)})
}

type Context struct {
	*ir.Block
	parent *Context
	vars   map[string]value.Value
	leave  *ir.Block
}

func NewContext(b *ir.Block) *Context {
	return &Context{
		Block:  b,
		parent: nil,
		vars:   make(map[string]value.Value),
		leave:  nil,
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

func (c *Context) localize(v value.Value) value.Value {
	switch v.(type) {
	case *ir.Global:
		r := v.(*ir.Global)
		return c.NewLoad(r.Typ.ElemType, r)
	}
	return v
}

func RetType(m *ir.Module, name *syntax.Name) types.Type {
	t := general.AsList(m.TypeDefs).Filter(func(i int, e types.Type) bool {
		return e.Name() == name.Value
	})
	if t == nil {

	}
	return nil
}

func (c *Compiler) CompileOper(space string, f *syntax.OperDecl) {
	r, l := c.QueryType(space, f.TypeR), c.QueryType(space, f.TypeL)
	rtn := c.QueryType(space, f.Return)
	if f.Oper.IsReversedOper() {
		r, l = l, r
		f.Oper = f.Oper - (syntax.Reversed_oper - syntax.Operator_beg)
	}
	c.Opers.Append(NewOperator(f, l, r, rtn))
}

func (c *Compiler) DefineFunc(space string, f *syntax.FuncDecl) *Context {
	m := c.Module
	name := f.Name.Value
	name = space + ".func." + name
	var ret types.Type
	ret = c.QueryType(space, f.Return)
	Params := make([]*ir.Param, 0)
	general.AsList(f.Param).Each(func(v *syntax.Field) {
		Params = append(Params, ir.NewParam(v.Name.Value, c.QueryType(space, v.Type)))
	})
	Func := m.NewFunc(name, ret, Params...)
	if Func.Name() == "main.func.main" {
		if !ret.Equal(types.Void) || f.Param != nil {
			fmt.Println("main function must have no arguments and no return values")
			os.Exit(1)

		}
		c.Main = Func
	}
	block := Func.NewBlock(Func.Name() + ".body")
	ctx := c.CurrentSpace.NewContext(block)
	general.AsList(Params).Each(func(v *ir.Param) {
		ctx.vars[v.Name()] = v
	})
	return ctx
}

func (c *Compiler) CompileFunc(space string, ctx *Context, def *syntax.FuncDecl) *ir.Func {
	last, ok := def.Body.StmtList[len(def.Body.StmtList)-1].(*syntax.ReturnStmt)
	if !ok {
		last = &syntax.ReturnStmt{Return: nil}
		def.Body.StmtList = append(def.Body.StmtList, last)
	}

	c.CompileBody(ctx, def.Body, false, nil)

	ret := ctx.Block.Parent.Sig.RetType

	got := c.QueryType(space, def.Return)
	if !ret.Equal(got) {
		println("expected return type " + ret.Name() + " but got " + got.Name())
		os.Exit(1)
	}

	return ctx.Block.Parent
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
		return builtin.Void
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

func (c *Compiler) CompileBody(ctx *Context, body *syntax.BlockStmt, inline bool, inlineRtn *value.Value) *Context {
	for _, v := range body.StmtList {
		c.CompileStmt(ctx, v, inline, inlineRtn)
	}
	return ctx
}

func (c *Compiler) EvalOperation(ctx *Context, oper *syntax.Operation) value.Value {
	if oper.Y == nil {
		return c.EvalUnary(ctx, oper)

	}
	x := c.CompileExpr(ctx, oper.X)
	y := c.CompileExpr(ctx, oper.Y)

	swaped := false
	op := oper.Op
	switch oper.Op {
	case syntax.Lss, syntax.Leq:
		x, y = y, x
		oper.Op += 2
		swaped = true
	}

	if oper.Op == syntax.Geq {

	}

	// TODO
	operation := c.QueryOperator(oper.Op, x.Type(), y.Type())
	if operation == nil {
		errStr := x.Type().Name() + " and " + y.Type().Name()
		if swaped {
			errStr = y.Type().Name() + " and " + x.Type().Name()
		}
		if x.Type().Equal(y.Type()) {
			errStr = x.Type().Name()
		}

		fmt.Println(oper.Pos().String() + " Invalid Operation: the operator " + op.String() + " is not defined for " + errStr)
		os.Exit(1)
	}
	return operation.Operation(c, ctx, x, y)

}

func (c *Compiler) EvalUnary(ctx *Context, op *syntax.Operation) value.Value {
	x := c.CompileExpr(ctx, op.X)
	oper := c.QueryOperator(op.Op, x.Type(), nil)
	if oper == nil {
		fmt.Println(op.Pos().String() + " Invalid Operation: the operator " + op.Op.String() + " is not defined for " + x.Type().Name())
		os.Exit(1)
	}
	return oper.Operation(c, ctx, x, nil)
}
func (c *Compiler) QueryOperator(oper syntax.Operator, L, R types.Type) *Operator {
	op := c.Opers.Filter(func(i int, v *Operator) bool {
		return v.Op == oper && v.TypeL.Equal(L) && (R == nil || v.TypeR.Equal(R))
	}).First()
	if op != nil {
		return *op
	}
	return nil
}

func (c *Compiler) CompileExpr(ctx *Context, s syntax.Expr) value.Value {
	// TODO
	switch s.(type) {
	case *syntax.BadExpr:
		expr := s.(*syntax.BadExpr)
		_ = expr

	case *syntax.Name:
		expr := s.(*syntax.Name)
		_var := ctx.lookupVariable(expr.Value)
		switch _var.(type) {
		case *ir.Global:
			v := _var.(*ir.Global)
			return ctx.NewLoad(v.Typ.ElemType, v)
		}
		return _var

	case *syntax.BasicLit:
		expr := s.(*syntax.BasicLit)
		switch expr.Kind {
		case syntax.StringLit:
			return c.NewGlobalString(ctx.Block, expr.Value)
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
			return constant.NewFloat(builtin.Float, v)
		}

	case *syntax.Operation:
		expr := s.(*syntax.Operation)
		return c.EvalOperation(ctx, expr)

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
		general.AsList(expr.ArgList).Each(func(v syntax.Expr) { args = append(args, c.CompileExpr(ctx, v)) })
		ctx.NewCall(callee, args...)

	case *syntax.Field:
		expr := s.(*syntax.Field)
		_ = expr

	}
	return nil
}

func (c *Compiler) CompileStmt(ctx *Context, s syntax.Stmt, inline bool, inlineRtn *value.Value) value.Value {
	switch s.(type) {
	case *syntax.ExprStmt:
		stmt := s.(*syntax.ExprStmt)
		return c.CompileExpr(ctx, stmt.X)
	case *syntax.EmptyStmt:
		return nil
	case *syntax.IncDecStmt:
		stmt := s.(*syntax.IncDecStmt)
		_ = stmt

	case *syntax.ContinueStmt:
		stmt := s.(*syntax.ContinueStmt)
		_ = stmt
	//case *syntax.WhileStmt:
	//	stmt := s.(*syntax.WhileStmt)
	//	c.CompileWhile(ctx, stmt)

	case *syntax.BreakStmt:
		ctx.NewBr(ctx.parent.Block)

	case *syntax.ReturnStmt:
		stmt := s.(*syntax.ReturnStmt)
		if inline {
			rtn := c.CompileExpr(ctx, stmt.Return)
			str := ctx.NewStore(rtn, *inlineRtn)
			return ctx.NewLoad(rtn.Type(), str.Dst)
		}
		return ctx.Block.NewRet(c.CompileExpr(ctx, stmt.Return)).X

	case *syntax.DeclStmt:
		stmt := s.(*syntax.DeclStmt)
		_ = stmt

	case *syntax.AssignStmt:
		stmt := s.(*syntax.AssignStmt)
		_ = stmt

	case *syntax.IfStmt:
		stmt := s.(*syntax.IfStmt)
		c.CompileIf(ctx, stmt, inline, inlineRtn)

	//case *syntax.ForStmt:
	//	stmt := s.(*syntax.ForStmt)
	//	c.CompileFor(ctx, stmt)

	case *syntax.BlockStmt:
		stmt := s.(*syntax.BlockStmt)
		c.CompileBody(ctx, stmt, inline, inlineRtn)
	}
	return nil
}

var (
	ifCntr    = 0
	forCntr   = 0
	whileCntr = 0
)

func (c *Compiler) CompileIf(ctx *Context, stmt *syntax.IfStmt, inline bool, inlineRtn *value.Value) {
	f := ctx.Block.Parent
	_cntr := ifCntr
	ifCntr++
	thenCtx := ctx.NewContext(f.NewBlock("if.then." + strconv.Itoa(_cntr)))
	block := thenCtx.Block
	c.CompileStmt(thenCtx, stmt.Block, inline, inlineRtn)
	elseB := f.NewBlock("if.else." + strconv.Itoa(_cntr))
	c.CompileStmt(ctx.NewContext(elseB), stmt.Else, inline, inlineRtn)
	//ctx.NewCondBr(constant.NewInt(builtin.Bool, 0), thenCtx.Block, elseB)
	ctx.NewCondBr(c.CompileExpr(ctx, stmt.Cond), block, elseB)
	var leaveB *ir.Block

	if thenCtx.Block.Term == nil {
		leaveB = f.NewBlock("leave.if." + strconv.Itoa(_cntr))
		thenCtx.NewBr(leaveB)
		ctx.Block = leaveB
	} else if inline {
		thenCtx.NewBr(ctx.leave)
	}
	//if elseB.Term == nil {
	//	if leaveB == nil {
	//		leaveB = f.NewBlock("leave.if." + strconv.Itoa(_cntr))
	//	}
	//	elseB.NewBr(leaveB)
	//	ctx.Block = leaveB
	//}
	ifCntr++
}

//func (c *Compiler) CompileFor(ctx *Context, stmt *syntax.ForStmt) {
//	f := ctx.Block.Parent
//	loopCtx := ctx.NewContext(f.NewBlock("for.loop.body"))
//	ctx.NewBr(loopCtx.Block)
//	firstAppear := loopCtx.NewPhi(ir.NewIncoming(c.CompileStmt(loopCtx, stmt.Init), ctx.Block))
//	loopCtx.vars[stmt.Init.Pos().String()] = firstAppear
//	step := c.CompileStmt(loopCtx, stmt.Post)
//	firstAppear.Incs = append(firstAppear.Incs, ir.NewIncoming(step, loopCtx.Block))
//	loopCtx.vars[stmt.Init.Pos().String()] = step
//	leaveB := f.NewBlock("leave.for.loop")
//	loopCtx.Block = leaveB
//	c.CompileStmt(loopCtx, stmt.Body)
//	loopCtx.NewCondBr(c.CompileExpr(loopCtx, stmt.Cond), loopCtx.Block, leaveB)
//}

//func (c *Compiler) CompileWhile(ctx *Context, stmt *syntax.WhileStmt) {
//	f := ctx.Block.Parent
//	loopCtx := ctx.NewContext(f.NewBlock("while.loop.body"))
//	ctx.NewBr(loopCtx.Block)
//	leaveB := f.NewBlock("leave.while.loop")
//	ctx.Block = leaveB
//	c.CompileStmt(loopCtx, stmt.Body)
//	if loopCtx.Block.Term != nil {
//		return
//	}
//	loopCtx.NewCondBr(c.CompileExpr(loopCtx, stmt.Cond), loopCtx.Block, leaveB)
//}
