package compiler

import (
	"bytes"
	"fmt"
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
	doc  CompiledDoc
	opts CompilerOpts
	err error
}

func (self *DefaultCompiler) Compile(input p.ParsedDoc, opts CompilerOpts) (doc CompiledDoc, e error) {
	self.doc = CompiledDoc{}
	self.opts = opts
	input.Accept(self)
	doc = self.doc
	return
}

func (self *DefaultCompiler) VisitDoctype(node *p.DoctypeNode) {
	decl := "unknown"
	switch self.opts.Format {
	case "xhtml":
		switch node.Specifier {
		case "XML":
			decl = "<?xml version='1.0' encoding='utf-8' ?>"
		case "1.1":
			decl = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.1//EN\" \"http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd\">"
		case "basic":
			decl = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML Basic 1.1//EN\" \"http://www.w3.org/TR/xhtml-basic/xhtml-basic11.dtd\">"
		case "frameset":
			decl = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Frameset//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-frameset.dtd\">"
		case "5":
			decl = "<!DOCTYPE html>"
		case "mobile":
			decl = "<!DOCTYPE html PUBLIC \"-//WAPFORUM//DTD XHTML Mobile 1.2//EN\" \"http://www.openmobilealliance.org/tech/DTD/xhtml-mobile12.dtd\">"
		case "":
			decl = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">"
		}
	case "html5":
		switch node.Specifier {
		case "XML":
			decl = ""
		case "":
			decl = "<!DOCTYPE html>"
		}
	case "html4":
		switch node.Specifier {
		case "XML":
			decl = ""
		case "frameset":
			decl = "<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 4.01 Frameset//EN\" \"http://www.w3.org/TR/html4/frameset.dtd\">"
		case "strict":
			decl = "<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 4.01//EN\" \"http://www.w3.org/TR/html4/strict.dtd\">"
		case "":
			decl = "<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 4.01 Transitional//EN\" \"http://www.w3.org/TR/html4/loose.dtd\">"
		}
	}
	self.doc.Outputs = append(self.doc.Outputs, &StaticOutput{decl})
}

func (self *DefaultCompiler) VisitTag(node *p.TagNode) {
	val := fmt.Sprintf("<%s></%s>", node.Name, node.Name)
	output := &StaticOutput{val}
	self.doc.Outputs = append(self.doc.Outputs, output)
}

type StaticOutput struct {
	Content string
}

func (self *StaticOutput) Render(scope map[string]interface{}) (output string, err error) {
	output = self.Content
	return
}
