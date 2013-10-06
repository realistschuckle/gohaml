package compiler

import (
	"bytes"
	p "github.com/realistschuckle/gohaml/parser"
)

type HamlCompiler interface {
	Compile(p.ParsedDoc, CompilerOpts) (CompiledDoc, error)
}

type CompiledOutput interface {
	Render(map[string]interface{}) (string, error)
}

type CompiledDoc struct {
	Outputs []CompiledOutput
}

func (self *CompiledDoc) Render(scope map[string]interface{}) (output string, err error) {
	buf := bytes.Buffer{}
	for i := 0; i < len(self.Outputs); i += 1 {
		o, _ := self.Outputs[i].Render(scope)
		buf.WriteString(o)
	}
	output = buf.String()
	return
}

type DefaultCompiler struct {
	doc CompiledDoc
}

func (self *DefaultCompiler) Compile(input p.ParsedDoc, opts CompilerOpts) (doc CompiledDoc, e error) {
	self.doc = CompiledDoc{}
	input.Accept(self)
	doc = self.doc
	return
}

func (self *DefaultCompiler) VisitDoctype(node *p.DoctypeNode) {
	self.doc.Outputs = append(self.doc.Outputs, &StaticOutput{"<?xml version='1.0' encoding='utf-8' ?>"})
}

type StaticOutput struct {
	Content string
}

func (self *StaticOutput) Render(scope map[string]interface{}) (output string, err error) {
	output = self.Content
	return
}
