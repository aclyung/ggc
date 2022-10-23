package main

import (
	"almeng.com/glang/vm"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var (
	pMode   vm.Mode
	verbose *bool
	file    *os.File
	defers  = make([]func(), 0)
)

func parseFlag() {
	verbose = new(bool)
	*verbose = false
	args := os.Args[1:]
	for i, v := range args {
		if v[0] == '-' {
			if len(args) < i+1 {
				args = args[:i]
			} else {

				args = append(args[:i], args[i+1:]...)
			}
			switch v {
			case "-v", "--verbose":
				*verbose = true
			}
		}
	}
	if len(args) != 2 {
		fmt.Println("invalid args")
		os.Exit(1)
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
		if *verbose {

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
			size := Bitsize(len(code))

			fmt.Printf("BC Size: %s\n", size)
		}
		wd, _ := os.Getwd()
		temp := wd + "/build"
		err = os.MkdirAll(temp, os.ModePerm)
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(temp)
		bcsrc := "[]byte{\n"
		for i, v := range code {
			if i%16 == 0 {
				if i != 0 {
					bcsrc += "\n"
				}
				bcsrc += "\t"
			}
			bcsrc += fmt.Sprintf("%d,", v)
		}
		bcsrc += "\n\t}"
		mod := `module main
	
	go 1.18
	
	replace almeng.com/glang/vm => /Users/seungyeoplee/Workspace/glang/vm
	
	require almeng.com/glang/vm v0.0.0`
		f := `package main
	
import (
	"almeng.com/glang/vm"
)

func main() {
	vm.Execute(` + bcsrc + `)
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
		//gmod := exec.Command("go", "mod", "tidy")
		//gmod.Dir = temp
		//gmod.Stdout = os.Stdout
		//gmod.Stderr = os.Stderr
		//err = gmod.Run()
		////fmt.Fprintf(gmod.Stderr, "%s\n", gmod.String())
		//if err != nil {
		//	panic(err.Error())
		//}

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

const (
	Byte = iota
	KB
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

var bsize = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

func Bitsize(size int) string {
	b := 0
	for size > 1024 && b < YB {
		size = size >> 10
		b++
	}
	return fmt.Sprintf("%d%s", size, bsize[b])
}
