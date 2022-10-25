package ir

import (
	vm "almeng.com/glang-vm"
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
	}

	InstCMP struct {
	}
	InstStore struct {
		Ident string
	}
)

func (i InstUnary) String() string {
	//TODO implement me
	panic("implement me")
}

func (i InstUnary) BCString() string {
	return fmt.Sprintf("%s", i.inst)
}

func (i InstStore) String() string {
	//TODO implement me
	panic("implement me")
}

func (i InstStore) BCString() string {
	return fmt.Sprintf("STORE %s", i.Ident)
}

func (i InstCALL) String() string {
	//TODO implement me
	panic("implement me")

}

func (i InstCALL) BCString() string {
	return "CALL " + i.Callee.Ident + "\n"

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
	return InstCALL{callee}
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
