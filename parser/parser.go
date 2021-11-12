package parser

import (
	"bufio"
	"os"

	"almeng.com/glang/general"
)



func Parser(file os.File) string {
	str := ""
	buf := newBuffer(&file)
	for buf.Scan() {
		text := buf.Text()
		str += text + ";\n"
		general.ErrCheck(buf.Err())
	}

	return str
}

func newBuffer(file *os.File) *bufio.Scanner {
	return bufio.NewScanner(file)
}
