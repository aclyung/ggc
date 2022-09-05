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
	_const{"true", "bool", constant.NewInt(Bool, 1)},
	{"false", "bool", constant.NewInt(Bool, 0)},
}
