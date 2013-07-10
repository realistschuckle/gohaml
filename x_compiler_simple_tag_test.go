package gohaml

import (
	"container/list"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func tagTest(t *testing.T, tag string, forceClose bool, format string, expected string) {
	n := TagNode{tag, forceClose}
	l := list.New()
	l.PushFront(&n)
	opts := DefaultEngineOptions()
	opts.Format = format
	pdoc := parsedDoc{l}
	scope := make(map[string]interface{})
	c := compiler{&pdoc, &opts, nil}

	cdoc, e1 := c.Compile()
	output, e2 := cdoc.Render(scope)

	assert.Nil(t, e1)
	assert.Nil(t, e2)
	assert.Equal(t, output, expected)
}

func TestSimpleEmptyTag(t *testing.T) {
	tagTest(t, "p", false, "html5", "<p></p>")
}

func TestSimpleEmptySelfClosingTagInXhtml(t *testing.T) {
	tags := []string{"area", "base", "br", "col", "hr", "img", "input", "link", "meta", "param"}
	for _, tag := range tags {
		expected := fmt.Sprintf("<%s />", tag)
		tagTest(t, tag, false, "xhtml", expected)
	}
}

func TestSimpleEmptySelfClosingTagInHtml4(t *testing.T) {
	tags := []string{"area", "base", "br", "col", "hr", "img", "input", "link", "meta", "param"}
	for _, tag := range tags {
		expected := fmt.Sprintf("<%s>", tag)
		tagTest(t, tag, false, "html4", expected)
	}
}

func TestSimpleEmptySelfClosingTagInHtml5(t *testing.T) {
	tags := []string{"area", "base", "br", "col", "hr", "img", "input", "link", "meta", "param"}
	for _, tag := range tags {
		expected := fmt.Sprintf("<%s>", tag)
		tagTest(t, tag, false, "html5", expected)
	}
}

func TestEmptySelfClosingTagWithModifierInXhtml(t *testing.T) {
	tagTest(t, "zzz", true, "xhtml", "<zzz />")
}

func TestEmptySelfClosingTagWithModifierInHtml5(t *testing.T) {
	tagTest(t, "zzz", true, "html5", "<zzz>")
}
