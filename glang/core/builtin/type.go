package builtin

import (
	"github.com/llir/llvm/ir/types"
)

var (
	Int8   = types.I8
	Int    = types.I64
	Bool   = types.NewInt(8)
	String = types.NewStruct(types.I64, types.NewPointer(types.I8))
	Float  = types.Float
	Void   = types.Void
)

type TypeDef struct {
	N string
	T types.Type
}

var ITypes = []TypeDef{
	TypeDef{"int8", Int8},
	TypeDef{"int", Int},
	TypeDef{"bool", Bool},
	TypeDef{"string", String},
	TypeDef{"float", Float},
}
