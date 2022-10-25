package types

import "strconv"

type Type interface {
	Name() string
	SetName(name string)
	BCString() string
}

type (
	_type struct {
		name string
	}
	FuncType struct {
		_type
		// Return type.
		RetType Type
		// Function parameters.
		Params []Type
		// Variable number of function arguments.
		Variadic bool
	}
	VoidType struct {
		_type
	}
	IntType struct {
		_type
		BitSize uint64
	}
	FloatType struct {
		_type
		BitSize uint64
	}
	SliceType struct {
		_type
		Elem Type
	}
	StructType struct {
		_type
		Fields []Type
	}
	PtrType struct {
		_type
		Elem Type
	}
)

func (t *_type) Name() string {
	return t.name
}
func (t *_type) SetName(name string) {
	t.name = name
}

func (v VoidType) BCString() string {
	return "void"
}

func (s *SliceType) BCString() (str string) {
	if s == String {
		return "string"
	}
	return "slice " + s.Elem.BCString()
}

func (i IntType) BCString() string {
	return "i" + strconv.Itoa(int(i.BitSize))
}

var (
	I8     = NewIntType(8)
	I16    = NewIntType(16)
	I32    = NewIntType(32)
	I64    = NewIntType(64)
	I8Ptr  = NewPtrType(I8)
	I16Ptr = NewPtrType(I16)
	I32Ptr = NewPtrType(I32)
	I64ptr = NewPtrType(I64)

	String = &SliceType{Elem: I8}
	Void   = &VoidType{}
)

func (v VoidType) SetName(name string) {
	//TODO implement me
	panic("implement me")
}

func NewIntType(bitSize uint64) *IntType {
	i := new(IntType)
	switch bitSize {
	case 8, 16, 32, 64:
		break
	default:
		panic("invalid bit size")
	}
	i.BitSize = bitSize
	return i
}

func NewPtrType(typ Type) *PtrType {
	p := new(PtrType)
	p.Elem = typ
	return p
}
