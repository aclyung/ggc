package main

import (
	"almeng.com/glang/ast"
	"almeng.com/glang/binding/boundNode"
	"almeng.com/glang/compile"
	"almeng.com/glang/general"
	"almeng.com/glang/token"
	"fmt"
	"github.com/c-bata/go-prompt"
	"os"
	"runtime/debug"
	"time"
)

const config = "./config.json"

func completer(d prompt.Document) []prompt.Suggest {
	return prompt.FilterHasPrefix(nil, d.GetWordBeforeCursor(), true)
}

func main() {
	debug.SetGCPercent(2000)
	show := true
	tsPromt := map[bool]string{true: "Showing parse trees", false: "Not showing parse trees"}

	vars := &map[general.VariableSymbol]boundNode.BoundExpression{}
	for {

		line := prompt.Input(">", completer)
		if line == "/exit" {
			fmt.Println("Existing...")
			os.Exit(0)
		}
		if line == "/show" {
			show = !show
			fmt.Println(tsPromt[show])
			continue
		}
		if line == "/cls" {
			fmt.Print("\033[H\033[2J")
			continue
		}
		start := time.Now()
		tree := ast.ParseTree(line)
		compiler := compile.NewCompiler(tree)
		result := compiler.Evaluate(vars)
		diag := result.Diags
		root := tree.Root
		kind := root.Kind()
		if root != nil && show && kind != token.EOF && kind != token.ILLEGAL {
			fmt.Println("Syntax")
			fmt.Println(tree.Root)
		}
		if len(diag.Notions) > 0 {
			general.Alert(diag, line)
		} else {
			fmt.Println("result: " + fmt.Sprint(result.Type()) + " | " + fmt.Sprint(result.Value))
		}
		fmt.Println(time.Since(start))
	}

}

//1. load file 2.
// conf := file.GetConfig(config)
// conf.PrintInfo()
// path := conf.Path
// ext := conf.Ext

// fileMan := file.NewFile()
// compfiles := fileMan.WalkPath(path).ExtractExt(ext)
// files := *compfiles.Open()
// for _, v := range files {
// 	fmt.Println("\nFile name: ", v.Name())
// 	// fileMan.ReadLine(v)
// 	fmt.Println(parser.Parser(*v))
