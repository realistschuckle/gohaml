package gohaml

import (
	"bytes"
	"fmt"
)

type compiledDoc struct {
	buf *bytes.Buffer
}

func (self *compiledDoc) Render(scope map[string]interface{}) (s string, e error) {
	s = self.buf.String()
	return
}

type compiler struct {
	doc     ParsedDocument
	options *EngineOptions
	out     *compiledDoc
}

func (self *compiler) Compile() (c CompiledDocument, e error) {
	self.out = &compiledDoc{&bytes.Buffer{}}
	self.doc.Accept(self)
	c = self.out
	self.out = nil
	return
}

func (self *compiler) VisitDocType(n *DocTypeNode) {
	f := "unknown doctype"
	if self.options.Format == "xhtml" {
		switch n.Specification {
		case "XML":
			f = "<?xml version='1.0' encoding='utf-8' ?>"
		case "1.1":
			f = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.1//EN\" \"http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd\">"
		case "5":
			f = "<!DOCTYPE html>"
		case "basic":
			f = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML Basic 1.1//EN\" \"http://www.w3.org/TR/xhtml-basic/xhtml-basic11.dtd\">"
		case "frameset":
			f = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Frameset//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-frameset.dtd\">"
		case "mobile":
			f = "<!DOCTYPE html PUBLIC \"-//WAPFORUM//DTD XHTML Mobile 1.2//EN\" \"http://www.openmobilealliance.org/tech/DTD/xhtml-mobile12.dtd\">"
		case "":
			f = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">"
		}
	}
	if self.options.Format == "html5" {
		switch n.Specification {
		case "XML":
			f = ""
		case "":
			f = "<!DOCTYPE html>"
		}
	}

	if self.options.Format == "html4" {
		switch n.Specification {
		case "XML":
			f = ""
		case "frameset":
			f = "<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 4.01 Frameset//EN\" \"http://www.w3.org/TR/html4/frameset.dtd\">"
		case "strict":
			f = "<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 4.01//EN\" \"http://www.w3.org/TR/html4/strict.dtd\">"
		case "":
			f = "<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 4.01 Transitional//EN\" \"http://www.w3.org/TR/html4/loose.dtd\">"
		}
	}

	self.out.buf.WriteString(f)
}

func (self *compiler) VisitTag(n *TagNode) {
	s := ""
	if self.options.Format == "xhtml" {
		switch n.Name {
			case "area", "base", "br", "col", "hr", "img", "input", "link", "meta", "param":
				s = fmt.Sprintf("<%s />", n.Name)
			default:
				if n.ForceClose {
					s = fmt.Sprintf("<%s />", n.Name)
				}
		}
	}
	if self.options.Format == "html5" {
		switch n.Name {
			case "area", "base", "br", "col", "hr", "img", "input", "link", "meta", "param":
				s = fmt.Sprintf("<%s>", n.Name)
			default:
				if n.ForceClose {
					s = fmt.Sprintf("<%s>", n.Name)
				} else {
					s = fmt.Sprintf("<%s></%s>", n.Name, n.Name)
				}
		}
	}
	if self.options.Format == "html4" {
		switch n.Name {
			case "area", "base", "br", "col", "hr", "img", "input", "link", "meta", "param":
				s = fmt.Sprintf("<%s>", n.Name)
		}
	}
	self.out.buf.WriteString(s)
}
