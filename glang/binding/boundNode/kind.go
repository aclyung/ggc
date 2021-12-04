package boundNode

type BoundNodeKind int

const (
	Literal BoundNodeKind = iota
	Unary
	Binary
	Variable
	Assign
)
