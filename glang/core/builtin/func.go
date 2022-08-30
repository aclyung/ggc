package buitin

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func RegisterFunc(m *ir.Module, f *ir.Func) {
	f.Parent = m
	m.Funcs = append(m.Funcs, f)
}

var Funcs = map[string]*ir.Func{
	"printf": Printf,
	"print":  Print,
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

func _println() *ir.Func {
	f := ir.NewFunc(
		"println",
		types.Void,
		ir.NewParam("", types.NewPointer(types.I8)),
	)
	b := f.NewBlock("")
	var prams []value.Value
	for _, v := range f.Params {
		prams = append(prams, v)
	}
	b.NewCall(Printf, prams...)
	c := b.NewGetElementPtr(NewLine.ContentType, NewLine, constant.NewInt(Int, 0), constant.NewInt(Int, 0))
	b.NewCall(Printf, c)
	b.NewRet(nil)
	return f
}

var Print = func() *ir.Func {
	f := ir.NewFunc(
		"print",
		types.Void,
		ir.NewParam("", types.NewPointer(types.I8)),
	)
	b := f.NewBlock("")
	var prams []value.Value
	for _, v := range f.Params {
		prams = append(prams, v)
	}
	b.NewCall(Printf, prams...)
	b.NewRet(nil)
	return f
}()
