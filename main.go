package main

import (
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
		fileMan.ReadLine(v)
	}
}
