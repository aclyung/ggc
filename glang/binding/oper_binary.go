package binding

type BoundBinaryOperKind = int

const (
	ILLEGAL BoundBinaryOperKind = iota
	ADD
	SUB
	MUL
	QUO
)
