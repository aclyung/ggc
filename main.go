package main

import (
	"bufio"
	"fmt"
	"os"

	_ "almeng.com/glang/file"
	"almeng.com/glang/general"
	"almeng.com/glang/lexer"
	_ "almeng.com/glang/parser"
	"almeng.com/glang/token"
)

const config = "./config.json"

func main() {
	for {
		fmt.Print(">")
		input := bufio.NewReader((os.Stdin))
		line, err := input.ReadString('\n')
		general.ErrCheck(err)
		lex := lexer.NewLexer(line)
		for {
			tok := lex.NextToken()
			if tok == nil {
				continue
			}
			if tok.Token == token.EOF {
				fmt.Println("EOF")
				break
			}
			fmt.Print(tok.Token, ": ", tok.Text)
			if tok.Value != nil {
				fmt.Print(" | ", tok.Value)
			}
			fmt.Println()
		}
	}
	// 1. load file 2.
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
	// a("asda")
}

// func a(g interface{}) { fmt.Print(g) }
