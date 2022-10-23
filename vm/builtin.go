package vm

type BuiltinFunction struct {
	Name string
	Args []string
	Func func() any
}

var a = `
println(...any) // fmt.Println
printf(format string, ...any) // fmt.Printf
print(...any) // fmt.Print
sprintf(format string, ...any) // fmt.Sprintf
sprintln(...any) // fmt.Sprintln
sprint(...any) // fmt.Sprint

len( slice | map | string ) // len
cap( slice ) // cap
append( slice[T], ...T ) // append
panic( any ) // panic
`

func FunctionCall(funcName string, args []any) any {
	// TODO
	panic(0)
}

func builtinCall(name string, s Stack[uint64]) {
	// TODO QUERY BUILTIN FUNCTION
	panic(0)
	_ = BuiltinFunction{
		Name: "println",
		Args: nil,
		Func: func() any {

			num_args := int64(s.Pop())
			args := make([]any, num_args)
			for num_args > 0 {
				typ := s.Pop()
				switch typ {
				case 0:
					// string
					size := s.Pop()
					str := make([]byte, size)
					for i := int64(0); i < int64(size); i++ {
						str[i] = byte(s.Pop())
					}
					args[num_args-1] = CastBytesToString(str)
				}
			}
			println(args)
			return nil
		},
	}
}
