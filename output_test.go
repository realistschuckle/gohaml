package gohaml

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStaticOuptutReturnsWhatsInContent(t *testing.T) {
	contents := []string{"one", "two", "three"}
	scope := make(map[string]interface{})
	for i := 0; i < len(contents); i += 1 {
		output := StaticOutput{contents[i]}
		o, e := output.Render(scope)
		assert.Nil(t, e)
		assert.Equal(t, contents[i], o)
	}
}
