package gohaml

import (
	"bytes"
	"container/vector"
	"fmt"
)

type res struct {
	value string
	needsResolution bool
}

type resPair struct {
	key res
	value res
}

type node struct {
	remainder res
	name string
	attrs vector.Vector
	noNewline bool
	indentLevel int
	children vector.Vector
}

type tree struct {
	nodes vector.Vector
}

func newTree() (output *tree) {
	output = &tree{}
	return
}

func (self res) resolve(scope map[string]interface{}) (output string) {
	output = self.value
	if self.needsResolution {output = fmt.Sprint(scope[self.value])}
	return
}

func (self tree) resolve(scope map[string]interface{}) (output string) {
	treeLen := self.nodes.Len()
	buf := bytes.NewBuffer(make([]byte, 0))
	for i, n := range self.nodes {
		node := n.(*node)
		node.resolve(scope, buf)
		if i != treeLen - 1 && !node.noNewline {
			buf.WriteString("\n")
		}
	}
	output = buf.String()
	return
}

func (self *node) addAttr(key string, value string) {
	keyLookup, valueLookup := true, true
	if key[0] == ':' {
		keyLookup = false
		key = key[1:]
	}
	if value[0] == '"' {
		valueLookup = false
		value = value[1:len(value) - 1]
	}
	self.attrs.Push(&resPair{res{key, keyLookup}, res{value, valueLookup}})
}

func (self *node) addAttrNoLookup(key string, value string) {
	self.attrs.Push(&resPair{res{key, false}, res{value, false}})
}

func (self node) resolve(scope map[string]interface{}, buf *bytes.Buffer) {
	remainder := self.remainder.resolve(scope)
	if self.attrs.Len() > 0 && len(remainder) > 0 {
		if len(self.name) == 0 {self.name = "div"}
		buf.WriteString("<")
		buf.WriteString(self.name)
		self.resolveAttrs(scope, buf)
		buf.WriteString(">")
		buf.WriteString(remainder)
		buf.WriteString("</")
		buf.WriteString(self.name)
		buf.WriteString(">")
	} else if self.attrs.Len() > 0 {
		if len(self.name) == 0 {self.name = "div"}
		buf.WriteString("<")
		buf.WriteString(self.name)
		self.resolveAttrs(scope, buf)
		childLen := self.children.Len()
		if childLen > 0 {
			buf.WriteString(">")
			for i, n := range self.children {
				node := n.(*node)
				node.resolve(scope, buf)
				if i != childLen - 1 && !node.noNewline {
					buf.WriteString("\n")
				}
			}
			buf.WriteString("</")
			buf.WriteString(self.name)
			buf.WriteString(">")
		} else {
			buf.WriteString(" />")
		}
	} else if len(self.name) > 0 && len(remainder) > 0 {
		buf.WriteString("<")
		buf.WriteString(self.name)
		buf.WriteString(">")
		buf.WriteString(remainder)
		buf.WriteString("</")
		buf.WriteString(self.name)
		buf.WriteString(">")
	} else if len(self.name) > 0 {
		buf.WriteString("<")
		buf.WriteString(self.name)
		buf.WriteString(" />")
	} else {
		buf.WriteString(remainder)
	}
}

func (self node) resolveAttrs(scope map[string]interface{}, buf *bytes.Buffer) {
	attrMap := make(map[string]string)
	for i := 0; i < self.attrs.Len(); i++ {
		resPair := self.attrs.At(i).(*resPair)
		key, value := resPair.key.resolve(scope), resPair.value.resolve(scope)
		if _, ok := attrMap[key]; ok {
			attrMap[key] += " " + value
		} else {
			attrMap[key] = value
		}
	}
	for key, value := range attrMap {
		buf.WriteString(" ")
		buf.WriteString(key)
		buf.WriteString("=\"")
		buf.WriteString(value)
		buf.WriteString("\"")
	}
}
