package vm

import (
	"encoding/binary"
	"reflect"
)

type Identifier struct {
	counter uint64
	ident   map[string]uint64
}

func NewIdentifier() *Identifier {
	return &Identifier{counter: 0, ident: make(map[string]uint64)}
}

func ExtractValue(b []byte) (typ DataType, data []byte) {
	typ = DataType(b[0])
	data = b[1:]
	if typ.IsDynamic() {
		size := binary.BigEndian.Uint64(data[:8])
		data = data[8 : 8+size]
	}
	return
}

type VM struct {
	program []byte
	index   int
	stack   *Stack[byte]
	vars    map[uint64][]byte
	retAddr *Stack[uint64]
	Idents  map[string]*Identifier
}

func NewVM(code []byte) *VM {
	return &VM{
		code,
		0,
		NewStack[byte](),
		make(map[uint64][]byte),
		NewStack[uint64](),
		make(map[string]*Identifier),
	}
}

func (vm *VM) Next() byte {
	if vm.index >= len(vm.program) {
		panic("index out of bounds")
	}
	defer func() { vm.index++ }()
	return vm.program[vm.index]
}

func (vm *VM) InstJump() {
	vm.index = int(vm.Int64())
}

func (vm *VM) StackPopValue() []byte {
	return PopValue(vm.stack)
}

func (vm *VM) ProgramReadValue() (b []byte) {
	typ := DataType(vm.Next())
	d_size := typ.Size()
	if typ.IsDynamic() {
		//Dynamically sized data
		size := vm.Int64()
		b = vm.ReadBuffer(int(size))
		for i, v := range b {
			b[i] = v //^ byte(i) ^ 0x80
		}
		b = append(WriteIntToBytes(size, 64), b...)
	} else {
		b = vm.ReadBuffer(d_size)
	}
	b = append([]byte{byte(typ)}, b...)
	return
}

func (vm *VM) InstRet() {
	vm.index = int(vm.retAddr.Pop())
}

func (vm *VM) InstCall() {
	vm.retAddr.Push(uint64(vm.index + 8))
	vm.InstJump()
}

func (vm *VM) InstPush() {
	vm.StackPushValue(vm.ProgramReadValue())
}

func (vm *VM) StackPushValue(v []byte) {
	vm.stack.PushAll(true, v...)
}

func (vm *VM) InstStore() {
	addr := vm.Int64()
	vm.vars[addr] = vm.StackPopValue()
}

func (vm *VM) InstCmp() {
	cond, a, b := vm.Inst(), vm.StackPopValue(), vm.StackPopValue()
	switch cond {
	case EQ:
		if reflect.DeepEqual(a, b) {
			vm.StackPushValue(TRUE)
		} else {
			vm.StackPushValue(FALSE)
		}
	default:
		panic("invalid comparison instruction")
	}
}

func (vm *VM) InstLoad() {
	addr := vm.Int64()
	vm.stack.PushAll(true, vm.vars[addr]...)
}

func (vm *VM) ReadBuffer(size int) []byte {
	defer func() { vm.index += size }()
	return vm.program[vm.index : vm.index+size]
}

func (vm *VM) Inst() InstSet {
	r := InstSet(binary.BigEndian.Uint16(vm.ReadBuffer(2)))
	if r >= 0xff00 && r < inst_end {
		return r
	}
	panic("invalid instruction")
}

func (vm *VM) Int64() uint64 {
	r := binary.BigEndian.Uint64(vm.ReadBuffer(8))
	return r
}

func (vm *VM) InstBranch() {
	cond := Const(vm.StackPopValue())
	if cond.Equal(TRUE) {
		vm.InstJump()
	} else {
		vm.Int64()
		vm.InstJump()
	}
}
