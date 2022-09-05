package main

import (
	"almeng.com/glang/core/compiler"
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	debug    *bool
	verbose  *bool
	no_build *bool
	asm      *bool
	run      *bool
	target   *string
	path     = "./main.gg"
)

func main() {
	parseFlag()
	c := compiler.Compile(path, *verbose, *target)
	compiled := c.GetIR()
	if *debug {
		print("Build result:\n")
		println(compiled)
	}

	wdir, err := os.Getwd()
	lldir := wdir

	if !*asm {
		tdir, terr := ioutil.TempDir("", "glang")
		if terr != nil {
			panic(terr)
		}
		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {

			}
		}(tdir)
		lldir = tdir
	}

	err = ioutil.WriteFile(lldir+"/main.ll", []byte(compiled), 0644)
	if err != nil {
		panic(err)
	}

	if err != nil {
		return
	}
	outputName := "exec"
	if !*no_build {
		if *run {
			wdir, err = ioutil.TempDir("", "glang")
			if err != nil {
				panic(err)
			}
			defer os.RemoveAll(wdir)
		}
		outputName = wdir + "/" + outputName
		clangArgs := []string{
			c.Target.String(),
			"-Wno-override-module",
			lldir + "/main.ll",
			"-o", outputName, "-O3",
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
		if *run {
			execute := exec.Command(outputName)
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
	}

	if *verbose {
		fmt.Println("Operator declarations:")
		c.Opers.Each(func(o *compiler.Operator) { println("\t", o.Name()) })
	}
}

func parseFlag() {
	debug = flag.Bool("d", false, "enable debug diagnosis")
	verbose = flag.Bool("v", false, "enable verbose debug diagnosis")
	no_build = flag.Bool("no-build", false, "make no binary build output. debug mode required")
	asm = flag.Bool("S", false, "build to llvm ir file")
	run = flag.Bool("run", false, "runs glang program")
	target = flag.String("t", "", "set target triple")
	flag.Parse()

	*debug = *debug || *verbose

	// Check if debug mode is enabled for 'no-build' option
	if !*debug && *no_build {
		fmt.Println("debug is required for '-no-build' option")
	}

	// Enable 'no-build' option when llvm assembly build option is enabled
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
