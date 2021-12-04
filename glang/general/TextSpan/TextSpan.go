package TextSpan

type TextSpan struct {
	Beg int
	End int
}

func Span(beg int, end int) TextSpan {
	return TextSpan{beg, end}
}
