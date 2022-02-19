package syntax

import "fmt"

type Pos struct {
	line, col int
}

func NewPos(line, col int) Pos {
	return Pos{line, col}
}

func (p Pos) String() string {
	return fmt.Sprintf("[%d:%d]", p.line, p.col)
}

// func (pos Pos) IsKnown() bool  { return pos.line > 0 }

func (p Pos) Pos() Pos  { return p }
func (p Pos) Line() int { return p.line }
func (p Pos) Col() int  { return p.col }
