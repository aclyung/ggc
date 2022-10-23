package vm

import "strconv"

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

func NewAssembler(code []string) *Assembler {
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

func (a *Assembler) PutInst(inst InstSet) {
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
	for !a.isEOF() {
		if inst, ok := words[a.Next()]; ok {
			switch inst {
			case META:
				a.PutInst(inst)
				str := a.Next()
				str = handleEscape(str)
				size := uint64(len(str))
				a.Put(Uint64ToBytes(size + 6)...)
				a.Put([]byte("!META " + str)...)
				continue
			case LABEL:
				a.NewLabelIdent(a.Next())
				continue
			case JMP, CALL:
				a.PutInst(inst)
				a.NewJump()
				continue
			case BR:
				a.PutInst(inst)
				a.NewJump()
				a.NewJump()
				continue
			case RET:
				a.PutInst(inst)
				continue
			case CMP:
				a.PutInst(inst)
				a.Put(Uint16ToBytes(uint16(words[a.Next()]))...)
				continue
			case LOAD, STORE:
				a.PutInst(inst)
				ident := a.NewIdent("VAR", a.Next())
				a.Put(Uint64ToBytes(ident)...)
				continue
			case PUSH:
				a.PutInst(inst)
				if typ, typeValid := LitTypes[a.Next()]; typeValid {
					a.Put(byte(typ))
					if typ == STRING {
						str := a.Next()
						str = handleEscape(str)
						size := uint64(len(str))
						a.Put(Uint64ToBytes(size)...)
						ss := []byte(str)
						for i, v := range ss {
							ss[i] = v ^ byte(i) ^ 0x80
						}
						a.Put(ss...)
						continue
					}
					n, err := strconv.ParseInt(a.Next(), 10, 64)
					if err == nil {
						bitsize := LitDataByteSize[typ] * 8
						a.Put(WriteIntToBytes(uint64(n), bitsize)...)
						continue
					}
				}
			default:
				a.PutInst(inst)
				continue
			}
		}
	}
	a.CompileJump()
	return a.bc
}

func (a *Assembler) CompileJump() {
	for k, v := range a.jumps {
		a.bc = append(a.bc[:k], append(Uint64ToBytes(a.NewIdent("LABEL", v)), a.bc[k+8:]...)...)
	}
}
