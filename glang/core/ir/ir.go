package ir

import (
	"almeng.com/glang/core/ir/types"
	"fmt"
)

type Func struct {
	Ident      string
	Params     []*Param
	Sig        *types.FuncType
	Blocks     []*Block
	Parent     *Module
	IsVariadic bool
}

type Param struct {
	Ident string
	Type  types.Type
}

type Block struct {
	Parent *Func
	Ident  string
	Inst   []Instruction
}

func (b *Block) NewStore(name string) Instruction {
	i := NewStore(name)
	b.Inst = append(b.Inst, i)
	return i
}

func NewStore(name string) Instruction {
	return &InstStore{name}
}

func NewParam(ident string, t types.Type) *Param {
	return &Param{Ident: ident, Type: t}
}

func NewFunc(ident string, ret any, param ...*Param) *Func {

	return &Func{Ident: ident, Blocks: make([]*Block, 0), Params: param}
}

func (f *Func) NewBlock(name string) *Block {
	b := &Block{Parent: f, Ident: f.Ident + "." + name + ".block", Inst: make([]Instruction, 0)}
	f.Blocks = append(f.Blocks, b)
	return b
}

func (f *Func) BCString() string {
	str := "!META __FUNCTION__" + f.Ident + "\n"

	for _, b := range f.Blocks {
		str += "\n" + b.BCString()
	}
	return str
}

//func (b *Block) NewLoad() string {
//}

func (b *Block) NewEOF() Instruction {
	i := &InstUnary{"EOF"}
	b.Inst = append(b.Inst, i)
	return i
}

func (b *Block) NewPop() Instruction {
	i := &InstUnary{inst: "POP"}
	b.Inst = append(b.Inst, i)
	return i
}

func (b *Block) NewCall(callee *Func, args ...Value) Instruction {
	for i, _ := range args {
		// PUSHES the arguments
		b.NewPush(args[len(args)-i-1])
	}
	// PUSHES the number of arguments
	// Callee will pop this number of arguments
	//if len(args) > 0 {
	numArgs := NewIntValue(types.I64, int64(len(args)))
	b.NewPush(numArgs)

	i := NewCall(callee, args...)
	b.Inst = append(b.Inst, i)
	return i
}

func (b *Block) NewReturn() Instruction {
	i := &InstUnary{"RET"}
	b.Inst = append(b.Inst, i)
	return i

}

func (b *Block) NewPrint() Instruction {
	i := NewPrint()
	b.Inst = append(b.Inst, i)
	return i
}

func (b *Block) NewPrintln() Instruction {
	i := NewPrintln()
	b.Inst = append(b.Inst, i)
	return i
}

func NewPrint() Instruction {
	return InstUnary{inst: "PRINT"}
}

func NewPrintln() Instruction {
	return InstUnary{inst: "PRINTLN"}
}

func (b *Block) NewPush(v Value) Instruction {
	i := NewPush(v)
	b.Inst = append(b.Inst, i)
	return i
}

func (b *Block) BCString() string {
	str := fmt.Sprintf("LABEL %s\n", b.Ident)
	for _, i := range b.Inst {
		str += fmt.Sprintf("\t%s\n", i.BCString())
	}
	return str
}
