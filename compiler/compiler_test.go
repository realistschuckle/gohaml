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
