package compiler

import (
	"github.com/realistschuckle/gohaml/parser"
)

type CompiledDocument struct {
}

type HamlCompiler interface {
	Compile(parser.ParsedDocument, *CompilerOptions) (CompiledDocument, error)
}
