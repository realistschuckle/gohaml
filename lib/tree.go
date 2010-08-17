package gohaml

import (
	"bytes"
	"container/vector"
	"fmt"
	"strings"
	"unicode"
)

type attrMap map[string]string

func (m attrMap) String() (s string) {
	buf := bytes.NewBuffer([]byte{})
	for k, v := range m {
		buf.WriteString(" ")
		buf.WriteString(k)
		buf.WriteString("=\"")
		buf.WriteString(v)
		buf.WriteString("\"")
	}
	s = buf.String()
	return
}

func (self attrMap) appendAttr(name string, value string) {
	if _, ok := self[name]; ok {
		self[name] += " " + value
	} else {
		self[name] = value
	}
	return
}

type node struct {
	name string
	attrs attrMap
	remainder string
	indentCount int
	parent *node
	children *vector.Vector
	closeTag string
	noNewline bool
	remainderLookup bool
}

type tree struct {
	node
	indent string
}

func newTree() (t *tree) {
	t = &tree{node{"tree", nil, "", -1, nil, new(vector.Vector), " /", false, false}, "\t"}
	return
}

func (self *node) createChild(name string, remainder string, indentCount int) (out *node) {
	out = &node{name, make(attrMap), remainder, indentCount, self, new(vector.Vector), " /", false, false}
	self.children.Push(out)
	return
}

func (self *node) topLevel() (isTopLevel bool) {
	return nil == self.parent.parent
	return
}

func (self *node) childAt(i int) (child *node) {
	child = self.children.At(i).(*node)
	return;
}

func (self *node) setAutocloseOff() {
	self.closeTag = ""
}

func (self *node) setNoNewline() {
	self.noNewline = true
}

func (self *node) needsLookup() {
	self.remainderLookup = true
}

func (self *node) appendAttr(name string, value string) {
	self.attrs.appendAttr(name, value)
	return
}

func (self *tree) String(lookup func(input string) (output string)) (output string) {
	buf := bytes.NewBuffer(make([]byte, 0))
	output = ""
	for i := 0; i < self.children.Len(); i++ {
		node := self.childAt(i)
		node.String(i == self.children.Len() - 1, "", self.indent, buf, lookup)
	}
	output = buf.String()
	return
}

func (self *node) String(last bool, indent string, customIndent string, buf *bytes.Buffer, lookup func(input string) (output string)) (output string) {
	remainder := self.remainder
	if self.remainderLookup {remainder = lookup(remainder)}
	lineEnd := "\n"
	if last || self.noNewline {lineEnd = ""}
	if len(self.name) > 0 && len(remainder) > 0 {
		buf.WriteString(indent)
		buf.WriteString(fmt.Sprintf("<%s%s>%s</%s>%s", self.name, self.attrs, remainder, self.name, lineEnd))
	} else if len(self.name) > 0 && self.children.Len() > 0 {
		if self.noNewline {
			buf.WriteString(indent)
			buf.WriteString(fmt.Sprintf("<%s%s>", self.name, self.attrs))
		} else {
			buf.WriteString(indent)
			buf.WriteString(fmt.Sprintf("<%s%s>\n", self.name, self.attrs))
		}
		childIndent := indent + customIndent
		childrenLen := self.children.Len()
		for i := 0; i < childrenLen; i++ {
			nextIndent := childIndent
			if 0 == i && self.noNewline {
				nextIndent = ""
			} else if self.noNewline {
				nextIndent = indent
			}
			node := self.childAt(i)
			lastNodeNeedsNoNewline := self.noNewline && i == childrenLen - 1
			node.String(lastNodeNeedsNoNewline, nextIndent, customIndent, buf, lookup)
		}
		if self.noNewline {
			buf.WriteString(strings.TrimRightFunc(output, unicode.IsSpace))
			indent = ""
			if !last {lineEnd = "\n"}
		}
		buf.WriteString(indent)
		buf.WriteString(fmt.Sprintf("</%s>%s", self.name, lineEnd))
	} else if len(self.name) > 0 {
		buf.WriteString(indent)
		buf.WriteString(fmt.Sprintf("<%s%s%s>%s", self.name, self.attrs, self.closeTag, lineEnd))
	} else if len(remainder) > 0 {
		buf.WriteString(indent)
		buf.WriteString(remainder)
		buf.WriteString(lineEnd)
	}
	return
}
