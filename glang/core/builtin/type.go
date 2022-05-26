package buitin

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

var (
	Bool   types.Type = types.NewInt(8)
	String types.Type = types.NewStruct(types.I64, types.NewPointer(types.I8))
	Int    types.Type = types.I64
	Float  types.Type = types.Float
	Type   types.Type
)

func InitModule(m *ir.Module) {
	m.NewTypeDef("bool", Bool)
	m.NewTypeDef("string", String)
	m.NewTypeDef("int", Int)
	m.NewTypeDef("float", Float)
}
