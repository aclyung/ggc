// Package global is unsafe
package global

import (
	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

var _mod *ir.Module
var cntr = 0

func Module() *ir.Module {
	return _mod
}

// Init function should be called at least once.
func Init(m *ir.Module) {
	_mod = m
}

func NewGlobalCharArrayConstant(s string) *ir.Global {
	c := constant.NewCharArrayFromString(s + "\000")
	n := fmt.Sprint(".str.", cntr)
	cntr += 1
	str := _mod.NewGlobalDef(n, c)
	return str
}

func NewGlobalString(b *ir.Block, s string) *ir.InstGetElementPtr {
	c := constant.NewCharArrayFromString(s + "\000")
	n := fmt.Sprint(".str.", cntr)
	cntr += 1
	str := _mod.NewGlobalDef(n, c)
	strPtr := b.NewGetElementPtr(str.ContentType, str, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
	return strPtr
}

func NewLocalString(b *ir.Block, s string) *ir.InstGetElementPtr {
	c := constant.NewCharArrayFromString(s + "\000")
	str := b.NewAlloca(c.Typ)
	b.NewStore(c, str)
	strPtr := b.NewGetElementPtr(str.ElemType, str, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
	return strPtr
}
