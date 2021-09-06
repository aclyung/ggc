package file

import (
	"os"
	"path/filepath"

	"almeng.com/ggc/general"
)

type File []string

func NewFile() File {
	return File{}
}

// ================= Recursions ==================

func (f File) WalkPath(path string) File {
	return walkPath(path)
}

func (f File) FileGet(ext string) File {
	fil := *new(File)
	for _, v := range f {
		if filterExt(v, ext) {
			fil = append(fil, v)
		}
	}
	return fil
}

// =============== Pure Functions =================

func walkPath(path string) []string {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	general.ErrCheck(err)
	return files
}

func filterExt(path string, ext string) bool {
	if filepath.Ext(path) != ext {
		return false
	}
	return true
}

//
// func FileGet() []string {
// 	var file_ext []string
// 	files := walkPath("./")
// 	for _, file := range files {
// 		if filterExt(file, ".ggc") {
// 			file_ext = append(file_ext, file)
// 			fmt.Println(file)
// 		}
// 	}
// 	print(file_ext)
// 	return file_ext
// }
