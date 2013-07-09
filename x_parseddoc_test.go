package gohaml

import (
	"container/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccpetPassesCompilerToEachNodeInList(t *testing.T) {
	c := fakeCompiler{}
	n1 := fakeNode{}
	n2 := fakeNode{}
	n3 := fakeNode{}
	l := list.New()
	l.PushBack(&n1)
	l.PushBack(&n2)
	l.PushBack(&n3)
	doc := parsedDoc{l}

	doc.Accept(&c)

	assert.Equal(t, n1.calledCount, 1)
	assert.Equal(t, n1.compiler, &c)
	assert.Equal(t, n2.calledCount, 1)
	assert.Equal(t, n2.compiler, &c)
	assert.Equal(t, n3.calledCount, 1)
	assert.Equal(t, n3.compiler, &c)
}
