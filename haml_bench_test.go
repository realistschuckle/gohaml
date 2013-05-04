package gohaml

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

const fn = "test/test.haml"

var template string

func loadTemplate() {
	if template != "" {
		return
	}
	var (
		file *os.File
		err  error
	)

	if file, err = os.Open(fn); err != nil {
		fmt.Printf("couldn't open testfile: %s (%s)\n", fn, err)
		os.Exit(1)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	bytes := make([]byte, 0)

	if bytes, err = ioutil.ReadAll(reader); err != nil {
		fmt.Printf("couldn't read testfile: %s (%s)\n", fn, err)
	}
	template = string(bytes)
}

func BenchmarkRender(b *testing.B) {
	var (
		engine *Engine
		err    error
	)
	b.StopTimer()
	loadTemplate()

	if engine, err = NewEngine(template); err != nil {
		fmt.Printf("couldn't init engine from testfile: %s (%s)\n", fn, err)
		os.Exit(1)
	}
	scope := make(map[string]interface{})
	b.StartTimer()
	for i := 0; i != b.N; i++ {
		engine.Render(scope)
	}
}
