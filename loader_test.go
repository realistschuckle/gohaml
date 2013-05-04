package gohaml

import (
	"bytes"
	"net/http"
	"net/url"
	"os"
	"testing"
)

const test_dir = "./test"
const simple_haml = "simple.haml"
const simple_html = "simple.html"

func TestLoadFile(t *testing.T) {
	var fsl Loader
	var err error

	if fsl, err = NewFileSystemLoader(test_dir+"/"); err != nil {
		t.Errorf("couldn't create fileSystemLoader: %s", err)
	}

	if _, err = fsl.Load(1); err == nil {
		t.Errorf("rats! expected error")
	}

	if _, err = fsl.Load(simple_haml); err != nil {
		t.Errorf("couldn't load: test.haml: %s", err)
	}

	if fsl, err = NewFileSystemLoader(test_dir); err != nil {
		t.Errorf("couldn't create fileSystemLoader: %s", err)
	}

	if _, err = fsl.Load(1); err == nil {
		t.Errorf("rats! expected error")
	}

	if _, err = fsl.Load(simple_haml); err != nil {
		t.Errorf("couldn't load: test.haml: %s", err)
	}
	if fsl, err = NewFileSystemLoader("blsadfasdf"); err == nil {
		t.Errorf("rats! expected error for non existing dir ...  ")
	}
}

func readFile(t *testing.T, fn string) ([]byte, error) {
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data := make([]byte, 1024)
	if count, err2 := file.Read(data); err2 != nil {
		return nil, err2
	} else {
		return data[:count], nil
	}
	// dead code
	return nil, nil
}

func TestHttp(t *testing.T) {
	if httpHandler, err := NewHamlHandler(test_dir); err != nil {
		t.Errorf("couldn't create HamlHandler: %s", err)
	} else {
		writer := &TestResponseWriter{bytes.NewBufferString(""), nil, 0}
		request := http.Request{}
		request.URL, _ = url.Parse("http://localhost/simple.html")
		if expected, err2 := readFile(t, test_dir + "/" +simple_html); err2 != nil {
			t.Errorf("couldn't load result: %s", err2)
		} else {
			httpHandler.ServeHTTP(writer, &request)
			if !bytes.Equal(writer.b.Bytes(), expected) {
				t.Errorf("unexpected result. <%s> >%s<", writer.b.Bytes(), expected)
			}
		}
		writer.h = nil
		request.URL, _ = url.Parse("http://localhost/index.html")
		httpHandler.ServeHTTP(writer, &request)
		if 301 != writer.s {
			t.Errorf("incorrect status: %d", writer.s)
		}
		if "./" != writer.h["Location"][0] {
			t.Errorf("incorrect redirect: %s", writer.h["Location"][0])
		}

		// shouldn't redirect
		writer.h = nil
		request.URL, _ = url.Parse("http://localhost/")
		httpHandler.ServeHTTP(writer, &request)
		if 404 != writer.s {
			t.Errorf("incorrect status: %d", writer.s)
		}
		if nil != writer.h["Location"] {
			t.Errorf("redirect but shouldn't: %s", writer.h["Location"][0])
		}
	}
}

type TestResponseWriter struct {
	b *bytes.Buffer
	h http.Header
	s int
}

func (t *TestResponseWriter) Header() http.Header {
	if t.h == nil {
		t.h = make(map[string][]string)
	}
	return t.h
}

func (t *TestResponseWriter) Write(bytes []byte) (int, error) {
	return t.b.Write(bytes)
}
func (t *TestResponseWriter) WriteHeader(i int) {
	t.s = i
}
