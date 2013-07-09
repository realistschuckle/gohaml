package gohaml

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDocTypeNodeAcceptCallsVisitDocTypeOnCompiler(t *testing.T) {
	c := fakeCompiler{}
	n := DocTypeNode{}

	n.Accept(&c)

	assert.True(t, c.visitDocTypeCalled)
	assert.Equal(t, c.docTypeNode, &n)
}

func TestTagNodeAcceptCallsVisitTagOnCompiler(t *testing.T) {
	c := fakeCompiler{}
	n := TagNode{}

	n.Accept(&c)

	assert.True(t, c.visitTagCalled)
	assert.Equal(t, c.tagNode, &n)
}
