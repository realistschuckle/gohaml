package gohaml

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestDoesNotParseSomethingThatIsntAnExclamationMark(t *testing.T) {
	input := strings.NewReader("abcdefg")
	p := &docTypeParser{}
	node, err := p.Parse(input)

	assert.Nil(t, node)
	assert.NotNil(t, err)
	assert.Equal(t, input.Len(), 7)
}

func TestDoesNotParseSingleExclamationPointAndLeavesReaderUnread(t *testing.T) {
	input := strings.NewReader("!")
	p := &docTypeParser{}
	node, err := p.Parse(input)

	assert.Nil(t, node)
	assert.NotNil(t, err)
	assert.Equal(t, input.Len(), 1)
}

func TestDoesNotParseDoubleExclamationPoint(t *testing.T) {
	input := strings.NewReader("!!")
	p := &docTypeParser{}
	node, err := p.Parse(input)

	assert.Nil(t, node)
	assert.NotNil(t, err)
	assert.Equal(t, input.Len(), 2)
}

func TestDoesParsesTripleExclamationPointWithEmptySpecification(t *testing.T) {
	input := strings.NewReader("!!!")
	p := &docTypeParser{}
	node, err := p.Parse(input)
	docTypeNode := node.(*DocTypeNode)

	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, input.Len(), 0)
	assert.Empty(t, docTypeNode.Specification)
}

func TestDoesParsesTripleExclamationPointWithSpecification(t *testing.T) {
	input := strings.NewReader("!!! some goofy thing")
	p := &docTypeParser{}
	node, err := p.Parse(input)
	docTypeNode := node.(*DocTypeNode)

	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, input.Len(), 0)
	assert.Equal(t, docTypeNode.Specification, "some goofy thing")
}

func TestDoesNotConsumeNewline(t *testing.T) {
	input := strings.NewReader("!!! some goofy thing\n")
	p := &docTypeParser{}
	node, err := p.Parse(input)
	docTypeNode := node.(*DocTypeNode)

	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, input.Len(), 1, "Should not consume newline")
	assert.Equal(t, docTypeNode.Specification, "some goofy thing")
}

func TestNextReturnsAnEmptySlice(t *testing.T) {
	p := &docTypeParser{}

	assert.Nil(t, p.Next())
}
