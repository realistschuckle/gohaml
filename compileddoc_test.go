package gohaml

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestCompiledDocRenderCallsRenderOnOutputs(t *testing.T) {
	scope := make(map[string]interface{})
	cdoc := CompiledDoc{}
	cdoc.Outputs = []CompiledOutput{
		&mockCompiledOutput{},
		&mockCompiledOutput{},
		&mockCompiledOutput{},
		&mockCompiledOutput{},
	}
	for i := 0; i < len(cdoc.Outputs); i += 1 {
		mock := cdoc.Outputs[i].(*mockCompiledOutput)
		mock.On("Render", scope).Return(string(i), nil)
	}

	cdoc.Render(scope)

	for i := 0; i < len(cdoc.Outputs); i += 1 {
		mock := cdoc.Outputs[i].(*mockCompiledOutput)
		mock.AssertCalled(t, "Render", scope)
	}
}

func TestCompiledDocRenderConcatenatesOutputAndReturnsIt(t *testing.T) {
	scope := make(map[string]interface{})
	cdoc := CompiledDoc{}
	cdoc.Outputs = []CompiledOutput{
		&mockCompiledOutput{},
		&mockCompiledOutput{},
		&mockCompiledOutput{},
		&mockCompiledOutput{},
	}
	for i := 0; i < len(cdoc.Outputs); i += 1 {
		mock := cdoc.Outputs[i].(*mockCompiledOutput)
		mock.On("Render", scope).Return(strconv.Itoa(i), nil)
	}

	output, _ := cdoc.Render(scope)

	assert.Equal(t, "0123", output)
}
