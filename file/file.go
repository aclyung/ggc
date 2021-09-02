package file

import (
	"fmt"
	"os"
	"path/filepath"

	"almeng.com/ggc/general"
)

var print func(interface{}) = func(i interface{}) {
	fmt.Println(i)
}

func WalkPath(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func FilterExt(path string, ext string) bool {
	if filepath.Ext(path) != ext {
		return false
	}
	return true
}

func FileGet() []string {
	var file_ext []string
	files, err := WalkPath("./")
	general.ErrCheck(err)
	for _, file := range files {
		if FilterExt(file, ".ggc") {
			file_ext = append(file_ext, file)
			fmt.Println(file)
		}
	}
	print(file_ext)
	return file_ext
}
