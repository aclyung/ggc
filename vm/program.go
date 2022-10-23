package vm

import (
	"encoding/binary"
	"reflect"
)

type Program struct {
	l       []byte
	index   int
	stack   *Stack[byte]
	vars    map[uint64][]byte
	retAddr *Stack[uint64]
	Idents  map[string]*Identifier
}

func NewProgram(l []byte) *Program {
	return &Program{
		l,
		0,
		NewStack[byte](),
		make(map[uint64][]byte),
		NewStack[uint64](),
		make(map[string]*Identifier),
	}
}

func (p *Program) Next() byte {
	if p.index >= len(p.l) {
		panic("index out of bounds")
	}
	defer func() { p.index++ }()
	return p.l[p.index]
}

func (p *Program) InstJump() {
	p.index = int(p.Int64())
}

func (p *Program) StackPopValue() []byte {
	return PopValue(p.stack)
}

func (p *Program) ProgramReadValue() (b []byte) {
	typ := DataType(p.Next())
	d_size := typ.Size()
	if typ.IsDynamic() {
		//Dynamically sized data
		size := p.Int64()
		b = p.ReadBuffer(int(size))
		for i, v := range b {
			b[i] = v ^ byte(i) ^ 0x80
		}
		b = append(WriteIntToBytes(size, 64), b...)
	} else {
		b = p.ReadBuffer(d_size)
	}
	b = append([]byte{byte(typ)}, b...)
	return
}

func (p *Program) InstRet() {
	p.index = int(p.retAddr.Pop())
}

func (p *Program) InstCall() {
	p.retAddr.Push(uint64(p.index + 8))
	p.InstJump()
}

func (p *Program) InstPush() {
	p.StackPushValue(p.ProgramReadValue())
	//typ := DataType(p.Next())
	//d_size := typ.Size()
	//if typ.IsDynamic() {
	//	//Dynamically sized data
	//	size := p.Int64()
	//	p.stack.PushAll(true, p.ReadBuffer(int(size))...)
	//	p.stack.PushAll(true, WriteIntToBytes(size, 64)...)
	//} else {
	//	p.stack.PushAll(true, p.ReadBuffer(d_size)...)
	//}
	//p.stack.Push(byte(typ))
}

func (p *Program) StackPushValue(v []byte) {
	p.stack.PushAll(true, v...)
}

func (p *Program) InstStore() {
	addr := p.Int64()
	p.vars[addr] = p.StackPopValue()
}

func (p *Program) InstCmp() {
	cond, a, b := p.Inst(), p.StackPopValue(), p.StackPopValue()
	switch cond {
	case EQ:
		if reflect.DeepEqual(a, b) {
			p.StackPushValue(TRUE)
		} else {
			p.StackPushValue(FALSE)
		}
	default:
		panic("invalid comparison instruction")
	}
}

func (p *Program) InstLoad() {
	addr := p.Int64()
	p.stack.PushAll(true, p.vars[addr]...)
}

func (p *Program) ReadBuffer(size int) []byte {
	defer func() { p.index += size }()
	return p.l[p.index : p.index+size]
}

func (p *Program) Inst() InstSet {
	r := InstSet(binary.BigEndian.Uint16(p.ReadBuffer(2)))
	if r >= 0xff00 && r < inst_end {
		return r
	}
	panic("invalid instruction")
}

func (p *Program) Int64() uint64 {
	r := binary.BigEndian.Uint64(p.ReadBuffer(8))
	return r
}

func (p *Program) InstBranch() {
	cond := Const(p.StackPopValue())
	if cond.Equal(TRUE) {
		p.InstJump()
	} else {
		p.Int64()
		p.InstJump()
	}
}
