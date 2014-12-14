package compiler

import (
	"bytes"
	"fmt"
	p "github.com/realistschuckle/gohaml/parser"
	"sort"
	"strings"
	"unicode"
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
		if i == len(self.Outputs)-1 && len(o) > 1 {
			o = strings.TrimRightFunc(o, func(r rune) bool {
				return unicode.IsSpace(r)
			})
		}
		buf.WriteString(o)
	}
	output = buf.String()
	return
}

func (self *CompiledDoc) Compress() {
	if len(self.Outputs) == 0 {
		return
	}
	outputs := []CompiledOutput{self.Outputs[0]}
	for i := 1; i < len(self.Outputs); i += 1 {
		output := self.Outputs[i]
		lastOutput := outputs[len(outputs)-1]
		if lastStatic, ok := lastOutput.(*StaticOutput); ok {
			if static, ok := output.(*StaticOutput); ok {
				lastStatic.Content += static.Content
			} else {
				outputs = append(outputs, output)
			}
		}
	}
	self.Outputs = outputs
}

type DefaultCompiler struct {
	doc  CompiledDoc
	opts CompilerOpts
	err  error
}

func (self *DefaultCompiler) Compile(input p.ParsedDoc, opts CompilerOpts) (doc CompiledDoc, e error) {
	self.doc = CompiledDoc{}
	self.opts = opts
	sort.Strings(self.opts.Autoclose)
	input.Accept(self)
	doc = self.doc
	doc.Compress()
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
	var val string
	i := sort.SearchStrings(self.opts.Autoclose, node.Name)
	autoClose := len(self.opts.Autoclose) > 0 && i < len(self.opts.Autoclose) && self.opts.Autoclose[i] == node.Name
	shouldClose := node.Close || autoClose

	classes := ""
	if len(node.Classes) > 0 {
		classes = fmt.Sprintf(" class='%v'", strings.Join(node.Classes, " "))
	}

	id := ""
	if len(node.Id) > 0 {
		id = fmt.Sprintf(" id='%v'", node.Id)
	}

	attrs := ""
	if len(node.Attrs) > 0 {
		attributes := []string{}
		for _, attr := range node.Attrs {
			content := attr.Value.(*p.StaticNode).Content
			a := fmt.Sprintf("%s='%s'", attr.Name, content)
			attributes = append(attributes, a)
		}
		attrs = " " + strings.Join(attributes, " ")
	}

	if len(node.Children) == 0 {
		switch {
		case self.opts.Format == "xhtml" && shouldClose:
			val = fmt.Sprintf("<%s%s%s%s />", node.Name, classes, id, attrs)
		case self.opts.Format == "html4" && shouldClose:
			val = fmt.Sprintf("<%s%s%s%s>", node.Name, classes, id, attrs)
		case self.opts.Format == "html5" && shouldClose:
			val = fmt.Sprintf("<%s%s%s%s>", node.Name, classes, id, attrs)
		default:
			val = fmt.Sprintf("%s<%s%s%s%s></%s>%s", node.Indent, node.Name, classes, id, attrs, node.Name, node.LineBreak)
		}
		output := &StaticOutput{val}
		self.doc.Outputs = append(self.doc.Outputs, output)
	} else {
		sn, ok := node.Children[0].(*p.StaticNode)
		if ok && len(node.Children) == 1 {
			content := []string{
				fmt.Sprintf("%s<%s%s%s%s>", node.Indent, node.Name, classes, id, attrs),
				sn.Content,
				fmt.Sprintf("</%s>%s", node.Name, node.LineBreak),
			}
			output := &StaticOutput{strings.Join(content, "")}
			self.doc.Outputs = append(self.doc.Outputs, output)
		} else {
			output := &StaticOutput{fmt.Sprintf("%s<%s%s%s%s>%s", node.Indent, node.Name, classes, id, attrs, node.LineBreak)}
			self.doc.Outputs = append(self.doc.Outputs, output)

			for _, child := range node.Children {
				child.Accept(self)
			}

			output = &StaticOutput{fmt.Sprintf("%s</%s>%s", node.Indent, node.Name, node.LineBreak)}
			self.doc.Outputs = append(self.doc.Outputs, output)
		}
	}
}

func (self *DefaultCompiler) VisitStaticLine(node *p.StaticLineNode) {
	val := fmt.Sprintf("%s%s\n", node.Indent, node.Content)
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
