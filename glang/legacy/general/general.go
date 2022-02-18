package general

import (
	"fmt"
	"log"
	"strings"

	"almeng.com/glang/legacy/general/Text"
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

func Alert(source *Text.Source, diagnotics Diags, line string) {
	for _, v := range diagnotics.Notions {
		line_index := source.LineIndex(v.Span.Beg)
		line_num := line_index + 1
		char := v.Span.Beg - source.Lines[line_index].Start + 1

		beg, end := v.Span.Beg, v.Span.End
		print_msg(v.Msg, v.Lev)
		pre, err, suf := line[:beg], line[beg:end], line[end:]
		fmt.Print("[", line_num, ":", char, "]")
		fmt.Println("\t" + pre + color[v.Lev] + err + "\033[0m" + suf)
	}
}

func IsEmpty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func print_msg(msg string, l Level) {
	fmt.Printf(color[l])
	fmt.Println(fmt.Sprint(l) + ": " + msg + "\033[0m")

}
