package binding

type BoundUnaryOperKind = int

const (
	Identity BoundUnaryOperKind = iota
	Negation
	NOT
)
