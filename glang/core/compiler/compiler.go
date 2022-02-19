package compiler

import (
	"almeng.com/glang/core/syntax"
	"os"
)

var errh = func(err error) { println(err) }

func Compile(filename string) {
	f, _ := os.Open(filename)

	// Node
	_ = syntax.Parse(filename, f, errh)

	// TODO: Node to llvm IR

	// TODO: link write file

	return
}

const src = 1
const f = src + 1

func CompileString(src string) {
	_ = syntax.ParseString(src, errh)
	return
}

func _compile() {

}
