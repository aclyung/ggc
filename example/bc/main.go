package main

import (
	"almeng.com/glang/vm"
	"fmt"
	"github.com/inhies/go-bytesize"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var (
	pMode  vm.Mode
	file   *os.File
	defers = make([]func(), 0)
)

func parseFlag() {

	args := os.Args[1:]
	if len(args) != 2 {
		panic("invalid args")
	}
	if mode, ok := vm.ValidMode(args[0]); ok {
		pMode = mode
	} else {
		panic("invalid arg")
	}
	open, err := os.Open(args[1])
	if err == nil {
		file = open
		defers = append(defers, func() { file.Close() })
		return
	}
	panic("invalid arg")
}

func ClearTmp() {
	for _, v := range defers {
		v()
	}
}

func main() {

	defer ClearTmp()
	parseFlag()
	//b, _ := ioutil.ReadFile("out.o")
	//Execute(b)
	src, err := io.ReadAll(file)
	if err != nil {
		return
	}
	if pMode == vm.Build || pMode == vm.BC {
		s := string(src)
		s = strings.ReplaceAll(s, "\n", " ")
		s = strings.ReplaceAll(s, "\t", "")
		split := strings.Split(s, " ")
		asm := make([]string, 0)
		for _, v := range split {
			if v != "" {
				asm = append(asm, v)
			}
		}
		code := vm.NewAssembler(asm).GenBC()

		fmt.Println("Byte Code:")
		for i, v := range code {
			if i%16 == 0 {
				if i != 0 {
					fmt.Println()
				}
				fmt.Printf("%08x: ", i)
			}
			fmt.Printf("%02x ", v)

		}
		fmt.Println()
		size := bytesize.New(float64(len(code)))

		fmt.Printf("BC Size: %s\n", size)
		wd, _ := os.Getwd()
		temp := wd + "/build"
		err = os.MkdirAll(temp, os.ModePerm)
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(temp)
		mod := `module main

go 1.18

replace almeng.com/glang/vm => /Users/seungyeoplee/Workspace/glang/vm

require almeng.com/glang/vm v0.0.0-00010101000000-000000000000

require github.com/inhies/go-bytesize v0.0.0-20220417184213-4913239db9cf // indirect`
		f := `package main

import (
	"almeng.com/glang/vm"
	_ "embed"
)

//go:embed out.bc
var src []byte

func main() {
	vm.Execute(src)
}`
		if pMode == vm.BC {
			temp = wd
		}
		err = ioutil.WriteFile(temp+"/out.bc", code, 0644)
		if err != nil {
			panic(err.Error())
		}
		if pMode == vm.BC {
			return
		}
		err = ioutil.WriteFile(temp+"/go.mod", []byte(mod), 0644)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(temp+"/main.go", []byte(f), 0644)
		if err != nil {
			panic(err)
		}
		gmod := exec.Command("go", "mod", "tidy")
		gmod.Dir = temp
		gmod.Stdout = os.Stdout
		gmod.Stderr = os.Stderr
		err = gmod.Run()
		//fmt.Fprintf(gmod.Stderr, "%s\n", gmod.String())
		if err != nil {
			panic(err.Error())
		}

		cmd := exec.Command("go", "build", "-o", wd+"/out", temp+"/main.go")
		cmd.Dir = temp
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()

		if err != nil {
			panic(err.Error())
		}
	} else {
		vm.Execute(src)
	}
	//for _, v := range code {
	//	fmt.Printf("%b", v)
	//
	//}
	//ioutil.WriteFile("out.o", code, os.ModePerm)
}
