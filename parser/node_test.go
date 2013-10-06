package parser

import (
	"testing"
)

func TestDoctypeNodeAcceptsVisitor(t *testing.T) {
	node := DoctypeNode{}
	visitor := new(mockVisitor)
	visitor.On("VisitDoctype", &node).Return()

	node.Accept(visitor)

	visitor.AssertCalled(t, "VisitDoctype", &node)
}
