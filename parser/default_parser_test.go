package parser

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
)

func DefaultParserReadsUntilErrorReturnedFromReadRune(t *testing.T) {
	reader := &mockRuneReader{}
	parser := &DefaultParser{}
	content := "html\n  head\n    title Hello\n  body This is great!"
	i := 0

	for ; i < len(content); i += 1 {
		reader.On("ReadRune").Return(content[i], 1, nil).Once()
	}
	reader.On("ReadRune").Return('\000', 0, errors.New(""))

	parser.Parse(reader)

	reader.AssertExpectations(t)
	assert.Equal(t, len(content) + 1, i)
}
