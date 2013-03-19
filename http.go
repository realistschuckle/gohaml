package gohaml

import (
	"bytes"
	"net/http"
	"strings"
	"time"
)

// create an http.Handler that loads haml files from locations relative
// to base dir, taking into account that they won't end in *haml in the
// http request. Translates URLs such that:
//
// /bla.html              -> ${base}/bla.haml
// /bla/bla/dingdong.html -> ${base}/bla/bla/dingdong.haml
// /bla/bla/              -> ${base}/bla/bla/index.haml
func NewHamlHandler(base string) (hndl http.Handler, err error) {
	var l Loader
	if l, err = NewFileSystemLoader(base); err != nil {
		return
	}
	return &httpHamlHandler{l, make(map[string]htmlEntry), base}, nil
}

type htmlEntry struct {
	ts   time.Time
	html *bytes.Buffer
}
type httpHamlHandler struct {
	loader  Loader
	cache   map[string]htmlEntry
	baseDir string
}

var defaultScope map[string]interface{}

func (h *httpHamlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const indexPage = "/index.html"
	path := r.URL.Path

	// borrowed from net/http/fs.go
	// redirect .../index.html to .../
	// can't use Redirect() because that would make the path absolute,
	// which would be a problem running under StripPrefix
	if strings.HasSuffix(path, indexPage) {
		newPath := "./"
		if q := r.URL.RawQuery; q != "" {
			newPath += "?" + q
		}
		w.Header().Set("Location", newPath)
		w.WriteHeader(http.StatusMovedPermanently)
		return
	}

	if entry, err := h.loader.Load(path); err != nil {
		http.NotFound(w, r)
	} else {
		var hentry htmlEntry
		var ok bool
		if hentry, ok = h.cache[path]; ok {
			if entry.ts.After(hentry.ts) {
				hentry.html = bytes.NewBufferString(entry.Engine.Render(defaultScope))
				hentry.ts = time.Now()
			}
		} else {
			hentry = htmlEntry{time.Now(), bytes.NewBufferString(entry.Engine.Render(defaultScope))}
			h.cache[path] = hentry
		}
		w.Write(hentry.html.Bytes())
	}
}
