package gohaml

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Loader, Entry are not particularly nice and custom tailored to the http handlers
// needs. Probably should be made private.

// The whole convolute is a victim of premature optimization ...

type Loader interface {
	Load(id interface{}) (entry *Entry, err error)
}

type Entry struct {
	Engine *Engine
	ts     time.Time
}

type fileSystemLoader struct {
	baseDir      string
	cache        map[interface{}]*Entry
	checkFSAfter time.Duration
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

	return &fileSystemLoader{dir, make(map[interface{}]*Entry), 2 * time.Second}, nil
}

func (l *fileSystemLoader) adjustSuffix(path string) string {
	const htmlExt = ".html"

	//fmt.Printf("oPath: >%s<\n", path)

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if strings.HasSuffix(path, "/") {
		path += "index"
	} else {
		// check if it's a dir
		if f, err := os.Open(l.baseDir + path); err == nil {
			defer f.Close()
			if fi, err := f.Stat(); err == nil {
				if fi.IsDir() {
					path += "/index"
				} // isDir
			} // stat 
		} //open
	}

	// swap .html extension ...
	if strings.HasSuffix(path, htmlExt) {
		path = path[:len(path)-len(htmlExt)]
	}

	// ... for haml
	path += ".haml"
	path = l.baseDir + path
	//fmt.Printf("Path: >%s<\n", path)
	return path
}

func (l *fileSystemLoader) Load(id_string interface{}) (entry *Entry, err error) {
	// totally prematurely optimized, cached filessystem loader
	// check
	id, ok := id_string.(string)
	if !ok {
		err = fmt.Errorf("id: %s is not a string", id)
		return
	}


	if entry, ok = l.cache[id]; ok {
		// if less than 2 seconds have passed, don't check fs for newer version.
		if time.Since(entry.ts) < l.checkFSAfter {
			return
		}
	}

	var file *os.File
	// check fs
	var path = l.adjustSuffix(id)
	if file, err = os.Open(path); err != nil {
		return
	}

	defer file.Close()

	if ok {
		var fi os.FileInfo
		if fi, err = file.Stat(); err != nil {
			return
		}

		if fi.ModTime().Before(entry.ts) {
			// fmt.Printf("cache hit: %s %s\n", id, entry.fsTs)
			entry.ts = time.Now()
			// fmt.Printf("cache new ts: %s\n", entry.fsTs)
			return
		}
	}

	// either no cache entry or the entry is stale

	var bb bytes.Buffer
	if _, err = io.Copy(&bb, file); err != nil {
		return
	}

	var engine *Engine
	if engine, err = NewEngine(bb.String()); err != nil {
		return
	}
	entry = &Entry{engine, time.Now()}
	l.cache[id] = entry

	return
}
