package parser

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultParserReadsUntilErrorReturnedFromReadRune(t *testing.T) {
	reader := &mockRuneReader{}
	parser := &DefaultParser{}
	content := []rune("%html\n  %head\n    %title Hello\n  %body This is great!")
	i := 0
	width := int(1)

	for ; i < len(content); i += 1 {
		reader.On("ReadRune").Return(content[i], width, nil).Once()
	}
	reader.On("ReadRune").Return('\000', width, errors.New(""))

	parser.Parse(reader)

	reader.AssertExpectations(t)
	assert.Equal(t, i, len(content))
}

func TestDefaultParserReturnsParsedDocumentWithDoctype(t *testing.T) {
	reader := &mockRuneReader{}
	parser := &DefaultParser{}
	content := []rune("!!! my_specification\n")
	i := 0

	for ; i < len(content); i += 1 {
		reader.On("ReadRune").Return(content[i], 1, nil).Once()
	}
	reader.On("ReadRune").Return('\000', 0, errors.New(""))

	doc, _ := parser.Parse(reader)
	assert.Equal(t, 1, len(doc.Nodes))

	dn := doc.Nodes[0].(*DoctypeNode)
	assert.Equal(t, "my_specification", dn.Specifier)
}

func TestDefaultParserReturnsTag(t *testing.T) {
	reader := &mockRuneReader{}
	parser := &DefaultParser{}
	content := []rune("%p")
	i := 0

	for ; i < len(content); i += 1 {
		reader.On("ReadRune").Return(content[i], 1, nil).Once()
	}
	reader.On("ReadRune").Return('\000', 0, errors.New(""))

	doc, e := parser.Parse(reader)
	if ok := assert.Nil(t, e); !ok {
		return
	}

	assert.Equal(t, 1, len(doc.Nodes))

	dn := doc.Nodes[0].(*TagNode)
	assert.Equal(t, "p", dn.Name)
	assert.Equal(t, "", dn.Id)
	assert.Equal(t, 0, len(dn.Classes))
	assert.Equal(t, 0, len(dn.Attrs))
	assert.Equal(t, 0, len(dn.Children))
}
