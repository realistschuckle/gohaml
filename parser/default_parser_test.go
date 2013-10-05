package parser

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultParserReadsUntilErrorReturnedFromReadRune(t *testing.T) {
	reader := &mockRuneReader{}
	parser := &DefaultParser{}
	content := []rune("html\n  head\n    title Hello\n  body This is great!")
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
