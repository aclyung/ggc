package syntax

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Error struct {
	Pos Pos
	Msg string
}

var _ error = Error{}

func (e Error) Error() string {
	return fmt.Sprintf("%s %s", e.Pos, e.Msg)
}

type ErrHandler func(err error)

var file rune = 'a' - 1

func ParseString(src string, errh ErrHandler) error {
	r := strings.NewReader(src)
	file++
	return Parse(string(file), r, errh)
}

func Parse(filename string, src io.Reader, errh ErrHandler) error {
	var p parser
	p.init(src, errh)
	p.next()
	fmt.Println("File " + filename)
	return p.EOF()
}

func ParseFile(filename string, errh ErrHandler) error {
	f, err := os.Open(filename)
	if err != nil {
		if errh != nil {
			errh(err)
		}
		return err
	}
	defer f.Close()
	return Parse(filename, f, errh)
}
