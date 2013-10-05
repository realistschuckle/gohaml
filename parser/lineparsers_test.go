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
