package syntax

import (
	"io"
	"unicode/utf8"
)

type source struct {
	in        string
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
	s.ch = ' '
	s.line, s.col = 0, 0
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
redo:
	if s.r >= len(s.in) {
		if s.ch == -1 {
			return
		}
		s.r++
		s.ch = -1
		return
	}
	s.col++
	if s.ch == '\n' {
		s.line++
		s.col = 0
	}
	if s.ch = rune(s.in[s.r]); s.ch < utf8.RuneSelf {
		s.r++
		if s.ch == 0 {
			s.error("invalid NUL character")
			goto redo
		}
		return
	}

}

func (s *source) start() { s.b = s.r }
func (s *source) segment() string {
	lit := s.in[s.b : s.r-1]
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
