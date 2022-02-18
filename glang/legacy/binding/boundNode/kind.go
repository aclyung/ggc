package boundNode

// Node Kind

type Kind int

const (
	EOF Kind = iota
	Literal
	Unary
	Binary
	Variable
	Assign
)
