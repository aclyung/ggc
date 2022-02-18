package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"almeng.com/glang/legacy/general"
)

type File []string

type Config struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Path    string `json:"path"`
	Ext     string `json:"extension"`
}

func GetConfig(path string) *Config {
	conf := &Config{}
	confFile := OpenFile(path)
	json.NewDecoder(confFile).Decode(&conf)
	return conf
}

func (c Config) PrintInfo() {
	fmt.Println(c.Name, "version", c.Version)
}

func NewFile() *File {
	return &File{}
}

// ================= Recursions ==================
func (f File) ReadLine(file *os.File) {
	buf := newBuffer(file)
	for buf.Scan() {
		fmt.Println(buf.Text())
		general.ErrCheck(buf.Err())
	}
}

func (f File) OpenFile(path string) *os.File {
	return OpenFile(path)
}

func (f *File) Open() *[]*os.File {
	files := *new([]*os.File)
	for _, v := range *f {
		files = append(files, OpenFile(v))
	}
	return &files
}

func (f File) WalkPath(path string) File {
	return walkPath(path)
}

func (f File) ExtractExt(ext string) File {
	fil := *new(File)
	for _, v := range f {
		if filterExt(v, ext) {
			fil = append(fil, v)
		}
	}
	return fil
}

// =============== Pure Functions =================
func newBuffer(file *os.File) *bufio.Scanner {
	return bufio.NewScanner(file)
}

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

func OpenFile(path string) *os.File {
	file, err := os.Open(path)
	general.ErrCheck(err)
	return file
}

func filterExt(path string, ext string) bool {
	if filepath.Ext(path) != ext {
		return false
	}
	return true
}
