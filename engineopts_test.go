package gohaml

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestDefaultEngineOptionsHasApostropheAsAttributeWrapper(t *testing.T) {
	opts := DefaultEngineOptions()
	assert.Equal(t, '\'', opts.AttributeWrapper)
}

func TestDefaultEngineOptionsHasAppropriateAutocloseTagList(t *testing.T) {
	opts := DefaultEngineOptions()
	tags := []string{"meta", "img", "link", "br", "hr", "input", "area", "param", "col", "base"}

	assert.Equal(t, 10, len(opts.Autoclose))
	for _, expected := range tags {
		for _, value := range opts.Autoclose {
			if value == expected {
				goto Found
			}
		}
		assert.Fail(t, "Failed to find expected tag", expected)
	Found:
	}
}

func TestDefaultEngineOptionsHasFalseCdataFlag(t *testing.T) {
	opts := DefaultEngineOptions()
	assert.False(t, opts.Cdata)
}

func TestDefaultEngineOptionsHasDefaultCompilerClassAsCompilerClass(t *testing.T) {
	var c DefaultCompiler
	cType := reflect.TypeOf(c)
	opts := DefaultEngineOptions()
	assert.Equal(t, cType, opts.CompilerClass)
}

func TestDefaultEngineOptionsHasUtf8Encoding(t *testing.T) {
	opts := DefaultEngineOptions()
	assert.Equal(t, "UTF-8", opts.Encoding)
}

func TestDefaultEngineOptionsHasEscapesAttributes(t *testing.T) {
	opts := DefaultEngineOptions()
	assert.True(t, opts.EscapeAttributes)
}

func TestDefaultEngineOptionsDoesNotHaveEscapeHtml(t *testing.T) {
	opts := DefaultEngineOptions()
	assert.False(t, opts.EscapeHtml)
}

func TestDefaultEngineOptionsHasHtml5Format(t *testing.T) {
	opts := DefaultEngineOptions()
	assert.Equal(t, "html5", opts.Format)
}

func TestDefaultEngineOptionsHasDataAttributeHyphenation(t *testing.T) {
	opts := DefaultEngineOptions()
	assert.True(t, opts.HyphenateDataAttributes)
}

func TestDefaultEngineOptionsHasDefaultParserClassAsParserClass(t *testing.T) {
	var p DefaultParser
	pType := reflect.TypeOf(p)
	opts := DefaultEngineOptions()
	assert.Equal(t, pType, opts.ParserClass)
}

func TestDefaultEngineOptionsDoesNotRemoveWhitespace(t *testing.T) {
	opts := DefaultEngineOptions()
	assert.False(t, opts.RemoveWhitespace)
}

func TestDefaultEngineOptionsDoesNotSuppressEvaluation(t *testing.T) {
	opts := DefaultEngineOptions()
	assert.False(t, opts.SuppressEval)
}

func TestDefaultEngineOptionsIsUgly(t *testing.T) {
	opts := DefaultEngineOptions()
	assert.True(t, opts.Ugly)
}
