package gohaml

import (
	"net/http"
	"strings"
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
	return &httpHamlHandler{l}, nil
}

type httpHamlHandler struct {
	loader Loader
}

var defaultScope map[string]interface{}

func adjustSuffix(path string) string {
	const htmlExt = ".html"
	const htmExt = ".htm"

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// swap .html extension ...
	if strings.HasSuffix(path, htmlExt) || strings.HasSuffix(path, htmExt) {
		path = path[:strings.LastIndex(path, ".")]
		return path + ".haml"
	}

	if strings.HasSuffix(path, "/") {
		return path + "index.haml"
	}

	return path + "/index.haml"
}

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
	path = adjustSuffix(path)
	if engine, err := h.loader.Load(path); err != nil {
		http.NotFound(w, r)
	} else {
		w.Write(([]byte)(engine.Render(defaultScope)))
	}
}
