package parser

import (
	"testing"
	"github.com/stretchr/testify/assert"
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
