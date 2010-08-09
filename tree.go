package gohaml

import (
	"bytes"
	"container/vector"
	"fmt"
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
}

type tree struct {
	node
}

func newTree() (t *tree) {
	t = &tree{node{"tree", nil, "", -1, nil, new(vector.Vector), " /"}}
	return
}

func (self *node) createChild(name string, remainder string, indentCount int) (out *node) {
	out = &node{name, make(attrMap), remainder, indentCount, self, new(vector.Vector), " /"}
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

func (self *node) appendAttr(name string, value string) {
	self.attrs.appendAttr(name, value)
	return
}

func (self *tree) String() (output string) {
	output = ""
	for i := 0; i < self.children.Len(); i++ {
		node := self.childAt(i)
		output += node.String(i == self.children.Len() - 1, "")
	}
	return
}

func (self *node) String(last bool, indent string) (output string) {
	lineEnd := "\n"
	if last {lineEnd = ""}
	if len(self.name) > 0 && len(self.remainder) > 0 {
		output += indent + fmt.Sprintf("<%s%s>%s</%s>%s", self.name, self.attrs, self.remainder, self.name, lineEnd)
	} else if len(self.name) > 0 && self.children.Len() > 0 {
		output += indent + fmt.Sprintf("<%s%s>\n", self.name, self.attrs)
		childIndent := indent + "\t"
		for i := 0; i < self.children.Len(); i++ {
			node := self.childAt(i)
			output += node.String(false, childIndent)
		}
		output += indent + fmt.Sprintf("</%s>%s", self.name, lineEnd)
	} else if len(self.name) > 0 {
		output += indent + fmt.Sprintf("<%s%s%s>%s", self.name, self.attrs, self.closeTag, lineEnd)
	} else if len(self.remainder) > 0 {
		output += indent + self.remainder + lineEnd
	}
	return
}
