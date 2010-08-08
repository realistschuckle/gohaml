package gohaml

import (
	"bytes"
	"container/vector"
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

type Node struct {
	name string
	attrs attrMap
	remainder string
	indentCount int
	parent *Node
	children *vector.Vector
}

type Tree struct {
	Node
}

func newTree() (tree *Tree) {
	tree = &Tree{}
	return
}

func (self *Node) createChild(name string, remainder string, indentCount int) (node *Node) {
	node = &Node{name, make(attrMap), remainder, indentCount, self, new(vector.Vector)}
	return
}

func (self *Node) topLevel() (isTopLevel bool) {
	return nil == self.parent.parent
	return
}

func (self *Node) childAt(i int) (child *Node) {
	child = self.children.At(i).(*Node)
	return;
}
