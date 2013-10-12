package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"fmt"
	"strings"
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

func TestTagParserReturnsDivWithProvidedContentAsChild(t *testing.T) {
	input := []rune("%h1 Hello, World!")
	parser := TagParser{}

	n, e := parser.Parse("", input)

	if ok := assert.Nil(t, e); !ok {
		return
	}
	if ok := assert.NotNil(t, n); !ok {
		return
	}

	dn := n.(*TagNode)
	assert.Equal(t, "h1", dn.Name)
	assert.Equal(t, "", dn.Id)
	assert.Equal(t, 0, len(dn.Attrs))
	assert.Equal(t, 0, len(dn.Classes))
	assert.False(t, dn.Close)
	assert.Empty(t, dn.LineBreak)

	if ok := assert.Equal(t, 1, len(dn.Children), "no static content"); !ok {
		return
	}

	sn := dn.Children[0].(*StaticNode)
	assert.Equal(t, "Hello, World!", sn.Content)
}

func TestStaticParserReturnsStaticLineNodeWithContentAndIndent(t *testing.T) {
	input := []rune("Here's some content.")
	indent := "    "
	parser := StaticParser{}

	n, e := parser.Parse(indent, input)

	if ok := assert.Nil(t, e); !ok {
		return
	}
	if ok := assert.NotNil(t, n); !ok {
		return
	}

	sn := n.(*StaticLineNode)
	assert.Equal(t, string(input), sn.Content)
	assert.Equal(t, indent, sn.Indent)
}

func TestTagParserReturnsDivWithProvidedHtmlStyleAttributes(t *testing.T) {
	input := []rune("%h1(style='background-color:red;')")
	parser := TagParser{}

	n, e := parser.Parse("", input)

	if ok := assert.Nil(t, e); !ok {
		return
	}
	if ok := assert.NotNil(t, n); !ok {
		return
	}

	dn := n.(*TagNode)
	assert.Equal(t, "h1", dn.Name)
	assert.Equal(t, "", dn.Id)
	assert.Equal(t, 0, len(dn.Children))
	assert.Equal(t, 0, len(dn.Classes))
	assert.False(t, dn.Close)
	assert.Empty(t, dn.LineBreak)

	if ok := assert.Equal(t, 1, len(dn.Attrs), "no HTML attributes"); !ok {
		return
	}

	attr := dn.Attrs[0]
	assert.Equal(t, "style", attr.Name)

	if ok := assert.NotNil(t, attr.Value); !ok {
		return
	}

	val := attr.Value.(*StaticNode)
	assert.Equal(t, "background-color:red;", val.Content)
}

func TestTagParserReturnsDivWithManyProvidedHtmlStyleAttributes(t *testing.T) {
	keys := []string{"alpha", "beta", "gamma", "delta"}
	values := []string{"a", "b", "c", "d"}
	attrs := []string{}
	for i := range keys {
		key := keys[i]
		value := values[i]
		attrs = append(attrs, fmt.Sprintf("%s='%s'", key, value))
	}
	rawInput := fmt.Sprintf("%%h1(%s)", strings.Join(attrs, " "))
	input := []rune(rawInput)
	parser := TagParser{}

	n, e := parser.Parse("", input)

	if ok := assert.Nil(t, e); !ok {
		return
	}
	if ok := assert.NotNil(t, n); !ok {
		return
	}

	dn := n.(*TagNode)
	assert.Equal(t, "h1", dn.Name)
	assert.Equal(t, "", dn.Id)
	assert.Equal(t, 0, len(dn.Children))
	assert.Equal(t, 0, len(dn.Classes))
	assert.False(t, dn.Close)
	assert.Empty(t, dn.LineBreak)

	if ok := assert.Equal(t, 4, len(dn.Attrs), "no HTML attributes"); !ok {
		return
	}

	for i := 0; i < len(keys); i += 1 {
		attr := dn.Attrs[i]
		assert.Equal(t, keys[i], attr.Name)
		val := attr.Value.(*StaticNode)
		assert.Equal(t, values[i], val.Content)
	}
}

