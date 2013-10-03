// The gohaml package contains a HAML parser similar to the one found at
// http://www.haml-lang.com.
//
// You can find the specifics about this implementation at
// http://github.com/realistschuckle/gohaml.
package gohaml

type ParsedDocument struct {
}

type CompiledDocument struct {
}

type HamlParser interface {
	Parse(string, *EngineSettings) (ParsedDocument, error)
}

type HamlCompiler interface {
	Compile(ParsedDocument, *EngineOptions) (CompiledDocument, error)
}

/*
Engine provides the template interpretation functionality to convert a HAML
template into its corresponding tag-based representation.
*/
type Engine struct {
	settings *EngineSettings
	options  *EngineOptions
}

/*
NewEngine returns a new Engine configured by the specified options that will
render the input.

If EngineOptions is nil, then the method will configure the returned Engine
with the result of DefaultEngineOptions.
*/
func NewEngine(input string, options *EngineOptions) (e *Engine, err error) {

	return
}

/*
Render interprets the HAML supplied to the NewEngine method.

If scope is nil, then the Engine will render without any local bindings.
*/
func (self *Engine) Render(scope map[string]interface{}) (s string, e error) {
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
