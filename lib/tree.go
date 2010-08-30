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

type inode interface {
	parent() inode
	indentLevel() int
	setIndentLevel(i int)
	setName(name string)
	setRemainder(value string, needsResolution bool)
	setNoNewline(b bool)
	setAutoclose(b bool)
	addAttr(key string, value string)
	addAttrNoLookup(key string, value string)
	addChild(n inode)
	noNewline() bool
	resolve(scope map[string]interface{}, buf *bytes.Buffer, curIndent string, indent string, autoclose bool)
	setParent(n inode)
}

type node struct {
	_parent inode
	_remainder res
	_name string
	_attrs vector.Vector
	_noNewline bool
	_autoclose bool
	_indentLevel int
	_children vector.Vector
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
		node := n.(*inode)
		node.resolve(scope, buf, "", indent, autoclose)
		if i != treeLen - 1 && !node.noNewline() {
			buf.WriteString("\n")
		}
	}
	output = buf.String()
	return
}

func (self node) resolve(scope map[string]interface{}, buf *bytes.Buffer, curIndent string, indent string, autoclose bool) {
	remainder := self._remainder.resolve(scope)
	if self._attrs.Len() > 0 && len(remainder) > 0 {
		if len(self._name) == 0 {self._name = "div"}
		buf.WriteString("<")
		buf.WriteString(self._name)
		self.resolveAttrs(scope, buf)
		buf.WriteString(">")
		buf.WriteString(remainder)
		buf.WriteString("</")
		buf.WriteString(self._name)
		buf.WriteString(">")
	} else if self._attrs.Len() > 0 {
		if len(self._name) == 0 {self._name = "div"}
		buf.WriteString("<")
		buf.WriteString(self._name)
		self.resolveAttrs(scope, buf)
		self.outputChildren(scope, buf, curIndent, indent, autoclose)
	} else if len(self._name) > 0 && len(remainder) > 0 {
		buf.WriteString("<")
		buf.WriteString(self._name)
		buf.WriteString(">")
		buf.WriteString(remainder)
		buf.WriteString("</")
		buf.WriteString(self._name)
		buf.WriteString(">")
	} else if len(self._name) > 0 {
		buf.WriteString("<")
		buf.WriteString(self._name)
		self.outputChildren(scope, buf, curIndent, indent, autoclose)
	} else {
		buf.WriteString(remainder)
	}
}

func (self node) outputChildren(scope map[string]interface{}, buf *bytes.Buffer, curIndent string, indent string, autoclose bool) {
	ind := curIndent + indent
	if self._noNewline {ind = curIndent}
	childLen := self._children.Len()
	if childLen > 0 {
		buf.WriteString(">")
		for i, n := range self._children {
			node := n.(*inode)
		if i != 0 || !self._noNewline {
				buf.WriteString("\n")
				buf.WriteString(ind)
			}
			node.resolve(scope, buf, ind, indent, autoclose)
		}
		if !self._noNewline {
			buf.WriteString("\n")
			buf.WriteString(curIndent)
		}
		buf.WriteString("</")
		buf.WriteString(self._name)
		buf.WriteString(">")
	} else {
	if autoclose || self._autoclose {
			buf.WriteString(" />")
		} else {
			buf.WriteString(">")
		}
	}
}

func (self node) resolveAttrs(scope map[string]interface{}, buf *bytes.Buffer) {
	attrMap := make(map[string]string)
	for i := 0; i < self._attrs.Len(); i++ {
		resPair := self._attrs.At(i).(*resPair)
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

func (self *node) addChild(n inode) {
	n.setParent(self)
	self._children.Push(n)
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
	self._attrs.Push(&resPair{res{key, keyLookup}, res{value, valueLookup}})
}

func (self *node) addAttrNoLookup(key string, value string) {
	self._attrs.Push(&resPair{res{key, false}, res{value, false}})
}

func (self *node) parent() inode {
	return self._parent
}

func (self *node) indentLevel() int {
	return self._indentLevel
}

func (self *node) setIndentLevel(i int) {
	self._indentLevel = i
}

func (self *node) setName(name string) {
	self._name = name
}

func (self *node) setRemainder(value string, needsResolution bool) {
	self._remainder = res{value, needsResolution}
}

func (self *node) setNoNewline(b bool) {
	self._noNewline = b
}

func (self *node) setAutoclose(b bool) {
	self._autoclose = b
}

func (self *node) noNewline() bool {
	return self._noNewline;
}

func (self *node) setParent(n inode) {
	self._parent = n
}
