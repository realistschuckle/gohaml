package parser

import (
	"github.com/realistschuckle/testify/mock"
)

type mockNode struct {
	mock.Mock
}

func (self *mockNode) Accept(visitor NodeVisitor) {
	self.Mock.Called(visitor)
}

func (self *mockNode) AddChild(child Node) bool {
	args := self.Mock.Called(child)
	return args.Bool(0)
}

type mockVisitor struct {
	mock.Mock
}

func (self *mockVisitor) VisitDoctype(node *DoctypeNode) {
	self.Mock.Called(node)
}

func (self *mockVisitor) VisitTag(node *TagNode) {
	self.Mock.Called(node)
}

func (self *mockVisitor) VisitStatic(node *StaticNode) {
	self.Mock.Called(node)
}

type mockRuneReader struct {
	mock.Mock
}

func (self *mockRuneReader) ReadRune() (rune, int, error) {
	args := self.Mock.Called()
	return args.Get(0).(rune), args.Int(1), args.Error(2)
}
