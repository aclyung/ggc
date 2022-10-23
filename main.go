package main

import (
	"almeng.com/glang/core/compiler/backend"
	"almeng.com/glang/core/compiler/backend/types"
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"almeng.com/glang/core/compiler"
)

var (
	debug, verbose, no_build, asm, run *bool
	target                             *string
	path                               = "./main.gg"
	defers                             = make([]func(), 0)
)

func main() {
	//defer ClearTmp()
	//parseFlag()
	//// TODO: support llvm version
	//c := compiler.Compile(path, *verbose, *target)
	//Build(c)

	//c := compiler.CompileBC(path, *verbose)
	//bc := compiler.BuildBC(c)
	//vm.Execute(bc)
	insts := []backend.Value{
		backend.NewStringValue("Hello, World!"),
		backend.NewIntValue(types.I64, 8000),
	}
	c := backend.NewCall(backend.NewFunc("print", types.Void, types.String), insts...)
	println(c.BCString())
}

func ClearTmp() {
	for _, f := range defers {
		f()
	}
}

func BuildIR(lldir string, compiled string) {
	err := ioutil.WriteFile(lldir+"/main.ll", []byte(compiled), 0644)
	if err != nil {
		panic(err)
	}
}

func remove(path string) func() {
	return func() {
		err := os.RemoveAll(path)
		if err != nil {
			panic(err)
		}
	}
}

func GetDir() (wdir, lldir string) {
	var err error
	wdir, err = os.Getwd()
	lldir = wdir
	if err != nil {
		panic(err)
	}
	if !*asm {
		tdir, terr := ioutil.TempDir("", "glang")
		if terr != nil {
			panic(terr)
		}
		defers = append(defers, remove(tdir))
		lldir = tdir
	}
	return
}
func Build(c *compiler.Compiler) {
	ir := c.GetIR()
	wdir, lldir := GetDir()
	if *debug {
		print("Build result:\n")
		println(ir)
	}
	BuildIR(lldir, ir)
	BuildExec(c.Target.String(), wdir, lldir)
	if *verbose {
		fmt.Println("Operator declarations:")
		c.Opers.Each(func(o *compiler.Operator) { println("\t", o.Name()) })
	}
}

func BuildExec(target, wdir, lldir string) {
	output_name := "exec"
	if *no_build {
		return
	}
	if *run {
		var err error
		wdir, err = ioutil.TempDir("", "glang")
		if err != nil {
			panic(err)
		}
		defers = append(defers, remove(wdir))
	}
	output_name = wdir + "/" + output_name
	clangArgs := []string{
		target,
		"-Wno-override-module",
		lldir + "/main.ll",
		"-o", output_name, "-O3",
	}
	cmd := exec.Command("clang", clangArgs...)
	output, _err := cmd.CombinedOutput()
	if _err != nil {
		println(string(output))
		panic(_err)
	}
	if len(output) > 0 {
		fmt.Println(string(output))
	}
	Execute(output_name)
}

func Execute(path string) {
	if !*run {
		return
	}
	execute := exec.Command(path)
	stdout, exec_err := execute.StdoutPipe()
	if exec_err != nil {
		log.Fatal(exec_err)
	}
	exec_err = execute.Start()
	if exec_err != nil {
		log.Fatalf("cmd.Start() failed with %s\n", exec_err)
	}

	stdin := bufio.NewScanner(stdout)
	for stdin.Scan() {
		fmt.Println(stdin.Text())
	}
	err := execute.Wait()
	if err != nil {
		panic(err)
	}
}

func parseFlag() {
	debug = flag.Bool("d", false, "enable debug diagnosis")
	verbose = flag.Bool("v", false, "enable verbose debug diagnosis")
	no_build = flag.Bool("no-Build", false, "make no binary Build output. debug mode required")
	asm = flag.Bool("S", false, "Build to llvm ir file")
	run = flag.Bool("run", false, "runs glang program")
	target = flag.String("t", "", "set target triple")
	flag.Parse()

	*debug = *debug || *verbose

	// Check if debug mode is enabled for 'no-Build' option
	if !*debug && *no_build {
		fmt.Println("debug is required for '-no-Build' option")
	}

	// Enable 'no-Build' option when llvm assembly Build option is enabled
	*no_build = *no_build || *asm

	path = "./main.gg"
	n := flag.NArg()
	if n == 1 {
		path = flag.Arg(0)
	}

	if !strings.HasSuffix(path, ".gg") {
		fmt.Println("file '" + path + "' isn't appropriate glang file")
		os.Exit(-1)
	}

	if !strings.HasPrefix(path, "/") && !strings.HasPrefix(path, "./") {
		path = "./" + path
	}
}
