package compiler

import (
	"almeng.com/glang-vm"
	"almeng.com/glang/core/compiler/backend"
	"almeng.com/glang/core/compiler/backend/types"
	"almeng.com/glang/core/syntax"

	"github.com/almenglee/general"
	"io"
	"os"
	"strconv"
)

type BCCompiler struct {
	Module       *backend.Module
	Spaces       general.List[*BCContext]
	CurrentSpace *BCContext
	Global       *BCContext
	Opers        general.List[*Operator]
	verbose      bool
}

func BuildBC(c *BCCompiler) []byte {
	ir := c.GetIR()
	return vm.CompileIR(ir)
}

func (c *BCCompiler) GetIR() string {
	return c.Module.String()
}

func CompileBC(filename string, verbose bool) *BCCompiler {
	f, ferr := os.Open(filename)
	if ferr != nil {
		println(ferr.Error())
		os.Exit(-1)
	}
	c := NewBCCompiler(verbose)
	if c.verbose {
		var src string
		_file, _err := os.Open(filename)
		if _err != nil {
			println(_err.Error())
			os.Exit(-1)
		}
		if b, err := io.ReadAll(_file); err != nil {
			println(err.Error())
			os.Exit(-1)
		} else {
			src = string(b)
		}
		syntax.TokenizingTest(filename, src)
	}
	c.Module.Funcs = append(c.Module.Funcs, c.Module.NewFunc("println", nil, nil))
	file := syntax.Parse(f, errh, c.verbose)
	c.CompileFile(file)
	return c
}

func (c *BCCompiler) QueryType(space string, t syntax.Expr) types.Type {
	if t == nil {
		return types.Void
	}
	name := ""
	q := _qName(space, t)
	if q.IsInternal {
		name_internal := space + ".type." + q.Name
		typ := c.Module.TypeDecl.Filter(func(i int, v types.Type) bool {
			//TODO: make sure that the type is not nil
			//return v.Name() == name_internal
			_ = name_internal
			return false
		}).First()
		if typ != nil {
			return *typ
		}
	}
	name += q.Name
	typ := c.Module.TypeDecl.Filter(func(i int, v types.Type) bool { return false }).First()
	if typ == nil {
		panic("no type name " + name + " was defined")
	}
	return *typ
}

func QueryName(space string, t syntax.Expr) Query {
	return _qName(space, t)
}

//func (c *BCCompiler) CompileType(space string, t *syntax.TypeDecl) {
//	name := space + ".type." + t.Name.Value
//	var typ types.Type
//	if t.Type == nil {
//		typ = types.Void
//	} else {
//		typ = c.QueryType(space, t.Type)
//	}
//}

func NewBCCompiler(verbose bool) *BCCompiler {
	return &BCCompiler{
		Module:  backend.NewModule(),
		Spaces:  nil,
		verbose: verbose,
	}
}

func (c *BCCompiler) InitGlobal() {
	m := c.Module
	m.TypeDecl.Append(types.NewIntType(64))
}

type BCContext struct {
	*backend.Block
	Parent *BCContext
	vars   map[string]backend.Value
}

func NewBCContext(b *backend.Block) *BCContext {
	return &BCContext{
		Block:  b,
		Parent: nil,
	}
}

func (c *BCContext) NewBCContext(b *backend.Block) *BCContext {
	ctx := NewBCContext(b)
	ctx.Parent = c
	return ctx
}

func (c *BCCompiler) CompileFile(f *syntax.File) {
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
	if c.CurrentSpace == nil || c.CurrentSpace.Ident != space {
		//TODO: make sure that the current space is not nil
		ctx := c.Global.NewBCContext(nil)
		c.CurrentSpace = ctx
		c.Spaces.Append(ctx)
	}
	if space == "main" {
		if Funcs.Filter(func(i int, v *syntax.FuncDecl) bool { return v.Name.Value == "main" }).First() == nil {
			println("program entry point was not found")
			os.Exit(1)
		}
	}
	//Types.Each(func(v *syntax.TypeDecl) { c.CompileType(space, v) })
	var FuncDefs = make([]*BCContext, 0)
	Funcs.Each(func(v *syntax.FuncDecl) { FuncDefs = append(FuncDefs, c.DefineFunc(space, v)) })
	Funcs.Iter(func(i int, v *syntax.FuncDecl) { c.CompileFunc(space, FuncDefs[i], v) })
}

func (c *BCCompiler) CompileFunc(space string, ctx *BCContext, def *syntax.FuncDecl) *backend.Func {
	for _, v := range def.Param {
		ctx.vars[v.Name.Value] = ctx.NewStore()
	}
	last, ok := def.Body.StmtList[len(def.Body.StmtList)-1].(*syntax.ReturnStmt)
	if !ok {
		last = &syntax.ReturnStmt{Return: nil}
		def.Body.StmtList = append(def.Body.StmtList, last)
	}

	c.CompileBody(ctx, def.Body)
	fn := ctx.Block.Parent
	ret := fn.Sig.RetType
	_ = ret
	return fn
}

func (c *BCCompiler) CompileBody(ctx *BCContext, body *syntax.BlockStmt) *BCContext {
	for _, v := range body.StmtList {
		c.CompileStmt(ctx, v)
	}
	return ctx
}

func (c *BCCompiler) CompileStmt(ctx *BCContext, s syntax.Stmt) backend.Value {
	switch s.(type) {
	case *syntax.ExprStmt:
		stmt := s.(*syntax.ExprStmt)
		return c.CompileExpr(ctx, stmt.X)
	case *syntax.EmptyStmt:
		//return nil
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
		//ctx.NewBr(ctx.parent.Block)

	case *syntax.ReturnStmt:
		//stmt := s.(*syntax.ReturnStmt)
		//return ctx.Block.NewRet(c.CompileExpr(ctx, stmt.Return)).X

	case *syntax.DeclStmt:
		stmt := s.(*syntax.DeclStmt)
		_ = stmt

	case *syntax.AssignStmt:
		stmt := s.(*syntax.AssignStmt)
		_ = stmt

	//case *syntax.IfStmt:
	//	stmt := s.(*syntax.IfStmt)
	//	c.CompileIf(ctx, stmt)

	//case *syntax.ForStmt:
	//	stmt := s.(*syntax.ForStmt)
	//	c.CompileFor(ctx, stmt)

	case *syntax.BlockStmt:
		stmt := s.(*syntax.BlockStmt)
		c.CompileBody(ctx, stmt)
	}
	return nil
}

func (c *BCCompiler) DefineFunc(space string, f *syntax.FuncDecl) *BCContext {
	m := c.Module
	name := f.Name.Value
	name = space + ".func." + name

	fn := m.NewFunc(name, nil, nil)
	return c.CurrentSpace.NewBCContext(fn.NewBlock("entry"))
}

func (c *BCCompiler) CompileExpr(ctx *BCContext, s syntax.Expr) backend.Value {
	switch s.(type) {
	case *syntax.CallExpr:
		expr := s.(*syntax.CallExpr)
		name := expr.Func.(*syntax.Name).Value
		callees := general.AsList(c.Module.Funcs).Filter(func(i int, v *backend.Func) bool { return v.Ident == name })
		if *callees == nil {
			println("function not found:", name)
			os.Exit(1)
		}
		callee := *(callees.First())
		var args []backend.Value
		general.AsList(expr.ArgList).Each(func(v syntax.Expr) { args = append(args, c.CompileExpr(ctx, v)) })
		ctx.NewCall(callee, args...)

	case *syntax.BasicLit:
		expr := s.(*syntax.BasicLit)
		switch expr.Kind {
		case syntax.StringLit:
			return backend.NewStringValue(expr.Value)
		case syntax.IntLit:
			v, err := strconv.ParseInt(expr.Value, 10, 64)
			if err != nil {
				panic("wrong value")
			}
			return backend.NewIntValue(types.I64, v)
		}
	}
	return nil
}

//func (c *BCCompiler) CompileType(space string, t *syntax.TypeDecl) {
//	name := space + ".type." + t.Name.Value
//	var typ types.Type
//	if t.Type == nil {
//		typ = types.Void
//	} else {
//		typ = c.QueryType(space, t.Type)
//	}
//	switch typ.(type) {
//	case *types.IntType:
//		c.Module.NewTypeDef(name, typ.(*types.IntType))
//	}
//}

//func (c *BCCompiler) CompileIf(ctx *BCContext, stmt *syntax.IfStmt) {
////	f := ctx.Block.Parent
////	thenCtx := ctx.NewBCContext(f.NewBlock("if"))
////	thenBlock := thenCtx.Block
////	elseCtx := ctx.NewBCContext(f.NewBlock("else"))
////	elseBlock := elseCtx.Block
////}
