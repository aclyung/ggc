package vm

type DataType uint8

const (
	UI8 DataType = iota
	UI64
	BOOL

	dynamic_beg
	STRING
	dynamic_end
)

func (t DataType) Size() int {
	return LitDataByteSize[t]
}

func (t DataType) IsDynamic() bool {
	return t > dynamic_beg && t < dynamic_end
}

var LitTypes = map[string]DataType{
	"i8":     UI8,
	"i64":    UI64,
	"bool":   BOOL,
	"string": STRING,
}

var LitDataByteSize = map[DataType]int{
	UI8:  1,
	UI64: 8,
	BOOL: 1,

	//DYNAMIC
	STRING: -1,
}
