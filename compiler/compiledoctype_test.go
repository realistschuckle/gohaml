package compiler

import (
	p "github.com/realistschuckle/gohaml/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDoctypeXmlSpecWithXhtmlFormat(t *testing.T) {
	opts := CompilerOpts{}
	opts.Format = "xhtml"
	nodes := []p.Node{&p.DoctypeNode{"XML"}}
	pdoc := p.ParsedDoc{}
	pdoc.Nodes = nodes
	compiler := DefaultCompiler{}

	cdoc, e := compiler.Compile(pdoc, opts)

	assert.Nil(t, e)
	assert.NotNil(t, cdoc)
	assert.Equal(t, len(cdoc.Outputs), 1)

	output := cdoc.Outputs[0].(*StaticOutput)
	assert.Equal(t, output.Content, "<?xml version='1.0' encoding='utf-8' ?>")
}

func TestDoctypeEmptySpecWithXhtmlFormat(t *testing.T) {
	opts := CompilerOpts{}
	opts.Format = "xhtml"
	nodes := []p.Node{&p.DoctypeNode{""}}
	pdoc := p.ParsedDoc{}
	pdoc.Nodes = nodes
	compiler := DefaultCompiler{}

	cdoc, e := compiler.Compile(pdoc, opts)

	assert.Nil(t, e)
	assert.NotNil(t, cdoc)
	assert.Equal(t, len(cdoc.Outputs), 1)

	output := cdoc.Outputs[0].(*StaticOutput)
	assert.Equal(t, output.Content, "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">")
}

func TestDoctype1Point1SpecWithXhtmlFormat(t *testing.T) {
	opts := CompilerOpts{}
	opts.Format = "xhtml"
	nodes := []p.Node{&p.DoctypeNode{"1.1"}}
	pdoc := p.ParsedDoc{}
	pdoc.Nodes = nodes
	compiler := DefaultCompiler{}

	cdoc, e := compiler.Compile(pdoc, opts)

	assert.Nil(t, e)
	assert.NotNil(t, cdoc)
	assert.Equal(t, len(cdoc.Outputs), 1)

	output := cdoc.Outputs[0].(*StaticOutput)
	assert.Equal(t, output.Content, "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.1//EN\" \"http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd\">")
}

func TestDoctypeMobileSpecWithXhtmlFormat(t *testing.T) {
	opts := CompilerOpts{}
	opts.Format = "xhtml"
	nodes := []p.Node{&p.DoctypeNode{"mobile"}}
	pdoc := p.ParsedDoc{}
	pdoc.Nodes = nodes
	compiler := DefaultCompiler{}

	cdoc, e := compiler.Compile(pdoc, opts)

	assert.Nil(t, e)
	assert.NotNil(t, cdoc)
	assert.Equal(t, len(cdoc.Outputs), 1)

	output := cdoc.Outputs[0].(*StaticOutput)
	assert.Equal(t, output.Content, "<!DOCTYPE html PUBLIC \"-//WAPFORUM//DTD XHTML Mobile 1.2//EN\" \"http://www.openmobilealliance.org/tech/DTD/xhtml-mobile12.dtd\">")
}
