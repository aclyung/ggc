package builtin

import (
	"github.com/llir/llvm/ir/constant"
)

type _const struct {
	Name   string
	Type   string
	IConst constant.Constant
}

var Consts = []_const{
	True,
	False,
}

var (
	True  = _const{"true", "bool", constant.NewInt(Bool, 1)}
	False = _const{"false", "bool", constant.NewInt(Bool, 0)}
)
