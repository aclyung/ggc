package backend

import "almeng.com/glang/core/compiler/backend/types"

type Func struct {
	Ident  string
	Params []types.Type
	Sig    *types.FuncType
	Blocks []*Block
	Parent *Module
}

type Block struct {
	Parent *Func
	Ident  string
	Inst   []Instruction
}

func NewFunc(ident string, ret any, param ...types.Type) *Func {

	return &Func{Ident: ident, Blocks: make([]*Block, 0), Params: param}
}

func (f *Func) NewBlock(name string) *Block {
	b := &Block{Parent: f, Ident: f.Ident + name + ".block", Inst: make([]Instruction, 0)}
	f.Blocks = append(f.Blocks, b)
	return b
}

func (f *Func) BCString() string {
	str := "LABEL " + f.Ident
	for _, b := range f.Blocks {
		str += "\n" + b.BCString()
	}
	//TODO: implement BCString
	panic("not implemented")
}

func (b *Block) NewLoad() string {
}

func (b *Block) NewCall(callee *Func, args ...Value) Instruction {
	return NewCall(callee, args...)
}

func (b *Block) NewPrint() Instruction {
	return NewPrint()
}

func NewPrint() Instruction {
	panic("")
}

func (b *Block) NewPush(v Value) Instruction {
	return NewPush(v)
}

func (b *Block) BCString() string {
	str := ""
	for _, i := range b.Inst {
		str += i.BCString() + "\n"
	}
	return str
}
