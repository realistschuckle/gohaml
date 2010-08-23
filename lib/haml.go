// The gohaml package contains a HAML parser similar to the one found at http://www.haml-lang.com.
//
//You can find the specifics about this implementation at http://github.com/realistschuckle/gohaml.
package gohaml

import (
)

/*
Engine provides the template interpretation functionality to convert a HAML template into its
corresponding tag-based representation.

Available options are:
  engine.Options["autoclose"] = true|false, default true

The Options field contains the values to modify the way that the engine produces the markup.

The Indentation field contains the string used by the engine to perform indentation.

The IncludeCallback field contains the callback invoked by the gohaml engine to process other files
included through the %include extension.
*/
type Engine struct {
	Autoclose bool
	Indentation string
	IncludeCallback func(string, map[string]interface{}) string
	ast *tree
}

// NewEngine returns a new Engine with the given input.
func NewEngine(input string) (engine *Engine) {
	engine = &Engine{true, "\t", nil, parser.parse(input)}
	return
}

// Render interprets the HAML supplied to the NewEngine method.
func (self *Engine) Render(scope map[string]interface{}) (output string) {
	output = self.ast.resolve(scope, self.Indentation, self.Autoclose)
	return
}
