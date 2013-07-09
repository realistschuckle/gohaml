package gohaml

import (
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
	inputs := []string{"%p", "%br", "%pre", "%abbr"}
	for _, input := range inputs {
		p := &tagParser{}
		r := strings.NewReader(input)
		node, err := p.Parse(r)

		assert.Nil(t, err)
		assert.NotNil(t, node)
		assert.Equal(t, r.Len(), 0)
	}
}
