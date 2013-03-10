package gohaml;

import (
	"os"
	"bytes"
	"io"
	"fmt"
)

type Loader interface {
	Load(id interface{}) (hamlEngine *Engine, err error)
}

type fileSystemLoader struct {
	baseDir string
}

func NewFileSystemLoader(dir string) (loader Loader, err error) {
	var f *os.File
	if f, err = os.Open(dir); err != nil {
		return
	}

	defer f.Close()

	var fi os.FileInfo
	if fi, err = f.Stat(); err != nil {
		return
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("%s: not a directory", fi.Name())
	}

	return &fileSystemLoader{dir}, nil
}

func (l  * fileSystemLoader) Load(id_string interface{}) (eng *Engine, err error){

	id, ok := id_string.(string)
	if !ok {
		err = fmt.Errorf("id: %s is not a string", id)
		return
	}

	var file *os.File
	if file, err = os.Open(l.baseDir + "/" + id) ; err != nil {
		return
	}
	defer file.Close()

	var bb bytes.Buffer
	if _, err = io.Copy(&bb, file); err!= nil {
		return
	}

	return NewEngine(bb.String())
}

