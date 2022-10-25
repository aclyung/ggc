//vm is Glang Virtual Machine
package vm

import (
	"encoding/binary"
	"fmt"
	"os"
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

func init() {
	InstString = make(map[InstSet]string)
	for k, v := range Words {
		InstString[v] = k
	}
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
	PRINTLN
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
	STORE
	ARRAY
	UINT64
	INT64
	inst_end
)

var InstString map[InstSet]string

var Words = map[string]InstSet{
	"EOF":     EOF,
	"EQ":      EQ,
	"!META":   META,
	"TYPE":    TYPE,
	"PUSH":    PUSH,
	"POP":     POP,
	"CONST":   CONST,
	"LOAD":    LOAD,
	"STORE":   STORE,
	"LOCAL":   LOCAL,
	"LABEL":   LABEL,
	"ADD":     ADD,
	"NOT":     NOT,
	"SUB":     SUB,
	"MUL":     MUL,
	"DIV":     DIV,
	"REM":     REM,
	"CALL":    CALL,
	"RET":     RET,
	"JMP":     JMP,
	"BR":      BR,
	"CMP":     CMP,
	"PRINT":   PRINT,
	"PRINTLN": PRINTLN,
	//"STRING": STRING,
	"ARRAY":  ARRAY,
	"UINT64": UINT64,
	"INT64":  INT64,
}

func (i InstSet) String() string {
	return InstString[i]
}

func (vm *VM) Execute() {
	length := len(vm.program)
loop:
	for vm.index < length {
		inst := vm.Inst()
		switch inst {
		case EOF:
			if len(vm.stack.l) == 0 {
				vm.stack.Push(0)
			}
			break loop
		case META:
			vm.index += int(vm.Int64())
		case POP:
			vm.StackPopValue()
		case PUSH:
			vm.InstPush()
		case CALL:
			vm.InstCall()
		case RET:
			vm.InstRet()
		case STORE:
			vm.InstStore()
		case LOAD:
			vm.InstLoad()

		case JMP:
			vm.InstJump()
		case CMP:
			vm.InstCmp()
		case PRINT, PRINTLN:
			num := binary.BigEndian.Uint64(vm.StackPopValue()[1:])
			str := ""
			for i := 0; i < int(num); i++ {
				typ, data := ExtractValue(vm.StackPopValue())
				switch typ {
				case UI8:
					str += fmt.Sprint(data[0])
				case UI64:
					str += fmt.Sprint(binary.BigEndian.Uint64(data))
				case BOOL:
					if TRUE.Equal(data) {
						str += "true"
					} else {
						str += "false"
					}
				case STRING:
					str += string(data)
				}
			}
			if inst == PRINTLN {
				fmt.Println(str)
			} else {
				fmt.Print(str)
			}
		case BR:
			vm.InstBranch()
		case NOT:
			if DataType(vm.stack.Head()) != BOOL {
				panic("not bool")
			}
			vm.stack.l[len(vm.stack.l)-2] = ^vm.stack.l[len(vm.stack.l)-2]
		case ADD, SUB, MUL, DIV, REM:
			b, a := vm.StackPopValue(), vm.StackPopValue()
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
				vm.stack.PushAll(true, res...)
				vm.stack.Push((byte)(STRING))
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
			//vm.stack.PushAll(true, res...)
			panic("error")

		}
	}
	fmt.Println("\nVM Diagnosis:")
	fmt.Println("Stack:", vm.stack.l)
	fmt.Println("Stack Length:", len(vm.stack.l))
	if len(vm.stack.l) != 1 && vm.stack.Head() != 0 {
		fmt.Println("exit code 1(stack not empty)")
		os.Exit(1)
	}
	fmt.Println("Stack Empty: true")
	os.Exit(0)
}

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
