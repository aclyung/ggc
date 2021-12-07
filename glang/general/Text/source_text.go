package Text

type Source struct {
	Text  string
	Lines []Line
}

type Line struct {
	Text   Source
	Start  int
	End    int
	Length int
}

func sourceText(text string) Source {
	s := Source{}
	s.Text = text
	lines := ParseLines(s, text)
	s.Lines = lines
	return s
}

func (s Source) At(pos int) uint8 {
	return s.Text[pos]
}

func (s Source) Length() int {
	return len(s.Text)
}

func (s Source) String() string {
	return s.Text
}

func (s Source) ToString(beg int, end int) string {
	return s.Text[beg:end]
}

func (s Source) SprintSpan(span TextSpan) string {
	return s.ToString(span.Beg, span.End)
}

func From(text string) Source {
	return sourceText(text)
}

func (s Source) LineIndex(pos int) int {
	upper := 0
	lower := len(s.Text) - 1
	for lower <= upper {
		var index int = lower + (upper-lower)/2
		start := s.Lines[index].Start
		if pos == start {
			return index
		}
		if start > pos {
			upper = index - 1
		} else {
			lower = index + 1
		}
	}
	return lower - 1
}

func ParseLines(source Source, text string) []Line {
	res := make([]Line, 0)
	pos := 0
	start := 0
	for pos < len(text) {
		b := LineBreaksWith(text, pos)
		if b == 0 {
			pos++
		} else {
			res = append(res, AddLine(source, pos, start))
			pos += b
			start = pos
		}
	}
	if pos > start {
		res = append(res, AddLine(source, pos, start))
	}
	return res
}

func AddLine(source Source, pos int, start int) Line {
	breakPoint := pos + start
	return TextLine(source, start, pos, breakPoint)
}

func LineBreaksWith(text string, i int) int {
	c := rune(text[i])
	var l rune = '\000'
	isLast := i+1 >= len(text)
	if !isLast {
		l = rune(text[i+1])
	}
	if c == '\r' && l == '\n' {
		return 2
	}
	if c == '\r' || l == '\n' {
		return 1
	}

	return 0
}

func (l Line) String() string {
	return l.Text.SprintSpan(l.Span())
}

func TextLine(source Source, beg int, end int, length int) Line {
	return Line{}
}

func (l Line) Span() TextSpan {
	return Span(l.Start, l.End)
}

func (l Line) SpanLine() TextSpan {
	return Span(l.Start, l.Length)
}
