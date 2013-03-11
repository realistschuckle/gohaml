package gohaml;

import (
	"os"
	"bytes"
	"io"
	"fmt"
	"time"
	"strings"
)

type Loader interface {
	Load(id interface{}) (entry Entry, err error)
}


type Entry struct {
	ts  time.Time
	engine * Engine
}

type fileSystemLoader struct {
	baseDir string;
	cache   map[interface{}]Entry
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

	return &fileSystemLoader{dir, make(map[interface{}]Entry)}, nil
}

func (l * fileSystemLoader) adjustSuffix(path string) string {
	fmt.Printf("oPath: >%s<\n", path)
	if !strings.HasPrefix(path, "/") {
		path = "/"+path
	}
	if strings.HasSuffix(path, "/") {
		path += "index"
	} else {
		// check id dir
		if f, err := os.Open(l.baseDir + path); err == nil {
			defer f.Close();
			if fi, err := f.Stat(); err == nil {
				if fi.IsDir() {
					path += "/index"
				} // isDir
			} // stat 
		} //open
	}

	if strings.HasSuffix(path, ".html") {
		path = path[:len(path)-len(".html")]
	}

	path += ".haml"
	path = l.baseDir + path
	fmt.Printf("Path: >%s<\n", path)
	return path
}

func (l  * fileSystemLoader) Load(id_string interface{}) (entry Entry, err error){

	id, ok := id_string.(string)
	if !ok {
		err = fmt.Errorf("id: %s is not a string", id)
		return
	}

	var path = l.adjustSuffix(id)
	var file *os.File
	if file, err = os.Open(path) ; err != nil {
		return
	}
	defer file.Close()

	if entry, ok = l.cache[path]; ok {
		var fi os.FileInfo
		if fi, err = file.Stat(); err != nil {
			return
		}
		if fi.ModTime().Before(entry.ts) {
			fmt.Printf("cache hit\n")
			return
		} 
	}

	// either no cache entry or the entry is stale

	var bb bytes.Buffer
	if _, err = io.Copy(&bb, file); err!= nil {
		return
	}

	var engine *Engine
	if engine, err = NewEngine(bb.String()); err != nil {
		return
	}
	entry = Entry{time.Now(), engine}
	l.cache[path] = entry

	return
}

