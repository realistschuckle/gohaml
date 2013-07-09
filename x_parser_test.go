package gohaml

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParseReturnsErrorAndDoesNotAdvanceReaderForNoNodeParsers(t *testing.T) {
	r := strings.NewReader("1234567890")
	p := parser{r, nil}
	_, e := p.Parse()

	assert.NotNil(t, e)
}

func TestParseSucceedsForFullyParsedInputByOneNodeParser(t *testing.T) {
	r := strings.NewReader("123")
	n := &fakeNode{}
	np := &fakeNodeParser{nil, 3, n, nil}
	p := parser{r, []NodeParser{np}}

	pdoc, e := p.Parse()
	doc := pdoc.(*parsedDoc)

	assert.Nil(t, e)
	assert.NotNil(t, doc)
	assert.Equal(t, r.Len(), 0)
	assert.Equal(t, doc.nodes.Len(), 1)
	assert.Equal(t, doc.nodes.Front().Value, n)
}

func TestParseSucceedsForFullyParsedInputByTwoNodeParsersSerially(t *testing.T) {
	r := strings.NewReader("1234")
	n1 := &fakeNode{}
	n2 := &fakeNode{}
	np2 := &fakeNodeParser{nil, 1, n2, nil}
	np1 := &fakeNodeParser{[]NodeParser{np2}, 3, n1, nil}
	p := parser{r, []NodeParser{np1}}

	pdoc, e := p.Parse()
	doc := pdoc.(*parsedDoc)

	assert.Nil(t, e)
	assert.NotNil(t, doc)
	assert.Equal(t, r.Len(), 0)
	assert.Equal(t, doc.nodes.Len(), 2)
	assert.Equal(t, doc.nodes.Front().Value, n1)
	assert.Equal(t, doc.nodes.Back().Value, n2)
}

func TestParseSucceedsForFullyParsedInputByTwoNodeParsersInSameStep(t *testing.T) {
	r := strings.NewReader("123")
	n1 := &fakeNode{}
	n2 := &fakeNode{}
	np1 := &fakeNodeParser{nil, 0, n1, errors.New("")}
	np2 := &fakeNodeParser{nil, 3, n2, nil}
	p := parser{r, []NodeParser{np1, np2}}

	pdoc, e := p.Parse()
	doc := pdoc.(*parsedDoc)

	assert.Nil(t, e)
	assert.NotNil(t, doc)
	assert.Equal(t, r.Len(), 0)
	assert.Equal(t, doc.nodes.Len(), 1)
	assert.Equal(t, doc.nodes.Front().Value, n2)
}

func TestParseFailsForPartiallyParsedInputAndNoNextState(t *testing.T) {
	r := strings.NewReader("1234")
	n := &fakeNode{}
	np := &fakeNodeParser{nil, 3, n, nil}
	p := parser{r, []NodeParser{np}}

	_, e := p.Parse()

	assert.NotNil(t, e)
}
