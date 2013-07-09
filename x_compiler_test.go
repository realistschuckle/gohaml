package gohaml

import (
	"container/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompileCallsAcceptOnParsedDocument(t *testing.T) {
	doc := fakeParsedDoc{}
	opts := DefaultEngineOptions()

	c := compiler{&doc, &opts, nil}

	c.Compile()

	assert.True(t, doc.acceptCalled)
	assert.Equal(t, doc.compiler, &c)
}

func assertDoctype(t *testing.T, format string, spec string, expected string) {
	n := DocTypeNode{spec}
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

func assertXhtmlDoctype(t *testing.T, spec string, expected string) {
	assertDoctype(t, "xhtml", spec, expected)
}

func assertHtml5Doctype(t *testing.T, spec string, expected string) {
	assertDoctype(t, "html5", spec, expected)
}

func assertHtml4Doctype(t *testing.T, spec string, expected string) {
	assertDoctype(t, "html4", spec, expected)
}

func TestForXmlDeclarationVisitDocTypeNodePutsXmlDeclarationIntoRenderedValue(t *testing.T) {
	spec := "XML"
	expected := "<?xml version='1.0' encoding='utf-8' ?>"
	assertXhtmlDoctype(t, spec, expected)
}

func TestForEmptyDeclarationWithXhtmlFormatVisitDocTypeNodePutsXmlDeclarationIntoRenderedValue(t *testing.T) {
	spec := ""
	expected := "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">"
	assertXhtmlDoctype(t, spec, expected)
}

func TestFor1Point1DeclarationWithXhtmlFormatVisitDocTypeNodePutsXmlDeclarationIntoRenderedValue(t *testing.T) {
	spec := "1.1"
	expected := "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.1//EN\" \"http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd\">"
	assertXhtmlDoctype(t, spec, expected)
}

func TestForMobileDeclarationWithXhtmlFormatVisitDocTypeNodePutsXmlDeclarationIntoRenderedValue(t *testing.T) {
	spec := "mobile"
	expected := "<!DOCTYPE html PUBLIC \"-//WAPFORUM//DTD XHTML Mobile 1.2//EN\" \"http://www.openmobilealliance.org/tech/DTD/xhtml-mobile12.dtd\">"
	assertXhtmlDoctype(t, spec, expected)
}

func TestForBasicDeclarationWithXhtmlFormatVisitDocTypeNodePutsXmlDeclarationIntoRenderedValue(t *testing.T) {
	spec := "basic"
	expected := "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML Basic 1.1//EN\" \"http://www.w3.org/TR/xhtml-basic/xhtml-basic11.dtd\">"
	assertXhtmlDoctype(t, spec, expected)
}

func TestForFramesetDeclarationWithXhtmlFormatVisitDocTypeNodePutsXmlDeclarationIntoRenderedValue(t *testing.T) {
	spec := "frameset"
	expected := "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Frameset//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-frameset.dtd\">"
	assertXhtmlDoctype(t, spec, expected)
}

func TestFor5DeclarationWithXhtmlFormatVisitDocTypeNodePutsXmlDeclarationIntoRenderedValue(t *testing.T) {
	spec := "5"
	expected := "<!DOCTYPE html>"
	assertXhtmlDoctype(t, spec, expected)
}

func TestForSilentXmlDeclarationWithHtml5FormatVisitDocTypeNodePutsXmlDeclarationIntoRenderedValue(t *testing.T) {
	spec := "XML"
	expected := ""
	assertHtml5Doctype(t, spec, expected)
}

func TestForHtml5DeclarationWithHtml5FormatVisitDocTypeNodePutsXmlDeclarationIntoRenderedValue(t *testing.T) {
	spec := ""
	expected := "<!DOCTYPE html>"
	assertHtml5Doctype(t, spec, expected)
}

func TestForSilentXmlDeclarationWithHtml4FormatVisitDocTypeNodePutsXmlDeclarationIntoRenderedValue(t *testing.T) {
	spec := "XML"
	expected := ""
	assertHtml4Doctype(t, spec, expected)
}

func TestForHtml4TransitionalDeclarationWithHtml4FormatVisitDocTypeNodePutsXmlDeclarationIntoRenderedValue(t *testing.T) {
	spec := ""
	expected := "<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 4.01 Transitional//EN\" \"http://www.w3.org/TR/html4/loose.dtd\">"
	assertHtml4Doctype(t, spec, expected)
}

func TestForHtml4FramesetDeclarationWithHtml4FormatVisitDocTypeNodePutsXmlDeclarationIntoRenderedValue(t *testing.T) {
	spec := "frameset"
	expected := "<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 4.01 Frameset//EN\" \"http://www.w3.org/TR/html4/frameset.dtd\">"
	assertHtml4Doctype(t, spec, expected)
}

func TestForHtml4StrictDeclarationWithHtml4FormatVisitDocTypeNodePutsXmlDeclarationIntoRenderedValue(t *testing.T) {
	spec := "strict"
	expected := "<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 4.01//EN\" \"http://www.w3.org/TR/html4/strict.dtd\">"
	assertHtml4Doctype(t, spec, expected)
}
