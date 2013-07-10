package gohaml

import (
	"container/list"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleEmptyTag(t *testing.T) {
	n := TagNode{"p", false}
	l := list.New()
	l.PushFront(&n)
	opts := DefaultEngineOptions()
	pdoc := parsedDoc{l}
	scope := make(map[string]interface{})
	c := compiler{&pdoc, &opts, nil}
	expected := "<p></p>"

	cdoc, e1 := c.Compile()
	output, e2 := cdoc.Render(scope)

	assert.Nil(t, e1)
	assert.Nil(t, e2)
	assert.Equal(t, output, expected)
}

func TestSimpleEmptySelfClosingTagInXhtml(t *testing.T) {
	tags := []string{"area", "base", "br", "col", "hr", "img", "input", "link", "meta", "param"}
	for _, tag := range tags {
		n := TagNode{tag, false}
		l := list.New()
		l.PushFront(&n)
		opts := DefaultEngineOptions()
		opts.Format = "xhtml"
		pdoc := parsedDoc{l}
		scope := make(map[string]interface{})
		c := compiler{&pdoc, &opts, nil}
		expected := fmt.Sprintf("<%s />", tag)

		cdoc, e1 := c.Compile()
		output, e2 := cdoc.Render(scope)

		assert.Nil(t, e1)
		assert.Nil(t, e2)
		assert.Equal(t, output, expected)
	}
}

func TestSimpleEmptySelfClosingTagInHtml4(t *testing.T) {
	tags := []string{"area", "base", "br", "col", "hr", "img", "input", "link", "meta", "param"}
	for _, tag := range tags {
		n := TagNode{tag, false}
		l := list.New()
		l.PushFront(&n)
		opts := DefaultEngineOptions()
		opts.Format = "html4"
		pdoc := parsedDoc{l}
		scope := make(map[string]interface{})
		c := compiler{&pdoc, &opts, nil}
		expected := fmt.Sprintf("<%s>", tag)

		cdoc, e1 := c.Compile()
		output, e2 := cdoc.Render(scope)

		assert.Nil(t, e1)
		assert.Nil(t, e2)
		assert.Equal(t, output, expected)
	}
}

func TestSimpleEmptySelfClosingTagInHtml5(t *testing.T) {
	tags := []string{"area", "base", "br", "col", "hr", "img", "input", "link", "meta", "param"}
	for _, tag := range tags {
		n := TagNode{tag, false}
		l := list.New()
		l.PushFront(&n)
		opts := DefaultEngineOptions()
		opts.Format = "html5"
		pdoc := parsedDoc{l}
		scope := make(map[string]interface{})
		c := compiler{&pdoc, &opts, nil}
		expected := fmt.Sprintf("<%s>", tag)

		cdoc, e1 := c.Compile()
		output, e2 := cdoc.Render(scope)

		assert.Nil(t, e1)
		assert.Nil(t, e2)
		assert.Equal(t, output, expected)
	}
}

func TestEmptySelfClosingTagWithModifierInXhtml(t *testing.T) {
	n := TagNode{"zzz", true}
	l := list.New()
	l.PushFront(&n)
	opts := DefaultEngineOptions()
	opts.Format = "xhtml"
	pdoc := parsedDoc{l}
	scope := make(map[string]interface{})
	c := compiler{&pdoc, &opts, nil}
	expected := "<zzz />"

	cdoc, e1 := c.Compile()
	output, e2 := cdoc.Render(scope)

	assert.Nil(t, e1)
	assert.Nil(t, e2)
	assert.Equal(t, output, expected)
}

