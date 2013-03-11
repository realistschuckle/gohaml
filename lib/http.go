package gohaml;


import (
	"net/http"
	"bytes"
	"time"
)

type htmlEntry struct {
	ts time.Time
	html * bytes.Buffer
}
type httpHamlHandler struct {
	loader Loader
	cache map[string]htmlEntry
	baseDir string
}

var defaultScope map[string]interface{}


func (h * httpHamlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if entry, err := h.loader.Load(path); err != nil {
		http.NotFound(w,r)
	} else {
		var hentry htmlEntry
		var ok bool
		if hentry, ok = h.cache[path]; ok {
			if entry.ts.After(hentry.ts) {
				hentry.html = bytes.NewBufferString(entry.engine.Render(defaultScope))
				hentry.ts   = time.Now()
			}
		} else {
			hentry = htmlEntry{time.Now(), bytes.NewBufferString(entry.engine.Render(defaultScope))}
			h.cache[path] = hentry
		}
		w.Write(hentry.html.Bytes())
	}
}

func NewHamlHandler(base string)(hndl http.Handler, err error) {
	var l Loader
	if l, err = NewFileSystemLoader(base); err != nil {
		return
	}
	return &httpHamlHandler{l, make(map[string]htmlEntry), base}, nil
}
