package main

import (
	"almeng.com/glang/glang/ast"
	"almeng.com/glang/glang/binding"
	"almeng.com/glang/glang/evaluator"
	"almeng.com/glang/glang/expression"
	"almeng.com/glang/glang/general"
	"almeng.com/glang/glang/syntax"
	"bufio"
	"fmt"
	"os"
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

func main() {
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

	// }
	show := true
	tsPromt := map[bool]string{true: "Showing parse trees", false: "Not showing parse trees"}
	for {
		fmt.Print(">")
		input := bufio.NewReader((os.Stdin))
		line, _ := input.ReadString('\n')
		//		general.ErrCheck(err)
		//line = strings.Replace(line,string(rune(0)),"",-1)
		if line == "/show\n" {
			show = !show
			fmt.Println(tsPromt[show])
			continue
		}
		if line == "cls\n" {
			fmt.Print("\033[H\033[2J")
			continue
		}
		start := time.Now()
		tree := ast.ParseTree(line)
		fmt.Println("Expression")
		binder := binding.NewBinder()
		boundExp := binder.Bind(tree.Root)
		notions := make([]general.Diag, 0)
		notions = append(notions, tree.Diagnostics.Notions...)
		notions = append(notions, binder.Diag.Notions...)
		diag := general.NewDiag()
		diag.Notions = notions
		if tree.Root != nil && show {
			pprint(tree.Root, "", true)
		}
		if len(diag.Notions) > 0 {
			general.Alert(diag)
		} else {
			eval := evaluator.NewEvaluator(boundExp)
			res := eval.Evaluate().(*binding.BoundLiteralExpression)
			fmt.Println("result: " + fmt.Sprint(res.Type()) + " | " + fmt.Sprint(res.Value))

		}
		fmt.Println(time.Since(start))
	}

}
