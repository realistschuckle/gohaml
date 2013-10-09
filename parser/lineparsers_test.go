package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDoctypeParserReturnsErrorWhenDoesNotStartWithTripleBang(t *testing.T) {
	input := []rune("Not a valid doctype")
	parser := DoctypeParser{}

	_, e := parser.Parse(input)

	assert.NotNil(t, e)
}

func TestDoctypeParserReturnsDoctypeNodeWithDoctypeSpecifier(t *testing.T) {
	input := []rune("!!! some_specifier  \r\n")
	parser := DoctypeParser{}

	n, e := parser.Parse(input)

	assert.Nil(t, e)
	assert.NotNil(t, n)

	dn := n.(*DoctypeNode)
	assert.Equal(t, "some_specifier", dn.Specifier)
}

func TestTagParserReturnsTagNodeWithNameForPercentSignInput(t *testing.T) {
	input := []rune("%html")
	parser := TagParser{}

	n, e := parser.Parse(input)

	if ok := assert.Nil(t, e); !ok {
		return
	}
	if ok := assert.NotNil(t, n); !ok {
		return
	}

	dn := n.(*TagNode)
	assert.Equal(t, "html", dn.Name)
	assert.Equal(t, "", dn.Id)
	assert.Equal(t, 0, len(dn.Classes))
	assert.Equal(t, 0, len(dn.Attrs))
	assert.Equal(t, 0, len(dn.Children))
	assert.False(t, dn.Close)
}

func TestTagParserReturnsCloseFlagTrueForIndicator(t *testing.T) {
	input := []rune("%giggety/")
	parser := TagParser{}

	n, e := parser.Parse(input)

	if ok := assert.Nil(t, e); !ok {
		return
	}
	if ok := assert.NotNil(t, n); !ok {
		return
	}

	dn := n.(*TagNode)
	assert.Equal(t, "giggety", dn.Name)
	assert.Equal(t, "", dn.Id)
	assert.Equal(t, 0, len(dn.Classes))
	assert.Equal(t, 0, len(dn.Attrs))
	assert.Equal(t, 0, len(dn.Children))
	assert.True(t, dn.Close)
}
