package gohaml

import (
	"bytes"
	"container/vector"
	"strings"
	"reflect"
	"fmt"
)

type res struct {
	value           string
	needsResolution bool
}

type resPair struct {
	key   res
	value res
}

type inode interface {
	parent() inode
	indentLevel() int
	setIndentLevel(i int)
	addChild(n inode)
	noNewline() bool
	resolve(scope map[string]interface{}, buf *bytes.Buffer, curIndent string, indent string, autoclose bool)
	setParent(n inode)
	nil() bool
}

type icodenode interface {
	setLHS(s string)
	parent() inode
	indentLevel() int
	setIndentLevel(i int)
	addChild(n inode)
	noNewline() bool
	resolve(scope map[string]interface{}, buf *bytes.Buffer, curIndent string, indent string, autoclose bool)
	setParent(n inode)
	nil() bool
}

type node struct {
	_parent      inode
	_remainder   res
	_name        string
	_attrs       vector.Vector
	_noNewline   bool
	_autoclose   bool
	_indentLevel int
	_children    vector.Vector
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
		curr := self.resolveValue(scope)

	OutputSwitch:
		switch t := curr; t.Kind() {
		case reflect.String:
			output = t.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			output = fmt.Sprint(t.Int())
		case reflect.Float32, reflect.Float64:
			output = fmt.Sprint(t.Float())
		case reflect.Ptr:
			if !t.IsNil() {
				curr = t.Elem()
				goto OutputSwitch
			}
			output = ""
		case reflect.Interface:
			curr = t.Elem()
			goto OutputSwitch
		default:
			output = fmt.Sprint(curr)
		}
	}
	return
}

func (self res) resolveValue(scope map[string]interface{}) (value reflect.Value) {
	keyPath := strings.Split(self.value, ".")
	curr := reflect.ValueOf(scope[keyPath[0]])
	for _, key := range keyPath[1:] {
	TypeSwitch:
		switch t := curr; t.Kind() {
		case reflect.Ptr:
			curr = t.Elem()
			goto TypeSwitch
		case reflect.Struct:
			curr = t.FieldByName(key)
		case reflect.Map:
			curr = t.MapIndex(reflect.ValueOf(key))
		}
	}
	value = curr
	return
}

func (self tree) resolve(scope map[string]interface{}, indent string, autoclose bool) (output string) {
	treeLen := self.nodes.Len()
	buf := bytes.NewBuffer(make([]byte, 0))
	for i, n := range self.nodes {
		node := n.(inode)
		node.resolve(scope, buf, "", indent, autoclose)
		if i != treeLen-1 && !node.noNewline() {
			buf.WriteString("\n")
		}
	}
	output = buf.String()
	return
}

func (self node) resolve(scope map[string]interface{}, buf *bytes.Buffer, curIndent string, indent string, autoclose bool) {
	remainder := self._remainder.resolve(scope)
	if self._attrs.Len() > 0 && len(remainder) > 0 {
		if len(self._name) == 0 {
			self._name = "div"
		}
		buf.WriteString("<")
		buf.WriteString(self._name)
		self.resolveAttrs(scope, buf)
		buf.WriteString(">")
		buf.WriteString(remainder)
		buf.WriteString("</")
		buf.WriteString(self._name)
		buf.WriteString(">")
	} else if self._attrs.Len() > 0 {
		if len(self._name) == 0 {
			self._name = "div"
		}
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
	if self._noNewline {
		ind = curIndent
	}
	childLen := self._children.Len()
	if childLen > 0 {
		buf.WriteString(">")
		for i, n := range self._children {
			node := n.(inode)
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
		if value == "false" {
			continue
		}
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
		value = value[1 : len(value)-1]
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

func (self *node) setRemainder(value string, needsResolution bool) {
	self._remainder = res{value, needsResolution}
}

func (self *node) setNoNewline(b bool) {
	self._noNewline = b
}

func (self *node) noNewline() bool {
	return self._noNewline
}

func (self *node) setParent(n inode) {
	self._parent = n
}

func (self *node) nil() bool {
	return self == nil
}

type rangenode struct {
	_parent      inode
	_indentLevel int
	_children    vector.Vector

	_lhs1, _lhs2 string
	_rhs         res
}

func (self *rangenode) parent() inode {
	return self._parent
}

func (self *rangenode) indentLevel() int {
	return self._indentLevel
}

func (self *rangenode) setIndentLevel(i int) {
	self._indentLevel = i
}

func (self *rangenode) addChild(n inode) {
	n.setParent(self)
	self._children.Push(n)
}

func (self *rangenode) noNewline() bool {
	return false
}

func (self *rangenode) resolve(scope map[string]interface{}, buf *bytes.Buffer, curIndent string, indent string, autoclose bool) {
	oldlhs1, oklhs1 := scope[self._lhs1]
	oldlhs2, oklhs2 := scope[self._lhs2]

	value := self._rhs.resolveValue(scope)

	switch t := value; t.Kind() {
	case reflect.Slice:
		for i := 0; i < t.Len(); i++ {
			v := t.Index(i)
			var iv interface{}

			switch t := v; t.Kind() {
			case reflect.String:
				iv = t.String()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				iv = fmt.Sprint(t.Int())
			case reflect.Float32, reflect.Float64:
				iv = fmt.Sprint(t.Float())
			}

			scope[self._lhs1] = i
			scope[self._lhs2] = iv

			for _, n := range self._children {
				node := n.(inode)
				node.resolve(scope, buf, curIndent, indent, autoclose)
				if i != t.Len()-1 && !node.noNewline() {
					buf.WriteString("\n")
					buf.WriteString(curIndent)
				}
			}
		}
	case reflect.Array:
		for i := 0; i < t.Len(); i++ {
			v := t.Index(i)
			var iv interface{}

			switch t := v; t.Kind() {
			case reflect.String:
				iv = t.String()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				iv = fmt.Sprint(t.Int())
			case reflect.Float32, reflect.Float64:
				iv = fmt.Sprint(t.Float())
			}

			scope[self._lhs1] = i
			scope[self._lhs2] = iv

			for _, n := range self._children {
				node := n.(inode)
				node.resolve(scope, buf, curIndent, indent, autoclose)
				if i != t.Len()-1 && !node.noNewline() {
					buf.WriteString("\n")
					buf.WriteString(curIndent)
				}
			}
		}
	case reflect.Map:
		for i, k := range t.MapKeys() {
			v := t.MapIndex(k)

			var iv interface{}

			switch t := k; t.Kind() {
			case reflect.String:
				iv = t.String()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				iv = fmt.Sprint(t.Int())
			case reflect.Float32, reflect.Float64:
				iv = fmt.Sprint(t.Float())
			}

			scope[self._lhs1] = iv

			switch t := v; t.Kind() {
			case reflect.String:
				iv = t.String()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				iv = fmt.Sprint(t.Int())
			case reflect.Float32, reflect.Float64:
				iv = fmt.Sprint(t.Float())
			}

			scope[self._lhs2] = iv

			for _, n := range self._children {
				node := n.(inode)
				node.resolve(scope, buf, curIndent, indent, autoclose)
				if i != t.Len()-1 && !node.noNewline() {
					buf.WriteString("\n")
					buf.WriteString(curIndent)
				}
			}
		}
	}

	scope[self._lhs1] = oldlhs1, oklhs1
	scope[self._lhs2] = oldlhs2, oklhs2
}

func (self *rangenode) setParent(n inode) {
	self._parent = n
}

func (self *rangenode) nil() bool {
	return self == nil
}

type declassnode struct {
	_parent      inode
	_indentLevel int
	_children    vector.Vector

	_lhs string
	_rhs interface{}
}

func (self *declassnode) parent() inode {
	return self._parent
}

func (self *declassnode) indentLevel() int {
	return self._indentLevel
}

func (self *declassnode) setIndentLevel(i int) {
	self._indentLevel = i
}

func (self *declassnode) addChild(n inode) {
	n.setParent(self)
	self._children.Push(n)
}

func (self *declassnode) noNewline() bool {
	return false
}

func (self *declassnode) resolve(scope map[string]interface{}, buf *bytes.Buffer, curIndent string, indent string, autoclose bool) {
	scope[self._lhs] = self._rhs
}

func (self *declassnode) setParent(n inode) {
	self._parent = n
}

func (self *declassnode) nil() bool {
	return self == nil
}

func (self *declassnode) setLHS(s string) {
	self._lhs = s
}

type vdeclassnode struct {
	_parent      inode
	_indentLevel int
	_children    vector.Vector

	_lhs string
	_rhs res
}

func (self *vdeclassnode) parent() inode {
	return self._parent
}

func (self *vdeclassnode) indentLevel() int {
	return self._indentLevel
}

func (self *vdeclassnode) setIndentLevel(i int) {
	self._indentLevel = i
}

func (self *vdeclassnode) addChild(n inode) {
	n.setParent(self)
	self._children.Push(n)
}

func (self *vdeclassnode) noNewline() bool {
	return false
}

func (self *vdeclassnode) resolve(scope map[string]interface{}, buf *bytes.Buffer, curIndent string, indent string, autoclose bool) {
	scope[self._lhs] = self._rhs.resolve(scope)
}

func (self *vdeclassnode) setParent(n inode) {
	self._parent = n
}

func (self *vdeclassnode) nil() bool {
	return self == nil
}

func (self *vdeclassnode) setLHS(s string) {
	self._lhs = s
}
