// The gohaml package contains a HAML parser similar to the one found at
// http://www.haml-lang.com.
//
// You can find the specifics about this implementation at
// http://github.com/realistschuckle/gohaml.
package gohaml

import (
	"strings"
)

/*
Engine provides the template interpretation functionality to convert a HAML
template into its corresponding tag-based representation.
*/
type Engine struct {
	options *EngineOptions
	doc     CompiledDocument
}

/*
Node defines the common interface for nodes that exist in a ParsedDocument.
*/
type Node interface {
	Accept(HamlCompiler)
}

/*
ParsedDocument defines the interface required for the functionality to contain
a parsed document generated from the HamlParser.
*/
type ParsedDocument interface {
	Accept(HamlCompiler)
}

/*
HamlParser defines the interface required for the functionality to parse HAML
input.
*/
type HamlParser interface {
	Parse() (ParsedDocument, error)
}

/*
NodeParser defines the interface required for functionality to parse nodes for
a HamlParser.
*/
type NodeParser interface {
	Parse(*strings.Reader) (Node, error)
	Next() []NodeParser
}

/*
CompiledDocument defines the interface required for the functionality that
provides the compiled HAML document.
*/
type CompiledDocument interface {
	Render(map[string]interface{}) (string, error)
}

/*
HamlCompiler defines the interface required for the functionality to compile
ParsedDocuments into CompiledDocuments.
*/
type HamlCompiler interface {
	Compile() (CompiledDocument, error)
	VisitDocType(*DocTypeNode)
	VisitTag(*TagNode)
}

/*
NewEngine returns a new Engine configured by the specified options that will
render the input.

If EngineOptions is nil, then the method will configure the returned Engine
with the result of DefaultEngineOptions.
*/
func NewEngine(input string, options *EngineOptions) (e *Engine, err error) {
	startStates := []NodeParser{&docTypeParser{}}
	p := &parser{strings.NewReader(input), startStates}

	var pdoc ParsedDocument
	if pdoc, err = p.Parse(); err != nil {
		return
	}

	c := &compiler{pdoc, options, nil}

	var cdoc CompiledDocument
	if cdoc, err = c.Compile(); err != nil {
		return
	}

	e = &Engine{options, cdoc}
	return
}

/*
Render interprets the HAML supplied to the NewEngine method.

If scope is nil, then the Engine will render without any local bindings.
*/
func (self *Engine) Render(scope map[string]interface{}) (s string, e error) {
	if scope == nil {
		scope = make(map[string]interface{})
	}
	s, e = self.doc.Render(scope)
	return
}

// func (self *Engine) Compiler() (compiler HamlCompiler) {
// 	compiler = self.options.compiler
// 	return
// }

// func (self *Engine) Indentation() (indentation string) {
// 	indentation = ""
// 	return
// }

// func (self *Engine) Options() (options *EngineOptions) {
// 	options = self.options
// 	return
// }

// func (self *Parser) Parser() (parser HamlParser) {
// 	parser = self.options.parser
// 	return
// }
