package gohaml

import (
	"container/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleEmptyTag(t *testing.T) {
	n := TagNode{"p"}
	l := list.New()
	l.PushFront(&n)
	opts := DefaultEngineOptions()
	pdoc := parsedDoc{l}
	scope := make(map[string]interface{})
	c := compiler{&pdoc, &opts, nil}
	expected := "<p></p>"

	cdoc, e1 := c.Compile()
	output, e2 := cdoc.Render(scope)

	return
	assert.Nil(t, e1)
	assert.Nil(t, e2)
	assert.Equal(t, output, expected)
}
