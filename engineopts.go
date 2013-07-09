package gohaml

import "reflect"

/*
EngineOptions provides the structure of options available to customize the
workings of the gohaml Engine.

Has most of the same properties found in HAML::Options. You can read that
documentation at http://haml.info/docs/yardoc/Haml/Options.html.
*/
type EngineOptions struct {
	// Contains the rune that should wrap element attributes.
	AttributeWrapper rune

	// Contains a list of tag names that should be automatically self-closed if
	// they have no content.
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

	// The name of the Haml file being parsed.
	Filename string

	// Determines the output format.
	Format string

	// If set to true, Haml will convert underscores to hyphens in all Custom
	// Data Attributes.
	HyphenateDataAttributes bool

	// The line offset of the Haml template being parsed.
	Line int

	// The mime type that the rendered document will be served with.
	MimeType string

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
	Compiler: nil
	Encoding: "UTF-8"
	EscapeAttributes: true
	EscapeHtml: false
	FileName: ""
	Format: "html5"
	HyphenateDataAttributes: true
	Line: 0
	MimeType: "text/html"
	Parser: nil
	RemoveWhitespace: false
	SuppressEval: false
	Ugly: true
*/
func DefaultEngineOptions() (opt EngineOptions) {
	closers := []string{"meta", "img", "link", "br", "hr", "input", "area", "param", "col", "base"}
	opt = EngineOptions{'\'', closers, false, nil, "UTF-8", true, false, "", "html5", true, 0, "text/html", nil, false, false, true}
	return
}
