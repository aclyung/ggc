package binding

type BoundBinaryOperKind = int

const (
	ILLEGAL BoundBinaryOperKind = iota
	ADD
	SUB
	MUL
	QUO

	log_beg
	LAND
	LOR
	log_end
)

func IsLogical(kind BoundBinaryOperKind) bool {
	return log_beg < kind && kind < log_end
}
