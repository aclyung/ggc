package assert

import (
	"log"
	"reflect"
)

func Equal(expected interface{}, actual interface{}) {
	if expected != actual {
		log.Fatalf("Expected: %v as %T, Got: %v as %T", expected, expected, actual, actual)
	}
}

func Single(in interface{}) {
	val := reflect.ValueOf(in)

	if reflect.TypeOf(in).Kind() != reflect.Slice {
		log.Fatal("The interface is not a slice")
	}
	if val.Len() != 1 {
		log.Fatalf("%v as %T has %v elements not 1", val, in, val.Len())
	}
}
