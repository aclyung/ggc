package Stack

import (
	"log"
	"reflect"
)

type stack struct {
	val        map[int]reflect.Value
	value_type reflect.Type
	size       int
}

func Stack(t reflect.Type) *stack {
	val := map[int]reflect.Value{0: reflect.ValueOf(reflect.Invalid)}
	return &stack{val, t, 0}
}

func (s *stack) Size() int {
	return s.size
}

func (s *stack) Push(elem interface{}) {
	if reflect.TypeOf(reflect.ValueOf(elem)) != s.value_type {
		log.Printf("Type Mismatch. Abort Pushing")
		return
	}
	s.val[s.size] = reflect.ValueOf(elem)
	s.size++
}

func (s *stack) Pop() interface{} {
	if s.size == 0 {
		log.Println("Invalid approach. The stack contains nothing")
	}
	rtn := s.val[s.Size()-1].Interface()
	s.size--
	return rtn
}
