package main

import (
	"fmt"
	"os"
	"path/filepath"

	"almeng.com/ggc/file"
)

func main() {
	fmt.Println("ggc compiler version 0.1")
	// args := os.Args[1:]
	// ar := make([]string, 0)
	// help := flag.String("help", "", "help")
	// for _, v := range args {
	// 	if v[0] != '-' {
	// 		ar = append(ar, v)
	// 	}
	// }
	// fmt.Println(os.Args, ar, *help)
	file.FileGet()
}
func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
