package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDoctypeParserReturnsErrorWhenDoesNotStartWithTripleBang(t *testing.T) {
	input := []rune("Not a valid doctype")
	parser := DoctypeParser{}

	_, e := parser.Parse("", input)

	assert.NotNil(t, e)
}

func TestDoctypeParserReturnsDoctypeNodeWithDoctypeSpecifier(t *testing.T) {
	input := []rune("!!! some_specifier  \r\n")
	parser := DoctypeParser{}

	n, e := parser.Parse("", input)

	assert.Nil(t, e)
	assert.NotNil(t, n)

	dn := n.(*DoctypeNode)
	assert.Equal(t, "some_specifier", dn.Specifier)
}

func TestTagParserReturnsErrorWhenDoesNotStartWithPercentSignOrOctothorpe(t *testing.T) {
	input := []rune("blah")
	parser := TagParser{}

	_, e := parser.Parse("", input)

	assert.NotNil(t, e)
}

func TestTagParserReturnsTagNodeWithNameForPercentSignInput(t *testing.T) {
	input := []rune("%html")
	parser := TagParser{}

	n, e := parser.Parse("", input)

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

	n, e := parser.Parse("", input)

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

func TestTagParserReturnsClassNameForTagWithClass(t *testing.T) {
	input := []rune("%sup.ui-helper-hidden")
	parser := TagParser{}

	n, e := parser.Parse("", input)

	if ok := assert.Nil(t, e); !ok {
		return
	}
	if ok := assert.NotNil(t, n); !ok {
		return
	}

	dn := n.(*TagNode)
	assert.Equal(t, "sup", dn.Name)
	assert.Equal(t, "", dn.Id)
	assert.Equal(t, 0, len(dn.Attrs))
	assert.Equal(t, 0, len(dn.Children))
	assert.False(t, dn.Close)

	if ok := assert.Equal(t, 1, len(dn.Classes)); !ok {
		return
	}

	assert.Equal(t, "ui-helper-hidden", dn.Classes[0])
}

func TestTagParserReturnsIdForTagWithId(t *testing.T) {
	input := []rune("%video#vid43")
	parser := TagParser{}

	n, e := parser.Parse("", input)

	if ok := assert.Nil(t, e); !ok {
		return
	}
	if ok := assert.NotNil(t, n); !ok {
		return
	}

	dn := n.(*TagNode)
	assert.Equal(t, "video", dn.Name)
	assert.Equal(t, "vid43", dn.Id)
	assert.Equal(t, 0, len(dn.Attrs))
	assert.Equal(t, 0, len(dn.Children))
	assert.Equal(t, 0, len(dn.Classes))
	assert.False(t, dn.Close)
}

func TestTagParserReturnsDivForJustCssId(t *testing.T) {
	input := []rune("#you-wish")
	parser := TagParser{}

	n, e := parser.Parse("", input)

	if ok := assert.Nil(t, e); !ok {
		return
	}
	if ok := assert.NotNil(t, n); !ok {
		return
	}

	dn := n.(*TagNode)
	assert.Equal(t, "div", dn.Name)
	assert.Equal(t, "you-wish", dn.Id)
	assert.Equal(t, 0, len(dn.Attrs))
	assert.Equal(t, 0, len(dn.Children))
	assert.Equal(t, 0, len(dn.Classes))
	assert.False(t, dn.Close)
}

func TestTagParserReturnsDivForJustCssClass(t *testing.T) {
	input := []rune(".i_am_legend")
	parser := TagParser{}

	n, e := parser.Parse("", input)

	if ok := assert.Nil(t, e); !ok {
		return
	}
	if ok := assert.NotNil(t, n); !ok {
		return
	}

	dn := n.(*TagNode)
	assert.Equal(t, "div", dn.Name)
	assert.Equal(t, "", dn.Id)
	assert.Equal(t, 0, len(dn.Attrs))
	assert.Equal(t, 0, len(dn.Children))
	assert.False(t, dn.Close)

	if ok := assert.Equal(t, 1, len(dn.Classes)); !ok {
		return
	}

	assert.Equal(t, "i_am_legend", dn.Classes[0])
}
