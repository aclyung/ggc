package ir

import (
	"almeng.com/glang/core/ir/types"
	"fmt"
	"strconv"
	"strings"
)

type Value interface {
	// BCString returns Type Len DATA
	BCString() string
}

type _value struct {
	Typ types.Type
	Val string
}

type StringValue struct {
	_value
}

type IntValue struct {
	_value
}

type Float64Value struct {
}

type BoolValue struct {
}

type SliceValue struct {
	Typ *types.SliceType
	// len of the slice will be chosen by size of the ElemType
	Val []Value
}

func StringEscape(s string) string {
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\t", `\t`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	s = strings.ReplaceAll(s, "\000", `\0`)
	s = strings.ReplaceAll(s, " ", `\20`)
	return s
}

func (v _value) BCString() string {
	return fmt.Sprintf("%s %s", v.Typ.BCString(), v.Val)
}

func NewStringValue(str string) (s *StringValue) {
	s = new(StringValue)
	s.Typ = types.String
	s.Val = StringEscape(str)
	return
}

func NewIntValue(typ *types.IntType, x int64) (i *IntValue) {
	i = new(IntValue)
	i.Typ = typ
	i.Val = strconv.Itoa(int(x))
	return
}
