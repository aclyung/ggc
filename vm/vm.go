//vm is Glang Virtual Machine
package vm

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var builtins = `
num => type number

LEN => int
APPEND num(both slice and elements) => slice

PRINTF num => void //string, ...any
PRINTLN  num => void //...any
PRINT num => void //...any
SPRINTF num => string // string,  ...any
SPRINT num => string // ...any
SPRINTLN num =>  string // ...any
SCANF num => void // string, ...any
`

type Mode uint8

const (
	Build Mode = iota
	ASM        //not yet
	BC         // .ir to .bc
	Run
)

var Modes = map[string]Mode{
	"build": Build,
	"run":   Run,
	"bc":    BC,
}

func init() {
	InstString = make(map[InstSet]string)
	for k, v := range Words {
		InstString[v] = k
	}
}

func ValidMode(m string) (mod Mode, b bool) {
	mod, b = Modes[m]
	return
}

var (
	pMode  Mode
	file   *os.File
	defers = make([]func(), 0)
)

func parseFlag() {
	args := os.Args[1:]
	if len(args) != 2 {
		panic("invalid args")
	}
	if mode, ok := Modes[args[0]]; ok {
		pMode = mode
	} else {
		panic("invalid arg")
	}
	open, err := os.Open(args[1])
	if err == nil {
		file = open
		defers = append(defers, func() { file.Close() })
		return
	}
	panic("invalid arg")
}

func ClearTmp() {
	for _, v := range defers {
		v()
	}
}

func CompileIR(ir string) []byte {
	ir = strings.ReplaceAll(ir, "\n", " ")
	asm := strings.Split(ir, " ")
	code := GenBC(asm)
	return code
}

func GenBC(code []string) []byte {
	s := code
	var r []byte
	addresses := make(map[string][]byte)
	jumps := make(map[int64]string)
	typedefs := make(map[string]string)
	for idx, v := range s {
		if v == "!skip" {
			continue
		}
		if i, ok := Words[v]; ok {
			switch i {
			case META:
				r = append(r, Uint16ToBytes(uint16(i))...)
				str := s[idx+1]
				s[idx+1] = "!skip"
				str = handleEscape(str)
				size := uint64(len(str))
				r = append(r, Uint64ToBytes(size+6)...)
				r = append(r, []byte("!META "+str)...)
				continue
			case TYPE:
				typedefs[s[idx+1]] = s[idx+2]
				s[idx+1] = "!skip"
				s[idx+2] = "!skip"
			case LABEL:
				addresses[s[idx+1]] = Uint64ToBytes(uint64(len(r)))
				s[idx+1] = "!skip"
				continue
			case JMP, CALL:
				r = append(r, Uint16ToBytes(uint16(i))...)
				jumps[int64(len(r))] = s[idx+1]
				r = append(r, Uint64ToBytes(0)...)
				s[idx+1] = "!skip"
				continue
			case BR:
				r = append(r, Uint16ToBytes(uint16(BR))...)
				jumps[int64(len(r))] = s[idx+1]
				r = append(r, Uint64ToBytes(0)...)
				s[idx+1] = "!skip"
				jumps[int64(len(r))] = s[idx+2]
				r = append(r, Uint64ToBytes(0)...)
				s[idx+2] = "!skip"
				continue
			case RET:
				r = append(r, Uint16ToBytes(uint16(RET))...)
				continue
			case CMP:
				r = append(r, Uint16ToBytes(uint16(CMP))...)
				// TODO: add support for other comparison operators
				r = append(r, Uint16ToBytes(uint16(Words[s[idx+1]]))...)
				s[idx+1] = "!skip"
				continue
			case LOAD, STORE:
				r = append(r, Uint16ToBytes(uint16(i))...)
				num, err := strconv.ParseInt(s[idx+1], 10, 64)
				if err != nil {
					panic(err)
				}
				r = append(r, Uint64ToBytes(uint64(num))...)
				s[idx+1] = "!skip"
				continue

			case PUSH:
				// PRE: PUSH TYPE VALUE

				// PUSH
				r = append(r, Uint16ToBytes(uint16(i))...)
				if typ, typeValid := LitTypes[s[idx+1]]; typeValid {
					// PUSH TYPE
					r = append(r, byte(typ))
					s[idx+1] = "!skip"
					if typ == STRING {
						// POST: PUSH TYPE SIZE STRING
						str := s[idx+2]
						s[idx+2] = "!skip"
						str = handleEscape(str)
						size := uint64(len(str))
						r = append(r, Uint64ToBytes(size)...)
						ss := []byte(str)
						for i, v := range ss {
							ss[i] = v ^ byte(i) ^ 0x80
						}
						r = append(r, ss...)
						continue
					}
					// PUSH TYPE VALUE(BIT_SIZE)
					n, err := strconv.ParseInt(s[idx+2], 10, 64)
					if err == nil {
						bitsize := LitDataByteSize[typ] * 8
						r = append(r, WriteIntToBytes(uint64(n), bitsize)...)
						s[idx+2] = "!skip"
						continue
					}
				}

				continue

			default:
				r = append(r, Uint16ToBytes(uint16(i))...)
				continue
			}

		}
		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			r = append(r, Uint64ToBytes(uint64(i))...)
			continue
		}

		r = append(r, []byte(v)...)
	}
	for k, v := range jumps {
		r = append(r[:k], append(addresses[v], r[k+8:]...)...)
	}
	return r
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

type ByteStack struct {
	Stack[byte]
}

type Stack[V any] struct {
	l []V
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

func Uint16ToBytes(v uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, v)
	return b
}

func NewStack[V any]() *Stack[V] {
	return &Stack[V]{l: make([]V, 0)}
}

func (s *Stack[V]) Head() V {
	return s.l[len(s.l)-1]
}

func (s *Stack[V]) Push(v V) {
	s.l = append(s.l, v)
}

func Reverse[V any](v []V) []V {
	rtn := make([]V, len(v))
	for i := len(v) - 1; i >= 0; i-- {
		rtn[len(rtn)-i-1] = v[i]
	}
	return rtn
}

func (s *Stack[V]) PopSized(reversed bool, size int) []V {
	defer func() { s.l = s.l[:len(s.l)-size] }()
	slice := s.l[len(s.l)-size:]
	if reversed {
		slice = Reverse(slice)
	}
	return slice
}

func (s *Stack[V]) PushAll(stackReverse bool, v ...V) {
	if stackReverse {
		v = Reverse(v)
	}
	s.l = append(s.l, v...)
}

func (s *Stack[V]) Pop() V {
	defer func() { s.l = s.l[:len(s.l)-1] }()
	if len(s.l) == 0 {
		panic("stack empty")
	}
	return s.l[len(s.l)-1]
}

type InstSet uint16

const (
	EOF InstSet = 0xff00 + iota
	META
	TYPE
	PUSH
	POP
	CONST
	LOAD
	LOCAL
	LABEL
	PRINT
	ADD
	SUB
	MUL
	DIV
	REM
	CALL
	RET
	JMP // unconditional jump
	BR  // conditional jump
	CMP
	EQ
	NOT
	GTR
	// STRING uft8 string indicator
	// STRING i64 len []byte
	//STRING
	STORE
	ARRAY
	UINT64
	INT64
	inst_end
)

type Type uint8

const (
	i8 Type = 1 + iota
	i16
	i32
	i64
	u8
	u16
	u32
	u64
)

var InstString map[InstSet]string

var Words = map[string]InstSet{
	"EOF":   EOF,
	"EQ":    EQ,
	"!META": META,
	"TYPE":  TYPE,
	"PUSH":  PUSH,
	"POP":   POP,
	"CONST": CONST,
	"LOAD":  LOAD,
	"STORE": STORE,
	"LOCAL": LOCAL,
	"LABEL": LABEL,
	"ADD":   ADD,
	"NOT":   NOT,
	"SUB":   SUB,
	"MUL":   MUL,
	"DIV":   DIV,
	"REM":   REM,
	"CALL":  CALL,
	"RET":   RET,
	"JMP":   JMP,
	"BR":    BR,
	"CMP":   CMP,
	"PRINT": PRINT,
	//"STRING": STRING,
	"ARRAY":  ARRAY,
	"UINT64": UINT64,
	"INT64":  INT64,
}

func (i InstSet) String() string {
	return InstString[i]
}

func readInt64(bytes []byte, index int64) uint64 {
	v := bytes[index : index+8]
	return binary.BigEndian.Uint64(v)
}

func Execute(hex []byte) {
	program := NewProgram(hex)
loop:
	for program.index < len(hex) {
		inst := program.Inst()
		switch inst {
		case EOF:
			if len(program.stack.l) == 0 {
				program.stack.Push(0)
			}
			break loop
		case META:
			program.index += int(program.Int64())
		case PUSH:
			program.InstPush()
		case CALL:
			program.InstCall()
		case RET:
			program.InstRet()
		case STORE:
			program.InstStore()
		case LOAD:
			program.InstLoad()

		case JMP:
			program.InstJump()
		case CMP:
			program.InstCmp()
		case PRINT:
			val := program.StackPopValue()

			str := string(val[9:])
			fmt.Print(str)
		case BR:
			program.InstBranch()
		case NOT:
			if DataType(program.stack.Head()) != BOOL {
				panic("not bool")
			}
			program.stack.l[len(program.stack.l)-2] = ^program.stack.l[len(program.stack.l)-2]
		case ADD, SUB, MUL, DIV, REM:
			b, a := program.StackPopValue(), program.StackPopValue()
			typ_a, typ_b := a[0], b[0]

			var res []byte
			if typ_b != typ_a {
				panic("Type mismatch")
			}

			if DataType(typ_a) == STRING {
				size_a, size_b := a[1:9], b[1:9]
				s_a, s_b := binary.BigEndian.Uint64(size_a), binary.BigEndian.Uint64(size_b)
				size := s_a + s_b
				a, b = a[9:], b[9:]

				if inst != ADD {
					panic("Operation not defined")
				}
				b_size := WriteIntToBytes(size, 64)
				_ = b_size
				res = append(b_size, append(a, b...)...)
				program.stack.PushAll(true, res...)
				program.stack.Push((byte)(STRING))
				continue
			}

			//switch inst {
			//case ADD:
			//	res = a + b
			//case SUB:
			//	res = a - b
			//case MUL:
			//	res = a * b
			//case DIV:
			//	res = a / b
			//case REM:
			//	res = a % b
			//}
			//program.stack.PushAll(true, res...)
			panic("error")

		}
	}
	if len(program.stack.l) != 1 && program.stack.Head() == 0 {
		//fmt.Println("exit code 1(stack not empty)")
		os.Exit(1)
	}
	os.Exit(0)
}

func FormatProgramData(program *Program, typ DataType) (b []byte) {
	switch typ {
	case UI8:
		b = []byte{program.Next()}
	case UI64:
		val := program.Int64()
		b = WriteIntToBytes(val, 64)
	case STRING:
		size := program.Int64()
		str := program.ReadBuffer(int(size))
		b = append(str, WriteIntToBytes(size, 64)...)
	default:
		panic("unexpected type")
	}
	return
}

//func SetData(v []byte, typ DataType) (b []byte) {
//	switch typ {
//	case UI8:
//		b = program.Next()
//	case UI64:
//		val := program.Int64()
//		program.stack.PushAll(true, WriteIntToBytes(val, 64)...)
//	case STRING:
//		size := program.Int64()
//		str := program.ReadBuffer(int(size))
//		program.stack.PushAll(true, str...)
//		program.stack.PushAll(true, WriteIntToBytes(size, 64)...)
//	default:
//		panic("unexpected type")
//	}
//}

func PopValue(s *Stack[byte]) (b []byte) {
	typ := s.Pop()
	return append([]byte{typ}, GetValue(s, DataType(typ))...)
}

func GetValue(s *Stack[byte], t DataType) (b []byte) {
	d_size := t.Size()
	if t.IsDynamic() {
		// read size of the data as int 64
		size_b := s.PopSized(true, 8)
		// case to int64
		size := binary.BigEndian.Uint64(size_b)
		// read data from stack with size
		data := s.PopSized(true, int(size))
		b = append(size_b, data...)
	} else {
		b = s.PopSized(true, d_size)
	}
	return
}

//func Run(insts []InstSet) {
//	programCounter := 0
//	stack := make([]int64, 0)
//	for programCounter < len(insts) {
//		ident := insts[programCounter]
//		switch ident {
//		case PUSH:
//			val := insts[programCounter+1]
//			stack = append(stack, val)
//			programCounter++
//		case ADD:
//			val := stack[len(stack)-1] + stack[len(stack)-2]
//			stack = append(stack, val)
//			programCounter++
//		}
//	}
//	println(stack[len(stack)-1])
//}

//type Instruction interface {
//	VMString() string
//}
//
//type InstOperation interface {
//}
//
//type InstFunction struct {
//	Name         string
//	Params       []*Param
//	Instructions []*Instruction
//	RtnType      string
//}
//
//func (f *InstFunction) VMString() string {
//	str := "FUCN"
//
//}
//
//type InstCAll struct {
//	Callee string
//	Args[]
//}
//
//type Param struct {
//	Name string
//	Type string
//}
