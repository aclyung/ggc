package compiler

import (
	"almeng.com/glang/core/syntax"
	"fmt"
	"github.com/almenglee/general"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

var errh = func(err error) { println(err.Error()) }

func Compile(filename string) {
	f, _ := os.Open(filename)
	// Node
	t := NewTarget(AARCH64, APPLE, DARWIN)
	c := NewCompiler(t)
	c.InitGlobal()
	file := syntax.Parse(filename, f, errh)
	// TODO: Node to llvm IR
	c.CompileFile(file)
	// TODO: link write file
	compiled := c.GetIR()

	println(compiled)
	tmpDir, err := ioutil.TempDir("", "glang")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(tmpDir+"/main.ll", []byte(compiled), 0644)
	if err != nil {
		panic(err)
	}

	out, err := os.UserHomeDir()
	if err != nil {
		return
	}

	clangArgs := []string{
		t.String(),
		"-Wno-override-module",
		tmpDir + "/main.ll",
		"-o", out + "/Desktop/exec", "-O3",
	}

	cmd := exec.Command("clang", clangArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		println(string(output))
		panic(err)
	}
	if len(output) > 0 {
		fmt.Println(string(output))
	}
	return
}

func (c *Compiler) CompileFile(f *syntax.File) {

}

func (c *Compiler) GetIR() string {
	return c.Module.String()
}

func NewCompiler(t Target) (c Compiler) {
	mod := ir.NewModule()
	c = Compiler{t, mod, nil, nil}
	c.InitGlobal()
	return
}

type Compiler struct {
	Target Target
	Module *ir.Module
	Spaces general.List[*Space]
	Global *Space
}

func (c *Compiler) InitGlobal() {

}

type Space struct {
	Name *syntax.Name
	Decl *general.List[syntax.Decl]
}

func CodeGen(node *syntax.File) {
	s := &Space{Name: node.SpaceName, Decl: general.NewList(node.DeclList...)}
	m := ir.NewModule()
	decl := func(d syntax.Decl) {
		ParseDecl(s, m, d)
	}
	s.Decl.Each(decl)
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
