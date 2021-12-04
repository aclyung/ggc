package general

import (
	"fmt"
	"log"
)

type Level int

const (
	WARN = iota
	CAUTION
	ERROR
)

var Levels = [...]string{
	WARN:    "WARN",
	CAUTION: "CAUTION",
	ERROR:   "ERROR",
}

func (l Level) String() string {
	return Levels[l]
}

var color = [...]string{
	WARN:    "\033[30m",
	CAUTION: "\033[33m",
	ERROR:   "\033[31m",
}

func ErrCheck(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Alert(diagnotics Diags, line string) {
	for _, v := range diagnotics.Notions {
		beg, end := v.Span.Beg, v.Span.End
		print_msg(v.Msg, v.Lev)
		pre, err, suf := line[:beg], line[beg:end], line[end:]
		fmt.Println("\t" + pre + color[v.Lev] + err + "\033[0m" + suf)
	}
}

func print_msg(msg string, l Level) {
	fmt.Printf(color[l])
	fmt.Println(fmt.Sprint(l) + ": " + msg + "\033[0m")

}

func Err(msg string) string {
	return "ERROR: " + msg
}
