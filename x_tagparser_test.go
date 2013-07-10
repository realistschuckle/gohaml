package gohaml

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestDoesNotParseSomethingThatDoesNotStartWithPercentSign(t *testing.T) {
	input := strings.NewReader("abcdefg")
	p := &tagParser{}
	node, err := p.Parse(input)

	assert.Nil(t, node)
	assert.NotNil(t, err)
	assert.Equal(t, input.Len(), 7)
}

func TestForOnlyTagNamesParsesAlphaCharactersUntilEndOfInput(t *testing.T) {
	inputs := []string{"p", "br", "pre", "abbr"}
	for _, input := range inputs {
		p := &tagParser{}
		s := fmt.Sprintf("%%%s", input)
		r := strings.NewReader(s)
		node, err := p.Parse(r)
		tag := node.(*TagNode)

		assert.Nil(t, err)
		assert.NotNil(t, node)
		assert.Equal(t, tag.Name, input)
		assert.Equal(t, tag.ForceClose, false)
		assert.Equal(t, r.Len(), 0)
	}
}

func TestForceCloseTrueForTagNameEndingInSlash(t *testing.T) {
	p := &tagParser{}
	tagName := "foo"
	input := "%foo/"
	r := strings.NewReader(input)
	node, err := p.Parse(r)
	tag := node.(*TagNode)

	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, tag.Name, tagName)
	assert.Equal(t, tag.ForceClose, true)
	assert.Equal(t, r.Len(), 0)
}
