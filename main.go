package main

import (
	"fmt"

	"almeng.com/ggc/file"
)

func main() {
	fmt.Println("ggc compiler version 0.1")
	fileMan := file.NewFile()
	files := fileMan.WalkPath("./").FileGet(".go")
	fmt.Println(files)
}
