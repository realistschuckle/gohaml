package parser

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultParserReadsUntilErrorReturnedFromReadRune(t *testing.T) {
	reader := &mockRuneReader{}
	parser := &DefaultParser{}
	content := []rune("!!!\n%html")
	i := 0
	width := int(1)

	for ; i < len(content); i += 1 {
		reader.On("ReadRune").Return(content[i], width, nil).Once()
	}
	reader.On("ReadRune").Return('\000', width, errors.New(""))

	_, e := parser.Parse(reader)

	if ok := assert.Nil(t, e); !ok {
		assert.Fail(t, e.Error())
		return
	}

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

func TestDefaultParserUnderstandsIndentation(t *testing.T) {
	reader := &mockRuneReader{}
	parser := &DefaultParser{}
	content := []rune("%div\n  %div\n    %div")
	i := 0

	for ; i < len(content); i += 1 {
		reader.On("ReadRune").Return(content[i], 1, nil).Once()
	}
	reader.On("ReadRune").Return('\000', 0, errors.New(""))

	_, e := parser.Parse(reader)
	if ok := assert.Nil(t, e); !ok {
		return
	}

	assert.Equal(t, "  ", parser.Indentation())
}

func TestDefaultParserPutsChildrenInTheRightPlace(t *testing.T) {
	reader := &mockRuneReader{}
	parser := &DefaultParser{}
	content := []rune("%div\n  %div\n    %div\n  %div\n%div")
	i := 0

	for ; i < len(content); i += 1 {
		reader.On("ReadRune").Return(content[i], 1, nil).Once()
	}
	reader.On("ReadRune").Return('\000', 0, errors.New(""))

	doc, e := parser.Parse(reader)
	if ok := assert.Nil(t, e); !ok {
		return
	}

	if ok := assert.Equal(t, 2, len(doc.Nodes)); !ok {
		return
	}

	dn := doc.Nodes[0].(*TagNode)
	assert.Equal(t, "div", dn.Name)
	assert.Equal(t, "", dn.Id)
	assert.Equal(t, 0, len(dn.Classes))
	assert.Equal(t, 0, len(dn.Attrs))
	if ok := assert.Equal(t, 2, len(dn.Children)); !ok {
		return
	}

	dn = doc.Nodes[0].(*TagNode).Children[0].(*TagNode)
	assert.Equal(t, "div", dn.Name)
	assert.Equal(t, "", dn.Id)
	assert.Equal(t, 0, len(dn.Classes))
	assert.Equal(t, 0, len(dn.Attrs))
	if ok := assert.Equal(t, 1, len(dn.Children)); !ok {
		return
	}
}

func TestDefaultParserRecognizesSaticContentOnLine(t *testing.T) {
	reader := &mockRuneReader{}
	parser := &DefaultParser{}
	content := []rune("Bite my shiny, metal ass!")
	i := 0

	for ; i < len(content); i += 1 {
		reader.On("ReadRune").Return(content[i], 1, nil).Once()
	}
	reader.On("ReadRune").Return('\000', 0, errors.New(""))

	doc, e := parser.Parse(reader)
	if ok := assert.Nil(t, e); !ok {
		return
	}

	if ok := assert.Equal(t, 1, len(doc.Nodes)); !ok {
		return
	}

	sn := doc.Nodes[0].(*StaticLineNode)
	assert.Equal(t, "Bite my shiny, metal ass!", sn.Content)
}
