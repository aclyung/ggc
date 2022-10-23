package types

import "strconv"

type Type interface {
	SetName(name string)
	BCString() string
}

type (
	FuncType struct {
		// Type name; or empty if not present.
		TypeName string
		// Return type.
		RetType Type
		// Function parameters.
		Params []Type
		// Variable number of function arguments.
		Variadic bool
	}
	VoidType struct{}
	IntType  struct {
		BitSize uint64
	}
	FloatType struct {
		BitSize uint64
	}
	SliceType struct {
		Elem Type
	}
	StructType struct {
		TypeName string
		Fields   []Type
	}
	PtrType struct {
		Elem Type
	}
)

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

func (s *SliceType) SetName(name string) {
	//TODO implement me
	panic("implement me")
}

func (i IntType) SetName(name string) {
	//TODO implement me
	panic("implement me")
}

var (
	I8     = &IntType{8}
	I16    = &IntType{16}
	I32    = &IntType{32}
	I64    = &IntType{64}
	I8Ptr  = &PtrType{I8}
	I16Ptr = &PtrType{I16}
	I32Ptr = &PtrType{I32}
	I64ptr = &PtrType{I64}

	String = &SliceType{Elem: I8}
	Void   = &VoidType{}
)

func (v VoidType) SetName(name string) {
	//TODO implement me
	panic("implement me")
}

func NewIntType(bitSize uint64) *IntType {
	switch bitSize {
	case 8, 16, 32, 64:
		break
	default:
		panic("invalid bit size")
	}
	return &IntType{bitSize}
}

func NewPtrType(typ Type) *PtrType {
	return &PtrType{typ}
}
