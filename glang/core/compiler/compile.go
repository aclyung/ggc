package compiler

import (
	buitin "almeng.com/glang/core/builtin"
	"almeng.com/glang/core/syntax"
	"almeng.com/glang/global"
	"fmt"
	"github.com/almenglee/general"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (c *Compiler) CompileFile(f *syntax.File) {

	for _, v := range f.DeclList {
		switch v.(type) {
		case *syntax.FuncDecl:
			c.CompileFunc(f.SpaceName.Value, v.(*syntax.FuncDecl))
		case *syntax.OperDecl:
			c.CompileOper(f.SpaceName.Value, v.(*syntax.OperDecl))
		case *syntax.VarDecl:

		case *syntax.TypeDecl:
			c.CompileType(f.SpaceName.Value, v.(*syntax.TypeDecl))
		}
	}
}

func ParseDecl(s *Space, m *ir.Module, decl syntax.Decl) {
	switch decl.(type) {
	case *syntax.FuncDecl:
		f := decl.(*syntax.FuncDecl)
		name := f.Name.Value
		var ret types.Type = types.Void
		if f.Return != nil {
			ret = RetType(m, f.Return.(*syntax.Name))
		}

		m.NewFunc(name, ret)
	case *syntax.VarDecl:
		_ = decl.(*syntax.VarDecl)
	case *syntax.OperDecl:
		_ = decl.(*syntax.OperDecl)

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

func (c *Compiler) CompileOper(space string, f *syntax.OperDecl) {
	name := space + ".oper." + f.TypeL.Type.(*syntax.Name).Value + "_" + f.TypeR.Type.(*syntax.Name).Value + "_" + f.Oper.String()
	c.Module.NewFunc(name, types.Void)
}

func (c *Compiler) CompileFunc(space string, f *syntax.FuncDecl) *ir.Func {
	m := c.Module
	decl := f
	name := decl.Name.Value
	if space != "main" || name != "main" {
		name = space + ".func." + name
	}
	//var ret types.Type
	//ret = c.QueryType(space, f.Return)
	Func := m.NewFunc(name, types.Void)
	entry := c.CompileBody(Func, decl.Body)
	entry.NewRet(nil)
	_ = entry
	return Func
}

func (c *Compiler) CompileType(space string, t *syntax.TypeDecl) {
	name := space + ".type." + t.Name.Value
	var tt types.Type
	if t.Type == nil {
		tt = types.Void
	} else {
		tt = c.QueryType(space, t.Type)
	}

	c.Module.NewTypeDef(name, tt)
}

func (c *Compiler) QueryType(space string, t syntax.Expr) types.Type {
	if t == nil {
		return types.Void
	}
	return types.NewInt(64)
}

var str *ir.InstGetElementPtr

func (c *Compiler) CompileBody(f *ir.Func, body *syntax.BlockStmt) *ir.Block {
	b := f.NewBlock("")
	str = global.NewGlobalString(b, "hello world!")
	for _, v := range body.StmtList {
		c.CompileStmt(b, v)
	}
	return b
}
func (c *Compiler) CompileExpr(p *ir.Block, s syntax.Expr) ir.Instruction {
	// TODO
	return nil
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

func (c *Compiler) CompileStmt(p *ir.Block, s syntax.Stmt) ir.Instruction {
	switch s.(type) {
	case *syntax.ExprStmt:
		stmt := s.(*syntax.ExprStmt)
		c.CompileExpr(p, stmt.X)
	case *syntax.EmptyStmt:
		stmt := s.(*syntax.EmptyStmt)
		_ = stmt
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
		_ = stmt

	case *syntax.DeclStmt:
		stmt := s.(*syntax.DeclStmt)
		_ = stmt

	case *syntax.AssignStmt:
		stmt := s.(*syntax.AssignStmt)
		_ = stmt

	case *syntax.IfStmt:
		stmt := s.(*syntax.IfStmt)
		_ = stmt

	case *syntax.ForStmt:
		stmt := s.(*syntax.ForStmt)
		_init, _cond, _post := stmt.Init, stmt.Cond, stmt.Post
		_, _, _ = _init, _cond, _post

	case *syntax.BlockStmt:
		stmt := s.(*syntax.BlockStmt)
		_ = stmt
	}
	return p.NewCall(buitin.Println, str)
}
