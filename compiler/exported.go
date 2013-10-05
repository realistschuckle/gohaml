package compiler

import (
	p "github.com/realistschuckle/gohaml/parser"
)

type CompiledDoc struct {
}

type HamlCompiler interface {
	Compile(p.ParsedDoc, CompilerOpts) (CompiledDoc, error)
}

type DefaultCompiler struct {
}

func (self *DefaultCompiler) Compile(input p.ParsedDoc, opts CompilerOpts) (doc CompiledDoc, e error) {
	doc = CompiledDoc{}
	return
}
