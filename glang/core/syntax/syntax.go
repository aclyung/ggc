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

var Reset string = "\033[0m"
var Cyan string = "\033[36m"
var Green string = "\033[32m"
var Yellow string = "\033[33m"
var White string = "\033[97m"
var Blue string = "\033[34m"
var Red string = "\033[31m"
var Purple string = "\033[35m"
var _ error = Error{}

func (e Error) Error() string {
	return fmt.Sprintf("%s %s", e.Pos, e.Msg)
}

type ErrHandler func(err error)

var file int = 0

func testPrint(msg string) {
	str := Cyan + "===== " + msg + Reset
	println("\n" + str + "\n")
}

func result(pass bool, time float64) {
	res := "FAIL"
	if pass {
		res = "PASS"
	}
	time = float64(int(time*100)) / 100

	println(Green + "--- " + res + ": " + fmt.Sprint(time) + "s" + Reset)
}

func TestParseString(filename string, src string, errh ErrHandler, testTok bool) *File {
	file++
	if testTok {
		TokenizingTest(filename, src)
	}
	testPrint("Testing Parsing")
	start := time.Now()
	r := strings.NewReader(src)
	e := Parse(filename, r, errh)
	took := time.Since(start).Seconds()
	testPrint("Testing Parsing complete")
	result(e != nil, took)
	return e
}

func Parse(filename string, src io.Reader, errh ErrHandler) *File {
	r := Reset
	defer func() { Reset = r }()
	var p parser
	p.init(src, errh)
	p.next()
	Reset = Blue
	fmt.Println(Blue + "File " + filename + "\n")
	return p.EOF()
}

func ParseFile(filename string, errh ErrHandler) *File {
	f, err := os.Open(filename)
	if err != nil {
		errh(err)
		return nil
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			if errh != nil {
				errh(err)
			}
		}
	}(f)
	return Parse(filename, f, errh)
}

func TestParseFile(filename string, errh ErrHandler, t bool) *File {
	f, err := os.Open(filename)
	if err != nil {
		if errh != nil {
			errh(err)
		}
		return nil
	}
	b, err := io.ReadAll(f)
	if err != nil {
		if errh != nil {
			errh(err)
		}
		return nil
	}
	return TestParseString(filename, string(b), errh, t)
}
