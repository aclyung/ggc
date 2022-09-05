package builtin

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

func RegisterFunc(m *ir.Module, f *ir.Func) {
	f.Parent = m
	m.Funcs = append(m.Funcs, f)
}

var Funcs = map[string]**ir.Func{
	"printf":  &Printf,
	"println": &Println,
	"print":   &Print,
}

var Printf = func() *ir.Func {
	f := ir.NewFunc(
		"printf",
		types.I32,
		ir.NewParam("", types.NewPointer(types.I8)),
	)
	f.Sig.Variadic = true
	return f
}()

var Println *ir.Func
var Print *ir.Func
