package syntax

import "fmt"

type Pos struct {
	line, col uint
}

func NewPos(line, col uint) Pos {
	return Pos{line, col}
}

func (p Pos) String() string {
	return fmt.Sprintf("[%d:%d]", p.line, p.col)
}

// func (pos Pos) IsKnown() bool  { return pos.line > 0 }

func (p Pos) Pos() Pos      { return p }
func (p Pos) Line() uint    { return p.line }
func (p Pos) Col() uint     { return p.col }
func (p Pos) IsKnown() bool { return p.line > 0 }
