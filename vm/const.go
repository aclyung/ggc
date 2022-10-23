package vm

import "reflect"

var (
	TRUE  = Const([]byte{byte(BOOL), 1})
	FALSE = Const([]byte{byte(BOOL), 0})
)

func (c Const) Equal(other Const) bool {
	return reflect.DeepEqual([]byte(c), []byte(other))
}

type Const []byte
