package main

import (
	"almeng.com/glang/ast"
	"almeng.com/glang/binding/boundNode"
	"almeng.com/glang/compile"
	"almeng.com/glang/expression"
	"almeng.com/glang/general"
	"almeng.com/glang/syntax"
	"almeng.com/glang/token"
	"encoding/csv"
	"github.com/c-bata/go-prompt"
	"runtime/debug"
	"strings"

	"fmt"
	"time"
)

const config = "./config.json"

func pprint(nod syntax.ExpressionSyntax, indent string, isLast bool) {
	var mark string

	if isLast {
		mark = "└────"
	} else {
		mark = "├────"
	}
	fmt.Print(indent)
	fmt.Print(mark)
	nodeType := nod.Type()
	fmt.Print(nodeType, " ")
	if nodeType == syntax.Token {
		val := nod.(expression.SyntaxToken)
		if val.Value != nil {
			fmt.Print(" ", val.Kind(), " | ", val.Value)
		} else {
			fmt.Print(" ", val.Kind())
		}
	}
	fmt.Println()
	if isLast {
		indent += "     "
	} else {
		indent += "│    "
	}

	for i, v := range nod.GetChildren() {
		pprint(v, indent, i == len(nod.GetChildren())-1)
	}
}

func completer(d prompt.Document) []prompt.Suggest {
	return prompt.FilterHasPrefix(nil, d.GetWordBeforeCursor(), true)
}

func main() {
	debug.SetGCPercent(2000)
	// Default prefix is "# ", but you can change it like so:
	// cli.SetPrefix("my-cli# ")

	// Executes a command directly, if one is given in arguments.
	// Otherwise creates a CLI.
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

	for _, each := range []string{
		"100 - 20 *a", "test1, test2, test3", "test1,test2,test3",
	} {
		r := csv.NewReader(strings.NewReader(each))
		r.TrimLeadingSpace = true
		s, e := r.Read()
		if e != nil {
			panic(e)
		}
		fmt.Printf("%q\n", s)
	}
	// }
	show := true
	tsPromt := map[bool]string{true: "Showing parse trees", false: "Not showing parse trees"}

	vars := &map[general.VariableSymbol]boundNode.BoundExpression{}
	for {

		line := prompt.Input(">", completer)
		//fmt.Print(">")
		//input := bufio.NewReader((os.Stdin))
		//line, _ := input.ReadString('\n')
		//		general.ErrCheck(err)
		//line = strings.Replace(line,string(rune(0)),"",-1)
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
			pprint(tree.Root, "", true)
		}
		if len(diag.Notions) > 0 {
			general.Alert(diag, line)
		} else {
			fmt.Println("result: " + fmt.Sprint(result.Type()) + " | " + fmt.Sprint(result.Value))
		}
		fmt.Println(time.Since(start))
	}

}
