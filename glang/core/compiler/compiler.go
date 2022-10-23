package compiler

import (
	"almeng.com/glang/core/builtin"
	"almeng.com/glang/core/syntax"
	"fmt"
	"github.com/almenglee/general"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"io"
	"os"
)

var errh = func(err error) { println(err.Error()); os.Exit(-1) }

type Compiler struct {
	Target       Target
	Module       *ir.Module
	Main         *ir.Func
	Spaces       general.List[*Context]
	CurrentSpace *Context
	Global       *Context
	Opers        general.List[*Operator]
	verbose      bool
}

var cntr = 0

func NewCompiler(t Target, verbose bool) *Compiler {
	mod := ir.NewModule()
	return &Compiler{t, mod, nil, nil, nil, NewContext(ir.NewBlock("global")), nil, verbose}
}

func (c *Compiler) NewGlobalCharArrayConstant(s string) *ir.Global {
	con := constant.NewCharArrayFromString(s + "\000")
	n := fmt.Sprint(".str.", cntr)
	cntr++
	_str := c.Module.NewGlobalDef(n, con)
	return _str
}

func (c *Compiler) NewGlobalString(b *ir.Block, s string) *ir.InstGetElementPtr {
	con := constant.NewCharArrayFromString(s + "\000")
	n := fmt.Sprint(".str.", cntr)
	cntr++
	_str := c.Module.NewGlobalDef(n, con)
	strPtr := b.NewGetElementPtr(_str.ContentType, _str, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))

	return strPtr
}

func (c *Compiler) NewLocalString(b *ir.Block, s string) *ir.InstGetElementPtr {
	con := constant.NewCharArrayFromString(s + "\000")
	_str := b.NewAlloca(con.Typ)
	b.NewStore(con, _str)
	strPtr := b.NewGetElementPtr(_str.ElemType, _str, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
	return strPtr
}

func (c *Compiler) InitGlobal() {
	m := c.Module
	for _, v := range builtin.ITypes {
		m.NewTypeDef(v.Name(), v)
	}
	for _, v := range builtin.Consts {
		c.Global.vars[v.Name] = m.NewGlobalDef(v.Name, v.IConst)
	}

	builtin.NewLine = c.NewGlobalCharArrayConstant("\n")
	builtin.Println = c.InitPrintln()
	builtin.Print = c.InitPrint()
	for _, v := range builtin.Funcs {
		builtin.RegisterFunc(m, *v)
	}
	c.InitOper()
}

type Space struct {
	Name *syntax.Name
	Decl *general.List[syntax.Decl]
}

func Compile(filename string, verbose bool, triple string) *Compiler {
	f, ferr := os.Open(filename)
	if ferr != nil {
		println(ferr.Error())
		os.Exit(-1)
	}
	// Node
	target := TargetFromTriple(triple)
	c := NewCompiler(target, verbose)
	c.InitGlobal()
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
	file := syntax.Parse(f, errh, c.verbose)

	c.CompileFile(file)
	entrypoint := c.Module.NewFunc("main", types.I32)
	block := entrypoint.NewBlock("")
	block.NewCall(c.Main)
	block.NewRet(constant.NewInt(types.I32, 0))

	return c
}

func (c *Compiler) InitPrintln() *ir.Func {
	f := ir.NewFunc(
		"println",
		types.Void,
		//ir.NewParam("str", types.I8Ptr),
	)
	f.Sig.Variadic = true
	b := f.NewBlock("")
	args := b.NewAlloca(types.NewArray(1, builtin.VAList))
	args.Align = 16
	ptr := b.NewGetElementPtr(args.ElemType, args, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	ptr.InBounds = true
	p := b.NewBitCast(ptr, types.I8Ptr)
	b.NewCall(builtin.VaStart, p)
	b.NewCall(builtin.Printf, c.NewGlobalString(b, "%d\n"), args)
	num := b.NewVAArg(args, types.I8)
	num = b.NewVAArg(args, types.I8)
	b.NewCall(builtin.Printf, c.NewGlobalString(b, "%d\n"), num)
	list_ptr := b.NewGetElementPtr(args.ElemType, args, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	list_ptr.InBounds = true

	println(list_ptr.Type().String())
	va := b.NewGetElementPtr(builtin.VAList, list_ptr, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	va.InBounds = true
	a := b.NewLoad(types.I8Ptr, va)
	n := b.NewBitCast(a, types.NewPointer(types.I8Ptr))
	l := b.NewLoad(types.I8Ptr, n)

	b.NewCall(builtin.Printf, l)
	b.NewCall(builtin.VaEnd, p)
	b.NewCall(builtin.Printf, c.NewGlobalString(b, "%d\n"), p)
	nl := b.NewGetElementPtr(builtin.NewLine.ContentType, builtin.NewLine, constant.NewInt(builtin.Int, 0), constant.NewInt(builtin.Int, 0))
	b.NewCall(builtin.Printf, nl)
	b.NewRet(nil)
	return f
}

func (c *Compiler) InitPrint() *ir.Func {
	f := ir.NewFunc(
		"print",
		types.Void,
		//ir.NewParam("str", types.I8Ptr),
	)
	f.Sig.Variadic = true
	b := f.NewBlock("")
	args := b.NewAlloca(types.NewArray(1, builtin.VAList))
	args.Align = 16
	ptr := b.NewGetElementPtr(args.ElemType, args, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	ptr.InBounds = true
	p := b.NewBitCast(ptr, types.I8Ptr)
	b.NewCall(builtin.VaStart, p)
	list_ptr := b.NewGetElementPtr(args.ElemType, args, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	list_ptr.InBounds = true
	println(list_ptr.Type().String())
	va := b.NewGetElementPtr(builtin.VAList, list_ptr, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	va.InBounds = true
	//b.NewCall(builtin.VaEnd, p)
	a := b.NewLoad(types.I8Ptr, va)
	n := b.NewBitCast(a, types.NewPointer(types.I8Ptr))
	l := b.NewLoad(types.I8Ptr, n)
	b.NewCall(builtin.VaEnd, p)
	b.NewCall(builtin.Printf, l)
	b.NewRet(nil)
	return f
}

func (c *Compiler) GetIR() string {
	return c.Module.String()
}
