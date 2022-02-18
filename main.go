package main

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"almeng.com/glang/legacy/ast"
	"almeng.com/glang/legacy/binding/boundNode"
	"almeng.com/glang/legacy/compile"
	"almeng.com/glang/legacy/general"
	"almeng.com/glang/legacy/general/Text"
	"almeng.com/glang/legacy/token"
	"github.com/c-bata/go-prompt"
)

const config = "./config.json"

func completer(d prompt.Document) []prompt.Suggest {
	return prompt.FilterHasPrefix(nil, d.GetWordBeforeCursor(), true)
}

func main() {
	debug.SetGCPercent(2000)
	show := true
	tsPromt := map[bool]string{true: "Showing parse trees", false: "Not showing parse trees"}
	str := ""
	lines := ""
	vars := &map[general.VariableSymbol]boundNode.BoundExpression{}
	source := &Text.Source{}
	l := 1
	for {
		line := ""
		if general.IsEmpty(str) {
			fmt.Print("\033[32mglang>\033[0m")
		} else {
			fmt.Print("\033[32m     |\033[0m")
		}
		line = prompt.Input("", completer)
		isBlank := general.IsEmpty(line)

		if len(str) == 0 {
			if isBlank {
				continue
			}
			if line == "/exit" {
				fmt.Println("Existing...")
				break
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
		}
		str = strings.TrimSuffix(str, "\n") + line
		start := time.Now()

		tree := ast.ParseTree(str, source)

		if !isBlank && len(tree.Diagnostics.Notions) != 0 {
			//lines += "\t"+"\033[31mInvalid Expression"+"\033[0m"+"\n"
			continue
		}

		compiler := compile.NewCompiler(tree)
		result := compiler.Evaluate(vars)
		diag := result.Diags
		root := tree.Root.Expression
		kind := root.Kind()
		if root != nil && show && kind != token.EOF && kind != token.ILLEGAL {
			fmt.Println("Syntax")
			fmt.Println(root)
		}
		lines += fmt.Sprint(l, "\t", str)
		l++
		if len(diag.Notions) > 0 {
			lines += "\t" + "\033[31mInvalid Expression" + "\033[0m" + "\n"
			general.Alert(tree.Source, diag, str+" ")
		} else {
			fmt.Println("result: " + fmt.Sprint(result.Type()) + " | " + fmt.Sprint(result.Value))
			lines += fmt.Sprintln("\t\033[32mEvaluation Result: " + fmt.Sprint(result.Type()) + " | " + fmt.Sprint(result.Value) + "\033[0m")
		}
		fmt.Println(time.Since(start))
		str = ""
	}
	fmt.Print(lines)
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
