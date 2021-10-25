package main

import (
	"fmt"

	"almeng.com/ggc/file"
)

const config = "./config.json"

func main() {

	conf := file.GetConfig(config)
	conf.PrintInfo()
	path := conf.Path
	ext := conf.Ext

	fileMan := file.NewFile()
	compfiles := fileMan.WalkPath(path).ExtractExt(ext)
	files := *compfiles.Open()
	for _, v := range files {
		fmt.Println("\nFile name: ",v.Name())
		fileMan.ReadLine(v)
	}
}


