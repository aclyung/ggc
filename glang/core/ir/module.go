package ir

import (
	"almeng.com/glang/core/ir/types"
	"fmt"
	"github.com/almenglee/general"
	"io"
	"strings"
)

type Module struct {
	TypeDecl general.List[types.Type]
	Funcs    []*Func
}

func (m *Module) NewTypeDef(name string, t types.Type) types.Type {
	t.SetName(name)
	m.TypeDecl.Append(t)
	return t
}

func NewModule() *Module {
	return &Module{Funcs: make([]*Func, 0), TypeDecl: general.List[types.Type]{}}
}

func (m *Module) NewFunc(name string, ret any, param ...*Param) *Func {
	f := NewFunc(name, ret, param...)
	m.Funcs = append(m.Funcs, f)
	return f
}

func (m *Module) String() string {
	buf := &strings.Builder{}
	if _, err := m.WriteTo(buf); err != nil {
		panic(fmt.Errorf("unable to write to string buffer; %v", err))
	}
	return buf.String()
}

func (m *Module) WriteTo(buf *strings.Builder) (n int64, err error) {
	fw := &fmtWriter{w: buf}

	if len(m.Funcs) > 0 && fw.size > 0 {
		fw.Fprint("\n")
	}
	for i, f := range m.Funcs {
		if i != 0 {
			fw.Fprint("\n")
		}
		fw.Fprintln(f.BCString())
	}
	return fw.size, fw.err
}

// Fprint formats using the default formats for its operands and writes to w.
// Spaces are added between operands when neither is a string. It returns the
// number of bytes written and any write error encountered.
func (fw *fmtWriter) Fprint(a ...interface{}) (n int, err error) {
	if fw.err != nil {
		// early return if a previous error has been encountered.
		return 0, nil
	}
	n, err = fmt.Fprint(fw.w, a...)
	fw.size += int64(n)
	fw.err = err
	return n, err
}

// Fprintf formats according to a format specifier and writes to w. It returns
// the number of bytes written and any write error encountered.
func (fw *fmtWriter) Fprintf(format string, a ...interface{}) (n int, err error) {
	if fw.err != nil {
		// early return if a previous error has been encountered.
		return 0, nil
	}
	n, err = fmt.Fprintf(fw.w, format, a...)
	fw.size += int64(n)
	fw.err = err
	return n, err
}

type fmtWriter struct {
	w    io.Writer
	err  error
	size int64
}

// Fprintln formats using the default formats for its operands and writes to w.
// Spaces are always added between operands and a newline is appended. It
// returns the number of bytes written and any write error encountered.
func (fw *fmtWriter) Fprintln(a ...interface{}) (n int, err error) {
	if fw.err != nil {
		// early return if a previous error has been encountered.
		return 0, nil
	}
	n, err = fmt.Fprintln(fw.w, a...)
	fw.size += int64(n)
	fw.err = err
	return n, err
}
