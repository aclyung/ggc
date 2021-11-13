package main

import (
	_ "almeng.com/glang/file"
	"almeng.com/glang/lexer"
	"almeng.com/glang/parser"
	_ "almeng.com/glang/parser"
	"almeng.com/glang/parser/node"
	"bufio"
	"fmt"
	"os"
)

const config = "./config.json"

func pprint(nod node.ExpressionSyntax, ident string) {
	indent := ""
	nodeType := nod.Type()
	fmt.Print(nodeType," ")
	if nodeType == node.SyntaxToken {
		val := nod.(lexer.SyntaxToken)
		if val.Value != nil {
			fmt.Print(" ", val.Value)
		} else {
			fmt.Print(" ", val.Token)
		}
	}
	fmt.Println()
	indent += "\t"
	for _, v := range nod.GetChildren() {
		pprint(v,indent)
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
	for {
		fmt.Print(">")
		input := bufio.NewReader((os.Stdin))
		line, _ := input.ReadString('\n')
		//		general.ErrCheck(err)
		pas := parser.Parser(line)
		exp := pas.Parse()
		fmt.Println(exp)
		pprint(exp, "")
		//	lex := lexer.NewLexer(line)
		//	for {
		//		tok := lex.NextToken()
		//		if tok == nil {
		//			continue
		//		}
		//		if tok.Token == token.EOF {
		//			fmt.Println("EOF")
		//			break
		//		}
		//		fmt.Print(tok.Token.String(), ": ", tok.Text)
		//		if tok.Value != nil {
		//			fmt.Print(" | ", tok.Value)
		//		}
		//		fmt.Println()
		//	}
		//}

		// a("asda")
	}

	// func a(g interface{}) { fmt.Print(g) }
}