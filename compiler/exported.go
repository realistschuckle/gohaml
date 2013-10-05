package compiler

import (
	"github.com/realistschuckle/gohaml/parser"
)

type CompiledDocument struct {
}

type HamlCompiler interface {
	Compile(parser.ParsedDocument, *CompilerOptions) (CompiledDocument, error)
}

type DefaultCompiler struct {
}

func (self *DefaultCompiler) Compile(parser.ParsedDocument, *CompilerOptions) (d CompiledDocument, e error) {
	return
}
