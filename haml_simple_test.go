package gohaml

import (
  "testing"
)

//const simple_haml = "%p(a='b')"
//const simple_haml = "%p{a=>'b'}"
//const simple_haml = "%p{:a=>'b'}"
const simple_haml = "%p{:a=>\"b\"}"
const simple_html = "<p a='b'></p>"

func TestSimple(t * testing.T) {
  engine, _ := NewEngine(simple_haml)
  html := engine.Render(defaultScope)
  print(html)
}
