package general

import (
	"fmt"
	"log"
)

type Level = int

type Diags struct {
	Notions []Diag
}

func NewDiag() Diags {
	return Diags{}
}

func (d *Diags) Diagnose(text string, l Level) {
	diag := Diag{text, l}
	d.Notions = append(d.Notions, diag)
}

type Diag struct {
	Msg string
	Lev Level
}

const (
	WARN = iota
	CAUTION
	ERROR
)

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

func Alert(diagnotics Diags) {
	for _, v := range diagnotics.Notions {
		print_msg(v.Msg, v.Lev)
	}
}

func print_msg(msg string, l Level) {
	fmt.Printf(color[l])
	fmt.Println(msg + "\033[0m")
}

func Err(msg string) string {
	return "ERROR: " + msg
}
