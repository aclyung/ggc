package backend

import (
	vm "almeng.com/glang-vm"
	"almeng.com/glang/core/compiler/backend/types"
	"fmt"
)

type Instruction interface {
	String() string
	BCString() string
}

type (
	// InstValue
	// PUSH | LOAD | STORE | CALL | JMP
	InstValue struct {
		vm.InstSet
		v Value
	}

	// InstUnary
	// RET | ADD | SUB | MUL | DIV | REM | POP
	InstUnary struct {
		inst string
	}

	// InstTernary
	// BR
	InstTernary struct {
		a, b uint64
	}

	InstLABEL struct {
	}

	InstCALL struct {
		Callee *Func
		Args   []Value
	}

	InstCMP struct {
	}
)

func (i InstCALL) String() string {
	//TODO implement me
	panic("implement me")

}

func (i InstCALL) BCString() string {
	str := ""
	for _, v := range i.Args {
		// PUSHES the arguments
		str += NewPush(v).BCString() + "\n"
	}
	numArgs := NewIntValue(types.I64, int64(len(i.Args)))
	// PUSHES the number of arguments
	// Callee will pop this number of arguments
	str += NewPush(numArgs).BCString() + "\n"
	str += "CALL " + i.Callee.Ident + "\n"
	return str
}

func (i InstValue) String() string {
	//TODO implement me
	panic("implement me")
}

func (i InstValue) BCString() string {
	return fmt.Sprintf("%s %s", i.InstSet.String(), i.v.BCString())
}

func NewPush(v Value) InstValue {
	return InstValue{vm.PUSH, v}
}

func NewLoad(ident string) InstValue {
	return InstValue{vm.LOAD, nil}
}

func NewCall(callee *Func, args ...Value) InstCALL {
	return InstCALL{callee, args}
}

func NewAdd() InstUnary {
	return InstUnary{"ADD"}
}

func NewSub() InstUnary {
	return InstUnary{"SUB"}
}

func NewMul() InstUnary {
	return InstUnary{"MUL"}
}

func NewDiv() InstUnary {
	return InstUnary{"DIV"}
}
