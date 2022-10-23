package builtin

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

var (
	Int8   = types.I8
	Int    = types.I64
	Bool   = types.NewInt(1)
	String = types.NewStruct(types.I64, types.NewPointer(types.I8))
	VAList = types.NewStruct(types.I8Ptr)
	Float  = types.Float
	Void   = types.Void
)

type Type struct {
	types.Type
	Method map[string]*ir.Func
}

func (t *Type) AttachMethod(name string, f *ir.Func) {
	if t.Method == nil {
		t.Method = make(map[string]*ir.Func)
	}
	t.Method[name] = f
}

var ITypes = []Type{
	{Int8, nil},
	{Int, nil},
	{Bool, nil},
	{VAList, nil},
	{String, nil},
	{Float, nil},
}
