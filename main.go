package main

import (
	"almeng.com/glang/glang/ast"
	"almeng.com/glang/glang/general"
	"almeng.com/glang/glang/lexer"
	"almeng.com/glang/glang/parser/evaluator"
	node2 "almeng.com/glang/glang/parser/node"
	"almeng.com/glang/glang/syntax/expression"
	"bufio"
	"fmt"
	"os"
	"time"
)

const config = "./config.json"

func pprint(nod node2.ExpressionSyntax, indent string, isLast bool) {

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
	if nodeType == expression.SyntaxToken {
		val := nod.(lexer.SyntaxToken)
		if val.Value != nil {
			fmt.Print(" ", val.Token, " | ", val.Value)
		} else {
			fmt.Print(" ", val.Token)
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
	show := false
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

		if tree.Root != nil && show {
			pprint(tree.Root, "", true)
		}
		if len(tree.Diagnostics) > 0 {
			general.Alert(tree.Diagnostics)
		}
		{
			//eval := evaluator.Evaluator(exp.Root)
			res := evaluator.ExpressionEvaluation(tree.Root).(lexer.SyntaxToken)
			fmt.Println("result: " + fmt.Sprint(res.Kind()) + " | " + fmt.Sprint(res.Value))

			fmt.Println(time.Since(start))
		}
	}

}
