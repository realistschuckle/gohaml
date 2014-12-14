package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDoctypeNodeAcceptsVisitor(t *testing.T) {
	node := DoctypeNode{}
	visitor := new(mockVisitor)
	visitor.On("VisitDoctype", &node).Return()

	node.Accept(visitor)

	visitor.AssertCalled(t, "VisitDoctype", &node)
}

func TestDoctypeNodeReturnsFalseForAddingChild(t *testing.T) {
	node := DoctypeNode{}
	assert.False(t, node.AddChild(nil))
	assert.False(t, node.AddChild(&mockNode{}))
}

func TestTagNodeAcceptsVisitor(t *testing.T) {
	node := TagNode{}
	visitor := new(mockVisitor)
	visitor.On("VisitTag", &node).Return()

	node.Accept(visitor)

	visitor.AssertCalled(t, "VisitTag", &node)
}

func TestDoctypeNodeReturnsTrueForAddingNonNilChild(t *testing.T) {
	node := TagNode{}
	child := &mockNode{}
	assert.True(t, node.AddChild(child))
	assert.Equal(t, 1, len(node.Children))
	assert.Equal(t, child, node.Children[0])
}

func TestDoctypeNodeReturnsFalseForAddingNilChild(t *testing.T) {
	node := TagNode{}
	assert.False(t, node.AddChild(nil))
}

func TestStaticNodeReturnsFalseForAddingChild(t *testing.T) {
	node := StaticNode{}
	assert.False(t, node.AddChild(nil))
	assert.False(t, node.AddChild(&mockNode{}))
}

func TestStaticNodeAcceptsVisitor(t *testing.T) {
	node := StaticNode{}
	visitor := new(mockVisitor)

	node.Accept(visitor)

	visitor.AssertNotCalled(t, "VisitStatic", &node)
}

func TestStaticLineNodeReturnsFalseForAddingChild(t *testing.T) {
	node := StaticLineNode{}
	assert.False(t, node.AddChild(nil))
	assert.False(t, node.AddChild(&mockNode{}))
}

func TestStaticLineNodeAcceptsVisitor(t *testing.T) {
	node := StaticLineNode{}
	visitor := new(mockVisitor)
	visitor.On("VisitStaticLine", &node).Return()

	node.Accept(visitor)

	visitor.AssertCalled(t, "VisitStaticLine", &node)
}
