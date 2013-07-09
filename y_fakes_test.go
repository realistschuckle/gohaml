package gohaml

import (
	"strings"
)

type fakeNode struct {
	calledCount int
	compiler    HamlCompiler
}

func (self *fakeNode) Accept(c HamlCompiler) {
	self.calledCount += 1
	self.compiler = c
	return
}

type fakeNodeParser struct {
	next        []NodeParser
	runesToRead int
	n           Node
	e           error
}

func (self *fakeNodeParser) Parse(i *strings.Reader) (n Node, e error) {
	for x := 0; x < self.runesToRead; x += 1 {
		i.ReadRune()
	}
	n = self.n
	e = self.e
	return
}

func (self *fakeNodeParser) Next() (n []NodeParser) {
	n = self.next
	return
}

type fakeCompiler struct {
	visitDocTypeCalled bool
	docTypeNode        *DocTypeNode
	visitTagCalled     bool
	tagNode            *TagNode
}

func (self *fakeCompiler) Compile() (c CompiledDocument, e error) {
	return
}

func (self *fakeCompiler) VisitDocType(n *DocTypeNode) {
	self.visitDocTypeCalled = true
	self.docTypeNode = n
}

func (self *fakeCompiler) VisitTag(n *TagNode) {
	self.visitTagCalled = true
	self.tagNode = n
}

type fakeParsedDoc struct {
	acceptCalled bool
	compiler     HamlCompiler
}

func (self *fakeParsedDoc) Accept(c HamlCompiler) {
	self.acceptCalled = true
	self.compiler = c
}
