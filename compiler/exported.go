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
	var decl string
	switch {
	case node.Specifier == "XML":
		decl = "<?xml version='1.0' encoding='utf-8' ?>"
	case node.Specifier == "1.1":
		decl = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.1//EN\" \"http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd\">"
	case node.Specifier == "basic":
		decl = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML Basic 1.1//EN\" \"http://www.w3.org/TR/xhtml-basic/xhtml-basic11.dtd\">"
	case node.Specifier == "frameset":
		decl = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Frameset//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-frameset.dtd\">"
	case node.Specifier == "5":
		decl = "<!DOCTYPE html>"
	case node.Specifier == "mobile":
		decl = "<!DOCTYPE html PUBLIC \"-//WAPFORUM//DTD XHTML Mobile 1.2//EN\" \"http://www.openmobilealliance.org/tech/DTD/xhtml-mobile12.dtd\">"
	case node.Specifier == "":
		decl = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">"
	}
	self.doc.Outputs = append(self.doc.Outputs, &StaticOutput{decl})
}

type StaticOutput struct {
	Content string
}

func (self *StaticOutput) Render(scope map[string]interface{}) (output string, err error) {
	output = self.Content
	return
}
