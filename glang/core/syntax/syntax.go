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
var Purple string = "\033[35m"
var _ error = Error{}

func (e Error) Error() string {
	return fmt.Sprintf("%s %s", e.Pos, e.Msg)
}

type ErrHandler func(err error)

var file rune = 'a' - 1

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

func TestParseString(src string, errh ErrHandler, testTok bool) error {
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
	fmt.Println(Blue + "File " + filename + "\n")
	return p.EOF()
}

func TestParseFile(filename string, errh ErrHandler, t bool) error {
	f, err := os.Open(filename)
	if err != nil {
		if errh != nil {
			errh(err)
		}
		return err
	}
	b, err := io.ReadAll(f)
	if err != nil {
		if errh != nil {
			errh(err)
		}
		return err
	}
	return TestParseString(string(b), errh, t)
}
