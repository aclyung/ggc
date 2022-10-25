package vm

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"unsafe"
)

var (
	TRUE  = Const([]byte{byte(BOOL), 1})
	FALSE = Const([]byte{byte(BOOL), 0})
)

func (c Const) Equal(other Const) bool {
	return reflect.DeepEqual([]byte(c), []byte(other))
}

type Const []byte

// abi.go

func CastBytesToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

func CastVMSliceToGoSlice(slice []byte, typ string) []any {
	switch typ {
	case "int64":
		rbuf := bytes.NewBuffer(slice)
		r64 := make([]int64, (len(slice)+7)/8)
		err := binary.Read(rbuf, binary.LittleEndian, &r64)
		if err != nil {
			panic("cast VM slice to Go slice failed: " + err.Error())
		}
		return *(*[]any)(unsafe.Pointer(&r64))
	//case "float":
	//	return CastVMSliceToGoSliceFloat(bytes)
	//case "string":
	//
	//	return CastBytesToSt
	//case "bool":
	//	return CastVMSliceToGoSliceBool(bytes)
	default:
		panic("invalid type")
	}
}
