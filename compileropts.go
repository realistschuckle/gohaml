package gohaml

/*
CompilerOptions provides the structure of options available to customize the
workings of the gohaml Compiler.
*/
type CompilerOpts struct {
	// The rune that should wrap element attributes.
	AttributeWrapper rune

	// A list of tag names that should be automatically self-closed if they
	// have no content.
	Autoclose []string

	// Whether to include CDATA sections around javascript and css blocks when
	// using the :javascript or :css filters.
	Cdata bool

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

	// If set to true, all tags are treated as if both whitespace removal
	// options were present.
	RemoveWhitespace bool

	// Whether or not attribute hashes and Go sections designated by = or ~
	// should be evaluated.
	SuppressEval bool

	// If set to true, Haml makes no attempt to properly indent or format the
	// HTML output.
	Ugly bool
}
