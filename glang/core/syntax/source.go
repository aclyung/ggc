package syntax

import (
	"io"
)

type source struct {
	in  string
	len int

	line, col int
	b, r, e   int
	ch        rune
	errh      func(line, col int, msg string)
}

func (s *source) init(r io.Reader, errh func(line, col int, msg string)) {
	if b, err := io.ReadAll(r); err != nil {
		panic(err)
	} else {
		s.in = string(b)
	}
	s.len = len(s.in)
	s.ch = ' '
	s.line, s.col = 0, 0
	s.b, s.r, s.e = 0, 0, 0
	s.errh = errh
}

const linebase = 1
const colbase = 1

// pos returns the (line, col) source position of s.ch.
func (s *source) pos() (line, col int) {
	return linebase + s.line, colbase + s.col
}

func (s *source) error(msg string) {
	line, col := s.pos()
	s.errh(line, col, msg)
}

// advance moves the lexer forward one rune.
func (s *source) nextch() {
	//redo:
	if s.ch == '\n' {
		s.line++
		s.col = 0
	}

	if s.r >= s.len {
		s.r = s.len + 1
		s.ch = -1
		return
	}
	s.col++

	s.ch = rune(s.in[s.r])
	s.r++
	return

}

func (s *source) start() {
	s.b = s.r
}
func (s *source) segment() string {
	lit := s.in[s.b-1 : s.r-1]
	return lit
}

func (s *source) current() rune {
	return s.ch
}

// lookahead returns the nextch rune in the text.
func (s *source) lookahead() rune {
	if s.col+1 >= len(s.in) {
		return EOFCHAR
	}
	return rune(s.in[s.r+1])
}
