package compiler

import (
	p "github.com/realistschuckle/gohaml/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagWithOnlyTagName(t *testing.T) {
	opts := CompilerOpts{}
	node := &p.TagNode{}
	node.Name = "p"
	nodes := []p.Node{node}
	pdoc := p.ParsedDoc{}
	pdoc.Nodes = nodes
	compiler := DefaultCompiler{}

	cdoc, e := compiler.Compile(pdoc, opts)

	if ok := assert.Nil(t, e); !ok {
		return
	}
	if ok := assert.NotNil(t, cdoc); !ok {
		return
	}
	if ok := assert.Equal(t, 1, len(cdoc.Outputs)); !ok {
		return
	}

	output := cdoc.Outputs[0].(*StaticOutput)
	assert.Equal(t, output.Content, "<p></p>")
}

