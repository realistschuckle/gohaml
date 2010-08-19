package gohaml

import (
	"bytes"
	"container/vector"
	"strings"
	"reflect"
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
	parent *node
	remainder res
	name string
	attrs vector.Vector
	noNewline bool
	autoclose bool
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
	if self.needsResolution {
		keyPath := strings.Split(self.value, ".", -1)
		curr := reflect.NewValue(scope[keyPath[0]])
		for _, key := range keyPath[1:] {
			TypeSwitch:
			switch t := curr.(type) {
			case *reflect.PtrValue:
				curr = t.Elem()
				goto TypeSwitch
			case *reflect.StructValue:
				curr = t.FieldByName(key)
			case *reflect.MapValue:
				curr = t.Elem(reflect.NewValue(key))
			}
		}
		
		OutputSwitch:
		switch t := curr.(type) {
		case *reflect.StringValue:
			output = t.Get()
		case *reflect.IntValue:
			output = fmt.Sprint(t.Get())
		case *reflect.FloatValue:
			output = fmt.Sprint(t.Get())
		case *reflect.PtrValue:
			if !t.IsNil() {goto OutputSwitch}
			output = ""
		case *reflect.InterfaceValue:
			curr = t.Elem()
			goto OutputSwitch
		default:
			output = fmt.Sprint(curr)
		}
	}
	return
}

func (self tree) resolve(scope map[string]interface{}, indent string, autoclose bool) (output string) {
	treeLen := self.nodes.Len()
	buf := bytes.NewBuffer(make([]byte, 0))
	for i, n := range self.nodes {
		node := n.(*node)
		node.resolve(scope, buf, "", indent, autoclose)
		if i != treeLen - 1 && !node.noNewline {
			buf.WriteString("\n")
		}
	}
	output = buf.String()
	return
}

func (self node) resolve(scope map[string]interface{}, buf *bytes.Buffer, curIndent string, indent string, autoclose bool) {
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
		self.outputChildren(scope, buf, curIndent, indent, autoclose)
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
		self.outputChildren(scope, buf, curIndent, indent, autoclose)
	} else {
		buf.WriteString(remainder)
	}
}

func (self node) outputChildren(scope map[string]interface{}, buf *bytes.Buffer, curIndent string, indent string, autoclose bool) {
	ind := curIndent + indent
	if self.noNewline {ind = curIndent}
	childLen := self.children.Len()
	if childLen > 0 {
		buf.WriteString(">")
		for i, n := range self.children {
			node := n.(*node)
			if i != 0 || !self.noNewline {
				buf.WriteString("\n")
				buf.WriteString(ind)
			}
			node.resolve(scope, buf, ind, indent, autoclose)
		}
		if !self.noNewline {
			buf.WriteString("\n")
			buf.WriteString(curIndent)
		}
		buf.WriteString("</")
		buf.WriteString(self.name)
		buf.WriteString(">")
	} else {
		if autoclose || self.autoclose {
			buf.WriteString(" />")
		} else {
			buf.WriteString(">")
		}
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
		if value == "false" {continue}
		buf.WriteString(" ")
		buf.WriteString(key)
		buf.WriteString("=\"")
		if value == "true" {
			buf.WriteString(key)
		} else {
			buf.WriteString(value)
		}
		buf.WriteString("\"")
	}
}

func (self *node) addChild(n *node) {
	n.parent = self
	self.children.Push(n)
}

func (self *node) addAttr(key string, value string) {
	keyLookup, valueLookup := true, true
	if key[0] == ':' {
		keyLookup = false
		key = key[1:]
	}
	if value == "true" || value == "false" {
		valueLookup = false
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
