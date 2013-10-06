package parser

import (
	"testing"
)

func TestParsedDocAcceptPassesVisitorToEachContainedNode(t *testing.T) {
	pdoc := ParsedDoc{}
	pdoc.Nodes = []Node{&mockNode{}, &mockNode{}, &mockNode{}, &mockNode{}}
	visitor := &mockVisitor{}
	for i := 0; i < len(pdoc.Nodes); i += 1 {
		mock := pdoc.Nodes[i].(*mockNode)
		mock.On("Accept", visitor).Return()
	}

	pdoc.Accept(visitor)

	for i := 0; i < len(pdoc.Nodes); i += 1 {
		mock := pdoc.Nodes[i].(*mockNode)
		mock.AssertCalled(t, "Accept", visitor)
	}
}
