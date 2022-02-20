package syntax

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
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

func testPrint(msg string) {
	str := "===== " + msg
	println("\n" + str + "\n")
}

func result(pass bool, time float64) {
	res := "FAIL"
	if pass {
		res = "PASS"
	}
	time = float64(int(time*100)) / 100

	println("--- " + res + ": " + fmt.Sprint(time) + "s")
}

func TestParse(src string, errh ErrHandler, testTok bool) error {
	file++
	fname := string(file)
	if testTok {
		TokenizingTest(fname, src)
	}
	testPrint("Testing Parsing")
	start := time.Now()
	r := strings.NewReader(src)
	e := Parse(string(file), r, errh)
	took := time.Since(start).Seconds()
	testPrint("Testing Parsing complete")
	result(e == nil, took)
	return e
}

func Parse(filename string, src io.Reader, errh ErrHandler) error {
	var p parser
	p.init(src, errh)
	p.next()
	fmt.Println("File " + filename + "\n")
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
