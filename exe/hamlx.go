package main

import (
	"fmt"
	"github.com/realistschuckle/gohaml/lib"
	"os"
	"io"
	"io/ioutil"
)

func main() {
	in, out := os.Stdin, os.Stdout
	var err os.Error
	defer func() {
		if in != os.Stdin {in.Close()}
		if out != os.Stdout {out.Close()}
	}()

	if len(os.Args) > 1 {
		in, err = os.Open(os.Args[1], os.O_RDONLY, 0)
		if nil != err {
			fmt.Println("Cannot", err)
			os.Exit(-1)
		}
	}
	if len(os.Args) > 2 {
		out, err = os.Open(os.Args[2], os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0666)
		if nil != err {
			fmt.Println("Cannot create output file", os.Args[2], ":", err)
			os.Exit(-1)
		}
	}
	
	bytes, _ := ioutil.ReadAll(in)
	
	engine, err := gohaml.NewEngine(string(bytes))
	if err != nil {
		os.Stderr.WriteString(err.String())
		os.Exit(1)
	}
	output := engine.Render(make(map[string]interface{}))
	
	io.WriteString(out, output)
	
	return
}
