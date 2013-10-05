package gohaml

import (
	"github.com/realistschuckle/gohaml/compiler"
	"github.com/realistschuckle/gohaml/parser"
	"reflect"
)

/*
EngineOptions provides the structure of options available to customize the
workings of the gohaml Engine.

Has most of the same properties found in HAML::Options. You can read that
documentation at http://haml.info/docs/yardoc/Haml/Options.html.
*/
type EngineOptions struct {
	// The rune that should wrap element attributes.
	AttributeWrapper rune

	// A list of tag names that should be automatically self-closed if they
	// have no content.
	Autoclose []string

	// Whether to include CDATA sections around javascript and css blocks when
	// using the :javascript or :css filters.
	Cdata bool

	// The compiler class to use. Must convert to HamlCompiler.
	CompilerClass reflect.Type

	// The encoding to use for the HTML output.
	Encoding string

	// Sets whether or not to escape HTML-sensitive characters in attributes.
	EscapeAttributes bool

	// Sets whether or not to escape HTML-sensitive characters in script.
	EscapeHtml bool

	// Determines the output format.
	Format string

	// If set to true, Haml will convert underscores to hyphens in all Custom
	// Data Attributes.
	HyphenateDataAttributes bool

	// The parser class to use. Must convert to HamlParser.
	ParserClass reflect.Type

	// If set to true, all tags are treated as if both whitespace removal
	// options were present.
	RemoveWhitespace bool

	// Whether or not attribute hashes and Ruby scripts designated by = or ~
	// should be evaluated.
	SuppressEval bool

	// If set to true, Haml makes no attempt to properly indent or format the
	// HTML output.
	Ugly bool
}

/*
Returns an EngineOptions with the default values set.

The default values are:

	AttributeWrapper: '\''
	Autoclose: ["meta", "img", "link", "br", "hr", "input", "area", "param", "col", "base"]
	Cdata: false
	CompilerClass: nil
	Encoding: "UTF-8"
	EscapeAttributes: true
	EscapeHtml: false
	Format: "html5"
	HyphenateDataAttributes: true
	ParserClass: nil
	RemoveWhitespace: false
	SuppressEval: false
	Ugly: true
*/
func DefaultEngineOptions() (opt EngineOptions) {
	var p parser.DefaultParser
	parserClass := reflect.TypeOf(p)
	var c compiler.DefaultCompiler
	compilerClass := reflect.TypeOf(c)
	closers := []string{"meta", "img", "link", "br", "hr", "input", "area", "param", "col", "base"}
	opt = EngineOptions{
		'\'',          // AttributeWrapper rune
		closers,       // Autoclose []string
		false,         // Cdata bool
		compilerClass, // CompilerClass reflect.Type
		"UTF-8",       // Encoding string
		true,          // EscapeAttributes bool
		false,         // EscapeHtml bool
		"html5",       // Format string
		true,          // HyphenateDataAttributes bool
		parserClass,   // ParserClass reflect.Type
		false,         // RemoveWhitespace bool
		false,         // SuppressEval bool
		true,          // Ugly bool
	}

	return
}
