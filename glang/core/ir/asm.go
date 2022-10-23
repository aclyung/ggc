package ir

import (
	vm "almeng.com/glang-vm"
	"encoding/binary"
	"strconv"
	"strings"
)

type Assembler struct {
	code   []string
	bc     []byte
	Idents map[string]*Identifier
	jumps  map[int64]string
}
type Identifier struct {
	counter uint64
	ident   map[string]uint64
}

func (i *Identifier) Assign(name string) {
	i.ident[name] = i.counter
	i.counter++
}

func NewIdentifier() *Identifier {
	return &Identifier{counter: 0, ident: make(map[string]uint64)}
}

func (a *Assembler) NewLabelIdent(name string) {
	a.NewIdent("LABEL", name)
	a.Idents["LABEL"].ident[name] = uint64(len(a.bc))
}

func (a *Assembler) NewIdent(typ string, name string) uint64 {
	var ident *Identifier
	if _, ok := a.Idents[typ]; !ok {
		a.Idents[typ] = NewIdentifier()
	}
	ident = a.Idents[typ]
	if _, ok := ident.ident[name]; !ok {
		ident.Assign(name)

	}
	return ident.ident[name]
}

func ParseIR(ir string) []string {
	ir = strings.ReplaceAll(ir, "\n", " ")
	ir = strings.ReplaceAll(ir, "\t", "")
	return strings.Split(ir, " ")
}

func NewAssembler(ir string) *Assembler {
	code := ParseIR(ir)
	return &Assembler{code, make([]byte, 0), make(map[string]*Identifier), make(map[int64]string)}
}

func (a *Assembler) Put(v ...byte) {
	a.bc = append(a.bc, v...)
}

func (a *Assembler) isEOF() bool {
	return len(a.code) == 0
}

func (a *Assembler) Next() string {
	v := a.code[0]
	a.code = a.code[1:]
	return v
}

func (a *Assembler) PutInst(inst vm.InstSet) {
	a.Put(Uint16ToBytes(uint16(inst))...)
}

func (a *Assembler) NewJump() {
	a.jumps[int64(len(a.bc))] = a.Next()
	a.Put(Uint64ToBytes(0)...)
}

func (a *Assembler) SetJump() {
	for k, v := range a.jumps {
		a.bc[k] = byte(a.Idents["LABEL"].ident[v])
	}
}

func (a *Assembler) GenBC() []byte {
	as := &Assembler{a.code, make([]byte, 0), make(map[string]*Identifier), make(map[int64]string)}
	for !as.isEOF() {
		if inst, ok := vm.Words[as.Next()]; ok {
			switch inst {
			case vm.META:
				as.PutInst(inst)
				str := as.Next()
				str = handleEscape(str)
				size := uint64(len(str))
				as.Put(Uint64ToBytes(size + 6)...)
				as.Put([]byte("!META " + str)...)
				continue
			case vm.LABEL:
				as.NewLabelIdent(as.Next())
				continue
			case vm.JMP, vm.CALL:
				as.PutInst(inst)
				as.NewJump()
				continue
			case vm.BR:
				as.PutInst(inst)
				as.NewJump()
				as.NewJump()
				continue
			case vm.RET:
				as.PutInst(inst)
				continue
			case vm.CMP:
				as.PutInst(inst)
				as.Put(vm.Uint16ToBytes(uint16(vm.Words[as.Next()]))...)
				continue
			case vm.LOAD, vm.STORE:
				as.PutInst(inst)
				ident := as.NewIdent("VAR", as.Next())
				as.Put(Uint64ToBytes(ident)...)
				continue
			case vm.PUSH:
				as.PutInst(inst)
				if typ, typeValid := vm.LitTypes[as.Next()]; typeValid {
					as.Put(byte(typ))
					if typ == vm.STRING {
						str := as.Next()
						str = handleEscape(str)
						size := uint64(len(str))
						as.Put(Uint64ToBytes(size)...)
						ss := []byte(str)
						for i, v := range ss {
							ss[i] = v ^ byte(i) ^ 0x80
						}
						as.Put(ss...)
						continue
					}
					n, err := strconv.ParseInt(as.Next(), 10, 64)
					if err == nil {
						bitsize := vm.LitDataByteSize[typ] * 8
						as.Put(WriteIntToBytes(uint64(n), bitsize)...)
						continue
					}
				}
			default:
				as.PutInst(inst)
				continue
			}
		}
	}
	as.CompileJump()
	return as.bc
}

func WriteIntToBytes(v uint64, bitsize int) []byte {
	switch bitsize {
	case 8:
		return []byte{byte(v)}
	case 16:
		return Uint64ToBytes(v)
	case 64:
		return Uint64ToBytes(v)
	}
	panic("unreachable")
}

func Uint64ToBytes(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func handleEscape(s string) string {
	s = strings.ReplaceAll(s, `\n`, "\n")
	s = strings.ReplaceAll(s, `\t`, "\t")
	s = strings.ReplaceAll(s, `\r`, "\r")
	s = strings.ReplaceAll(s, `\0`, "\000")
	s = strings.ReplaceAll(s, `\20`, " ")
	s = strings.ReplaceAll(s, `\\`, "\\")
	return s
}

func Uint16ToBytes(v uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, v)
	return b
}

func (a *Assembler) CompileJump() {
	for k, v := range a.jumps {
		a.bc = append(a.bc[:k], append(Uint64ToBytes(a.NewIdent("LABEL", v)), a.bc[k+8:]...)...)
	}
}
