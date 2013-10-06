package compiler

import (
	"github.com/realistschuckle/testify/mock"
)

type mockCompiledOutput struct {
	mock.Mock
}

func (self *mockCompiledOutput) Render(scope map[string]interface{}) (string, error) {
	args := self.Mock.Called(scope)
	return args.String(0), args.Error(1)
}
