package gohaml

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestDoesNotParseSomethingThatDoesNotStartWithDecimal(t *testing.T) {
	input := strings.NewReader("abcdefg")
	p := &classNameParser{}
	node, err := p.Parse(input)

	assert.Nil(t, node)
	assert.NotNil(t, err)
	assert.Equal(t, input.Len(), 7)
}

func TestForClassNamesParsesAlphaCharactersUntilEndOfInput(t *testing.T) {
	inputs := []string{"a", "ab", "ab1", "ab1c"}
	for _, input := range inputs {
		p := &classNameParser{}
		s := fmt.Sprintf(".%s", input)
		r := strings.NewReader(s)
		node, err := p.Parse(r)
		tag := node.(*ClassNameNode)

		assert.Nil(t, err)
		assert.NotNil(t, node)
		assert.Equal(t, tag.Name, input)
		assert.Equal(t, r.Len(), 0)
	}
}
