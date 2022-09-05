package compiler

import (
	"almeng.com/glang/core/builtin"
	"almeng.com/glang/core/syntax"
	"almeng.com/glang/global"
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

func NewCompiler(t Target, verbose bool) *Compiler {
	mod := ir.NewModule()
	return &Compiler{t, mod, nil, nil, nil, NewContext(ir.NewBlock("global")), nil, verbose}
}

func (c *Compiler) InitGlobal() {
	m := c.Module
	for _, v := range builtin.ITypes {
		m.NewTypeDef(v.N, v.T)
	}
	for _, v := range builtin.Consts {
		c.Global.vars[v.Name] = m.NewGlobalDef(v.Name, v.IConst)
	}

	builtin.NewLine = global.NewGlobalCharArrayConstant("\n")
	builtin.Println = builtin.InitPrintln()
	builtin.Print = builtin.InitPrint()
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
	global.Init(c.Module)
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

func (c *Compiler) GetIR() string {
	return c.Module.String()
}
