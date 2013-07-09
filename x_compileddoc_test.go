package gohaml

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompiledDocRenderReturnsValueInBuffer(t *testing.T) {
	s := "Ain't this grand?"
	buf := bytes.NewBufferString(s)
	doc := compiledDoc{buf}
	scope := make(map[string]interface{})

	output, e := doc.Render(scope)

	assert.Nil(t, e)
	assert.Equal(t, output, s)
}
