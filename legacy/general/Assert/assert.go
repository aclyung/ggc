package assert

import (
	"fmt"
	"reflect"
)

func Empty(in interface{}) {
	val := reflect.ValueOf(in)

	if reflect.TypeOf(in).Kind() != reflect.Slice {
		panic("The interface is not a slice")
	}
	if val.Len() != 0 {
		msg := fmt.Sprintf("%v as %T has %v elements, not empty", val, in, val.Len())
		panic(msg)
	}
}

func Equal(expected interface{}, actual interface{}) {
	if expected != actual {
		msg := fmt.Sprintf("Expected: %v as %T, Got: %v as %T", expected, expected, actual, actual)
		panic(msg)
	}
}

func NotEqual(expected interface{}, actual interface{}) {
	if expected == actual {
		msg := fmt.Sprintf("Expected to be not equal: %v as %T, Got: %v as %T", expected, expected, actual, actual)
		panic(msg)
	}
}

func False(condition bool) {
	if condition {
		panic("Expected False, but got True")
	}
}

func True(condition bool) {
	if !condition {
		panic("Expected True, but got False")
	}
}

func Single(in interface{}) {
	val := reflect.ValueOf(in)

	if reflect.TypeOf(in).Kind() != reflect.Slice {
		panic("The interface is not a slice")
	}
	if val.Len() != 1 {
		msg := fmt.Sprintf("%v as %T has %v elements, not 1", val, in, val.Len())
		panic(msg)
	}
}
